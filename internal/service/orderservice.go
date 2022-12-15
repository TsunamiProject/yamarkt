package service

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, login string, orderID string) (err error)
	OrderList(ctx context.Context, login string) (ol []models.OrderList, err error)
}

type OrderService struct {
	storage OrderStorage
}

func NewOrderService(os OrderStorage, accURL string) *OrderService {
	return &OrderService{
		storage: os,
	}
}

//CreateOrder service for creating orders to accrual point by authorized user
func (os *OrderService) CreateOrder(ctx context.Context, login string, orderID string) (err error) {
	//calling CreateOrder postgres storage method for creating new order
	err = os.storage.CreateOrder(ctx, login, orderID)
	return err
}

//OrderList service returns order list by authorized user
func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}
