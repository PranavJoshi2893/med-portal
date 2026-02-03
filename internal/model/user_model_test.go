package model

import (
	"errors"
	"testing"
)

func TestCreateUser_Validate(t *testing.T) {
	tests := []struct {
		name     string
		user     CreateUser
		wantErr  bool
		errField string
	}{
		{
			name: "valid",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Pass123!",
			},
			wantErr: false,
		},
		{
			name: "valid with hyphen and apostrophe",
			user: CreateUser{
				FirstName: "Mary-Jane",
				LastName:  "O'Brien",
				Email:     "mary@example.com",
				Password:  "Pass123!",
			},
			wantErr: false,
		},
		{
			name: "empty first name",
			user: CreateUser{
				FirstName: "",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Pass123!",
			},
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "invalid first name",
			user: CreateUser{
				FirstName: "John123",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Pass123!",
			},
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "empty last name",
			user: CreateUser{
				FirstName: "John",
				LastName:  "",
				Email:     "john@example.com",
				Password:  "Pass123!",
			},
			wantErr:  true,
			errField: "last_name",
		},
		{
			name: "empty email",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
				Password:  "Pass123!",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "invalid email",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "not-an-email",
				Password:  "Pass123!",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "empty password",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "",
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "password too short",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Pass1!",
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "password missing upper",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "pass123!",
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "password missing special",
			user: CreateUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Pass1234",
			},
			wantErr:  true,
			errField: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var vErrs ValidationErrors
				if errors.As(err, &vErrs) && tt.errField != "" {
					found := false
					for _, fe := range vErrs {
						if fe.Field == tt.errField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected field %q in errors, got %v", tt.errField, vErrs)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoginUser_Validate(t *testing.T) {
	tests := []struct {
		name     string
		user     LoginUser
		wantErr  bool
		errField string
	}{
		{
			name: "valid",
			user: LoginUser{
				Email:    "john@example.com",
				Password: "anypassword",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			user: LoginUser{
				Email:    "",
				Password: "password",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "invalid email",
			user: LoginUser{
				Email:    "@invalid",
				Password: "password",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "empty password",
			user: LoginUser{
				Email:    "john@example.com",
				Password: "",
			},
			wantErr:  true,
			errField: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var vErrs ValidationErrors
				if errors.As(err, &vErrs) && tt.errField != "" {
					found := false
					for _, fe := range vErrs {
						if fe.Field == tt.errField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected field %q in errors, got %v", tt.errField, vErrs)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestUpdateUser_Validate(t *testing.T) {
	validFirst := "John"
	validLast := "Doe"
	invalidName := "John123"
	empty := ""

	tests := []struct {
		name     string
		user     *UpdateUser
		wantErr  bool
		errField string
	}{
		{
			name: "valid",
			user: &UpdateUser{
				FirstName: &validFirst,
				LastName:  &validLast,
			},
			wantErr: false,
		},
		{
			name:     "nil",
			user:     nil,
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "nil first name",
			user: &UpdateUser{
				FirstName: nil,
				LastName:  &validLast,
			},
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "nil last name",
			user: &UpdateUser{
				FirstName: &validFirst,
				LastName:  nil,
			},
			wantErr:  true,
			errField: "last_name",
		},
		{
			name: "empty first name",
			user: &UpdateUser{
				FirstName: &empty,
				LastName:  &validLast,
			},
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "whitespace first name",
			user: &UpdateUser{
				FirstName: strPtr("   "),
				LastName:  &validLast,
			},
			wantErr:  true,
			errField: "first_name",
		},
		{
			name: "invalid first name",
			user: &UpdateUser{
				FirstName: &invalidName,
				LastName:  &validLast,
			},
			wantErr:  true,
			errField: "first_name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.user == nil {
				err = (*UpdateUser)(nil).Validate()
			} else {
				err = tt.user.Validate()
			}
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var vErrs ValidationErrors
				if errors.As(err, &vErrs) && tt.errField != "" {
					found := false
					for _, fe := range vErrs {
						if fe.Field == tt.errField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected field %q in errors, got %v", tt.errField, vErrs)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
