package service

import (
	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/internal/repository"
)

type accrualService interface {
	LoadOld()
	AddToFetchList(order *models.Order)
}

type Service struct {
	repo           repository.Repository
	accrualService accrualService
}

func NewService(repo repository.Repository, accrualService accrualService) *Service {
	return &Service{repo: repo, accrualService: accrualService}
}
