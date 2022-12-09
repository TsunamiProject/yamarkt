package handlers__test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/handler/servicemock"
)

func TestHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name               string
		inputMethod        string
		inputEndpoint      string
		inputBody          string
		expectedStatusCode int
		inputLogin         string
	}{
		{
			name:               "#1. CreateOrder handler. Positive",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test",
			inputBody:          "5871772181",
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "#2. CreateOrder handler. Bad orderID",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test",
			inputBody:          "12wrongorderid34",
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "#3. CreateOrder handler. Bad orderID: Luhn validation",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test",
			inputBody:          "123",
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "#4. CreateOrder handler. Order already exist",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test2",
			inputBody:          "5871772181",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "#5. CreateOrder handler. Order already created by another user",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test3",
			inputBody:          "5871772181",
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "#6. CreateOrder handler. Negative",
			inputMethod:        http.MethodPost,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "wrong",
			inputBody:          "5871772181",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	os := &servicemock.OrderServiceMock{}
	oh := handler.NewOrderHandler(os)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, bytes.NewBufferString(tCase.inputBody))
			jwtToken := jwt.New()
			jwtToken.Set(`login`, tCase.inputLogin)
			requestContext := jwtauth.NewContext(request.Context(), jwtToken, nil)
			request = request.WithContext(requestContext)
			w := httptest.NewRecorder()
			w.Header().Set("Authorization", "Bearer "+tCase.inputLogin)
			oh.CreateOrder(w, request)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_OrderList(t *testing.T) {
	tests := []struct {
		name               string
		inputMethod        string
		inputEndpoint      string
		expectedStatusCode int
		inputLogin         string
	}{
		{
			name:               "#1. OrderList handler. Positive",
			inputMethod:        http.MethodGet,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "#2. OrderList handler. No orders",
			inputMethod:        http.MethodGet,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "test2",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "#3. OrderList handler. Negative",
			inputMethod:        http.MethodGet,
			inputEndpoint:      "/api/user/orders",
			inputLogin:         "wrong",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	os := &servicemock.OrderServiceMock{}
	oh := handler.NewOrderHandler(os)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, nil)
			jwtToken := jwt.New()
			jwtToken.Set(`login`, tCase.inputLogin)
			requestContext := jwtauth.NewContext(request.Context(), jwtToken, nil)
			request = request.WithContext(requestContext)
			w := httptest.NewRecorder()
			w.Header().Set("Authorization", "Bearer "+tCase.inputLogin)
			oh.OrderList(w, request)
			body, _ := io.ReadAll(w.Body)
			log.Printf("%s", body)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
		})
	}
}
