package handler

import (
	"context"
	"net/http"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type OrderServiceProvider interface {
	CreateOrder(ctx context.Context, login string, orderID string) error
	OrderList(ctx context.Context, login string) (ol []models.OrderList, err error)
}

type OrderHandler struct {
	service OrderServiceProvider
}

func NewOrderHandler(osp OrderServiceProvider) *OrderHandler {
	return &OrderHandler{service: osp}
}

func (oh OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

}

func (oh OrderHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	
}