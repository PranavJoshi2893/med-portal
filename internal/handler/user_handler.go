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

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.service.GetAll(ctx)
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
		return
	}

	ctx := r.Context()

	var user *model.GetByID
	if user, err = h.service.GetByID(ctx, id); err != nil {
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

	ctx := r.Context()

	if err := h.service.DeleteByID(ctx, id); err != nil {
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

func (h *UserHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "INVALID_ID",
			Message: "Invalid User ID",
		})

		return
	}

	var user *model.UpdateUser

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

	ctx := r.Context()
	err = h.service.UpdateByID(ctx, id, user)
	if err != nil {
		responses.WriteError(w, responses.FromModelError(err, err.Error()))
		return
	}

	responses.WriteSuccess(
		w,
		http.StatusOK,
		"user updated successfully",
		nil,
	)
}
