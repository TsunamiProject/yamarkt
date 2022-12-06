package customerrs

import "errors"

var (
	ErrUserDoesNotExist           = errors.New("user doesn't exist")
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrOrderAlreadyExists         = errors.New("order already exists")
	ErrOrderCreatedByAnotherLogin = errors.New("order already created by another customer")
	ErrUnauthorizedUser           = errors.New("unauthorized user")
	ErrNoOrders                   = errors.New("no orders created by login")
	ErrNoFunds                    = errors.New("no funds on login balance")
	ErrWrongPassword              = errors.New("wrong password received")
)
