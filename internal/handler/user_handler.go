package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
	"github.com/PranavJoshi2893/med-portal/internal/service"
	"github.com/PranavJoshi2893/med-portal/pkg/responses"
	"github.com/google/uuid"
)

func getCallerFromContext(ctx interface{ Value(any) any }) (*uuid.UUID, string) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, ""
	}
	role, _ := ctx.Value("role").(string)
	if role == "" {
		role = "user"
	}
	return &userID, role
}

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
	callerID, callerRole := getCallerFromContext(ctx)
	if callerID == nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "Unauthorized",
		})
		return
	}
	users, err := h.service.GetAll(ctx, *callerID, callerRole)
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
	callerID, callerRole := getCallerFromContext(ctx)
	if callerID == nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "Unauthorized",
		})
		return
	}

	var user *model.GetByID
	if user, err = h.service.GetByID(ctx, id, *callerID, callerRole); err != nil {
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
	callerID, callerRole := getCallerFromContext(ctx)
	if callerID == nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "Unauthorized",
		})
		return
	}

	if err := h.service.DeleteByID(ctx, id, *callerID, callerRole); err != nil {
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

	if err := user.Validate(); err != nil {
		responses.WriteError(w, responses.FromModelError(err, ""))
		return
	}

	ctx := r.Context()
	callerID, callerRole := getCallerFromContext(ctx)
	if callerID == nil {
		responses.WriteError(w, responses.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "Unauthorized",
		})
		return
	}
	err = h.service.UpdateByID(ctx, id, user, *callerID, callerRole)
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
