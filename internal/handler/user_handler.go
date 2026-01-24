package handler

import (
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

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"success",
		users,
	)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "INVALID_ID",
			Message: "Invalid User ID",
		})
	}

	var user *model.GetByID
	if user, err = h.service.GetByID(id); err != nil {
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(w, http.StatusOK, "success", user)
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
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"user deleted successfully",
		nil,
	)

}
