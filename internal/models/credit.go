package models

import "time"

type Credit struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"userId"`
	AccountID       int64     `json:"accountId"`
	PrincipalAmount float64   `json:"principalAmount"`
	InterestRate    float64   `json:"interestRate"`
	TermMonths      int       `json:"termMonths"`
	MonthlyPayment  float64   `json:"monthlyPayment"`
	RemainingAmount float64   `json:"remainingAmount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type PaymentSchedule struct {
	ID            int64      `json:"id"`
	CreditID      int64      `json:"creditId"`
	PaymentNumber int        `json:"paymentNumber"`
	PaymentDate   time.Time  `json:"paymentDate"`
	Amount        float64    `json:"amount"`
	PrincipalPart float64    `json:"principalPart"`
	InterestPart  float64    `json:"interestPart"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paidAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type CreateCreditRequest struct {
	AccountID       int64   `json:"accountId"`
	PrincipalAmount float64 `json:"principalAmount"`
	InterestRate    float64 `json:"interestRate"`
	TermMonths      int     `json:"termMonths"`
}

type CreateCreditResponse struct {
	ID              int64   `json:"id"`
	AccountID       int64   `json:"accountId"`
	PrincipalAmount float64 `json:"principalAmount"`
	InterestRate    float64 `json:"interestRate"`
	TermMonths      int     `json:"termMonths"`
	MonthlyPayment  float64 `json:"monthlyPayment"`
	RemainingAmount float64 `json:"remainingAmount"`
	Status          string  `json:"status"`
	Message         string  `json:"message"`
}
