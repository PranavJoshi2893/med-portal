package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		return fmt.Errorf("failed to generate uuid: %w", err)
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

// Login performs authentication for a user based on provided credentials.
// Steps:
//   1. Retrieves the user by email from the repository.
//   2. Verifies the provided password against the stored hash.
//   3. Generates an access token and a refresh token if authentication is successful.
//   4. Stores the refresh token in the database for session management.
//   5. Returns the generated access and refresh tokens to the caller.

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

	if access_token, err = auth.GenerateAccessToken(s.accessTokenKey, data.ID, data.Role); err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	if refresh_token, err = auth.GenerateRefreshToken(s.refreshTokenKey, data.ID, data.Role); err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	tokenHash := encrypt.HashToken(refresh_token)
	err = s.repo.StoreRefreshToken(ctx, id, data.ID, tokenHash, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	resp := model.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return &resp, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	tokenHash := encrypt.HashToken(token)
	err := s.repo.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (s *AuthService) Refresh(ctx context.Context) (*model.LoginResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("internal server error")
	}

	role, _ := ctx.Value("role").(string)
	if role == "" {
		role = "user"
	}

	oldToken, _ := ctx.Value("refresh_token").(string)
	if oldToken != "" {
		tokenHash := encrypt.HashToken(oldToken)
		if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil, model.ErrUnauthorized
			}
			return nil, fmt.Errorf("internal server error")
		}
	}

	accessToken, err := auth.GenerateAccessToken(s.accessTokenKey, userID, role)
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	refreshToken, err := auth.GenerateRefreshToken(s.refreshTokenKey, userID, role)
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	tokenHash := encrypt.HashToken(refreshToken)
	err = s.repo.StoreRefreshToken(ctx, id, userID, tokenHash, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
