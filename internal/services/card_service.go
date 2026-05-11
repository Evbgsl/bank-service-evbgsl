package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCardData  = errors.New("invalid card data")
	ErrCardNotFound     = errors.New("card not found")
	ErrCardDataModified = errors.New("card data integrity check failed")
)

type CardService struct {
	cardRepo       *repositories.CardRepository
	accountRepo    *repositories.AccountRepository
	cardPGPKey     string
	cardHMACSecret string
}

func NewCardService(
	cardRepo *repositories.CardRepository,
	accountRepo *repositories.AccountRepository,
	cardPGPKey string,
	cardHMACSecret string,
) *CardService {
	return &CardService{
		cardRepo:       cardRepo,
		accountRepo:    accountRepo,
		cardPGPKey:     cardPGPKey,
		cardHMACSecret: cardHMACSecret,
	}
}

func (s *CardService) CreateCard(userID int64, req models.CreateCardRequest) (*models.CreateCardResponse, error) {
	if req.AccountID <= 0 {
		return nil, ErrInvalidCardData
	}

	account, err := s.accountRepo.FindByIDAndUserID(req.AccountID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	cardNumber, err := generateCardNumber()
	if err != nil {
		return nil, err
	}

	cvv, err := generateCVV()
	if err != nil {
		return nil, err
	}

	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	expiryTime := time.Now().AddDate(3, 0, 0)
	expiry := fmt.Sprintf("%02d/%d", int(expiryTime.Month()), expiryTime.Year())

	card := &models.Card{
		UserID:         userID,
		AccountID:      account.ID,
		CardNumberHMAC: computeHMAC(cardNumber, s.cardHMACSecret),
		MaskedNumber:   maskCardNumber(cardNumber),
		CVVHash:        string(cvvHash),
		Status:         "ACTIVE",
	}

	if err := s.cardRepo.Create(card, cardNumber, expiry, s.cardPGPKey); err != nil {
		return nil, err
	}

	return &models.CreateCardResponse{
		ID:           card.ID,
		AccountID:    card.AccountID,
		CardNumber:   cardNumber,
		MaskedNumber: card.MaskedNumber,
		Expiry:       expiry,
		CVV:          cvv,
		Message:      "card created successfully. Save card number and CVV now; CVV will not be shown again.",
	}, nil
}

func (s *CardService) GetUserCards(userID int64) ([]models.Card, error) {
	return s.cardRepo.FindByUserID(userID)
}

func (s *CardService) GetCardDetails(userID int64, cardID int64) (*models.CardDetailsResponse, error) {
	if cardID <= 0 {
		return nil, ErrInvalidCardData
	}

	card, err := s.cardRepo.FindByIDAndUserID(cardID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrCardNotFound) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	cardNumber, err := s.cardRepo.DecryptCardNumber(cardID, userID, s.cardPGPKey)
	if err != nil {
		if errors.Is(err, repositories.ErrCardNotFound) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	expectedHMAC := computeHMAC(cardNumber, s.cardHMACSecret)
	if !hmac.Equal([]byte(expectedHMAC), []byte(card.CardNumberHMAC)) {
		return nil, ErrCardDataModified
	}

	expiry, err := s.cardRepo.DecryptExpiry(cardID, userID, s.cardPGPKey)
	if err != nil {
		if errors.Is(err, repositories.ErrCardNotFound) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	return &models.CardDetailsResponse{
		ID:           card.ID,
		AccountID:    card.AccountID,
		CardNumber:   cardNumber,
		MaskedNumber: card.MaskedNumber,
		Expiry:       expiry,
		Status:       card.Status,
	}, nil
}

func (s *CardService) PayByCard(
	userID int64,
	cardID int64,
	req models.CardPaymentRequest,
) (*models.CardPaymentResponse, error) {
	if cardID <= 0 {
		return nil, ErrInvalidCardData
	}

	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	merchant := strings.TrimSpace(req.Merchant)
	if merchant == "" {
		merchant = "Unknown merchant"
	}

	resp, err := s.cardRepo.PayByCard(userID, cardID, req.Amount)
	if err != nil {
		if errors.Is(err, repositories.ErrCardNotFound) {
			return nil, ErrCardNotFound
		}

		if errors.Is(err, repositories.ErrInsufficientFunds) {
			return nil, ErrInsufficientFunds
		}

		return nil, err
	}

	resp.Merchant = merchant

	return resp, nil
}

func generateCardNumber() (string, error) {
	const prefix = "2202"

	randomPart, err := randomCardDigits(11)
	if err != nil {
		return "", err
	}

	numberWithoutCheckDigit := prefix + randomPart
	checkDigit := calculateLuhnCheckDigit(numberWithoutCheckDigit)

	return numberWithoutCheckDigit + strconv.Itoa(checkDigit), nil
}

func calculateLuhnCheckDigit(numberWithoutCheckDigit string) int {
	sum := 0
	double := true

	for i := len(numberWithoutCheckDigit) - 1; i >= 0; i-- {
		digit := int(numberWithoutCheckDigit[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return (10 - (sum % 10)) % 10
}

func generateCVV() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", n.Int64()), nil
}

func randomCardDigits(length int) (string, error) {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}

		result[i] = byte('0' + n.Int64())
	}

	return string(result), nil
}

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) < 10 {
		return cardNumber
	}

	return cardNumber[:6] + "******" + cardNumber[len(cardNumber)-4:]
}

func computeHMAC(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}
