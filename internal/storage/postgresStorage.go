package storage

import "database/sql"

type PostgresStorage struct {
	PostgresQL *sql.DB
}
