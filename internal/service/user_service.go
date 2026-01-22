package service

import (
	"errors"
	"fmt"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/repository"
	"github.com/PranavJoshi2893/med-portal/pkg/auth"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type UserService struct {
	repo            *repository.UserRepo
	hasher          *encrypt.PasswordHasher
	accessTokenKey  string
	refreshTokenKey string
}

func NewUserService(repo *repository.UserRepo, pepper string, accessTokenKey string, refreshTokenKey string) *UserService {
	return &UserService{
		repo:            repo,
		hasher:          encrypt.NewPasswordHasher(pepper),
		accessTokenKey:  accessTokenKey,
		refreshTokenKey: refreshTokenKey,
	}
}

func (s *UserService) Register(user *model.CreateUser) error {

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

	err = s.repo.Register(newUser)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return fmt.Errorf("email %w", err)
		}

		return err
	}

	return nil
}

func (s *UserService) Login(user *model.LoginUser) (*model.LoginResponse, error) {
	data, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("email %w", err)
		}
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
