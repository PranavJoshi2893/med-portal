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

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		if vErrs, ok := err.(model.ValidationErrors); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)

			json.NewEncoder(w).Encode(map[string]any{
				"message": "validation failed",
				"errors":  vErrs,
			})
			return
		}

		// unexpected validation error (should not happen)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.service.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusCreated, "New user registered", nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Login(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusOK, "Login successful", nil)

}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "Something bad happen when retriving data", http.StatusBadRequest)
	}

	responses.JSONResponse(w, http.StatusOK, "Ok", users)
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses.JSONResponse(w, http.StatusOK, "Ok", nil)

}
