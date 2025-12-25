package service

import "errors"

var (
	ErrUserExists        = errors.New("user exists")
	ErrUserDoesNotExists = errors.New("user doesn't exists")
	ErrIncorrectPassword = errors.New("incorrect password")

	ErrInternalError = errors.New("internal error password")

	ErrInvalidNumber      = errors.New("invalid number")
	ErrOrderAlreadyExists = errors.New("order already exists ")
)
