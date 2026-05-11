package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/evbgsl/bank-service-evbgsl/internal/middleware"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/evbgsl/bank-service-evbgsl/pkg/response"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) GetMonthlyAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	month := r.URL.Query().Get("month")

	analytics, err := h.analyticsService.GetMonthlyAnalytics(userID, month)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAnalyticsData) {
			response.Error(w, http.StatusBadRequest, "month must have format YYYY-MM")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)

	accountID, err := strconv.ParseInt(vars["accountId"], 10, 64)
	if err != nil || accountID <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	days := 30

	daysParam := r.URL.Query().Get("days")
	if daysParam != "" {
		days, err = strconv.Atoi(daysParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "days must be a number")
			return
		}
	}

	prediction, err := h.analyticsService.PredictBalance(userID, accountID, days)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAnalyticsData) {
			response.Error(w, http.StatusBadRequest, "invalid analytics data")
			return
		}

		if errors.Is(err, services.ErrInvalidPredictionDays) {
			response.Error(w, http.StatusBadRequest, "prediction days must be between 1 and 365")
			return
		}

		if errors.Is(err, services.ErrAccountNotFound) {
			response.Error(w, http.StatusNotFound, "account not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, prediction)
}
