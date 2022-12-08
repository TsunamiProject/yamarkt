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

func (bs *BalanceService) CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) (err error) {
	err = bs.balanceStorage.CreateWithdrawal(ctx, login, w)
	return err
}

func (bs *BalanceService) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	wl, err = bs.balanceStorage.GetWithdrawalList(ctx, login)
	return wl, err
}

func (bs *BalanceService) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	cb, err = bs.balanceStorage.GetCurrentBalance(ctx, login)
	return cb, nil
}
