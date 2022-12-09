package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
)

//Auth method finds user in user table from credentials
func (ps *PostgresStorage) Auth(ctx context.Context, login string, pass string) error {
	var dbPassword string
	//getting user password by login
	err := ps.PostgresQL.QueryRowContext(ctx, userPasswordQuery, login).Scan(&dbPassword)
	//errors equality means that login doesn't exist in table
	if err == sql.ErrNoRows {
		log.Printf("Auth. Error while getting user pass for auth: %s", customErr.ErrUserDoesNotExist)
		return customErr.ErrUserDoesNotExist
	} else if err != nil {
		log.Printf("Auth. Error: %s", err)
		return err
	}

	if dbPassword != pass {
		log.Printf("Auth. Wrong password received. Login: %s", login)
		return customErr.ErrWrongPassword
	}
	return nil
}

//Register method creating new user in db by credentials user credentials
func (ps *PostgresStorage) Register(ctx context.Context, login string, pass string) error {
	//collecting transaction instance
	tx, err := ps.PostgresQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error while initializing tx from user register: %s", err)
		return err
	}
	defer tx.Rollback()
	{
		//making create new user transaction
		_, err = tx.Exec(createNewUserQuery, login, pass)
		var txErr *pgconn.PgError
		//checking user already exists in db
		if errors.As(err, &txErr) && txErr.Code == pgerrcode.UniqueViolation {
			log.Printf("Register. Create new user tx err: %s. Login: %s", customErr.ErrUserAlreadyExists, login)
			err = tx.Rollback()
			if err != nil {
				log.Printf("Register. Error while rollback transaction after err in creating new user")
			}
			return customErr.ErrUserAlreadyExists
		}
		if err != nil {
			log.Printf("Register. Error while executing create new user transaction: %s", err)
			err = tx.Rollback()
			if err != nil {
				log.Printf("Register. Error on rollback create new user transaction")
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

	//committing transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Register. Error while committing create new user transaction: %s", err)
	}
	return err
}
