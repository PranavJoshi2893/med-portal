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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (m *CreateUser) Validate() error {

	// trim trailing and leading space
	m.FirstName = strings.TrimSpace(m.FirstName)
	m.LastName = strings.TrimSpace(m.LastName)
	m.Email = strings.TrimSpace(m.Email)

	if m.FirstName == "" {
		return fmt.Errorf("first_name is required")
	}

	if m.LastName == "" {
		return fmt.Errorf("last_name is required")
	}

	if m.Email == "" {
		return fmt.Errorf("email is required")
	}

	if m.Password == "" {
		return fmt.Errorf("password is required")
	}

	// Letters (Unicode)
	// Space
	// Hyphen -
	// Apostrophe '
	if !isValidName(m.FirstName) {
		return fmt.Errorf("Invalid firstname")
	}

	if !isValidName(m.LastName) {
		return fmt.Errorf("Invalid lastname")
	}

	addr, err := mail.ParseAddress(m.Email)
	if err != nil || addr.Address != m.Email {
		return fmt.Errorf("invalid email")
	}

	// normalize email
	m.Email = strings.ToLower(m.Email)

	if len(m.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if !isValidPassword(m.Password) {
		return fmt.Errorf("Password must contain at least upper, lower, digit and special character ")
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
