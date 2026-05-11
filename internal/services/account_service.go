package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
)

var (
	ErrInvalidAccountData = errors.New("invalid account data")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrAccountNotFound    = errors.New("account not found")
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) CreateAccount(userID int64, req models.CreateAccountRequest) (*models.CreateAccountResponse, error) {
	currency := strings.TrimSpace(strings.ToUpper(req.Currency))
	if currency == "" {
		currency = "RUB"
	}

	if currency != "RUB" {
		return nil, ErrInvalidAccountData
	}

	accountNumber, err := generateAccountNumber()
	if err != nil {
		return nil, err
	}

	account := &models.Account{
		UserID:        userID,
		AccountNumber: accountNumber,
		Balance:       0,
		Currency:      currency,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return &models.CreateAccountResponse{
		ID:            account.ID,
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
		Currency:      account.Currency,
		Message:       "account created successfully",
	}, nil
}

func (s *AccountService) GetUserAccounts(userID int64) ([]models.Account, error) {
	return s.accountRepo.FindByUserID(userID)
}

func (s *AccountService) Deposit(userID int64, accountID int64, req models.DepositRequest) (*models.DepositResponse, error) {
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	account, err := s.accountRepo.Deposit(accountID, userID, req.Amount)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return &models.DepositResponse{
		AccountID: account.ID,
		Balance:   account.Balance,
		Message:   "account deposited successfully",
	}, nil
}

func generateAccountNumber() (string, error) {
	const prefix = "40817810"

	randomPart, err := randomDigits(12)
	if err != nil {
		return "", err
	}

	return prefix + randomPart, nil
}

func randomDigits(length int) (string, error) {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}

		result[i] = byte('0' + n.Int64())
	}

	return fmt.Sprintf("%s", result), nil
}
