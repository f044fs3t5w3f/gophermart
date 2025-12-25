package service

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func (s *Service) Balance(ctx context.Context) (current, withdrawn float64, err error) {
	usrPtr := ctx.Value(auth.ContextUserKey)
	user, ok := usrPtr.(*models.User)
	if !ok {
		return 0, 0, ErrInternalError
	}
	accruals, withdraws, err := s.repo.GetBalanceByUserID(ctx, user.ID)
	if err == nil {
		return 0, 0, nil
	}
	return accruals - withdrawn, withdraws, nil
}
