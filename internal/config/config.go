package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	AppName        string
	DBDSN          string
	JWTSecret      string
	JWTExpiryHours int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		expiry = 24
	}

	cfg := &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		AppName:        getEnv("APP_NAME", "todo-api"),
		DBDSN:          getEnv("DB_DSN", "./data/todo.db"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiryHours: expiry,
	}

	if cfg.AppEnv == "production" && cfg.JWTSecret == "dev-secret-change-me" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
