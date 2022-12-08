package service__test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/models"
	"github.com/TsunamiProject/yamarkt/internal/service"
	"github.com/TsunamiProject/yamarkt/internal/service/storagemock"
)

func TestHandler_GetCurrentBalance(t *testing.T) {
	tests := []struct {
		name                 string
		inputLogin           string
		expectedStruct       models.CurrentBalance
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:       "#1. Balance service. GetCurrentBalance. Positive",
			inputLogin: "test",
			expectedStruct: models.CurrentBalance{
				Current:   decimal.NewFromFloatWithExponent(123.123, -2),
				Withdrawn: decimal.NewFromFloatWithExponent(12, -2),
			},
			expectedError: nil,
		},
		{
			name:           "#2. Balance service. GetCurrentBalance. Negative",
			inputLogin:     "test2",
			expectedStruct: models.CurrentBalance{},
			expectedError:  errors.New("internal server error"),
		},
	}

	s := &storagemock.BalanceStorage{}
	bs := service.NewBalanceService(s)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			ec, err := bs.GetCurrentBalance(ctx, tCase.inputLogin)
			log.Println(err)
			assert.Equal(t, tCase.expectedError, err)
			assert.Equal(t, tCase.expectedStruct, ec)
		})
	}
}

func TestHandler_CreateWithdrawal(t *testing.T) {
	tests := []struct {
		name                 string
		inputLogin           string
		inputStruct          models.Withdrawal
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:       "#1. Balance service. Create withdrawal. Positive",
			inputLogin: "test",
			inputStruct: models.Withdrawal{
				Order: "2377225624",
				Sum:   decimal.NewFromFloatWithExponent(42, -2),
			},
			expectedError: nil,
		},
		{
			name:       "#2. Balance service. Create withdrawal. Negative",
			inputLogin: "test",
			inputStruct: models.Withdrawal{
				Order: "",
				Sum:   decimal.NewFromFloatWithExponent(42, -2),
			},
			expectedError: errors.New("internal server error"),
		},
	}

	s := &storagemock.BalanceStorage{}
	bs := service.NewBalanceService(s)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			err := bs.CreateWithdrawal(ctx, tCase.inputLogin, tCase.inputStruct)
			assert.Equal(t, tCase.expectedError, err)
		})
	}
}

func TestHandler_GetWithdrawalsList(t *testing.T) {
	tests := []struct {
		name                 string
		inputLogin           string
		inputStruct          []models.WithdrawalList
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:       "#1. Balance service. GetWithdrawalList. Positive",
			inputLogin: "test",
			inputStruct: []models.WithdrawalList{
				{
					Order:       "123456",
					Sum:         decimal.NewFromFloatWithExponent(123.123, -2),
					ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
				},
				{
					Order:       "654321",
					Sum:         decimal.NewFromFloatWithExponent(312.321, -2),
					ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
				},
			},
			expectedError: nil,
		},
		{
			name:       "#2. Balance service. GetWithdrawalList. Negative",
			inputLogin: "test2",
			inputStruct: []models.WithdrawalList{
				{
					Order:       "123456",
					Sum:         decimal.NewFromFloatWithExponent(123.123, -2),
					ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
				},
				{
					Order:       "654321",
					Sum:         decimal.NewFromFloatWithExponent(312.321, -2),
					ProcessedAt: time.Date(2022, time.December, 7, 22, 44, 10, 0, time.Local),
				},
			},
			expectedError: errors.New("internal server error"),
		},
	}

	s := &storagemock.BalanceStorage{}
	bs := service.NewBalanceService(s)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			ec, err := bs.GetWithdrawalList(ctx, tCase.inputLogin)
			assert.Equal(t, tCase.expectedError, err)
			if tCase.inputLogin == "test" {
				assert.Equal(t, tCase.inputStruct, ec)
			}
		})
	}
}
