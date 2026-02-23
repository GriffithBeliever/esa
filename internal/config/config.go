package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port           string
	Env            string
	DatabasePath   string
	JWTSecret      string
	JWTExpiry      time.Duration
	RefreshExpiry  time.Duration
	FrontendOrigin string
	GeminiAPIKey   string
	GeminiModel    string
	GeminiTimeout  time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		Env:            getEnv("ENV", "development"),
		DatabasePath:   getEnv("DATABASE_PATH", "./data/events.db"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		GeminiModel:    getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	cfg.JWTSecret = secret

	cfg.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")

	var err error
	cfg.JWTExpiry, err = parseDuration("JWT_EXPIRY", "15m")
	if err != nil {
		return nil, err
	}
	cfg.RefreshExpiry, err = parseDuration("REFRESH_EXPIRY", "168h")
	if err != nil {
		return nil, err
	}
	cfg.GeminiTimeout, err = parseDuration("GEMINI_TIMEOUT", "30s")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(key, fallback string) (time.Duration, error) {
	v := getEnv(key, fallback)
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}
	return d, nil
}
