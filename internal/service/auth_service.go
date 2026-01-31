package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/repository"
	"github.com/PranavJoshi2893/med-portal/pkg/auth"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type AuthService struct {
	repo            repository.AuthRepository
	hasher          *encrypt.PasswordHasher
	accessTokenKey  string
	refreshTokenKey string
}

func NewAuthService(repo repository.AuthRepository, pepper string, accessTokenKey string, refreshTokenKey string) *AuthService {
	return &AuthService{
		repo:            repo,
		hasher:          encrypt.NewPasswordHasher(pepper),
		accessTokenKey:  accessTokenKey,
		refreshTokenKey: refreshTokenKey,
	}
}

func (s *AuthService) Register(ctx context.Context, user *model.CreateUser) error {

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to genereate uuid: %w", err)
	}

	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := model.User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	}

	err = s.repo.Register(ctx, newUser)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return fmt.Errorf("email %w", err)
		}

		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, user *model.LoginUser) (*model.LoginResponse, error) {
	data, err := s.repo.Login(ctx, user.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("email %w", err)
		}
		return nil, err
	}

	if ok := s.hasher.VerifyPassword(user.Password, data.Password); !ok {
		return nil, model.ErrUnauthorized
	}

	var access_token string
	var refresh_token string

	if access_token, err = auth.GenerateAccessToken(s.accessTokenKey, data.ID); err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	if refresh_token, err = auth.GenerateRefreshToken(s.refreshTokenKey, data.ID); err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	resp := model.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return &resp, nil
}
