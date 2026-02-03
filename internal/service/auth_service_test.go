package service

import (
	"context"
	"testing"
	"time"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type mockAuthRepo struct {
	registerFunc      func(ctx context.Context, user model.User) error
	loginFunc         func(ctx context.Context, email string) (*model.GetByEmail, error)
	storeRefreshFunc  func(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error
	revokeRefreshFunc func(ctx context.Context, token string) error
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

func (m *mockAuthRepo) StoreRefreshToken(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error {
	if m.storeRefreshFunc != nil {
		return m.storeRefreshFunc(ctx, id, userID, token, expiresAt)
	}
	return nil
}

func (m *mockAuthRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	if m.revokeRefreshFunc != nil {
		return m.revokeRefreshFunc(ctx, token)
	}
	return nil
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
		{
			name: "email already exists",
			mockFunc: func(ctx context.Context, user model.User) error {
				return model.ErrAlreadyExists
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
		name             string
		loginData        *model.LoginUser
		mockFunc         func(ctx context.Context, email string) (*model.GetByEmail, error)
		storeRefreshFunc func(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error
		expectErr        bool
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
					Role:     "user",
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
					Role:     "user",
				}, nil
			},
			expectErr: true,
		},
		{
			name: "store refresh token fails",
			loginData: &model.LoginUser{
				Email:    "johndoe@test.com",
				Password: "password123",
			},
			mockFunc: func(ctx context.Context, email string) (*model.GetByEmail, error) {
				return &model.GetByEmail{
					ID:       testID,
					Password: hashedPassword,
					Role:     "user",
				}, nil
			},
			storeRefreshFunc: func(ctx context.Context, id uuid.UUID, userID uuid.UUID, token string, expiresAt time.Time) error {
				return errRepo
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAuthRepo{
				loginFunc:        tt.mockFunc,
				storeRefreshFunc: tt.storeRefreshFunc,
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

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		mockFunc  func(ctx context.Context, token string) error
		expectErr bool
	}{
		{
			name:  "success",
			token: "valid-refresh-token",
			mockFunc: func(ctx context.Context, token string) error {
				return nil
			},
			expectErr: false,
		},
		{
			name:  "repo error",
			token: "valid-refresh-token",
			mockFunc: func(ctx context.Context, token string) error {
				return errRepo
			},
			expectErr: true,
		},
		{
			name:  "token not found - idempotent success",
			token: "invalid-token",
			mockFunc: func(ctx context.Context, token string) error {
				return model.ErrNotFound
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAuthRepo{
				revokeRefreshFunc: tt.mockFunc,
			}

			service := NewAuthService(
				mockRepo,
				"test-pepper",
				"test-access-key",
				"test-refresh-key",
			)

			err := service.Logout(context.Background(), tt.token)

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

func TestAuthService_Refresh(t *testing.T) {
	testID, _ := uuid.NewV7()

	tests := []struct {
		name              string
		ctx               context.Context
		revokeRefreshFunc func(ctx context.Context, token string) error
		expectErr         bool
	}{
		{
			name: "success",
			ctx: context.WithValue(
				context.WithValue(
					context.WithValue(context.Background(), "user_id", testID),
					"role", "user",
				),
				"refresh_token", "old-token",
			),
			revokeRefreshFunc: func(ctx context.Context, token string) error {
				return nil
			},
			expectErr: false,
		},
		{
			name: "token reuse - returns unauthorized",
			ctx: context.WithValue(
				context.WithValue(
					context.WithValue(context.Background(), "user_id", testID),
					"role", "user",
				),
				"refresh_token", "reused-token",
			),
			revokeRefreshFunc: func(ctx context.Context, token string) error {
				return model.ErrNotFound
			},
			expectErr: true,
		},
		{
			name:      "missing user_id in context",
			ctx:       context.Background(),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAuthRepo{revokeRefreshFunc: tt.revokeRefreshFunc}
			service := NewAuthService(mockRepo, "test-pepper", "test-access-key", "test-refresh-key")

			resp, err := service.Refresh(tt.ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil || resp.AccessToken == "" {
				t.Error("expected access token in response")
			}
			if resp != nil && resp.RefreshToken == "" {
				t.Error("expected refresh token in response")
			}
		})
	}
}
