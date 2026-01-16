package model

import "errors"

var (
	ErrAlreadyExists    = errors.New("already exists")        // 409
	ErrNotFound         = errors.New("not found")             // 404
	ErrAlreadyDeleted   = errors.New("already deleted")       // 410
	ErrUnauthorized     = errors.New("unauthorized")          // 401
	ErrForbidden        = errors.New("forbidden")             // 403
	ErrBadRequest       = errors.New("bad request")           // 400
	ErrInternal         = errors.New("internal server error") // 500
	ErrConflict         = errors.New("conflict")              // 409
	ErrValidationFailed = errors.New("validation failed")     // 422
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	return "validation failed"
}
