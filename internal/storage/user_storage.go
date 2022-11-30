package storage

import "context"

func (ps *PostgresStorage) Auth(ctx context.Context, login string, pass string) error {
	return nil
}

func (ps *PostgresStorage) Register(ctx context.Context, login string, pass string) error {
	return nil
}
