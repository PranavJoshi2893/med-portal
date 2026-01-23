package service

import (
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

func TestUserService_GetAll_Success(t *testing.T) {

	testID_1, _ := uuid.NewV7()
	testID_2, _ := uuid.NewV7()

	users := []model.GetAll{
		{ID: testID_1, FirstName: "John", LastName: "Doe", Email: "johndoe@test.com"},
		{ID: testID_2, FirstName: "Jane", LastName: "Doe", Email: "janedoe@test.com"},
	}

	mock := &mockUserRepo{
		getAllFunc: func() ([]model.GetAll, error) {
			return users, nil
		},
	}

	service := NewUserService(mock)

	resp, err := service.GetAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !slices.Equal(resp, users) {
		t.Errorf("got: %v want %v", resp, users)
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	testID, _ := uuid.NewV7()

	user := model.GetByID{
		ID:        testID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@test.com",
	}

	mock := &mockUserRepo{
		getByIDFunc: func(id uuid.UUID) (*model.GetByID, error) {
			return &user, nil
		},
	}

	service := NewUserService(mock)

	resp, err := service.GetByID(testID)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(&user, resp) {
		t.Errorf("got %v want %v", user, resp)
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
