package models

import "time"

type Transaction struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"userId"`
	FromAccountID   int64     `json:"fromAccountId"`
	ToAccountID     int64     `json:"toAccountId"`
	TransactionType string    `json:"transactionType"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"createdAt"`
}

type TransferRequest struct {
	FromAccountID int64   `json:"fromAccountId"`
	ToAccountID   int64   `json:"toAccountId"`
	Amount        float64 `json:"amount"`
}

type TransferResponse struct {
	TransactionID int64   `json:"transactionId"`
	FromAccountID int64   `json:"fromAccountId"`
	ToAccountID   int64   `json:"toAccountId"`
	Amount        float64 `json:"amount"`
	FromBalance   float64 `json:"fromBalance"`
	Message       string  `json:"message"`
}
