package customerrs

import "errors"

var (
	ErrUserDoesNotExist            = errors.New("user doesn't exist")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrOrderAlreadyExists          = errors.New("order already exists")
	ErrOrderCreatedByAnotherLogin  = errors.New("order already created by another customer")
	ErrWithdrawalOrderAlreadyExist = errors.New("withdrawal order already exists")
	ErrUnauthorizedUser            = errors.New("unauthorized user")
	ErrNoOrders                    = errors.New("no orders created by login")
	ErrNoFunds                     = errors.New("no funds on login balance")
	ErrNoWithdrawals               = errors.New("no withdrawals by login")
	ErrWrongPassword               = errors.New("wrong password received")
)
