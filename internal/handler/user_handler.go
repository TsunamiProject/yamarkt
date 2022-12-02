package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TsunamiProject/yamarkt/internal/config"
	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type UserServiceProvider interface {
	Auth(ctx context.Context, cred models.Credentials) error
	Register(ctx context.Context, cred models.Credentials) error
}

type UserHandler struct {
	service UserServiceProvider
}

func NewUserHandler(userHandler UserServiceProvider) *UserHandler {
	return &UserHandler{service: userHandler}
}

func (uh UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	credInstance := models.Credentials{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	err = uh.service.Auth(ctx, credInstance)
	if errors.As(err, &customErr.ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusUnauthorized)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	//TODO: create jwt token for auth bearer
	jwtToken := ""
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}

func (uh UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	credInstance := models.Credentials{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	err = uh.service.Register(ctx, credInstance)
	if errors.As(err, &customErr.ErrUserAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	//TODO: create jwt token for auth bearer
	jwtToken := ""
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}
