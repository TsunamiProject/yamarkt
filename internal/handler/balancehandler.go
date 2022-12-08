package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-chi/jwtauth/v5"

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

func (bh BalanceHandler) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	withdrawal := models.Withdrawal{}
	err := json.NewDecoder(r.Body).Decode(&withdrawal)
	if err != nil {
		log.Printf("Request body decoding error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	err = goluhn.Validate(withdrawal.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	tokenString := jwtauth.TokenFromHeader(r)
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("error while decoding token string to jwtToken in create withdrawal handler: %s",
			err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("error while getting login from claims in create withdrawal handler: %s", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	login := fmt.Sprintf("%v", claims)

	err = bh.service.CreateWithdrawal(ctx, login, withdrawal)
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoFunds):
		w.WriteHeader(http.StatusPaymentRequired)
	case err != nil && errors.Is(err, customErr.ErrWithdrawalOrderAlreadyExist):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case err != nil:
		log.Printf("create withdrawal service error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func (bh BalanceHandler) GetWithdrawalList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	tokenString := jwtauth.TokenFromHeader(r)
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("error while decoding token string to jwtToken in get withdrawal list handler: %s",
			err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("error while getting login from claims in get withdrawal list  handler: %s", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	login := fmt.Sprintf("%v", claims)

	withdrawalList, err := bh.service.GetWithdrawalList(ctx, login)
	switch {
	case err != nil && errors.Is(err, customErr.ErrNoWithdrawals):
		w.WriteHeader(http.StatusNoContent)
	case err != nil:
		log.Printf("get withdrawal list service error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(withdrawalList)
	}
}

func (bh BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	tokenString := jwtauth.TokenFromHeader(r)
	jwtToken, err := config.TokenAuth.Decode(tokenString)
	if err != nil {
		errString := fmt.Sprintf("error while decoding token string to jwtToken in get current balance handler: %s",
			err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	claims, ok := jwtToken.Get("login")
	if !ok {
		errString := fmt.Sprintf("error while getting login from claims in get current balance handler: %s", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	login := fmt.Sprintf("%v", claims)
	w.Header().Set("Content-Type", "application/json")
	currentBalance, err := bh.service.GetCurrentBalance(ctx, login)
	switch {
	case err != nil:
		log.Printf("get current balance service error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(currentBalance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
