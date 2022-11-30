package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TsunamiProject/yamarkt/internal/config"
)

type PostgresStorage struct {
	PostgresQL *sql.DB
}

func NewPostgresStorage(databaseDsn string) (*PostgresStorage, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, config.StorageContextTimeout)
	defer cancel()

	_, err := connectToPostgresKernel(databaseDsn)
	if err != nil {
		return nil, err
	}
	//...
	return nil, nil
}

func connectToPostgresKernel(databaseDsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, fmt.Errorf("connection error to Postgres kernel with creds: %s", databaseDsn)
	}
	return db, nil
}
