package service

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, login string, orderID string) error
	OrderList(ctx context.Context, login string) (ol []models.OrderList, err error)
}

type OrderService struct {
	storage    OrderStorage
	AccrualURL string
}

func NewOrderService(os OrderStorage, accURL string) *OrderService {
	return &OrderService{
		storage:    os,
		AccrualURL: accURL,
	}
}

func (os *OrderService) CreateOrder(ctx context.Context, login string, orderID string) error {
	return nil
}

func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	return nil, nil
}
