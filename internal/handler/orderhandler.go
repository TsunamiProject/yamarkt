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

// CreateOrder takes on enter user login from Authentication header and orderID from request payload
// and returns status codes: 200 - on order already exist,
// 202 - on order accepted to process,
// 400 - on bad request,
// 401 - on unauthorized user,
// 409 - on order already exist,
// 422 - on wrong orderID format,
// 500 - on internal server error
func (oh OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	//getting token string from Authentication header
	tokenString := jwtauth.TokenFromHeader(r)
	//decoding token string to jwtToken instance
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("CreateOrder handler. Error while decoding token string to jwtToken: %s",
			err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	//getting login from jwtToken
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("CreateOrder handler. Error while getting login from claims: %s", err)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	login := fmt.Sprintf("%v", claims)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("CreateOrder handler. Error while reading request body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	stringBody := fmt.Sprintf("%s", body)
	//getting orderID from response body string
	_, err = strconv.ParseInt(stringBody, 10, 0)
	if err != nil {
		errString := fmt.Sprintf("CreateOrder handler. Error while converting request body to int: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnprocessableEntity)
	}

	//validating orderID via Luhn algorithm
	err = goluhn.Validate(stringBody)
	if err != nil {
		errString := fmt.Sprintf("CreateOrder handler. Luhn validating error: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnprocessableEntity)
		return
	}

	//calling CreateOrder service method
	err = oh.service.CreateOrder(ctx, login, stringBody)
	switch {
	case err != nil && errors.Is(err, customErr.ErrOrderAlreadyExists):
		log.Printf("CreateOrder handler. Login: %s: Order already exists", login)
		w.WriteHeader(http.StatusOK)
	case err != nil && errors.Is(err, customErr.ErrOrderCreatedByAnotherLogin):
		log.Printf("CreateOrder handler. Login: %s: Order created by another login", login)
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		log.Printf("CreateOrder handler. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		log.Printf("CreateOrder. Order %s created by login: %s: Error: %s", stringBody, login, err)
		w.WriteHeader(http.StatusAccepted)
	}
}

// OrderList takes on enter user login from Authentication header
// and returns:
// status code 200: user withdrawals list: json [
//
//	    {
//	        "number": "9278923470",
//	        "status": "PROCESSED",
//	        "accrual": 500,
//	        "uploaded_at": "2020-12-10T15:15:45+03:00"
//	    },
//	...
//	]
//
// 204: on no orders from user,
// 401 - on unauthorized user,
// 500 - on internal server error
func (oh OrderHandler) OrderList(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	//getting token string from Authentication header
	tokenString := jwtauth.TokenFromHeader(r)
	//decoding token string to jwtToken instance
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("Order List handler. "+
			"Error while decoding token string to jwtToken: %s",
			err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	//getting login from jwtToken
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("OrderList handler. Error while getting login from claims: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	login := fmt.Sprintf("%v", claims)

	//calling OrderList service method
	orderList, err := oh.service.OrderList(ctx, login)
	w.Header().Set("Content-Type", "application/json")
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoOrders):
		log.Printf("OrderList handler. Login: %s: No orders", login)
		w.WriteHeader(http.StatusNoContent)
	case err != nil:
		log.Printf("OrderList handler. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		log.Printf("OrderList handler. Login: %s", login)
		log.Printf("OrderList. Output: %s", orderList)
		json.NewEncoder(w).Encode(orderList)
	}
}
