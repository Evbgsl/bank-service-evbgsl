package models

import "time"

type Account struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userId"`
	AccountNumber string    `json:"accountNumber"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	ID            int64   `json:"id"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	Message       string  `json:"message"`
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type DepositResponse struct {
	AccountID int64   `json:"accountId"`
	Balance   float64 `json:"balance"`
	Message   string  `json:"message"`
}
