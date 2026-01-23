package service

import (
	"errors"
	"fmt"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAll() ([]model.GetAll, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteByID(id uuid.UUID) error {
	err := s.repo.DeleteByID(id)
	if errors.Is(err, model.ErrAlreadyDeleted) {
		return fmt.Errorf("user %w", err)
	}
	return nil
}

func (s *UserService) GetByID(id uuid.UUID) (*model.GetByID, error) {

	user, err := s.repo.GetByID(id)
	if errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("user %w", err)
	}

	return user, nil

}
