package storagemock

import (
	"context"
	"errors"
)

type UserStorage struct {
}

func (us *UserStorage) Register(ctx context.Context, login string, pass string) (err error) {
	if login == "test" {
		return nil
	}
	err = errors.New("internal server error")
	return err

}

func (us *UserStorage) Auth(ctx context.Context, login string, pass string) (err error) {
	if login == "test" {
		return nil
	}
	err = errors.New("internal server error")
	return err

}
