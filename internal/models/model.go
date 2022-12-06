package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Credentials struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type Withdrawal struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

type WithdrawalList struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type CurrentBalance struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

type OrderInfo struct {
	Order   string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}

type OrderList struct {
	Number     string          `json:"number"`
	Status     string          `json:"status"`
	Accrual    decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt time.Time       `json:"uploaded_at"`
}

type AccrualJSON struct {
	Order string `json:"order"`
}
