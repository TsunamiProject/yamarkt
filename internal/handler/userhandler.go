package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"

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
		return
	}
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	err = uh.service.Auth(ctx, credInstance)
	if errors.Is(err, customErr.ErrUserDoesNotExist) || errors.Is(err, customErr.ErrWrongPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	claims := map[string]interface{}{
		"login": credInstance.Login,
	}
	jwtauth.SetExpiryIn(claims, config.TokenTTL)
	_, jwtToken, err := config.TokenAuth.Encode(claims)
	if err != nil {
		log.Printf("error while encoding auth token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}

func (uh UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	credInstance := models.Credentials{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	err = uh.service.Register(ctx, credInstance)
	if errors.Is(err, customErr.ErrUserAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	claims := map[string]interface{}{
		"login": credInstance.Login,
	}
	jwtauth.SetExpiryIn(claims, config.TokenTTL)
	_, jwtToken, err := config.TokenAuth.Encode(claims)
	if err != nil {
		log.Printf("error while encoding auth token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}
