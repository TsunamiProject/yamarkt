package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Credentials struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type Withdrawal struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type WithdrawalList struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type CurrentBalance struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn int             `json:"withdrawn"`
}

type OrderInfo struct {
	Order  string `json:"order"`
	Status string `json:"status"`
}

type OrderList struct {
	Accrual    decimal.Decimal `json:"accrual"`
	UploadedAt time.Time       `json:"uploaded_at"`
}

type AccrualJSON struct {
	Order string `json:"order"`
}
