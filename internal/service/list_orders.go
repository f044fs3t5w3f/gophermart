package service

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/auth"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func (s *Service) ListOrders(ctx context.Context) ([]*models.Order, error) {
	usrPtr := ctx.Value(auth.ContextUserKey)
	user, ok := usrPtr.(*models.User)
	if !ok {
		return nil, ErrInternalError
	}
	orders, err := s.repo.GetOrdersByUserID(ctx, user.ID)
	return orders, err
}
