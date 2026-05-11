package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Transfer(userID int64, fromAccountID int64, toAccountID int64, amount float64) (*models.Transaction, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, 0, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	fromAccount, err := getAccountForUpdate(ctx, tx, fromAccountID)
	if err != nil {
		return nil, 0, err
	}

	if fromAccount.UserID != userID {
		return nil, 0, ErrForbiddenAccount
	}

	toAccount, err := getAccountForUpdate(ctx, tx, toAccountID)
	if err != nil {
		return nil, 0, err
	}

	if fromAccount.Currency != toAccount.Currency {
		return nil, 0, ErrInvalidAccountCurrency
	}

	if fromAccount.Balance < amount {
		return nil, 0, ErrInsufficientFunds
	}

	var fromBalance float64

	err = tx.QueryRowContext(
		ctx,
		`
			UPDATE accounts
			SET balance = balance - $1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
			RETURNING balance
		`,
		amount,
		fromAccountID,
	).Scan(&fromBalance)
	if err != nil {
		return nil, 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE accounts
			SET balance = balance + $1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`,
		amount,
		toAccountID,
	)
	if err != nil {
		return nil, 0, err
	}

	transaction := &models.Transaction{
		UserID:          userID,
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		TransactionType: "TRANSFER",
		Amount:          amount,
	}

	err = tx.QueryRowContext(
		ctx,
		`
			INSERT INTO transactions (
			    user_id,
			    from_account_id,
			    to_account_id,
			    transaction_type,
			    amount
			)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at
		`,
		transaction.UserID,
		transaction.FromAccountID,
		transaction.ToAccountID,
		transaction.TransactionType,
		transaction.Amount,
	).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	return transaction, fromBalance, nil
}

func (r *TransactionRepository) FindByUserID(userID int64) ([]models.Transaction, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    COALESCE(from_account_id, 0),
		    COALESCE(to_account_id, 0),
		    transaction_type,
		    amount,
		    created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)

	for rows.Next() {
		var transaction models.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.TransactionType,
			&transaction.Amount,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func getAccountForUpdate(ctx context.Context, tx *sql.Tx, accountID int64) (*models.Account, error) {
	query := `
		SELECT id, user_id, account_number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	account := &models.Account{}

	err := tx.QueryRowContext(ctx, query, accountID).Scan(
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
