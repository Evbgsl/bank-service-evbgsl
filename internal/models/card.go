package models

import "time"

type Card struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"userId"`
	AccountID           int64     `json:"accountId"`
	CardNumberEncrypted []byte    `json:"-"`
	ExpiryEncrypted     []byte    `json:"-"`
	CardNumberHMAC      string    `json:"-"`
	MaskedNumber        string    `json:"maskedNumber"`
	CVVHash             string    `json:"-"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type CreateCardRequest struct {
	AccountID int64 `json:"accountId"`
}

type CreateCardResponse struct {
	ID           int64  `json:"id"`
	AccountID    int64  `json:"accountId"`
	CardNumber   string `json:"cardNumber"`
	MaskedNumber string `json:"maskedNumber"`
	Expiry       string `json:"expiry"`
	CVV          string `json:"cvv"`
	Message      string `json:"message"`
}

type CardDetailsResponse struct {
	ID           int64  `json:"id"`
	AccountID    int64  `json:"accountId"`
	CardNumber   string `json:"cardNumber"`
	MaskedNumber string `json:"maskedNumber"`
	Expiry       string `json:"expiry"`
	Status       string `json:"status"`
}
