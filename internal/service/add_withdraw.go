package service

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/pkg/luhn"
)

func (s *Service) AddWithdraw(ctx context.Context, order string, sum float64) error {
	usrPtr := ctx.Value(auth.ContextUserKey)
	user, ok := usrPtr.(*models.User)
	if !ok {
		return ErrInternalError
	}

	if !luhn.Check(order) {
		return ErrInvalidNumber
	}

	withdraw := &models.Withdraw{
		Order:  order,
		Sum:    sum,
		UserId: user.ID,
	}
	return s.repo.CreateWithdraw(ctx, withdraw)
}
