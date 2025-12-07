package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists        = errors.New("User exists")
	ErrUserDoesNotExists = errors.New("User doesn't exists")
	ErrIncorrectPassword = errors.New("Incorrect password")

	ErrInternalError = errors.New("Internal error password")
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func getPasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

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
		UserId: user.Id,
		Token:  token,
	}
	err = s.repo.CreateSession(ctx, session)
	if err != nil {
		return "", ErrInternalError
	}

	return token, nil
}
