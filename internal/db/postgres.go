package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(30 * time.Minute)

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, err
	}

	return database, nil
}
