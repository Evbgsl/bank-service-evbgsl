package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/evbgsl/bank-service-evbgsl/internal/middleware"
	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
)

type CreditHandler struct {
	creditService *services.CreditService
}

func NewCreditHandler(creditService *services.CreditService) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
	}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.creditService.CreateCredit(userID, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCreditData) {
			response.Error(w, http.StatusBadRequest, "invalid credit data")
			return
		}

		if errors.Is(err, services.ErrAccountNotFound) {
			response.Error(w, http.StatusNotFound, "account not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}

func (h *CreditHandler) GetUserCredits(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	credits, err := h.creditService.GetUserCredits(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, credits)
}

func (h *CreditHandler) GetCreditSchedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)

	creditID, err := strconv.ParseInt(vars["creditId"], 10, 64)
	if err != nil || creditID <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid credit id")
		return
	}

	schedule, err := h.creditService.GetCreditSchedule(userID, creditID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCreditData) {
			response.Error(w, http.StatusBadRequest, "invalid credit data")
			return
		}

		if errors.Is(err, services.ErrCreditNotFound) {
			response.Error(w, http.StatusNotFound, "credit not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, schedule)
}
