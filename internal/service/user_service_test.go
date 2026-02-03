package service

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	getAllFunc     func(ctx context.Context) ([]model.GetAll, error)
	deleteByIDFunc func(ctx context.Context, id uuid.UUID) error
	getByIDFunc    func(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
	updateByIDFunc func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error
}

func (m *mockUserRepo) GetAll(ctx context.Context) ([]model.GetAll, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockUserRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if m.deleteByIDFunc != nil {
		return m.deleteByIDFunc(ctx, id)
	}
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepo) UpdateByID(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
	if m.updateByIDFunc != nil {
		return m.updateByIDFunc(ctx, id, data)
	}
	return nil
}

var errRepo = errors.New("repo error")

func TestUserService_GetAll(t *testing.T) {

	testID_1, _ := uuid.NewV7()
	testID_2, _ := uuid.NewV7()

	users := []model.GetAll{
		{ID: testID_1, FirstName: "John", LastName: "Doe", Email: "johndoe@test.com"},
		{ID: testID_2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@test.com"},
	}

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context) ([]model.GetAll, error)
		expectErr bool
		expect    []model.GetAll
	}{
		{
			name: "sucess",
			mockFunc: func(ctx context.Context) ([]model.GetAll, error) {
				return users, nil
			},
			expectErr: false,
			expect:    users,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context) ([]model.GetAll, error) {
				return nil, errRepo
			},
			expectErr: true,
		},
		{
			name: "nil slice treated as empty",
			mockFunc: func(ctx context.Context) ([]model.GetAll, error) {
				return nil, nil
			},
			expectErr: false,
			expect:    []model.GetAll{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{getAllFunc: tt.mockFunc}
			service := NewUserService(mock)

			resp, err := service.GetAll(context.Background())

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(resp, tt.expect) {
				t.Errorf("got %v want %v", resp, tt.expect)
			}
		})
	}

}

func TestUserService_GetByID(t *testing.T) {
	testID, _ := uuid.NewV7()

	user := model.GetByID{
		ID:        testID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@test.com",
	}

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
		expectErr bool
		expect    *model.GetByID
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return &user, nil
			},
			expectErr: false,
			expect:    &user,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return nil, errRepo
			},
			expectErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return nil, model.ErrNotFound
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{getByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			resp, err := service.GetByID(context.Background(), testID)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(resp, tt.expect) {
				t.Errorf("got %+v want %+v", resp, tt.expect)
			}
		})
	}
}

func TestUserService_DeleteByID(t *testing.T) {

	testID, _ := uuid.NewV7()

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, id uuid.UUID) error
		expectErr bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			expectErr: false,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return errRepo
			},
			expectErr: true,
		},
		{
			name: "user already deleted",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return model.ErrAlreadyDeleted
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{deleteByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			err := service.DeleteByID(context.Background(), testID)

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

func TestUserService_UpdateByID(t *testing.T) {
	testID, _ := uuid.NewV7()

	firstName := "John"
	lastName := "Doe"
	updateData := &model.UpdateUser{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error
		expectErr bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return nil
			},
			expectErr: false,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return errRepo
			},
			expectErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrNotFound
			},
			expectErr: true,
		},
		{
			name: "user already deleted",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrAlreadyDeleted
			},
			expectErr: true,
		},
		{
			name: "bad request - nil fields",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrBadRequest
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{updateByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			err := service.UpdateByID(context.Background(), testID, updateData)

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
