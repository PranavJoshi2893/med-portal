package service

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	getAllFunc     func() ([]model.GetAll, error)
	deleteByIDFunc func(id uuid.UUID) error
	getByIDFunc    func(id uuid.UUID) (*model.GetByID, error)
}

func (m *mockUserRepo) GetAll() ([]model.GetAll, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil, nil
}

func (m *mockUserRepo) DeleteByID(id uuid.UUID) error {
	if m.deleteByIDFunc != nil {
		return m.deleteByIDFunc(id)
	}
	return nil
}

func (m *mockUserRepo) GetByID(id uuid.UUID) (*model.GetByID, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
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
		mockFunc  func() ([]model.GetAll, error)
		expectErr bool
		expect    []model.GetAll
	}{
		{
			name: "sucess",
			mockFunc: func() ([]model.GetAll, error) {
				return users, nil
			},
			expectErr: false,
			expect:    users,
		},
		{
			name: "repo error",
			mockFunc: func() ([]model.GetAll, error) {
				return nil, errRepo
			},
			expectErr: true,
		},
		{
			name: "nil slice treated as empty",
			mockFunc: func() ([]model.GetAll, error) {
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

			resp, err := service.GetAll()

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
		mockFunc  func(id uuid.UUID) (*model.GetByID, error)
		expectErr bool
		expect    *model.GetByID
	}{
		{
			name: "success",
			mockFunc: func(id uuid.UUID) (*model.GetByID, error) {
				return &user, nil
			},
			expectErr: false,
			expect:    &user,
		},
		{
			name: "repo error",
			mockFunc: func(id uuid.UUID) (*model.GetByID, error) {
				return nil, errRepo
			},
			expectErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(id uuid.UUID) (*model.GetByID, error) {
				return nil, model.ErrNotFound
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{getByIDFunc: tt.mockFunc}
			service := NewUserService(mock)

			resp, err := service.GetByID(testID)

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

func TestUserService_DeleteByID_Success(t *testing.T) {
	testID, _ := uuid.NewV7()

	mock := &mockUserRepo{
		deleteByIDFunc: func(id uuid.UUID) error {
			return nil
		},
	}

	service := NewUserService(mock)

	err := service.DeleteByID(testID)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
