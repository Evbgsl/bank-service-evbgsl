package repositories

import (
	"database/sql"
	"errors"
	"time"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (r *AnalyticsRepository) GetMonthlyIncome(userID int64, from time.Time, to time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		WHERE t.to_account_id IN (
		    SELECT id
		    FROM accounts
		    WHERE user_id = $1
		)
		  AND t.created_at >= $2
		  AND t.created_at < $3
	`

	var income float64

	err := r.db.QueryRow(query, userID, from, to).Scan(&income)
	if err != nil {
		return 0, err
	}

	return income, nil
}

func (r *AnalyticsRepository) GetMonthlyExpenses(userID int64, from time.Time, to time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		WHERE t.from_account_id IN (
		    SELECT id
		    FROM accounts
		    WHERE user_id = $1
		)
		  AND t.created_at >= $2
		  AND t.created_at < $3
	`

	var expenses float64

	err := r.db.QueryRow(query, userID, from, to).Scan(&expenses)
	if err != nil {
		return 0, err
	}

	return expenses, nil
}

func (r *AnalyticsRepository) GetActiveCreditsCount(userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM credits
		WHERE user_id = $1
		  AND status = 'ACTIVE'
	`

	var count int

	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AnalyticsRepository) GetTotalMonthlyCreditPayments(userID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(monthly_payment), 0)
		FROM credits
		WHERE user_id = $1
		  AND status = 'ACTIVE'
	`

	var total float64

	err := r.db.QueryRow(query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *AnalyticsRepository) GetTotalRemainingDebt(userID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(remaining_amount), 0)
		FROM credits
		WHERE user_id = $1
		  AND status = 'ACTIVE'
	`

	var total float64

	err := r.db.QueryRow(query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *AnalyticsRepository) GetAccountBalance(accountID int64, userID int64) (float64, error) {
	query := `
		SELECT balance
		FROM accounts
		WHERE id = $1 AND user_id = $2
	`

	var balance float64

	err := r.db.QueryRow(query, accountID, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrAccountNotFound
		}

		return 0, err
	}

	return balance, nil
}

func (r *AnalyticsRepository) GetPlannedPaymentsForAccount(accountID int64, userID int64, toDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(ps.amount), 0)
		FROM payment_schedules ps
		JOIN credits c ON c.id = ps.credit_id
		WHERE c.account_id = $1
		  AND c.user_id = $2
		  AND c.status = 'ACTIVE'
		  AND ps.status IN ('PLANNED', 'OVERDUE')
		  AND ps.payment_date <= $3
	`

	var total float64

	err := r.db.QueryRow(query, accountID, userID, toDate).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
