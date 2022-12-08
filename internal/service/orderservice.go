package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

	go func() {
		err = UpdateOrderStatus(os.storage, os.AccrualURL, login, orderID)
		if err != nil {
			log.Printf("error while updating order status: %s", err)
		}
	}()
	return err
}

func (os *OrderService) OrderList(ctx context.Context, login string) (ol []models.OrderList, err error) {
	orderList, err := os.storage.OrderList(ctx, login)
	return orderList, err
}

func UpdateOrderStatus(orderStorage OrderStorage, accrualURL string, login string, orderID string) (err error) {
	for {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("%s/api/orders/%s", accrualURL, orderID), nil)
		if err != nil {
			log.Printf("error while constructing request to accrual service :%s", err)
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
		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusConflict {
			log.Printf("recieved status code from accrual service: %v", resp.StatusCode)
			return err
		}
		log.Printf("status code %v recieved from accrual service", resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			oi := models.OrderInfo{}
			err = json.NewDecoder(resp.Body).Decode(&oi)
			if err != nil {
				log.Printf("error while unmarshalling resp from accrual service: %s", err)
				return err
			}

			ctx, cancel := context.WithTimeout(req.Context(), config.StorageContextTimeout)
			defer cancel()
			err = orderStorage.UpdateOrder(ctx, login, oi)
			if err != nil {
				log.Printf("error while updating order :%s", err)
				return err
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
				return err
			}
			time.Sleep(time.Duration(timeout) * 1000 * time.Millisecond)
		}
	}
}
