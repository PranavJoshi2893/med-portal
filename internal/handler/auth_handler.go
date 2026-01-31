package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/service"
	"github.com/PranavJoshi2893/med-portal/pkg/responses"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.CreateUser

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	defer r.Body.Close()

	if err := dec.Decode(&user); err != nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "INVALID_JSON",
			Message: "Invalid JSON payload",
		})
		return
	}

	if err := user.Validate(); err != nil {
		responses.WriteError(w, responses.FromModelError(err, ""))
		return
	}

	ctx := r.Context()

	if err := h.service.Register(ctx,&user); err != nil {
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusCreated,
		"user registered successfully",
		nil,
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.LoginUser

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	defer r.Body.Close()

	if err := dec.Decode(&user); err != nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "INVALID_JSON",
			Message: "Invalid JSON payload",
		})
		return
	}

	if err := user.Validate(); err != nil {
		responses.WriteError(w, responses.FromModelError(err, ""))
		return
	}

	ctx := r.Context()

	data, err := h.service.Login(ctx,&user)
	if err != nil {
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"login successful",
		&data,
	)

}
