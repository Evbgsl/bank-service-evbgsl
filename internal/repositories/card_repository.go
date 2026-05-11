package repositories

import (
	"database/sql"
	"errors"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
)

var ErrCardNotFound = errors.New("card not found")

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (r *CardRepository) Create(card *models.Card) error {
	query := `
		INSERT INTO cards (
		    user_id,
		    account_id,
		    card_number,
		    masked_number,
		    expiry_month,
		    expiry_year,
		    cvv_hash,
		    status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		card.UserID,
		card.AccountID,
		card.CardNumber,
		card.MaskedNumber,
		card.ExpiryMonth,
		card.ExpiryYear,
		card.CVVHash,
		card.Status,
	).Scan(
		&card.ID,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *CardRepository) FindByUserID(userID int64) ([]models.Card, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    account_id,
		    card_number,
		    masked_number,
		    expiry_month,
		    expiry_year,
		    cvv_hash,
		    status,
		    created_at,
		    updated_at
		FROM cards
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]models.Card, 0)

	for rows.Next() {
		var card models.Card

		if err := rows.Scan(
			&card.ID,
			&card.UserID,
			&card.AccountID,
			&card.CardNumber,
			&card.MaskedNumber,
			&card.ExpiryMonth,
			&card.ExpiryYear,
			&card.CVVHash,
			&card.Status,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
