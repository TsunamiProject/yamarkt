package storage

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

func (ps *PostgresStorage) CreateOrder(ctx context.Context, login string, orderID string) error {
	return nil
}

func (ps *PostgresStorage) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	return nil, nil
}
