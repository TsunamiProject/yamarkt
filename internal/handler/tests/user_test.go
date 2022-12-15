package handlers__test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/handler/servicemock"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name                  string
		inputMethod           string
		inputEndpoint         string
		inputBody             string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:                  "#1. Register handler. Positive",
			inputMethod:           "POST",
			inputEndpoint:         "/api/user/register",
			inputBody:             `{ "login": "test", "password": "qwerty" }`,
			expectedStatusCode:    http.StatusOK,
			expectedHeader:        "Authorization",
			expectedHeaderContent: "Bearer",
		},
		{
			name:               "#2. Register handler. Bad json received",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/register",
			inputBody:          `{ "logi, "passwor"" }`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "#3. Register handler. User already exist",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/register",
			inputBody:          `{ "login": "test2", "password": "qwerty" }`,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "#4. Register handler. Negative",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/register",
			inputBody:          `{ "login": "wrong", "password": "wrong" }`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	us := &servicemock.UserServiceMock{}
	uh := handler.NewUserHandler(us)

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, bytes.NewBufferString(tCase.inputBody))
			w := httptest.NewRecorder()
			uh.Register(w, request)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Contains(t, w.Header().Get(tCase.expectedHeader), tCase.expectedHeaderContent)
			}

		})
	}
}

func TestHandler_Auth(t *testing.T) {
	tests := []struct {
		name                  string
		inputMethod           string
		inputEndpoint         string
		inputBody             string
		expectedStatusCode    int
		expectedResponseBody  string
		expectedHeader        string
		expectedHeaderContent string
	}{
		{
			name:               "#1. Auth handler. Positive",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/login",
			inputBody:          `{ "login": "test", "password": "qwerty" }`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "#2. Auth handler. Wrong password",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/login",
			inputBody:          `{ "login": "test", "password": "wrong" }`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "#3. Auth handler. Wrong login",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/login",
			inputBody:          `{ "login": "test3", "password": "qwerty123" }`,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "#4. Auth handler. Negative",
			inputMethod:        "POST",
			inputEndpoint:      "/api/user/register",
			inputBody:          `{ "login": "wrong", "password": "wrong" }`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	us := &servicemock.UserServiceMock{}
	uh := handler.NewUserHandler(us)

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			request := httptest.NewRequest(tCase.inputMethod, tCase.inputEndpoint, bytes.NewBufferString(tCase.inputBody))
			w := httptest.NewRecorder()
			uh.Auth(w, request)
			assert.Equal(t, tCase.expectedStatusCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Contains(t, w.Header().Get(tCase.expectedHeader), tCase.expectedHeaderContent)
			}
		})
	}
}
