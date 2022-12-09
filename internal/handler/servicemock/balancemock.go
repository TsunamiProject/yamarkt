package servicemock

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type BalanceServiceMock struct {
}

func (bs *BalanceServiceMock) GetCurrentBalance(ctx context.Context, login string) (ec models.CurrentBalance, err error) {
	switch login {

	case "test":
		ec = models.CurrentBalance{
			Current:   decimal.NewFromFloatWithExponent(123.123, -2),
			Withdrawn: decimal.NewFromFloatWithExponent(12, -2),
		}
		return ec, nil

	default:
		return ec, errors.New("internal server error")
	}
}

func (bs *BalanceServiceMock) CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) (err error) {
	switch {
	case login == "test" && w.Order == "6532528541" && w.Sum.Equal(decimal.NewFromFloat(123)):
		return nil
	case login == "test" && w.Order == "6532528541" && w.Sum.GreaterThan(decimal.NewFromFloat(123)):
		return customErr.ErrNoFunds
	case login == "test2" && w.Order == "5660169110":
		return customErr.ErrWithdrawalOrderAlreadyExist
	default:
		return errors.New("internal server error")
	}
}

func (bs *BalanceServiceMock) GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error) {
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
	switch {
	case login == "test":
		return wl, nil
	case login == "no-records":
		return nil, customErr.ErrNoWithdrawals
	case login == "wrong":
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("internal server error")
	}
}
