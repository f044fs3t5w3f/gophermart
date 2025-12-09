package service

import (
	"context"
	"errors"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/pkg/luhn"
)

var (
	ErrInvalidNumber      = errors.New("Invalid number")
	ErrOrderAlreadyExists = errors.New("Order already exists ")
)

func (s *Service) AddOrder(ctx context.Context, token, orderNumber string) (bool, error) {
	usrPtr := ctx.Value(auth.ContextUserKey)
	user, ok := usrPtr.(*models.User)
	if !ok {
		return false, ErrInternalError
	}

	isNumberValid := luhn.Check(orderNumber)
	if !isNumberValid {
		return false, ErrInvalidNumber
	}

	existedOrder, err := s.repo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return false, ErrInternalError
	}
	if existedOrder != nil {
		if existedOrder.UserId != user.Id {
			return false, ErrOrderAlreadyExists
		}
		return false, nil
	}

	order := &models.Order{
		Number: orderNumber,
		UserId: user.Id,
		Status: "NEW",
	}

	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		return false, ErrInternalError
	}

	return true, nil
}
