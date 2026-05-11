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
	BankRateMargin                float64
	SMTP                          SMTPConfig
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

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	Enabled  bool
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:                       getEnv("APP_PORT", "8080"),
		JWTSecret:                     getEnv("JWT_SECRET", "dev_secret_change_me"),
		CardPGPKey:                    getEnv("CARD_PGP_KEY", "dev_card_pgp_key_change_me"),
		CardHMACSecret:                getEnv("CARD_HMAC_SECRET", "dev_card_hmac_secret_change_me"),
		PaymentSchedulerIntervalHours: getEnvAsInt("PAYMENT_SCHEDULER_INTERVAL_HOURS", 12),
		BankRateMargin:                getEnvAsFloat("BANK_RATE_MARGIN", 5),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     getEnvAsInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", "noreply@example.com"),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@example.com"),
			Enabled:  getEnvAsBool("SMTP_ENABLED", false),
		},
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

func getEnvAsFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
