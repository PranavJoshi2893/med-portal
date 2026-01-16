package model

import "errors"

var (
	ErrUserAlreadyExists = errors.New("already exists")
	ErrUserNotFound      = errors.New("not found")
	ErrAlreadyDeleted    = errors.New("already deleted or does not exist")
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	return "validation failed"
}
