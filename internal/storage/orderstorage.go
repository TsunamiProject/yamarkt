package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/TsunamiProject/yamarkt/internal/config"
	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

//CreateOrder storage method for placing new user order with user login and orderID
func (ps *PostgresStorage) CreateOrder(ctx context.Context, login string, orderID string) (err error) {
	//sending create new user order query
	_, err = ps.PostgresQL.ExecContext(ctx, createNewUserOrderQuery, orderID, login)
	if err != nil {
		var pgErr *pgconn.PgError
		//if order already exists in database
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			dbLogin := ""
			//sending get user by orderID query
			err = ps.PostgresQL.QueryRowContext(ctx, getUserByOrderIDQuery, orderID).Scan(&dbLogin)
			if err != nil {
				log.Printf("CreateOrder. Error while scanning get user by order query result: %s", err)
				return err
			}

			if dbLogin != login {
				return customErr.ErrOrderCreatedByAnotherLogin
			}
			return customErr.ErrOrderAlreadyExists
		}
		return err
	}
	return err
}

//OrderList storage method for getting orders from db by user login
func (ps *PostgresStorage) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	//sending get user order list query
	rows, err := ps.PostgresQL.QueryContext(ctx, getUserOrdersListQuery, login)
	if err != nil {
		log.Printf("OrderList. Error on get user orders list query: %s", err)
		return ol, err
	}
	defer rows.Close()
	orderList := models.OrderList{}

	//scanning rows
	for rows.Next() {
		err = rows.Scan(&orderList.Number, &orderList.Status, &orderList.Accrual, &orderList.UploadedAt)
		if err != nil {
			log.Printf("OrderList. Error while scanning row: %s", err)
			return ol, err
		}
		//appending order to orderList
		ol = append(ol, orderList)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("OrderList. Error on iteration scan in get user orders list query: %s", err)
		return ol, err
	}

	if len(ol) == 0 {
		log.Printf("OrderList. No orders by current login: %s", login)
		err = customErr.ErrNoOrders
	}
	return ol, err
}

//UpdateOrder storage method for update order status and accrual info in created order
func (ps *PostgresStorage) UpdateOrder(ctx context.Context, login string, oi models.OrderInfo) (err error) {
	//creating transaction instance
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("UpdateOrder. Error while creating tx instance: %s", err)
		return err
	}
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			log.Printf("UpdateOrder. Transaction rollback error: %s", err)
		}
	}(tx)

	//sending update user order query with new info about order status and accrual
	_, err = tx.Exec(updateUserOrderQuery, login, oi.Order, oi.Status, oi.Accrual)
	if err != nil {
		log.Printf("UpdateOrder. Error while updating user order info: %s", err)
		return err
	}

	if oi.Accrual.GreaterThan(decimal.NewFromInt(0)) {
		var dbBalance decimal.Decimal
		//sending query for getting user balance
		err = ps.PostgresQL.QueryRow(getUserBalanceQuery, login).Scan(&dbBalance)
		if err != nil {
			log.Printf("UpdateOrder. Error while scanning get user balance query result: %s", err)
			return err
		}
		log.Printf("UpdateOrder. Before. Login: %s: Balance: %s", login, dbBalance)
		dbBalance = oi.Accrual.Add(dbBalance)
		log.Printf("UpdateOrder. After. Login: %s: Balance : %s", login, dbBalance)

		//sending update user balance query to update user actual balance
		_, err = tx.Exec(updateUserBalanceQuery, login, dbBalance)
		if err != nil {
			log.Printf("UpdateOrder. Error while updating user balance: %s", err)
			return err
		}
	}
	//committing transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("UpdateOrder. Error while committing tx on update user order")
	}
	return err
}

//GetUnprocessedOrdersList storage method for getting unprocessed orders from db
func (ps *PostgresStorage) GetUnprocessedOrdersList(ctx context.Context) (ol []models.UnprocessedOrdersList, err error) {
	//sending get unprocessed orders list query
	rows, err := ps.PostgresQL.QueryContext(ctx, getUnprocessedOrdersQuery, config.InvalidOrderStatus,
		config.ProcessedOrderStatus)
	if err != nil {
		log.Printf("GetUnprocessedOrdersList. Error on get unprocessed orders list query: %s", err)
		return ol, err
	}
	defer rows.Close()
	unprocessedOrderList := models.UnprocessedOrdersList{}

	//scanning rows
	for rows.Next() {
		err = rows.Scan(&unprocessedOrderList.Number, &unprocessedOrderList.Login)
		if err != nil {
			log.Printf("GetUnprocessedOrdersList. Error while scanning row: %s", err)
			return ol, err
		}
		//appending order to orderList
		ol = append(ol, unprocessedOrderList)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("GetUnprocessedOrdersList. Error on iteration scan in get unprocessed orders list query: %s", err)
		return ol, err
	}

	if len(ol) == 0 {
		log.Printf("GetUnprocessedOrdersList. No unprocessed orders")
		err = customErr.ErrNoUnprocessedOrders
	}
	return ol, err
}
