package storage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
)

func (ps *PostgresStorage) Auth(ctx context.Context, login string, pass string) error {
	var dbPassword string
	err := ps.PostgresQL.QueryRowContext(ctx, userPasswordQuery, login).Scan(&dbPassword)
	if err != nil {
		log.Printf("error while getting user pass for auth: %s", err)
		return customErr.ErrUserDoesNotExist
	}

	if dbPassword != pass {
		log.Printf("wrong password recieved. login: %s", login)
		return customErr.ErrWrongPassword
	}
	return nil
}

func (ps *PostgresStorage) Register(ctx context.Context, login string, pass string) error {
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error while initializing tx from user register: %s", err)
		return err
	}
	defer tx.Rollback()

	{
		//create new user tx
		_, err = tx.Exec(createNewUserQuery, login, pass)
		var txErr *pgconn.PgError

		if errors.As(err, &txErr) && txErr.Code == pgerrcode.UniqueViolation {
			log.Printf("create new user tx err: %s. login: %s", customErr.ErrUserAlreadyExists, login)
			err = tx.Rollback()
			if err != nil {
				log.Printf("error while rollback transaction after err in creating new user")
			}
			return customErr.ErrUserAlreadyExists
		}

		if err != nil {
			log.Printf("error while executing create new user transaction: %s", err)
			err = tx.Rollback()
			if err != nil {
				log.Printf("error while rollback transaction after err in creating new user")
			}
			return err
		}
	}

	{
		//create balance row for new user
		_, err = tx.Exec(createUserBalanceQuery, login)
		var txErr *pgconn.PgError

		if errors.As(err, &txErr) && txErr.Code == pgerrcode.UniqueViolation {
			log.Printf("create new user balance err: %s. login: %s", customErr.ErrUserAlreadyExists, login)
			err = tx.Rollback()
			if err != nil {
				log.Printf("error while rollback transaction after err in creating new user balance")
			}
			return customErr.ErrUserAlreadyExists
		}

		if err != nil {
			log.Printf("error while executing create new user balance transaction: %s", err)
			err = tx.Rollback()
			if err != nil {
				log.Printf("error while rollback transaction after err in creating new user balance")
			}
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error while committing create new user transaction: %s", err)
	}

	return err
}
