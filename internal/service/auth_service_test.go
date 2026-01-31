package service

import (
	"context"
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type mockAuthRepo struct {
	registerFunc func(ctx context.Context, user model.User) error
	loginFunc    func(ctx context.Context, email string) (*model.GetByEmail, error)
}

func (m *mockAuthRepo) Register(ctx context.Context, user model.User) error {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, user)
	}
	return nil
}

func (m *mockAuthRepo) Login(ctx context.Context, email string) (*model.GetByEmail, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, email)
	}
	return nil, nil
}

func TestAuthService_Register(t *testing.T) {

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, user model.User) error
		expectErr bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, user model.User) error {
				return nil
			},
			expectErr: false,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, user model.User) error {
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

			err := service.Register(context.Background(), user)

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

func TestAuthService_Login(t *testing.T) {

	hasher := encrypt.NewPasswordHasher("test-pepper")

	testID, _ := uuid.NewV7()
	hashedPassword, _ := hasher.HashPassword("password123")

	tests := []struct {
		name      string
		loginData *model.LoginUser
		mockFunc  func(ctx context.Context, email string) (*model.GetByEmail, error)
		expectErr bool
	}{
		{
			name: "success",
			loginData: &model.LoginUser{
				Email:    "johndoe@test.com",
				Password: "password123",
			},
			mockFunc: func(ctx context.Context, email string) (*model.GetByEmail, error) {
				return &model.GetByEmail{
					ID:       testID,
					Password: hashedPassword,
				}, nil
			},
			expectErr: false,
		},
		{
			name: "user not found",
			loginData: &model.LoginUser{
				Email:    "missing@test.com",
				Password: "password123",
			},
			mockFunc: func(ctx context.Context, email string) (*model.GetByEmail, error) {
				return nil, model.ErrNotFound
			},
			expectErr: true,
		},
		{
			name: "wrong password",
			loginData: &model.LoginUser{
				Email:    "johndoe@test.com",
				Password: "wrong-password",
			},
			mockFunc: func(ctx context.Context, email string) (*model.GetByEmail, error) {
				return &model.GetByEmail{
					ID:       testID,
					Password: hashedPassword,
				}, nil
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAuthRepo{
				loginFunc: tt.mockFunc,
			}

			service := NewAuthService(
				mockRepo,
				"test-pepper",
				"test-access-key",
				"test-refresh-key",
			)

			resp, err := service.Login(context.Background(), tt.loginData)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if resp != nil {
					t.Fatalf("expected nil response, got %+v", resp)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("expected response, got nil")
			}

			if resp.AccessToken == "" {
				t.Error("expected access token, got empty string")
			}

			if resp.RefreshToken == "" {
				t.Error("expected refresh token, got empty string")
			}
		})
	}

}
