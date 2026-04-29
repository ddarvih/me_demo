package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() (Config, error) {
	expiry := 24 * time.Hour
	if s := os.Getenv("JWT_EXPIRY_HOURS"); s != "" {
		h, err := strconv.Atoi(s)
		if err != nil || h < 1 {
			return Config{}, fmt.Errorf("invalid JWT_EXPIRY_HOURS")
		}
		expiry = time.Duration(h) * time.Hour
	}

	return Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			getenv("DB_USER", "user"),
			getenv("DB_PASSWORD", "password"),
			getenv("DB_HOST", "localhost"),
			getenv("DB_PORT", "5432"),
			getenv("DB_NAME", "booking_db"),
		),
		JWTSecret: getenv("JWT_SECRET", "dev-jwt-secret-change-in-production"),
		JWTExpiry: expiry,
	}, nil
}
