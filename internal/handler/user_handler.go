package handler

import (
	"context"
	"net/http"

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

}

func (uh UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}
