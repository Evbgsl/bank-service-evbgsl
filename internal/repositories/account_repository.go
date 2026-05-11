package repositories

import (
	"database/sql"
	"errors"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/lib/pq"
)

var (
	ErrAccountNotFound        = errors.New("account not found")
	ErrAccountAlreadyExists   = errors.New("account already exists")
	ErrForbiddenAccount       = errors.New("account does not belong to user")
	ErrInvalidAccountCurrency = errors.New("invalid account currency")
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (user_id, account_number, balance, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		account.UserID,
		account.AccountNumber,
		account.Balance,
		account.Currency,
	).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrAccountAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *AccountRepository) FindByUserID(userID int64) ([]models.Account, error) {
	query := `
		SELECT id, user_id, account_number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]models.Account, 0)

	for rows.Next() {
		var account models.Account

		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountNumber,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) FindByIDAndUserID(accountID int64, userID int64) (*models.Account, error) {
	query := `
		SELECT id, user_id, account_number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1 AND user_id = $2
	`

	account := &models.Account{}

	err := r.db.QueryRow(query, accountID, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) Deposit(accountID int64, userID int64, amount float64) (*models.Account, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, account_number, balance, currency, created_at, updated_at
	`

	account := &models.Account{}

	err := r.db.QueryRow(query, amount, accountID, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountNumber,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	return account, nil
}
