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
