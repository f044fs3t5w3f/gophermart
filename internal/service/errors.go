package service

import (
	"errors"
	"fmt"

	"github.com/f044fs3t5w3f/gophermart/internal/repository"
)

var (
	ErrUserExists        = errors.New("user exists")
	ErrUserDoesNotExists = errors.New("user doesn't exists")
	ErrIncorrectPassword = errors.New("incorrect password")

	ErrInternalError = errors.New("internal error password")

	ErrInvalidNumber      = errors.New("invalid number")
	ErrOrderAlreadyExists = errors.New("order already exists ")

	ErrWithdrawAllreadyExistsForAnotherUser = fmt.Errorf("service: %w", repository.ErrWithdrawAllreadyExistsForAnotherUser)
	ErrWithdrawAllreadyExists               = fmt.Errorf("service: %w", repository.ErrWithdrawAllreadyExists)
	ErrNotEnough                            = fmt.Errorf("service: %w", repository.ErrNotEnough)
)
