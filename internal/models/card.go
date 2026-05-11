package models

import "time"

type Card struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	AccountID    int64     `json:"accountId"`
	CardNumber   string    `json:"-"`
	MaskedNumber string    `json:"maskedNumber"`
	ExpiryMonth  int       `json:"expiryMonth"`
	ExpiryYear   int       `json:"expiryYear"`
	CVVHash      string    `json:"-"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreateCardRequest struct {
	AccountID int64 `json:"accountId"`
}

type CreateCardResponse struct {
	ID           int64  `json:"id"`
	AccountID    int64  `json:"accountId"`
	CardNumber   string `json:"cardNumber"`
	MaskedNumber string `json:"maskedNumber"`
	ExpiryMonth  int    `json:"expiryMonth"`
	ExpiryYear   int    `json:"expiryYear"`
	CVV          string `json:"cvv"`
	Message      string `json:"message"`
}
