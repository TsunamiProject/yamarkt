package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/shopspring/decimal"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

func (ps *PostgresStorage) CreateWithdrawal(ctx context.Context, login string, withdrawal models.Withdrawal) error {
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error while creating tx instance: %s", err)
	}
	defer tx.Rollback()
	{
		var balanceWithdrawn decimal.Decimal
		var balanceCurrent decimal.Decimal
		err = ps.PostgresQL.QueryRow(getUserWithdrawnInfoQuery, login).Scan(&balanceCurrent, &balanceWithdrawn)
		if err != nil {
			log.Printf("error while scanning user withdrawn info query result: %s", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback error after getting user withdrawn info: %s", rollbackErr)
			}
			return err
		}
		if decimal.Decimal.Cmp(withdrawal.Sum, balanceCurrent) > 0 {
			log.Printf("%s", customErr.ErrNoFunds)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback error after getting user withdrawn info: %s", rollbackErr)
			}
			return customErr.ErrNoFunds
		}

		balanceWithdrawn = balanceWithdrawn.Add(withdrawal.Sum)
		balanceCurrent = balanceCurrent.Sub(withdrawal.Sum)

		_, err = ps.PostgresQL.Exec(updateUserWithdrawalBalanceQuery, login, balanceCurrent, balanceWithdrawn)
		if err != nil {
			log.Printf("error while updating user withdrawals info: %s", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback error after updating user withdrawn info: %s", rollbackErr)
			}
			return err
		}
	}
	{
		_, err = ps.PostgresQL.Exec(updateUserWithdrawalBalanceQuery, withdrawal.Order, login, withdrawal.Sum)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				log.Printf("%s", customErr.ErrWithdrawalOrderAlreadyExist)
				return customErr.ErrWithdrawalOrderAlreadyExist
			}
			log.Printf("error while updating user withdrawal balance: %s", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback error after updating user withdrawal balance: %s", rollbackErr)
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("error while committing withdrawal update tx: %s", err)
	}
	return err
}

func (ps *PostgresStorage) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	rows, err := ps.PostgresQL.QueryContext(ctx, getUserWithdrawalsQuery)
	if err != nil {
		log.Printf("err while getting user withdrawals: %s", err)
		return wl, err
	}
	defer rows.Close()
	withdrawalsList := models.WithdrawalList{}
	for rows.Next() {
		err = rows.Scan(&withdrawalsList.Order, &withdrawalsList.Sum, &withdrawalsList.ProcessedAt)
		if err != nil {
			log.Printf("error while scanning user withdrawals row: %s", err)
			return wl, err
		}
		wl = append(wl, withdrawalsList)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("error on iteration scan in get user withdrawals list: %s", err)
	}

	if len(wl) == 0 {
		log.Printf("%s", customErr.ErrNoWithdrawals)
		err = customErr.ErrNoWithdrawals
	}
	return wl, err
}

func (ps *PostgresStorage) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	err = ps.PostgresQL.QueryRowContext(ctx, getUserWithdrawnInfoQuery, login).Scan(&cb.Current, &cb.Withdrawn)
	if err != nil {
		log.Printf("error while getting user withdrawn info: %s", err)
	}
	return cb, err
}
