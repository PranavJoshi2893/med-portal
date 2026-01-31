package service

import (
	"context"
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

func (s *UserService) GetAll(ctx context.Context) ([]model.GetAll, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteByID(ctx, id)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyDeleted) {
			return fmt.Errorf("user %w", err)
		}
		return err
	}
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {

	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("user %w", err)
		}
		return nil, err
	}

	return user, nil

}
