package service

import (
	"context"

	"github.com/rs/zerolog/log"

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
	if err != nil {
		log.Printf("CreateOrder service. Create order storage error: %s", err)
		return err
	}
	//accOrder := models.AccrualJSON{Order: orderID}
	//accOrderJSON, err := json.Marshal(accOrder)
	//if err != nil {
	//	log.Printf("error while marshalling json for accrual service: %s", err)
	//	return err
	//}
	//req, err := http.NewRequest("POST", os.AccrualURL+"/api/orders", bytes.NewBuffer(accOrderJSON))
	//if err != nil {
	//	return err
	//}
	//client := http.Client{
	//	Timeout: 5 * time.Second,
	//}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Printf("error while making request to accrual service: %s", err)
	//	return err
	//}
	//defer resp.Body.Close()
	//log.Printf("%d", resp.StatusCode)

	return err
}

//OrderList service returns order list by authorized user
func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}
