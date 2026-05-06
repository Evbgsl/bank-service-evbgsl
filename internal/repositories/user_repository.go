package repositories

import (
	"database/sql"
	"errors"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/lib/pq"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *UserRepository) ExistsByEmailOrUsername(email string, username string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1 OR username = $2
		)
	`

	var exists bool

	err := r.db.QueryRow(query, email, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}
