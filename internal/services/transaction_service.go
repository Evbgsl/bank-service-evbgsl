package services

import (
	"errors"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
)

var (
	ErrInvalidTransferData = errors.New("invalid transfer data")
	ErrForbiddenTransfer   = errors.New("forbidden transfer")
	ErrInsufficientFunds   = errors.New("insufficient funds")
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Transfer(userID int64, req models.TransferRequest) (*models.TransferResponse, error) {
	if req.FromAccountID <= 0 || req.ToAccountID <= 0 {
		return nil, ErrInvalidTransferData
	}

	if req.FromAccountID == req.ToAccountID {
		return nil, ErrInvalidTransferData
	}

	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	transaction, fromBalance, err := s.transactionRepo.Transfer(
		userID,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
	)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}

		if errors.Is(err, repositories.ErrForbiddenAccount) {
			return nil, ErrForbiddenTransfer
		}

		if errors.Is(err, repositories.ErrInsufficientFunds) {
			return nil, ErrInsufficientFunds
		}

		if errors.Is(err, repositories.ErrInvalidAccountCurrency) {
			return nil, ErrInvalidAccountData
		}

		return nil, err
	}

	return &models.TransferResponse{
		TransactionID: transaction.ID,
		FromAccountID: transaction.FromAccountID,
		ToAccountID:   transaction.ToAccountID,
		Amount:        transaction.Amount,
		FromBalance:   fromBalance,
		Message:       "transfer completed successfully",
	}, nil
}

func (s *TransactionService) GetUserTransactions(userID int64) ([]models.Transaction, error) {
	return s.transactionRepo.FindByUserID(userID)
}
