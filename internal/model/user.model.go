package model

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type CreateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type GetAll struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type DeleteUser struct {
	ID uuid.UUID `json:"id"`
}

func (m *CreateUser) Validate() error {
	var errs ValidationErrors

	// normalize
	m.FirstName = strings.TrimSpace(m.FirstName)
	m.LastName = strings.TrimSpace(m.LastName)
	m.Email = strings.TrimSpace(m.Email)

	if m.FirstName == "" {
		errs = append(errs, FieldError{
			Field:   "first_name",
			Message: "first name is required",
		})
	} else if !isValidName(m.FirstName) {
		errs = append(errs, FieldError{
			Field:   "first_name",
			Message: "invalid first name",
		})
	}

	if m.LastName == "" {
		errs = append(errs, FieldError{
			Field:   "last_name",
			Message: "last name is required",
		})
	} else if !isValidName(m.LastName) {
		errs = append(errs, FieldError{
			Field:   "last_name",
			Message: "invalid last name",
		})
	}

	if m.Email == "" {
		errs = append(errs, FieldError{
			Field:   "email",
			Message: "email is required",
		})
	} else {
		addr, err := mail.ParseAddress(m.Email)
		if err != nil || addr.Address != m.Email {
			errs = append(errs, FieldError{
				Field:   "email",
				Message: "invalid email",
			})
		} else {
			m.Email = strings.ToLower(m.Email)
		}
	}

	if m.Password == "" {
		errs = append(errs, FieldError{
			Field:   "password",
			Message: "password is required",
		})
	} else {
		if len(m.Password) < 8 {
			errs = append(errs, FieldError{
				Field:   "password",
				Message: "password must be at least 8 characters",
			})
		}

		if !isValidPassword(m.Password) {
			errs = append(errs, FieldError{
				Field:   "password",
				Message: "password must contain upper, lower, digit and special character",
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (m *LoginUser) Validate() error {

	m.Email = strings.TrimSpace(m.Email)

	if m.Email == "" {
		return fmt.Errorf("email is required")
	}

	if m.Password == "" {
		return fmt.Errorf("password is required")
	}

	addr, err := mail.ParseAddress(m.Email)
	if err != nil || addr.Address != m.Email {
		return fmt.Errorf("invalid email")
	}

	// normalize email
	m.Email = strings.ToLower(m.Email)

	return nil
}

func isValidName(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && char != '-' && char != '\'' {
			return false
		}
	}
	return true
}

func isValidPassword(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecialCharacter bool

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}

		if unicode.IsLower(char) {
			hasLower = true
		}

		if unicode.IsDigit(char) {
			hasDigit = true
		}

		if unicode.IsPunct(char) {
			hasSpecialCharacter = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecialCharacter
}
