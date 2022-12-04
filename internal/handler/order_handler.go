package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/TsunamiProject/yamarkt/internal/config"
	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	stringBody := fmt.Sprintf("%s", body)
	_, err = strconv.ParseInt(stringBody, 10, 0)
	if err != nil {
		log.Printf("Error while converting request body to int: %s", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = goluhn.Validate(stringBody)
	if err != nil {
		log.Printf("Error while validating request body vai Luhn algo: %s", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	//TODO: get login from request
	login := ""

	err = oh.service.CreateOrder(ctx, login, stringBody)
	if errors.As(err, &customErr.ErrOrderAlreadyExists) {
		w.WriteHeader(http.StatusOK)
	} else if errors.As(err, &customErr.ErrOrderCreatedByAnotherLogin) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}

func (oh OrderHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	//TODO: get login from request
	login := ""
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	orderList, err := oh.service.OrderList(ctx, login)
	if errors.As(err, &customErr.ErrNoOrders) {
		log.Printf("no orders for %s login", login)
		w.WriteHeader(http.StatusNoContent)
	} else if errors.As(err, &customErr.ErrUnauthorizedUser) {
		log.Printf("unauthorized login: %s", login)
		w.WriteHeader(http.StatusUnauthorized)
	} else if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orderList)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}
