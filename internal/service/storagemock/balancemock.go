package storagemock

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type BalanceStorage struct {
}

func (bs *BalanceStorage) GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error) {
	switch {
	case login == "test":
		cb = models.CurrentBalance{
			Current:   decimal.NewFromFloatWithExponent(123.123, -2),
			Withdrawn: decimal.NewFromFloatWithExponent(12, -2),
		}
		return cb, nil
	default:
		err = errors.New("internal server error")
		return models.CurrentBalance{}, err
	}
}

func (bs *BalanceStorage) CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) (err error) {
	wch := models.Withdrawal{
		Order: "2377225624",
		Sum:   decimal.NewFromFloatWithExponent(42, -2),
	}
	if w.Sum.Equal(wch.Sum) && w.Order == wch.Order {
		return nil
	}
	err = errors.New("internal server error")
	return err
}

func (bs *BalanceStorage) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
	if login == "test" {
		wl = []models.WithdrawalList{
			{
				Order:       "123456",
				Sum:         decimal.NewFromFloatWithExponent(123.123, -2),
				ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
			},
			{
				Order:       "654321",
				Sum:         decimal.NewFromFloatWithExponent(312.321, -2),
				ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
			},
		}
		return wl, nil
	}
	err = errors.New("internal server error")
	return nil, err
}
