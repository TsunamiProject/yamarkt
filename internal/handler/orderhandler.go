package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"

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
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

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

	tokenString := jwtauth.TokenFromHeader(r)
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("error while decoding token string to jwtToken in create order handler: %s",
			err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("error while getting login from claims in create order handler: %s", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	login := fmt.Sprintf("%v", claims)

	err = oh.service.CreateOrder(ctx, login, stringBody)
	switch {
	case err != nil && errors.Is(err, customErr.ErrOrderAlreadyExists):
		w.WriteHeader(http.StatusOK)
	case err != nil && errors.Is(err, customErr.ErrOrderCreatedByAnotherLogin):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		log.Printf("create order service error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

func (oh OrderHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	tokenString := jwtauth.TokenFromHeader(r)
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("error while decoding token string to jwtToken in order list handler: %s",
			err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("error while getting login from claims in order list handler: %s", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	login := fmt.Sprintf("%v", claims)

	orderList, err := oh.service.OrderList(ctx, login)
	w.Header().Set("Content-Type", "application/json")
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoOrders):
		log.Printf("no orders for %s login", login)
		w.WriteHeader(http.StatusNoContent)
	case err != nil:
		log.Printf("order list service error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		log.Printf("order list: %s", orderList)
		err = json.NewEncoder(w).Encode(orderList)
		if err != nil {
			log.Printf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
