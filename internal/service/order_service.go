package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	err := os.storage.CreateOrder(ctx, login, orderID)
	if err != nil {
		return err
	}
	accOrder := models.AccrualJSON{Order: orderID}
	accOrderJSON, err := json.Marshal(accOrder)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", os.AccrualURL+"/api/orders", bytes.NewBuffer(accOrderJSON))
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error while making request to accrual service: %s", err)
		return err
	}
	defer resp.Body.Close()

	//TODO: worker for updating order info?

	return nil
}

func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}
