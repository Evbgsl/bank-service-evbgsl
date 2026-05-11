package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
)

var ErrCreditNotFound = errors.New("credit not found")

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{
		db: db,
	}
}

func (r *CreditRepository) CreateCreditWithSchedule(
	credit *models.Credit,
	schedule []models.PaymentSchedule,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = tx.QueryRowContext(
		ctx,
		`
			INSERT INTO credits (
			    user_id,
			    account_id,
			    principal_amount,
			    interest_rate,
			    term_months,
			    monthly_payment,
			    remaining_amount,
			    status
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at
		`,
		credit.UserID,
		credit.AccountID,
		credit.PrincipalAmount,
		credit.InterestRate,
		credit.TermMonths,
		credit.MonthlyPayment,
		credit.RemainingAmount,
		credit.Status,
	).Scan(
		&credit.ID,
		&credit.CreatedAt,
		&credit.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE accounts
			SET balance = balance + $1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND user_id = $3
		`,
		credit.PrincipalAmount,
		credit.AccountID,
		credit.UserID,
	)
	if err != nil {
		return err
	}

	for i := range schedule {
		schedule[i].CreditID = credit.ID

		err = tx.QueryRowContext(
			ctx,
			`
				INSERT INTO payment_schedules (
				    credit_id,
				    payment_number,
				    payment_date,
				    amount,
				    principal_part,
				    interest_part,
				    status
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id, created_at
			`,
			schedule[i].CreditID,
			schedule[i].PaymentNumber,
			schedule[i].PaymentDate,
			schedule[i].Amount,
			schedule[i].PrincipalPart,
			schedule[i].InterestPart,
			schedule[i].Status,
		).Scan(
			&schedule[i].ID,
			&schedule[i].CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CreditRepository) FindByUserID(userID int64) ([]models.Credit, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    account_id,
		    principal_amount,
		    interest_rate,
		    term_months,
		    monthly_payment,
		    remaining_amount,
		    status,
		    created_at,
		    updated_at
		FROM credits
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credits := make([]models.Credit, 0)

	for rows.Next() {
		var credit models.Credit

		if err := rows.Scan(
			&credit.ID,
			&credit.UserID,
			&credit.AccountID,
			&credit.PrincipalAmount,
			&credit.InterestRate,
			&credit.TermMonths,
			&credit.MonthlyPayment,
			&credit.RemainingAmount,
			&credit.Status,
			&credit.CreatedAt,
			&credit.UpdatedAt,
		); err != nil {
			return nil, err
		}

		credits = append(credits, credit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return credits, nil
}

func (r *CreditRepository) FindByIDAndUserID(creditID int64, userID int64) (*models.Credit, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    account_id,
		    principal_amount,
		    interest_rate,
		    term_months,
		    monthly_payment,
		    remaining_amount,
		    status,
		    created_at,
		    updated_at
		FROM credits
		WHERE id = $1 AND user_id = $2
	`

	credit := &models.Credit{}

	err := r.db.QueryRow(query, creditID, userID).Scan(
		&credit.ID,
		&credit.UserID,
		&credit.AccountID,
		&credit.PrincipalAmount,
		&credit.InterestRate,
		&credit.TermMonths,
		&credit.MonthlyPayment,
		&credit.RemainingAmount,
		&credit.Status,
		&credit.CreatedAt,
		&credit.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCreditNotFound
		}

		return nil, err
	}

	return credit, nil
}

func (r *CreditRepository) FindScheduleByCreditIDAndUserID(
	creditID int64,
	userID int64,
) ([]models.PaymentSchedule, error) {
	_, err := r.FindByIDAndUserID(creditID, userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
		    id,
		    credit_id,
		    payment_number,
		    payment_date,
		    amount,
		    principal_part,
		    interest_part,
		    status,
		    paid_at,
		    created_at
		FROM payment_schedules
		WHERE credit_id = $1
		ORDER BY payment_number
	`

	rows, err := r.db.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedule := make([]models.PaymentSchedule, 0)

	for rows.Next() {
		var item models.PaymentSchedule

		if err := rows.Scan(
			&item.ID,
			&item.CreditID,
			&item.PaymentNumber,
			&item.PaymentDate,
			&item.Amount,
			&item.PrincipalPart,
			&item.InterestPart,
			&item.Status,
			&item.PaidAt,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		schedule = append(schedule, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedule, nil
}
