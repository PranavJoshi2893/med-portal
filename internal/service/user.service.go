package service

import (
	"errors"
	"fmt"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(repo *repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(user *model.CreateUser) error {

	newUser := model.User{
		ID:        uuid.New(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}

	err := s.repo.Register(newUser)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return fmt.Errorf("email %w", err)
		}

		return err
	}

	return nil
}

func (s *UserService) Login(user *model.LoginUser) error {
	_, err := s.repo.Login(user)

	if err != nil {
		return err
	}

	return nil
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
