package models

type KeyRateResponse struct {
	CentralBankRate float64 `json:"centralBankRate"`
	BankMargin      float64 `json:"bankMargin"`
	FinalRate       float64 `json:"finalRate"`
	Message         string  `json:"message"`
}

type TestEmailRequest struct {
	To string `json:"to"`
}

type TestEmailResponse struct {
	Message string `json:"message"`
}
