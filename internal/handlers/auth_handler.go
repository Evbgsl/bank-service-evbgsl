package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRegisterData) {
			response.Error(w, http.StatusBadRequest, "invalid username, email or password")
			return
		}

		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Error(w, http.StatusConflict, "user with this email or username already exists")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}
