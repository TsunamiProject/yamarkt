package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"

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
	encodedPass, err := Encode(cred.Pass)
	if err != nil {
		log.Printf("encoding password error: %s", err)
		return err
	}
	err = us.userStorage.Register(ctx, cred.Login, encodedPass)
	return nil
}

func (us *UserService) Auth(ctx context.Context, cred models.Credentials) error {
	encodedPass, err := Encode(cred.Pass)
	if err != nil {
		log.Printf("encoding password error: %s", err)
		return err
	}
	err = us.userStorage.Auth(ctx, cred.Login, encodedPass)
	return nil
}

func Encode(src string) (encodedString string, err error) {
	crInst := sha256.New()
	crInst.Write([]byte(src))
	srcBytes := crInst.Sum(nil)
	encodedString = hex.EncodeToString(srcBytes)
	return encodedString, nil
}
