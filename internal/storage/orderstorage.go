package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

func (ps *PostgresStorage) CreateOrder(ctx context.Context, login string, orderID string) (err error) {
	_, err = ps.PostgresQL.ExecContext(ctx, createNewUserOrderQuery, orderID, login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			dbLogin := ""
			err = ps.PostgresQL.QueryRowContext(ctx, getUserByOrderIDQuery, orderID).Scan(&dbLogin)
			if err != nil {
				log.Printf("error while scanning get user by order query result: %s", err)
				return err
			}

			if dbLogin != login {
				log.Printf("%s", customErr.ErrOrderCreatedByAnotherLogin)
				return customErr.ErrOrderCreatedByAnotherLogin
			}

			log.Printf("login: %s, %s", login, customErr.ErrOrderAlreadyExists)
			return customErr.ErrOrderAlreadyExists
		}
	}
	return err
}

func (ps *PostgresStorage) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	rows, err := ps.PostgresQL.QueryContext(ctx, getUserOrdersListQuery, login)
	if err != nil {
		log.Printf("error on get user orders list query: %s", err)
		return ol, err
	}
	defer rows.Close()
	orderList := models.OrderList{}

	for rows.Next() {
		err = rows.Scan(&orderList.Number, &orderList.Status, &orderList.Accrual, &orderList.UploadedAt)
		if err != nil {
			log.Printf("error while scanning row: %s", err)
			return ol, err
		}
		ol = append(ol, orderList)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("error on iteration scan in get user orders list query: %s", err)
		return ol, err
	}

	if len(ol) == 0 {
		log.Printf("no orders by current login: %s", login)
		err = customErr.ErrNoOrders
	}

	return ol, err
}

func (ps *PostgresStorage) UpdateOrder(ctx context.Context, login string, oi models.OrderInfo) (err error) {
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error while creating tx instance: %s", err)
		return err
	}
	defer tx.Rollback()

	{
		_, err = ps.PostgresQL.Exec(updateUserOrderQuery, login, oi.Order, oi.Status, oi.Accrual)
		if err != nil {
			log.Printf("error while updating user order info: %s", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback error while updating user order info: %s", err)
			}
			return err
		}
	}
	{
		if decimal.Decimal.Cmp(oi.Accrual, decimal.NewFromInt(0)) > 0 {
			var dbBalance decimal.Decimal
			err = ps.PostgresQL.QueryRow(getUserBalanceQuery, login).Scan(&dbBalance)
			if err != nil {
				log.Printf("error while scanning get user balance query result: %s", err)
				rollbackErr := tx.Rollback()
				if err != nil {
					log.Printf("rollback error after getting user balance: %s", rollbackErr)
				}
				return err
			}

			_, err = ps.PostgresQL.Exec(updateUserBalanceQuery, login, oi.Accrual.Add(dbBalance))
			if err != nil {
				log.Printf("error while updating user balance: %s", err)
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					log.Printf("rollback error after updating user balance: %s", rollbackErr)
				}
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error while committing tx on update user order")
	}
	return err
}
