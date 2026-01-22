package service

import (
	"slices"
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/pkg/encrypt"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	registerFunc   func(user model.User) error
	getByEmailFunc func(email string) (*model.GetByEmail, error)
	getAllFunc     func() ([]model.GetAll, error)
	deleteByIDFunc func(id uuid.UUID) error
	getByIDFunc    func(id uuid.UUID) (*model.GetByID, error)
}

func (m *mockUserRepo) Register(user model.User) error {
	if m.registerFunc != nil {
		return m.registerFunc(user)
	}
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*model.GetByEmail, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	return nil, nil
}

// Method 3: GetAll
func (m *mockUserRepo) GetAll() ([]model.GetAll, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil, nil
}

// Method 4: DeleteByID - THIS WAS MISSING!
func (m *mockUserRepo) DeleteByID(id uuid.UUID) error {
	if m.deleteByIDFunc != nil {
		return m.deleteByIDFunc(id)
	}
	return nil
}

// Method 5: GetByID
func (m *mockUserRepo) GetByID(id uuid.UUID) (*model.GetByID, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func TestUserService_Register_Success(t *testing.T) {

	mock := &mockUserRepo{
		registerFunc: func(user model.User) error {
			return nil
		},
	}

	service := &UserService{
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

func TestUserService_Login_Success(t *testing.T) {

	hasher := encrypt.NewPasswordHasher("test-pepper")
	hashedPassword, _ := hasher.HashPassword("password123")

	testID, _ := uuid.NewV7()

	mock := &mockUserRepo{
		getByEmailFunc: func(email string) (*model.GetByEmail, error) {

			return &model.GetByEmail{
				ID:       testID,
				Password: hashedPassword,
			}, nil
		},
	}

	service := NewUserService(
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

	service := NewUserService(
		mock,
		"",
		"",
		"",
	)

	resp, err := service.GetAll()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !slices.Equal(resp, users) {
		t.Errorf("got: %v want %v", resp, users)
	}
}
