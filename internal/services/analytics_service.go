package services

import (
	"errors"
	"math"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
)

var (
	ErrInvalidAnalyticsData  = errors.New("invalid analytics data")
	ErrInvalidPredictionDays = errors.New("prediction days must be between 1 and 365")
)

type AnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
	}
}

func (s *AnalyticsService) GetMonthlyAnalytics(userID int64, month string) (*models.MonthlyAnalyticsResponse, error) {
	from, to, normalizedMonth, err := parseMonthPeriod(month)
	if err != nil {
		return nil, ErrInvalidAnalyticsData
	}

	income, err := s.analyticsRepo.GetMonthlyIncome(userID, from, to)
	if err != nil {
		return nil, err
	}

	expenses, err := s.analyticsRepo.GetMonthlyExpenses(userID, from, to)
	if err != nil {
		return nil, err
	}

	activeCreditsCount, err := s.analyticsRepo.GetActiveCreditsCount(userID)
	if err != nil {
		return nil, err
	}

	totalMonthlyPayments, err := s.analyticsRepo.GetTotalMonthlyCreditPayments(userID)
	if err != nil {
		return nil, err
	}

	totalRemainingDebt, err := s.analyticsRepo.GetTotalRemainingDebt(userID)
	if err != nil {
		return nil, err
	}

	income = roundAnalyticsMoney(income)
	expenses = roundAnalyticsMoney(expenses)

	return &models.MonthlyAnalyticsResponse{
		Month:    normalizedMonth,
		Income:   income,
		Expenses: expenses,
		Net:      roundAnalyticsMoney(income - expenses),
		CreditLoad: models.CreditLoadResponse{
			ActiveCreditsCount:   activeCreditsCount,
			TotalMonthlyPayments: roundAnalyticsMoney(totalMonthlyPayments),
			TotalRemainingDebt:   roundAnalyticsMoney(totalRemainingDebt),
		},
	}, nil
}

func (s *AnalyticsService) PredictBalance(userID int64, accountID int64, days int) (*models.BalancePredictionResponse, error) {
	if accountID <= 0 {
		return nil, ErrInvalidAnalyticsData
	}

	if days <= 0 || days > 365 {
		return nil, ErrInvalidPredictionDays
	}

	currentBalance, err := s.analyticsRepo.GetAccountBalance(accountID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	toDate := time.Now().AddDate(0, 0, days)

	plannedPayments, err := s.analyticsRepo.GetPlannedPaymentsForAccount(accountID, userID, toDate)
	if err != nil {
		return nil, err
	}

	currentBalance = roundAnalyticsMoney(currentBalance)
	plannedPayments = roundAnalyticsMoney(plannedPayments)
	predictedBalance := roundAnalyticsMoney(currentBalance - plannedPayments)

	return &models.BalancePredictionResponse{
		AccountID:        accountID,
		Days:             days,
		CurrentBalance:   currentBalance,
		PlannedPayments:  plannedPayments,
		PredictedBalance: predictedBalance,
		Message:          "prediction includes planned and overdue credit payments only",
	}, nil
}

func parseMonthPeriod(month string) (time.Time, time.Time, string, error) {
	if month == "" {
		now := time.Now()
		from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		to := from.AddDate(0, 1, 0)

		return from, to, from.Format("2006-01"), nil
	}

	from, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, "", err
	}

	to := from.AddDate(0, 1, 0)

	return from, to, from.Format("2006-01"), nil
}

func roundAnalyticsMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
