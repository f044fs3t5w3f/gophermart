package repository

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

type Repository interface {
	CreateUser(context.Context, *models.User) error
	GetUserById(context.Context, int64) (*models.User, error)
	GetUserByLogin(context.Context, string) (*models.User, error)
	DoesUserExist(context.Context, string) (bool, error)

	CreateSession(ctx context.Context, session *models.Session) error
	GetUserByToken(ctx context.Context, token string) (*models.User, error)

	GetOrderByNumber(ctx context.Context, number string) (*models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) error

	GetOrdersByUserId(ctx context.Context, userId int64) ([]*models.Order, error)
}
