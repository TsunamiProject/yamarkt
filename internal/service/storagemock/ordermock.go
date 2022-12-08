package storagemock

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderStorage struct {
}

func (os *OrderStorage) CreateOrder(ctx context.Context, login string, orderID string) (err error) {
	if login == "test" && orderID == "2377225624" {
		return nil
	}
	err = errors.New("internal server error")
	return err
}

func (os *OrderStorage) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	if login == "test" {
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
	}
	err = errors.New("internal server error")
	return nil, err
}

func (os *OrderStorage) UpdateOrder(ctx context.Context, login string, oi models.OrderInfo) (err error) {
	return nil
}
