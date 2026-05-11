package models

type MonthlyAnalyticsResponse struct {
	Month      string             `json:"month"`
	Income     float64            `json:"income"`
	Expenses   float64            `json:"expenses"`
	Net        float64            `json:"net"`
	CreditLoad CreditLoadResponse `json:"creditLoad"`
}

type CreditLoadResponse struct {
	ActiveCreditsCount   int     `json:"activeCreditsCount"`
	TotalMonthlyPayments float64 `json:"totalMonthlyPayments"`
	TotalRemainingDebt   float64 `json:"totalRemainingDebt"`
}

type BalancePredictionResponse struct {
	AccountID        int64   `json:"accountId"`
	Days             int     `json:"days"`
	CurrentBalance   float64 `json:"currentBalance"`
	PlannedPayments  float64 `json:"plannedPayments"`
	PredictedBalance float64 `json:"predictedBalance"`
	Message          string  `json:"message"`
}
