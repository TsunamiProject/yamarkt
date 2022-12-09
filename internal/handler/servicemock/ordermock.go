package servicemock

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderServiceMock struct {
}

func (os *OrderServiceMock) CreateOrder(ctx context.Context, login string, orderNum string) (err error) {
	switch {
	case login == "test" && orderNum == "5871772181":
		return nil
	case login == "test2" && orderNum == "5871772181":
		return customErr.ErrOrderAlreadyExists
	case login == "test3" && orderNum == "5871772181":
		return customErr.ErrOrderCreatedByAnotherLogin
	default:
		return errors.New("internal server error")
	}

}

func (os *OrderServiceMock) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	switch login {
	case "test":
		ol = []models.OrderList{
			{
				Number:     "4289742787",
				Status:     "PROCESSED",
				Accrual:    decimal.NewFromFloatWithExponent(123, -2),
				UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
			},
			{
				Number:     "5322351601",
				Status:     "PROCESSING",
				UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
			},
			{
				Number:     "6721797030",
				Status:     "INVALID",
				UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
			},
		}
		return ol, nil
	case "test2":
		return nil, customErr.ErrNoOrders
	default:
		return nil, errors.New("internal server error")
	}
}
