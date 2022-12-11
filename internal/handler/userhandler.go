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

// Auth takes on enter user login from Authentication header and credentials from request payload like
//{
//    "login": "<login>",
//    "password": "<password>"
//}
//and returns status codes: 200 - on successfully authentication,
//400 - on bad request,
//401 - on wrong credentials,
//500 - on internal server error
func (uh UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	credInstance := models.Credentials{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Auth handler. Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//unmarshalling request body to credentials struct
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Auth handler. Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
		return
	}
	//calling Auth service method
	err = uh.service.Auth(ctx, credInstance)
	if errors.Is(err, customErr.ErrUserDoesNotExist) || errors.Is(err, customErr.ErrWrongPassword) {
		log.Printf("Auth handler. Login: %s: Wrong credentials", credInstance.Login)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Auth handler. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//collecting claims map for setting up to jwtToken
	claims := map[string]interface{}{
		"login": credInstance.Login,
	}
	//setting up jwtToken time to live
	jwtauth.SetExpiryIn(claims, config.TokenTTL)
	//encoding claims to jwtToken
	_, jwtToken, err := config.TokenAuth.Encode(claims)
	if err != nil {
		errString := fmt.Sprintf("Auth handler. Error while encoding auth token: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}

// Register takes on enter user login from Authentication header and credentials from request payload like
//{
//    "login": "<login>",
//    "password": "<password>"
//}
//and returns status codes: 200 - on successfully register and authentication,
//400 - on bad request,
//409 - on user already exists,
//500 - on internal server error
func (uh UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	credInstance := models.Credentials{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Register handler. Error while reading body: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//unmarshalling request body to credentials struct
	err = json.Unmarshal(body, &credInstance)
	if err != nil {
		log.Printf("Register handler. Request body unmarshal error: %s", err)
		http.Error(w, "Invalid credentials JSON received", http.StatusBadRequest)
	}

	//creating context from request context
	ctx, cancel := context.WithTimeout(r.Context(), config.StorageContextTimeout)
	defer cancel()
	//calling Register service method
	err = uh.service.Register(ctx, credInstance)
	if errors.Is(err, customErr.ErrUserAlreadyExists) {
		log.Printf("Register handler. Login: %s: User already exist", credInstance.Login)
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Printf("Register handler. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//collecting claims map for setting up to jwtToken
	claims := map[string]interface{}{
		"login": credInstance.Login,
	}
	//setting up jwtToken time to live
	jwtauth.SetExpiryIn(claims, config.TokenTTL)
	//encoding claims to jwtToken
	_, jwtToken, err := config.TokenAuth.Encode(claims)
	if err != nil {
		errString := fmt.Sprintf("Register handler. Error while encoding auth token: %s", err)
		log.Printf(errString)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	w.WriteHeader(http.StatusOK)
}
