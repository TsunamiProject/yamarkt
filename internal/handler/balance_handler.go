package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"

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

func (bh BalanceHandler) NewWithdrawal(w http.ResponseWriter, r *http.Request) {
	//TODO: get login from request
	login := ""

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	withdrawal := models.Withdrawal{}
	err = json.Unmarshal(body, &withdrawal)
	if err != nil {
		log.Printf("Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	err = goluhn.Validate(withdrawal.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	err = bh.service.CreateWithdrawal(ctx, login, withdrawal)
	if errors.As(err, &customErr.ErrNoFunds) {
		w.WriteHeader(http.StatusPaymentRequired)
	} else if errors.As(err, &customErr.ErrUnauthorizedUser) {
		w.WriteHeader(http.StatusUnauthorized)
	} else if errors.As(err, &customErr.ErrOrderAlreadyExists) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (bh BalanceHandler) GetWithdrawalList(w http.ResponseWriter, r *http.Request) {
	//TODO: get login from request
	login := ""

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	withdrawalList, err := bh.service.GetWithdrawalList(ctx, login)
	if errors.As(err, &customErr.ErrNoOrders) {
		w.WriteHeader(http.StatusNoContent)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(withdrawalList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (bh BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	//TODO: get login from request
	login := ""

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	currentBalance, err := bh.service.GetCurrentBalance(ctx, login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(currentBalance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
