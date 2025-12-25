package service

import (
	"context"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func (s *Service) Register(ctx context.Context, login, password string) error {
	passwordHash, err := getPasswordHash(password)
	if err != nil {
		return err
	}
	userExists, err := s.repo.DoesUserExist(ctx, login)
	if err != nil {
		return err
	}
	if userExists {
		return ErrUserExists
	}
	user := &models.User{
		Login:        login,
		PasswordHash: passwordHash,
	}
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return err
}
