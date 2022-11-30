package service

import (
	"context"

	"github.com/TsunamiProject/yamarkt/internal/models"
)

type UserStorage interface {
	Auth(ctx context.Context, login string, pass string) error
	Register(ctx context.Context, login string, pass string) error
}

type UserService struct {
	userStorage UserStorage
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{userStorage: userStorage}
}

func (us *UserService) Register(ctx context.Context, cred models.Credentials) error {
	return nil
}

func (us *UserService) Auth(ctx context.Context, cred models.Credentials) error {
	return nil
}
