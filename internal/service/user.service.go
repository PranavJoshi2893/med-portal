package service

import (
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
