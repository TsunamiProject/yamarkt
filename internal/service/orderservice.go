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

func (os *OrderService) CreateOrder(ctx context.Context, login string, orderID string) (err error) {
	err = os.storage.CreateOrder(ctx, login, orderID)
	if err != nil {
		log.Printf("create order storage error: %s", err)
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = UpdateOrderStatus(&wg, os.storage, os.AccrualURL, login, orderID)
		if err != nil {
			log.Printf("error while updating order status: %s", err)
		}
	}()
	wg.Wait()
	return err
}

func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}

func UpdateOrderStatus(wg *sync.WaitGroup, orderStorage OrderStorage, accrualURL string, login string, orderID string) (err error) {
	defer wg.Done()
	for {
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s/api/orders/%s", accrualURL, orderID), nil)
		//if err != nil {
		//	log.Printf("error while constructing request to accrual service :%s", err)
		//	return err
		//}
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error while making request to accrual service: %s", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusConflict {
			log.Printf("received status code from accrual service: %v", resp.StatusCode)
			return nil
		}
		log.Printf("status code %v received from accrual service", resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			oi := models.OrderInfo{}
			err = json.NewDecoder(resp.Body).Decode(&oi)
			if err != nil {
				log.Printf("error while unmarshalling resp from accrual service: %s", err)
				continue
			}

			ctx, cancel := context.WithTimeout(req.Context(), config.StorageContextTimeout)
			defer cancel()
			log.Printf("accrual:", oi.Accrual)
			err = orderStorage.UpdateOrder(ctx, login, oi)
			if err != nil {
				log.Printf("error while updating order :%s", err)
				continue
			}
			if oi.Status == "INVALID" || oi.Status == "PROCESSED" {
				log.Printf("order %s has updated status to %s", oi.Order, oi.Status)
				return nil
			}
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				log.Printf("error converting Retry-After to int:%s", err)
			}
			time.Sleep(time.Duration(timeout) * 1000 * time.Millisecond)
		}
	}
}
