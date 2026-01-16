package responses

import (
	"errors"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

func FromModelError(err error) ErrorResponse {
	var vErrs model.ValidationErrors

	switch {
	case errors.As(err, &vErrs):
		return ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Status:  "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  vErrs,
		}

	case errors.Is(err, model.ErrUserAlreadyExists):
		return ErrorResponse{
			Code:    http.StatusConflict,
			Status:  "ALREADY_EXISTS",
			Message: "User already exists",
		}

	case errors.Is(err, model.ErrUserNotFound):
		return ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "NOT_FOUND",
			Message: "User not found",
		}

	case errors.Is(err, model.ErrAlreadyDeleted):
		return ErrorResponse{
			Code:    http.StatusGone,
			Status:  "ALREADY_DELETED",
			Message: "Already deleted or does not exist",
		}

	default:
		return ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "INTERNAL_ERROR",
			Message: "Internal server error",
		}
	}
}
