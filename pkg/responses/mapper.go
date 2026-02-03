package responses

import (
	"errors"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

func FromModelError(err error, message string) ErrorResponse {
	var vErrs model.ValidationErrors

	switch {
	case errors.As(err, &vErrs):
		return ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Status:  "VALIDATION_ERROR",
			Message: "validation error",
			Fields:  vErrs,
		}

	case errors.Is(err, model.ErrAlreadyExists):
		return ErrorResponse{
			Code:    http.StatusConflict,
			Status:  "ALREADY_EXISTS",
			Message: message,
		}

	case errors.Is(err, model.ErrNotFound):
		return ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "NOT_FOUND",
			Message: message,
		}

	case errors.Is(err, model.ErrAlreadyDeleted):
		return ErrorResponse{
			Code:    http.StatusGone,
			Status:  "ALREADY_DELETED",
			Message: message,
		}

	case errors.Is(err, model.ErrUnauthorized):
		return ErrorResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: message,
		}

	case errors.Is(err, model.ErrBadRequest):
		return ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Message: message,
		}

	default:
		return ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "INTERNAL_ERROR",
			Message: "Internal server error",
		}
	}
}
