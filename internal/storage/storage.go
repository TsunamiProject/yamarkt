package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/TsunamiProject/yamarkt/internal/config"
)

type PostgresStorage struct {
	PostgresQL *sql.DB
}

func NewPostgresStorage(databaseDsn string) (*PostgresStorage, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, config.StorageContextTimeout)
	defer cancel()

	pdb, err := connectToPostgresKernel(databaseDsn)
	if err != nil {
		return nil, err
	}

	_, err = pdb.Exec(usersTableQuery)
	if err != nil {
		log.Printf("error while creating users table: %s", err)
		return nil, err
	}

	_, err = pdb.Exec(balanceTableQuery)
	if err != nil {
		log.Printf("error while creating balance table: %s", err)
		return nil, err
	}

	_, err = pdb.Exec(ordersTableQuery)
	if err != nil {
		log.Printf("error while creating orders table: %s", err)
		return nil, err
	}

	_, err = pdb.Exec(withdrawalsTableQuery)
	if err != nil {
		log.Printf("error while creating withdrawals table: %s", err)
		return nil, err
	}
	return &PostgresStorage{PostgresQL: pdb}, nil
}

func connectToPostgresKernel(databaseDsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, fmt.Errorf("connection error to Postgres kernel with creds: %s. %s", databaseDsn, err)
	}
	return db, nil
}

func (ps *PostgresStorage) CloseConnection() error {
	err := ps.PostgresQL.Close()
	if err != nil {
		return err
	}
	return nil
}
