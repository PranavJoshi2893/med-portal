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
	getAllFunc     func(ctx context.Context, limit, offset int) ([]model.GetAll, error)
	getCountFunc   func(ctx context.Context) (int, error)
	deleteByIDFunc func(ctx context.Context, id uuid.UUID) error
	getByIDFunc    func(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
	updateByIDFunc func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error
}

func (m *mockUserRepo) GetAll(ctx context.Context, limit, offset int) ([]model.GetAll, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockUserRepo) GetCount(ctx context.Context) (int, error) {
	if m.getCountFunc != nil {
		return m.getCountFunc(ctx)
	}
	return 0, nil
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

	params := model.PaginationParams{Page: 1, Limit: 10}

	tests := []struct {
		name         string
		mockFunc     func(ctx context.Context, limit, offset int) ([]model.GetAll, error)
		getCountFunc func(ctx context.Context) (int, error)
		getByIDFunc  func(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
		callerID     *uuid.UUID
		callerRole   string
		expectErr    bool
		expectItems  []model.GetAll
		expectTotal  int
	}{
		{
			name: "success - admin gets all",
			mockFunc: func(ctx context.Context, limit, offset int) ([]model.GetAll, error) {
				return users, nil
			},
			getCountFunc: func(ctx context.Context) (int, error) { return 2, nil },
			callerID:     &testID_1,
			callerRole:   "admin",
			expectErr:    false,
			expectItems:  users,
			expectTotal:  2,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, limit, offset int) ([]model.GetAll, error) {
				return nil, errRepo
			},
			callerID:   &testID_1,
			callerRole: "admin",
			expectErr:  true,
		},
		{
			name: "nil slice treated as empty",
			mockFunc: func(ctx context.Context, limit, offset int) ([]model.GetAll, error) {
				return nil, nil
			},
			getCountFunc: func(ctx context.Context) (int, error) { return 0, nil },
			callerID:     &testID_1,
			callerRole:   "admin",
			expectErr:    false,
			expectItems:  []model.GetAll{},
			expectTotal:  0,
		},
		{
			name: "regular user gets only self",
			mockFunc: func(ctx context.Context, limit, offset int) ([]model.GetAll, error) {
				return nil, nil
			},
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return &model.GetByID{ID: testID_1, FirstName: "John", LastName: "Doe", Email: "johndoe@test.com"}, nil
			},
			callerID:    &testID_1,
			callerRole:  "user",
			expectErr:   false,
			expectItems: []model.GetAll{{ID: testID_1, FirstName: "John", LastName: "Doe", Email: "johndoe@test.com"}},
			expectTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				getAllFunc:   tt.mockFunc,
				getCountFunc: tt.getCountFunc,
				getByIDFunc:  tt.getByIDFunc,
			}
			service := NewUserService(mock)
			callerID := testID_1
			if tt.callerID != nil {
				callerID = *tt.callerID
			}
			callerRole := tt.callerRole
			if callerRole == "" {
				callerRole = "user"
			}

			resp, err := service.GetAll(context.Background(), callerID, callerRole, params)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(resp.Items, tt.expectItems) {
				t.Errorf("got items %v want %v", resp.Items, tt.expectItems)
			}
			if resp.Meta.Total != tt.expectTotal {
				t.Errorf("got total %d want %d", resp.Meta.Total, tt.expectTotal)
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
		name       string
		mockFunc   func(ctx context.Context, id uuid.UUID) (*model.GetByID, error)
		callerID   uuid.UUID
		callerRole string
		expectErr  bool
		expect     *model.GetByID
	}{
		{
			name: "success - same user",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return &user, nil
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  false,
			expect:     &user,
		},
		{
			name: "success - admin accessing other user",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return &user, nil
			},
			callerID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			callerRole: "admin",
			expectErr:  false,
			expect:     &user,
		},
		{
			name: "forbidden - user accessing other user",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return &user, nil
			},
			callerID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return nil, errRepo
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "user not found",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*model.GetByID, error) {
				return nil, model.ErrNotFound
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{getByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			resp, err := service.GetByID(context.Background(), testID, tt.callerID, tt.callerRole)

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
		name       string
		mockFunc   func(ctx context.Context, id uuid.UUID) error
		callerID   uuid.UUID
		callerRole string
		expectErr  bool
	}{
		{
			name: "success - same user",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  false,
		},
		{
			name: "success - admin deleting other user",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			callerID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			callerRole: "admin",
			expectErr:  false,
		},
		{
			name: "forbidden - user deleting other user",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
			callerID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return errRepo
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "user already deleted",
			mockFunc: func(ctx context.Context, id uuid.UUID) error {
				return model.ErrAlreadyDeleted
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{deleteByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			err := service.DeleteByID(context.Background(), testID, tt.callerID, tt.callerRole)

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
		name       string
		mockFunc   func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error
		callerID   uuid.UUID
		callerRole string
		expectErr  bool
	}{
		{
			name: "success - same user",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return nil
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  false,
		},
		{
			name: "forbidden - user updating other user",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return nil
			},
			callerID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "repo error",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return errRepo
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "user not found",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrNotFound
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "user already deleted",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrAlreadyDeleted
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
		{
			name: "bad request - nil fields",
			mockFunc: func(ctx context.Context, id uuid.UUID, data *model.UpdateUser) error {
				return model.ErrBadRequest
			},
			callerID:   testID,
			callerRole: "user",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{updateByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			err := service.UpdateByID(context.Background(), testID, updateData, tt.callerID, tt.callerRole)

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
