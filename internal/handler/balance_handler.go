package handler

import (
	"context"
	"net/http"

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

}

func (bh BalanceHandler) GetWithdrawalList(w http.ResponseWriter, r *http.Request) {

}

func (bh BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {

}
