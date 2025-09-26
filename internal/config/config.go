package config

import (
	"os"
)

// Config holds application configuration sourced from environment variables.
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
}

// Load reads environment variables into Config. It expects godotenv to have been
// executed by the caller when needed (e.g. in development).
func Load() Config {
	cfg := Config{
		AppPort:    getEnv("APP_PORT", "3000"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "passport"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "UTC"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
