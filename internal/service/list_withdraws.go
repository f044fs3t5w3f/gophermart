package service

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func (s *Service) ListWithdraws(ctx context.Context) ([]*models.Withdraw, error) {
	usrPtr := ctx.Value(auth.ContextUserKey)
	user, ok := usrPtr.(*models.User)
	if !ok {
		return nil, ErrInternalError
	}
	withdraws, err := s.repo.GetWithdrawsByUserId(ctx, user.Id)
	return withdraws, err
}
