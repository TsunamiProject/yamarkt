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

type UpdateOrderStorage interface {
	GetUnprocessedOrdersList(ctx context.Context) (ol []models.UnprocessedOrdersList, err error)
	UpdateOrder(ctx context.Context, login string, oi models.OrderInfo) (err error)
}

type UpdateOrderService struct {
	storage    UpdateOrderStorage
	AccrualURL string
}

func NewUpdateOrderService(os UpdateOrderStorage, accURL string) *UpdateOrderService {
	return &UpdateOrderService{
		storage:    os,
		AccrualURL: accURL,
	}
}

//UpdateOrderStatus service for sending requests to accrual system
func (uo *UpdateOrderService) UpdateOrderStatus(ctx context.Context, wg *sync.WaitGroup) {
	//collecting http client
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			time.Sleep(config.GetUnprocessedOrdersFrequency)
			//getting unprocessed orders from storage
			unprocessedOrderList, err := uo.GetUnprocessedOrdersList(ctx)
			if err != nil {
				log.Printf("UpdateOrderStatus. Error while getting unprocessed order list: %s", err)
				continue
			}
			for order := range unprocessedOrderList {
				//collecting request to accrual system
				req, err := http.NewRequestWithContext(ctx, http.MethodGet,
					fmt.Sprintf("%s%s%s", uo.AccrualURL, config.AccrualOrderStatusURN, unprocessedOrderList[order].Number), nil)
				if err != nil {
					log.Printf("UpdateOrderStatus service. Error while collecting http request: %s", err)
					continue
				}
				//making request to accrual system
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("UpdateOrderStatus service. Error while making request to accrual system: %s", err)
					continue
				}
				log.Printf("UpdateOrderStatus service. Received status code from accrual system: %d", resp.StatusCode)
				if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusConflict {
					continue
				}
				if resp.StatusCode == http.StatusTooManyRequests {
					//parsing Retry-After header to timeout variable for making sleep on Retry-After header value
					timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
					if err != nil {
						log.Printf("UpdateOrderStatus service. Error while converting Retry-After header to int: %s", err)
						time.Sleep(config.RetryAfterErrorDefaultTimeout)
					} else {
						time.Sleep(time.Duration(timeout) * time.Second)
					}
				}
				if resp.StatusCode == http.StatusOK {
					oi := models.OrderInfo{}
					//decoding response from accrual system to OrderInfo struct
					err = json.NewDecoder(resp.Body).Decode(&oi)
					if err != nil {
						log.Printf("UpdateOrderStatus service. Error while unmarshalling resp from accrual service: %s", err)
						continue
					}
					//log.Printf("UpdateOrderStatus service. Received order info: %s", oi)
					//creating context from parent context
					updateContext, cancel := context.WithTimeout(ctx, config.StorageContextTimeout)
					//calling UpdateOrder storage method
					err = uo.storage.UpdateOrder(updateContext, unprocessedOrderList[order].Login, oi)
					if err != nil {
						cancel()
						log.Printf("UpdateOrderStatus service. Error while updating order :%s", err)
						continue
					}
					cancel()
					if oi.Status == config.InvalidOrderStatus || oi.Status == config.ProcessedOrderStatus {
						log.Printf("UpdateOrderStatus service. Order %s has updated status to %s", oi.Order, oi.Status)
						continue
					}
				}
				err = resp.Body.Close()
				if err != nil {
					log.Printf("UpdateOrderStatus service. Error while closing response body: %s", err)
				}
			}
		}
	}
}

//GetUnprocessedOrdersList service for getting unprocessed orders from db
func (uo *UpdateOrderService) GetUnprocessedOrdersList(ctx context.Context) (ol []models.UnprocessedOrdersList, err error) {
	ol, err = uo.storage.GetUnprocessedOrdersList(ctx)
	return ol, err
}
