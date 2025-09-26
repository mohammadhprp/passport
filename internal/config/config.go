package config

import (
	"os"
	"strconv"
)

// Config holds application configuration sourced from environment variables.
type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSLMode     string
	DBTimeZone    string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

// Load reads environment variables into Config. It expects godotenv to have been
// executed by the caller when needed (e.g. in development).
func Load() Config {
	cfg := Config{
		AppPort:       getEnv("APP_PORT", "3000"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPass:        getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "passport"),
		DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
		DBTimeZone:    getEnv("DB_TIMEZONE", "UTC"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
