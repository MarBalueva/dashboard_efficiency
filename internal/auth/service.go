package auth

import (
	"errors"

	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(email, password, name string) error {
	_, err := s.repo.GetByEmail(email)
	if err == nil {
		return errors.New("user already exists")
	}

	hash, err := services.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: hash,
		Name:     name,
	}

	return s.repo.Create(user)
}

func (s *Service) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !services.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
