package storage

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

func (ps *PostgresStorage) CreateWithdrawal(ctx context.Context, login string, withdrawal models.Withdrawal) error {
	return nil
}

func (ps *PostgresStorage) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	return nil, nil
}

func (ps *PostgresStorage) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	return models.CurrentBalance{}, nil
}
