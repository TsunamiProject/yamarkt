package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	//_, claims, err := jwtauth.FromContext(r.Context())
	//if err != nil {
	//	log.Printf("error while getting claims from new withdrawal request context: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//login, ok := claims["login"].(string)
	//if !ok {
	//	log.Printf("error while getting login from claims in new withdrawal handler: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

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
	if errors.Is(err, customErr.ErrNoFunds) {
		w.WriteHeader(http.StatusPaymentRequired)
	} else if errors.Is(err, customErr.ErrWithdrawalOrderAlreadyExist) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (bh BalanceHandler) GetWithdrawalList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	//_, claims, err := jwtauth.FromContext(r.Context())
	//if err != nil {
	//	log.Printf("error while getting claims from get withdrawal list request context: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//login, ok := claims["login"].(string)
	//if !ok {
	//	log.Printf("error while getting login from claims in get withdrawals list handler: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
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
	if errors.Is(err, customErr.ErrNoWithdrawals) {
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
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()

	//_, claims, err := jwtauth.FromContext(r.Context())
	//if err != nil {
	//	log.Printf("error while getting claims from get current balance request context: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//login, ok := claims["login"].(string)
	//if !ok {
	//	log.Printf("error while getting login from claims in get current balance handler: %s", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

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
