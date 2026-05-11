package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
)

type IntegrationHandler struct {
	cbrService   *services.CBRService
	emailService *services.EmailService
}

func NewIntegrationHandler(
	cbrService *services.CBRService,
	emailService *services.EmailService,
) *IntegrationHandler {
	return &IntegrationHandler{
		cbrService:   cbrService,
		emailService: emailService,
	}
}

func (h *IntegrationHandler) GetKeyRate(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cbrService.GetKeyRate()
	if err != nil {
		if errors.Is(err, services.ErrKeyRateNotFound) {
			response.Error(w, http.StatusNotFound, "key rate not found")
			return
		}

		response.Error(w, http.StatusBadGateway, "central bank service unavailable")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *IntegrationHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	var req models.TestEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.emailService.SendTestEmail(req.To); err != nil {
		if errors.Is(err, services.ErrInvalidEmail) {
			response.Error(w, http.StatusBadRequest, "invalid email")
			return
		}

		response.Error(w, http.StatusBadGateway, "email sending failed")
		return
	}

	response.JSON(w, http.StatusOK, models.TestEmailResponse{
		Message: "test email processed successfully",
	})
}
