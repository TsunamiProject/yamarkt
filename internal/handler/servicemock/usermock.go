package servicemock

import (
	"context"
	"errors"

	customErr "github.com/TsunamiProject/yamarkt/internal/customerrs"
	"github.com/TsunamiProject/yamarkt/internal/models"
)

type UserServiceMock struct {
}

func (us *UserServiceMock) Register(ctx context.Context, cr models.Credentials) (err error) {
	switch {
	case cr.Login == "test" && cr.Pass == "qwerty":
		return nil
	case cr.Login == "test2" && cr.Pass == "qwerty":
		return customErr.ErrUserAlreadyExists
	default:
		return errors.New("internal server error")
	}
}
func (us *UserServiceMock) Auth(ctx context.Context, cr models.Credentials) (err error) {
	switch {
	case cr.Login == "test" && cr.Pass == "qwerty":
		return nil
	case cr.Login == "test" && cr.Pass == "wrong":
		return customErr.ErrWrongPassword
	case cr.Login == "test3" && cr.Pass == "qwerty123":
		return customErr.ErrUserDoesNotExist
	default:
		return errors.New("internal server error")
	}
}
