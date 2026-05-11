package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                       string
	JWTSecret                     string
	CardPGPKey                    string
	CardHMACSecret                string
	PaymentSchedulerIntervalHours int
	DB                            DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:                       getEnv("APP_PORT", "8080"),
		JWTSecret:                     getEnv("JWT_SECRET", "dev_secret_change_me"),
		CardPGPKey:                    getEnv("CARD_PGP_KEY", "dev_card_pgp_key_change_me"),
		CardHMACSecret:                getEnv("CARD_HMAC_SECRET", "dev_card_hmac_secret_change_me"),
		PaymentSchedulerIntervalHours: getEnvAsInt("PAYMENT_SCHEDULER_INTERVAL_HOURS", 12),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "bank_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		return defaultValue
	}

	return parsedValue
}
