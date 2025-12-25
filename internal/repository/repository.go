package repository

import (
	"context"
	"errors"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

var (
	ErrNotEnough                            = errors.New("not enough")
	ErrWithdrawAllreadyExists               = errors.New("withdraw allready exists")
	ErrWithdrawAllreadyExistsForAnotherUser = errors.New("withdraw allready exists for another user")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	DoesUserExist(ctx context.Context, login string) (bool, error)
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
	GetBalanceByUserID(ctx context.Context, userID int64) (float64, float64, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByNumber(ctx context.Context, number string) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
}

type WithdrawRepository interface {
	CreateWithdraw(ctx context.Context, order *models.Withdraw) error
	GetWithdrawsByUserID(ctx context.Context, userID int64) ([]*models.Withdraw, error)
}

type Repository interface {
	UserRepository
	SessionRepository
	OrderRepository
	WithdrawRepository
}
