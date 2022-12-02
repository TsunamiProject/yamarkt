package customerrs

import "errors"

var (
	ErrUserDoesNotExist  = errors.New("user doesn't exist")
	ErrUserAlreadyExists = errors.New("user already exists")
)
