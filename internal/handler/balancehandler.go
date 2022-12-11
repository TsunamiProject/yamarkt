package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"

	"github.com/TsunamiProject/yamarkt/internal/config"
	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type BalanceServiceProvider interface {
	CreateWithdrawal(ctx context.Context, login string, w models.Withdrawal) error
	GetWithdrawalList(ctx context.Context, login string) (wl []models.WithdrawalList, err error)
	GetCurrentBalance(ctx context.Context, login string) (cb models.CurrentBalance, err error)
}

type BalanceHandler struct {
	service BalanceServiceProvider
}

func NewBalanceHandler(bhp BalanceServiceProvider) *BalanceHandler {
	return &BalanceHandler{service: bhp}
}

// CreateWithdrawal takes on enter json like {"order":"2377225624", "sum":751}, login from Authentication header
//and returns: status codes: 200 - on correct withdrawal order processing,
//401 - on unauthorized user,
//402 - on insufficient funds on user balance,
//422 - on incorrect orderID received
//500 - on internal server error
func (bh BalanceHandler) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	withdrawal := models.Withdrawal{}
	//decoding request body to withdrawal instance
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		log.Printf("CreateWithdrawal handler. Request body decoding error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	//validating orderID via Luhn algorithm
	err = goluhn.Validate(withdrawal.Order)
	if err != nil {
		errString := fmt.Sprintf("CreateWithdrawal handler. Luhn validating error: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnprocessableEntity)
	}

	//getting token string from Authentication header
	tokenString := jwtauth.TokenFromHeader(r)
	//decoding token string to jwtToken instance
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("CreateWithdrawal handler. "+
			"Error while decoding token string to jwtToken in create withdrawal handler: %s",
			err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	//getting login from jwtToken
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("CreateWithdrawal handler. Error while getting login from claims: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	login := fmt.Sprintf("%v", claims)

	//calling CreateWithdrawal service method
	err = bh.service.CreateWithdrawal(ctx, login, withdrawal)
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoFunds):
		log.Printf("CreateWithdrawal handler. Login: %s: No funds", login)
		w.WriteHeader(http.StatusPaymentRequired)
	case err != nil && errors.Is(err, customErr.ErrWithdrawalOrderAlreadyExist):
		log.Printf("CreateWithdrawal handler. Login: %s: Withdrawal already exist", login)
		w.WriteHeader(http.StatusUnprocessableEntity)
	case err != nil:
		log.Printf("CreateWithdrawal handler. Login: %s: Error: %s", login, err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

// GetWithdrawalList takes on enter user login from Authentication header
//and returns:
//status code 200: user withdrawals list: json [
//      {
//          "order": "2377225624",
//          "sum": 500,
//          "processed_at": "2020-12-09T16:09:57+03:00"
//      },
//  ...
//  ]
//401 - on unauthorized user,
//500 - on internal server error
func (bh BalanceHandler) GetWithdrawalList(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	//getting token string from Authentication header
	tokenString := jwtauth.TokenFromHeader(r)
	//decoding token string to jwtToken instance
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("GetWithdrawalList handler. Error while decoding token string to jwtToken: %s",
			err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	//getting login from jwtToken
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("GetWithdrawalList handler. Error while getting login from claims: %s", err)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	login := fmt.Sprintf("%v", claims)
	//calling GetWithdrawalList service method
	w.Header().Set("Content-Type", "application/json")
	withdrawalList, err := bh.service.GetWithdrawalList(ctx, login)
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoWithdrawals):
		log.Printf("GetWithdrawalList handler. No orders from: %s", login)
		w.WriteHeader(http.StatusNoContent)
	case err != nil:
		log.Printf("GetWithdrawalList handler. Login: %s : Error: %s", login, err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(withdrawalList)
	}
}

// GetCurrentBalance takes on enter user login from Authentication header
//and returns: status code 200: user balance info: json {"current":123, "withdrawn":12},
//401 - on unauthorized user,
//500 - on internal server errors
func (bh BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	//getting token string from Authentication header
	tokenString := jwtauth.TokenFromHeader(r)
	//decoding token string to jwtToken instance
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("GetCurrentBalance handler. Error while decoding token string to jwtToken: %s",
			err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	//getting login from jwtToken
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("GetCurrentBalance handler. Error while getting login from claims: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusUnauthorized)
		return
	}
	login := fmt.Sprintf("%v", claims)
	w.Header().Set("Content-Type", "application/json")
	//calling GetCurrentBalance service method
	currentBalance, err := bh.service.GetCurrentBalance(ctx, login)
	switch {
	case err != nil:
		log.Printf("GetCurrentBalance handler. Login: %s: Error: %s", login, err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		log.Printf("GetCurrentBalance handler. Login: %s: User balance info: %s", login, currentBalance)
		err = json.NewEncoder(w).Encode(currentBalance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
