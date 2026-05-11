package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/evbgsl/bank-service-evbgsl/internal/middleware"
	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.transactionService.Transfer(userID, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidTransferData) {
			response.Error(w, http.StatusBadRequest, "invalid transfer data")
			return
		}

		if errors.Is(err, services.ErrInvalidAmount) {
			response.Error(w, http.StatusBadRequest, "amount must be greater than zero")
			return
		}

		if errors.Is(err, services.ErrAccountNotFound) {
			response.Error(w, http.StatusNotFound, "account not found")
			return
		}

		if errors.Is(err, services.ErrForbiddenTransfer) {
			response.Error(w, http.StatusForbidden, "source account does not belong to user")
			return
		}

		if errors.Is(err, services.ErrInsufficientFunds) {
			response.Error(w, http.StatusBadRequest, "insufficient funds")
			return
		}

		if errors.Is(err, services.ErrInvalidAccountData) {
			response.Error(w, http.StatusBadRequest, "invalid account data")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *TransactionHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, transactions)
}
