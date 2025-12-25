package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func generateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", ErrInternalError
	}
	if user == nil {
		return "", ErrUserDoesNotExists
	}
	passwordIsCorrect := checkPassword(user.PasswordHash, password)
	if !passwordIsCorrect {
		return "", ErrIncorrectPassword
	}
	token, err := generateToken()
	if err != nil {
		return "", ErrInternalError
	}

	session := &models.Session{
		UserID: user.ID,
		Token:  token,
	}
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return "", ErrInternalError
	}

	return token, nil
}
