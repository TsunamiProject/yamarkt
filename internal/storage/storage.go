package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/TsunamiProject/yamarkt/internal/config"
)

type PostgresStorage struct {
	PostgresQL *sql.DB
}

//NewPostgresStorage method for collecting postgres instance with given databaseDsn
func NewPostgresStorage(databaseDsn string) (*PostgresStorage, error) {
	//collecting context
	ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
	defer cancel()

	//connecting to postgres database
	pdb, err := connectToPostgres(databaseDsn)
	if err != nil {
		return nil, err
	}

	//checking connection with database
	err = pdb.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("connection with database not alive: %s", err)
	}

	//creating users table
	_, err = pdb.Exec(usersTableQuery)
	if err != nil {
		log.Printf("Storage. Error while creating users table: %s", err)
		return nil, err
	}

	//creating balance table
	_, err = pdb.Exec(balanceTableQuery)
	if err != nil {
		log.Printf("Storage. Error while creating balance table: %s", err)
		return nil, err
	}

	//creating orders table
	_, err = pdb.Exec(ordersTableQuery)
	if err != nil {
		log.Printf("Storage. Error while creating orders table: %s", err)
		return nil, err
	}

	//creating withdrawals table
	_, err = pdb.Exec(withdrawalsTableQuery)
	if err != nil {
		log.Printf("Storage. Error while creating withdrawals table: %s", err)
		return nil, err
	}
	return &PostgresStorage{PostgresQL: pdb}, nil
}

//connectToPostgres makes connection with postgres database via pgx driver
func connectToPostgres(databaseDsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, fmt.Errorf("connection error to Postgres kernel with creds: %s. %s", databaseDsn, err)
	}
	return db, nil
}

//CloseConnection closes connection with postgres database
func (ps *PostgresStorage) CloseConnection() error {
	err := ps.PostgresQL.Close()
	if err != nil {
		return err
	}
	return nil
}
