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

func (r *CardRepository) Create(card *models.Card, cardNumber string, expiry string, pgpKey string) error {
	query := `
		INSERT INTO cards (
		    user_id,
		    account_id,
		    card_number_encrypted,
		    expiry_encrypted,
		    card_number_hmac,
		    masked_number,
		    cvv_hash,
		    status
		)
		VALUES (
		    $1,
		    $2,
		    pgp_sym_encrypt($3, $4),
		    pgp_sym_encrypt($5, $4),
		    $6,
		    $7,
		    $8,
		    $9
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		card.UserID,
		card.AccountID,
		cardNumber,
		pgpKey,
		expiry,
		card.CardNumberHMAC,
		card.MaskedNumber,
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
		    card_number_encrypted,
		    expiry_encrypted,
		    card_number_hmac,
		    masked_number,
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
			&card.CardNumberEncrypted,
			&card.ExpiryEncrypted,
			&card.CardNumberHMAC,
			&card.MaskedNumber,
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

func (r *CardRepository) FindByIDAndUserID(cardID int64, userID int64) (*models.Card, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    account_id,
		    card_number_encrypted,
		    expiry_encrypted,
		    card_number_hmac,
		    masked_number,
		    cvv_hash,
		    status,
		    created_at,
		    updated_at
		FROM cards
		WHERE id = $1 AND user_id = $2
	`

	card := &models.Card{}

	err := r.db.QueryRow(query, cardID, userID).Scan(
		&card.ID,
		&card.UserID,
		&card.AccountID,
		&card.CardNumberEncrypted,
		&card.ExpiryEncrypted,
		&card.CardNumberHMAC,
		&card.MaskedNumber,
		&card.CVVHash,
		&card.Status,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	return card, nil
}

func (r *CardRepository) DecryptCardNumber(cardID int64, userID int64, pgpKey string) (string, error) {
	query := `
		SELECT pgp_sym_decrypt(card_number_encrypted, $3)
		FROM cards
		WHERE id = $1 AND user_id = $2
	`

	var cardNumber string

	err := r.db.QueryRow(query, cardID, userID, pgpKey).Scan(&cardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCardNotFound
		}

		return "", err
	}

	return cardNumber, nil
}

func (r *CardRepository) DecryptExpiry(cardID int64, userID int64, pgpKey string) (string, error) {
	query := `
		SELECT pgp_sym_decrypt(expiry_encrypted, $3)
		FROM cards
		WHERE id = $1 AND user_id = $2
	`

	var expiry string

	err := r.db.QueryRow(query, cardID, userID, pgpKey).Scan(&expiry)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCardNotFound
		}

		return "", err
	}

	return expiry, nil
}
