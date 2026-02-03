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

func isAdmin(role string) bool {
	return role == "admin" || role == "super_admin"
}

func (s *UserService) GetAll(ctx context.Context, callerID uuid.UUID, callerRole string, params model.PaginationParams) (*model.PaginatedUsersResponse, error) {
	if isAdmin(callerRole) {
		users, err := s.repo.GetAll(ctx, params.Limit, (params.Page-1)*params.Limit)
		if err != nil {
			return nil, err
		}
		total, err := s.repo.GetCount(ctx)
		if err != nil {
			return nil, err
		}
		totalPages := total / params.Limit
		if total%params.Limit > 0 {
			totalPages++
		}
		return &model.PaginatedUsersResponse{
			Items: users,
			Meta: model.PaginationMeta{
				Page:       params.Page,
				Limit:      params.Limit,
				Total:      total,
				TotalPages: totalPages,
			},
		}, nil
	}
	user, err := s.repo.GetByID(ctx, callerID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("user %w", err)
		}
		return nil, err
	}
	return &model.PaginatedUsersResponse{
		Items: []model.GetAll{{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}},
		Meta: model.PaginationMeta{
			Page:       1,
			Limit:      params.Limit,
			Total:      1,
			TotalPages: 1,
		},
	}, nil
}

func (s *UserService) DeleteByID(ctx context.Context, id uuid.UUID, callerID uuid.UUID, callerRole string) error {
	if id != callerID && !isAdmin(callerRole) {
		return model.ErrForbidden
	}
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrAlreadyDeleted) {
			return fmt.Errorf("user %w", err)
		}
		return err
	}
	return nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID, callerID uuid.UUID, callerRole string) (*model.GetByID, error) {
	if id != callerID && !isAdmin(callerRole) {
		return nil, model.ErrForbidden
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("user %w", err)
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateByID(ctx context.Context, id uuid.UUID, data *model.UpdateUser, callerID uuid.UUID, callerRole string) error {
	if id != callerID && !isAdmin(callerRole) {
		return model.ErrForbidden
	}
	err := s.repo.UpdateByID(ctx, id, data)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return fmt.Errorf("user %w", err)
		}
		if errors.Is(err, model.ErrAlreadyDeleted) {
			return fmt.Errorf("user %w", err)
		}
		return err
	}
	return nil
}
