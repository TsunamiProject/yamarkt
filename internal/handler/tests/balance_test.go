package handlers__test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/handler/servicemock"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

func TestHandler_GetCurrentBalance(t *testing.T) {
	tests := []struct {
		name                  string
		inputMethod           string
		inputEndpoint         string
		inputLogin            string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "#1. GetCurrentBalance handler. Positive test",
			inputMethod:           http.MethodGet,
			inputEndpoint:         "/api/user/balance",
			inputLogin:            "test",
			expectedStatusCode:    http.StatusOK,
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json",
		},
		{
			name:                  "#2. GetCurrentBalance handler. Negative test",
			inputMethod:           http.MethodGet,
			inputEndpoint:         "/api/user/balance",
			inputLogin:            "wrong",
			expectedStatusCode:    http.StatusInternalServerError,
			expectedHeader:        "Content-Type",
			expectedHeaderContent: "application/json",
		},
	}

	bs := &servicemock.BalanceServiceMock{}
	bh := handler.NewBalanceHandler(bs)

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, nil)
			claims := map[string]interface{}{
				"login": tCase.inputLogin,
			}
			_, jwtToken, _ := config.TokenAuth.Encode(claims)
			request.Header.Set("Authorization", "Bearer "+jwtToken)
			w := httptest.NewRecorder()
			bh.GetCurrentBalance(w, request)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
			assert.Equal(t, tCase.expectedHeaderContent, w.Header().Get(tCase.expectedHeader))
		})
	}
}

func TestHandler_CreateWithdrawal(t *testing.T) {
	tests := []struct {
		name                  string
		inputMethod           string
		inputEndpoint         string
		inputLogin            string
		inputBody             string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:               "#1. CreateWithdrawal handler. New withdrawal. Positive",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/balance/withdraw",
			inputLogin:         "test",
			inputBody:          `{"order": "6532528541", "sum":123}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "#2. CreateWithdrawal handler. New withdrawal. Negative: no funds",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/balance/withdraw",
			inputLogin:         "test",
			inputBody:          `{"order": "6532528541", "sum":12345}`,
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name:               "#3. CreateWithdrawal handler. New withdrawal. Negative: order already exists",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/balance/withdraw",
			inputLogin:         "test2",
			inputBody:          `{"order": "5660169110", "sum":1}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "#4. CreateWithdrawal handler. New withdrawal. Negative",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/balance/withdraw",
			inputLogin:         "wrong",
			inputBody:          `{"order": "5639730687", "sum":123}`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	bs := &servicemock.BalanceServiceMock{}
	bh := handler.NewBalanceHandler(bs)

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, bytes.NewBufferString(tCase.inputBody))
			claims := map[string]interface{}{
				"login": tCase.inputLogin,
			}
			_, jwtToken, _ := config.TokenAuth.Encode(claims)
			request.Header.Set("Authorization", "Bearer "+jwtToken)
			w := httptest.NewRecorder()
			w.Header().Set("Authorization", "Bearer "+tCase.inputLogin)
			bh.CreateWithdrawal(w, request)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_GetWithdrawalList(t *testing.T) {
	tests := []struct {
		name                   string
		inputMethod            string
		inputEndpoint          string
		inputLogin             string
		inputBody              string
		expectedStatusCode     int
		expectedStruct         []models.WithdrawalList
		expectedResponseBody   string
		expectedHeader1        string
		expectedHeaderContent1 string
	}{
		{
			name:                   "#1. GetWithdrawal handler. Positive",
			inputMethod:            http.MethodGet,
			inputEndpoint:          "/api/user/withdrawals",
			inputLogin:             "test",
			expectedStatusCode:     http.StatusOK,
			expectedHeader1:        "Content-Type",
			expectedHeaderContent1: "application/json",
			expectedStruct: []models.WithdrawalList{
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
		},
		{
			name:                   "#2. GetWithdrawal handler. Negative: no content",
			inputMethod:            http.MethodGet,
			inputEndpoint:          "/api/user/withdrawals",
			inputLogin:             "no-records",
			expectedStatusCode:     http.StatusNoContent,
			expectedHeader1:        "Content-Type",
			expectedHeaderContent1: "application/json",
		},
		{
			name:                   "#3. GetWithdrawal handler. Negative",
			inputMethod:            http.MethodGet,
			inputEndpoint:          "/api/user/withdrawals",
			inputLogin:             "wrong",
			expectedStatusCode:     http.StatusInternalServerError,
			expectedHeader1:        "Content-Type",
			expectedHeaderContent1: "application/json",
		},
	}

	bs := &servicemock.BalanceServiceMock{}
	bh := handler.NewBalanceHandler(bs)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			tCase := tCase
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, nil)
			claims := map[string]interface{}{
				"login": tCase.inputLogin,
			}
			_, jwtToken, _ := config.TokenAuth.Encode(claims)
			request.Header.Set("Authorization", "Bearer "+jwtToken)
			w := httptest.NewRecorder()
			w.Header().Set("Authorization", "Bearer "+tCase.inputLogin)
			bh.GetWithdrawalList(w, request)
			b, err := json.Marshal(tCase.expectedStruct)
			if err != nil {
				log.Printf("json.Marshal error: %v", err)
				return
			}
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
			assert.Equal(t, tCase.expectedHeaderContent1, w.Header().Get(tCase.expectedHeader1))
			if tCase.inputLogin == "test" {
				assert.Equal(t, string(b), strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
