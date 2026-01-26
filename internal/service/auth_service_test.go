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

func TestAuthService_Register(t *testing.T) {

	tests := []struct {
		name      string
		mockFunc  func(user model.User) error
		expectErr bool
	}{
		{
			name: "success",
			mockFunc: func(user model.User) error {
				return nil
			},
			expectErr: false,
		},
		{
			name: "repo error",
			mockFunc: func(user model.User) error {
				return errRepo
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAuthRepo{registerFunc: tt.mockFunc}
			service := &AuthService{
				mock,
				encrypt.NewPasswordHasher("test-pepper"),
				"test-access-key",
				"test=refresh-key",
			}

			user := &model.CreateUser{
				FirstName: "john",
				LastName:  "doe",
				Email:     "johndoe@test.com",
				Password:  "pass123",
			}

			err := service.Register(user)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
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
