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

type CardHandler struct {
	cardService *services.CardService
}

func NewCardHandler(cardService *services.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
	}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.cardService.CreateCard(userID, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCardData) {
			response.Error(w, http.StatusBadRequest, "invalid card data")
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

func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	cards, err := h.cardService.GetUserCards(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, cards)
}

func (h *CardHandler) GetCardDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)

	cardID, err := strconv.ParseInt(vars["cardId"], 10, 64)
	if err != nil || cardID <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid card id")
		return
	}

	cardDetails, err := h.cardService.GetCardDetails(userID, cardID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCardData) {
			response.Error(w, http.StatusBadRequest, "invalid card data")
			return
		}

		if errors.Is(err, services.ErrCardNotFound) {
			response.Error(w, http.StatusNotFound, "card not found")
			return
		}

		if errors.Is(err, services.ErrCardDataModified) {
			response.Error(w, http.StatusConflict, "card data integrity check failed")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, cardDetails)
}

func (h *CardHandler) PayByCard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)

	cardID, err := strconv.ParseInt(vars["cardId"], 10, 64)
	if err != nil || cardID <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid card id")
		return
	}

	var req models.CardPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.cardService.PayByCard(userID, cardID, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCardData) {
			response.Error(w, http.StatusBadRequest, "invalid card data")
			return
		}

		if errors.Is(err, services.ErrInvalidAmount) {
			response.Error(w, http.StatusBadRequest, "amount must be greater than zero")
			return
		}

		if errors.Is(err, services.ErrCardNotFound) {
			response.Error(w, http.StatusNotFound, "card not found")
			return
		}

		if errors.Is(err, services.ErrInsufficientFunds) {
			response.Error(w, http.StatusBadRequest, "insufficient funds")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
