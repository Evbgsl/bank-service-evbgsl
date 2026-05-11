package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCardData = errors.New("invalid card data")
	ErrCardNotFound    = errors.New("card not found")
)

type CardService struct {
	cardRepo    *repositories.CardRepository
	accountRepo *repositories.AccountRepository
}

func NewCardService(cardRepo *repositories.CardRepository, accountRepo *repositories.AccountRepository) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
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

	expiry := time.Now().AddDate(3, 0, 0)

	card := &models.Card{
		UserID:       userID,
		AccountID:    account.ID,
		CardNumber:   cardNumber,
		MaskedNumber: maskCardNumber(cardNumber),
		ExpiryMonth:  int(expiry.Month()),
		ExpiryYear:   expiry.Year(),
		CVVHash:      string(cvvHash),
		Status:       "ACTIVE",
	}

	if err := s.cardRepo.Create(card); err != nil {
		return nil, err
	}

	return &models.CreateCardResponse{
		ID:           card.ID,
		AccountID:    card.AccountID,
		CardNumber:   cardNumber,
		MaskedNumber: card.MaskedNumber,
		ExpiryMonth:  card.ExpiryMonth,
		ExpiryYear:   card.ExpiryYear,
		CVV:          cvv,
		Message:      "card created successfully. Save card number and CVV now; CVV will not be shown again.",
	}, nil
}

func (s *CardService) GetUserCards(userID int64) ([]models.Card, error) {
	return s.cardRepo.FindByUserID(userID)
}

func generateCardNumber() (string, error) {
	const prefix = "2202"

	base := prefix

	randomPart, err := randomCardDigits(11)
	if err != nil {
		return "", err
	}

	base += randomPart

	checkDigit := calculateLuhnCheckDigit(base)

	return base + strconv.Itoa(checkDigit), nil
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
