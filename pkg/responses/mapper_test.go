package responses

import (
	"errors"
	"net/http"
	"testing"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

func TestFromModelError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		message  string
		wantCode int
		wantStat string
	}{
		{
			name:     "validation errors",
			err:      model.ValidationErrors{{Field: "email", Message: "invalid"}},
			message:  "",
			wantCode: http.StatusUnprocessableEntity,
			wantStat: "VALIDATION_ERROR",
		},
		{
			name:     "already exists",
			err:      model.ErrAlreadyExists,
			message:  "email already exists",
			wantCode: http.StatusConflict,
			wantStat: "ALREADY_EXISTS",
		},
		{
			name:     "not found",
			err:      model.ErrNotFound,
			message:  "user not found",
			wantCode: http.StatusNotFound,
			wantStat: "NOT_FOUND",
		},
		{
			name:     "already deleted",
			err:      model.ErrAlreadyDeleted,
			message:  "already deleted",
			wantCode: http.StatusGone,
			wantStat: "ALREADY_DELETED",
		},
		{
			name:     "unauthorized",
			err:      model.ErrUnauthorized,
			message:  "unauthorized",
			wantCode: http.StatusUnauthorized,
			wantStat: "UNAUTHORIZED",
		},
		{
			name:     "bad request",
			err:      model.ErrBadRequest,
			message:  "bad request",
			wantCode: http.StatusBadRequest,
			wantStat: "BAD_REQUEST",
		},
		{
			name:     "forbidden",
			err:      model.ErrForbidden,
			message:  "forbidden",
			wantCode: http.StatusForbidden,
			wantStat: "FORBIDDEN",
		},
		{
			name:     "unknown error",
			err:      errors.New("unknown"),
			message:  "",
			wantCode: http.StatusInternalServerError,
			wantStat: "INTERNAL_ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromModelError(tt.err, tt.message)
			if got.Code != tt.wantCode {
				t.Errorf("Code: got %d want %d", got.Code, tt.wantCode)
			}
			if got.Status != tt.wantStat {
				t.Errorf("Status: got %q want %q", got.Status, tt.wantStat)
			}
		})
	}
}
