package service

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type BalanceStorage interface {
	CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) error
	GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error)
	GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error)
}

type BalanceService struct {
	balanceStorage BalanceStorage
}

func NewBalanceService(bs BalanceStorage) *BalanceService {
	return &BalanceService{balanceStorage: bs}
}

func (bs *BalanceService) CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) error {
	return nil
}

func (bs *BalanceService) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	return nil, nil
}

func (bs *BalanceService) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	return models.CurrentBalance{}, nil
}
