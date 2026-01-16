package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/service"
	"github.com/PranavJoshi2893/med-portal/pkg/responses"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
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
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	if err := h.service.Register(&user); err != nil {
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusCreated,
		"New user registered",
		nil,
	)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	if err := h.service.Login(&user); err != nil {
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"Login Successful",
		nil,
	)

}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"OK",
		users,
	)
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "INVALID_ID",
			Message: "Invalid User ID",
		})
		return
	}

	if err := h.service.DeleteByID(id); err != nil {
		responses.WriteError(w, responses.FromModelError(err))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"User Deleted Successfully",
		nil,
	)

}
