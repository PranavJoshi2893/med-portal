package service

import (
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type mockAuthRepo struct {
	registerFunc func(user model.User) error
	loginFunc    func(email string) (*model.GetByEmail, error)
}

func (m *mockAuthRepo) Register(user model.User) error {
	if m.registerFunc != nil {
		return m.registerFunc(user)
	}
	return nil
}

func (m *mockAuthRepo) Login(email string) (*model.GetByEmail, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email)
	}
	return nil, nil
}

func TestAuthService_Register_Success(t *testing.T) {

	mock := &mockAuthRepo{
		registerFunc: func(user model.User) error {
			return nil
		},
	}

	service := &AuthService{
		repo:   mock,
		hasher: encrypt.NewPasswordHasher("test-pepper"),
	}

	user := &model.CreateUser{
		FirstName: "john",
		LastName:  "doe",
		Email:     "johndoe@test.com",
		Password:  "pass123",
	}

	err := service.Register(user)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

}

func TestAuthService_Login_Success(t *testing.T) {

	hasher := encrypt.NewPasswordHasher("test-pepper")
	hashedPassword, _ := hasher.HashPassword("password123")

	testID, _ := uuid.NewV7()

	mock := &mockAuthRepo{
		loginFunc: func(email string) (*model.GetByEmail, error) {

			return &model.GetByEmail{
				ID:       testID,
				Password: hashedPassword,
			}, nil
		},
	}

	service := NewAuthService(
		mock,
		"test-pepper",
		"test-access-key",
		"test-refresh-key",
	)

	loginUser := &model.LoginUser{
		Email:    "johndoe@test.com",
		Password: "password123",
	}

	resp, err := service.Login(loginUser)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if resp.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

}
