package service__test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/models"
	"github.com/TsunamiProject/yamarkt/internal/service"
	"github.com/TsunamiProject/yamarkt/internal/service/storagemock"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name                 string
		inputLogin           string
		inputStruct          models.Credentials
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:       "#1. User service. Register. Positive",
			inputLogin: "test",
			inputStruct: models.Credentials{
				Login: "test",
				Pass:  "qwerty",
			},
			expectedError: nil,
		},
		{
			name:       "#2. User service. Register. Negative",
			inputLogin: "wrong",
			inputStruct: models.Credentials{
				Login: "wrong",
				Pass:  "wrong",
			},
			expectedError: errors.New("internal server error"),
		},
	}

	s := &storagemock.UserStorage{}
	us := service.NewUserService(s)

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			err := us.Register(ctx, tCase.inputStruct)
			assert.Equal(t, tCase.expectedError, err)
		})
	}
}

func TestHandler_Auth(t *testing.T) {
	tests := []struct {
		name                 string
		inputLogin           string
		inputStruct          models.Credentials
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:       "#1. User service. Auth. Positive",
			inputLogin: "test",
			inputStruct: models.Credentials{
				Login: "test",
				Pass:  "qwerty",
			},
			expectedError: nil,
		},
		{
			name:       "#2. User service. Auth. Negative",
			inputLogin: "wrong",
			inputStruct: models.Credentials{
				Login: "wrong",
				Pass:  "wrong",
			},
			expectedError: errors.New("internal server error"),
		},
	}

	s := &storagemock.UserStorage{}
	us := service.NewUserService(s)

	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.StorageContextTimeout)
			defer cancel()
			err := us.Auth(ctx, tCase.inputStruct)
			assert.Equal(t, tCase.expectedError, err)
		})
	}
}
