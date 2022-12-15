package service__test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/models"
	"github.com/TsunamiProject/yamarkt/internal/service"
	"github.com/TsunamiProject/yamarkt/internal/service/storagemock"
)

func TestHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name          string
		inputLogin    string
		inputOrderID  string
		expectedError error
	}{
		{
			name:          "#1. Order service. CreateOrder. Positive",
			inputLogin:    "test",
			inputOrderID:  "6721797030",
			expectedError: nil,
		},
	}

	s := &storagemock.OrderStorage{}
	os := service.NewOrderService(s, "")

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			err := os.CreateOrder(ctx, tCase.inputLogin, tCase.inputOrderID)
			assert.Equal(t, err, err)
		})
	}
}

func TestHandler_OrderList(t *testing.T) {
	tests := []struct {
		name           string
		inputLogin     string
		expectedStruct []models.OrderList
		expectedError  error
	}{
		{
			name:       "#1. Order service. Order list. Positive",
			inputLogin: "test",
			expectedStruct: []models.OrderList{
				{
					Number:     "4289742787",
					Status:     "PROCESSED",
					Accrual:    decimal.NewFromFloatWithExponent(123, -2),
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
				{
					Number:     "5322351601",
					Status:     "PROCESSING",
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
				{
					Number:     "6721797030",
					Status:     "INVALID",
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
			},
			expectedError: nil,
		},
		{
			name:       "#1. Order service. Order list. Negative",
			inputLogin: "test2",
			expectedStruct: []models.OrderList{
				{
					Number:     "4289742787",
					Status:     "PROCESSED",
					Accrual:    decimal.NewFromFloatWithExponent(123, -2),
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
				{
					Number:     "5322351601",
					Status:     "PROCESSING",
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
				{
					Number:     "6721797030",
					Status:     "INVALID",
					UploadedAt: time.Date(2022, time.December, 8, 01, 43, 12, 0, time.Local),
				},
			},
			expectedError: errors.New("internal server error"),
		},
	}

	s := &storagemock.OrderStorage{}
	os := service.NewOrderService(s, "")

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			ol, err := os.OrderList(ctx, tCase.inputLogin)
			assert.Equal(t, tCase.expectedError, err)
			if tCase.inputLogin == "test" {
				assert.Equal(t, tCase.expectedStruct, ol)
			}
		})
	}
}
