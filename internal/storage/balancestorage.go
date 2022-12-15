package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

//CreateWithdrawal storage method for creating new withdrawal by user login
func (ps *PostgresStorage) CreateWithdrawal(ctx context.Context, login string, withdrawal models.Withdrawal) error {
	//creating transaction instance
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("CreateWithdrawal. Error while creating tx instance: %s", err)
	}
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			log.Printf("CreateWithdrawal. Transaction rollback error: %s", err)
		}
	}(tx)

	var balanceWithdrawn decimal.Decimal
	var balanceCurrent decimal.Decimal
	//sending get user withdrawn query for get actual info about user balance and total withdrawn
	err = ps.PostgresQL.QueryRow(getUserWithdrawnInfoQuery, login).Scan(&balanceCurrent, &balanceWithdrawn)
	if err != nil {
		log.Printf("CreateWithdrawal. Error while scanning user withdrawn info query result: %s", err)
		return err
	}
	//if withdrawal order sum greater than actual user balance
	if withdrawal.Sum.GreaterThan(balanceCurrent) {
		//returning no funds error
		return customErr.ErrNoFunds
	}

	//increase total withdrawn info
	balanceWithdrawn = balanceWithdrawn.Add(withdrawal.Sum)
	log.Printf("CreateWithdrawal. Login: %s: Total withdrawn: %s", login, balanceWithdrawn)
	//decrease user current balance
	balanceCurrent = balanceCurrent.Sub(withdrawal.Sum)
	log.Printf("CreateWithdrawal. Login: %s: Current balance: %s", login, balanceCurrent)

	//sending create user withdrawal query for create new withdrawal if not exists
	_, err = tx.Exec(createUserWithdrawalQuery, withdrawal.Order, login, withdrawal.Sum)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			//returns withdrawal already exists error
			return customErr.ErrWithdrawalOrderAlreadyExist
		}
		log.Printf("CreateWithdrawal. Error while creating new user withdrawal: %s", err)
		return err
	}
	//sending update user balance after creating new withdrawal query
	_, err = tx.Exec(updateUserWithdrawalBalanceQuery, login, balanceCurrent, balanceWithdrawn)
	if err != nil {
		log.Printf("CreateWithdrawal. Error while updating user balance after creating new withdrawal: %s", err)
		return err
	}

	//committing transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateWithdrawal. Error while committing withdrawal update tx: %s", err)
	}
	return err
}

//GetWithdrawalList storage method for getting all user withdrawals
func (ps *PostgresStorage) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	//sending get user withdrawals query
	rows, err := ps.PostgresQL.QueryContext(ctx, getUserWithdrawalsQuery, login)
	if err != nil {
		log.Printf("GetWithdrawalList. Error while getting user withdrawals: %s", err)
		return wl, err
	}
	defer rows.Close()
	withdrawalsList := models.WithdrawalList{}
	//scanning rows
	for rows.Next() {
		err = rows.Scan(&withdrawalsList.Order, &withdrawalsList.Sum, &withdrawalsList.ProcessedAt)
		if err != nil {
			log.Printf("GetWithdrawalList. Error while scanning user withdrawals row: %s", err)
			return wl, err
		}
		//appending withdrawal to withdrawals list
		wl = append(wl, withdrawalsList)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("GetWithdrawalList. Error on iteration scan in get user withdrawals list: %s", err)
	}

	//checking that result not empty
	if len(wl) == 0 {
		err = customErr.ErrNoWithdrawals
	}
	return wl, err
}

//GetCurrentBalance storage method for getting user actual balance and total withdrawn info
func (ps *PostgresStorage) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	//sending get user withdrawn info query and scanning to struct
	err = ps.PostgresQL.QueryRowContext(ctx, getUserWithdrawnInfoQuery, login).Scan(&cb.Current, &cb.Withdrawn)
	if err != nil {
		log.Printf("GetCurrentBalance. Error while getting user withdrawn info: %s", err)
	}
	return cb, err
}
