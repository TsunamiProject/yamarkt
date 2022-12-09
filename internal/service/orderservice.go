package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, login string, orderID string) error
	OrderList(ctx context.Context, login string) (ol []models.OrderList, err error)
	UpdateOrder(ctx context.Context, login string, oi models.OrderInfo) (err error)
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

	//creating sync.WaitGroup instance for handling update order status goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go UpdateOrderStatus(&wg, os.storage, os.AccrualURL, login, orderID)
	wg.Wait()
	return err
}

//OrderList service returns order list by authorized user
func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}

//UpdateOrderStatus service for sending requests to accrual system
func UpdateOrderStatus(wg *sync.WaitGroup, orderStorage OrderStorage, accrualURL string, login string, orderID string) {
	defer wg.Done()
	//collecting request to accrual system
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/orders/%s", accrualURL, orderID), nil)
	for {
		//collecting http client
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		//makeing request to accrual system
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("UpdateOrderStatus service. Error while making request to accrual system: %s", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusConflict {
			return
		}
		if resp.StatusCode == http.StatusOK {
			oi := models.OrderInfo{}
			//decoding response from accrual system to OrderInfo struct
			err = json.NewDecoder(resp.Body).Decode(&oi)
			if err != nil {
				log.Printf("UpdateOrderStatus service. Error while unmarshalling resp from accrual service: %s", err)
				return
			}

			//creating context from background context
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			//calling UpdateOrder storage mehod
			err = orderStorage.UpdateOrder(ctx, login, oi)
			if err != nil {
				log.Printf("UpdateOrderStatus service. Error while updating order :%s", err)
				return
			}
			if oi.Status == "INVALID" || oi.Status == "PROCESSED" {
				log.Printf("UpdateOrderStatus service. Order %s has updated status to %s", oi.Order, oi.Status)
				return
			}
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			//parsing Retry-After header to timeout variable for making sleep on Retry-After header value
			timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				log.Printf("UpdateOrderStatus service. Error while converting Retry-After header to int: %s", err)
				return
			}
			time.Sleep(time.Duration(timeout) * 1000 * time.Millisecond)
		}
	}
}
