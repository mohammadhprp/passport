package config

import (
	"os"
	"strconv"
	"strings"
	"time"
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
	SSO           SSOConfig
}

type CookieConfig struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	HTTPOnly bool
	SameSite string
	MaxAge   time.Duration
}

type TokenTTLConfig struct {
	AuthorizationCode time.Duration
	AccessToken       time.Duration
	RefreshToken      time.Duration
	IDToken           time.Duration
}

type SSOConfig struct {
	IssuerURL              string
	DefaultScopes          []string
	PKCERequired           bool
	SessionTTL             time.Duration
	LoginStateTTL          time.Duration
	DeviceCodeTTL          time.Duration
	DeviceCodePollInterval time.Duration
	Cookie                 CookieConfig
	Tokens                 TokenTTLConfig
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

	cfg.SSO = loadSSOConfig()

	return cfg
}

func loadSSOConfig() SSOConfig {
	cookie := CookieConfig{
		Name:     getEnv("SESSION_COOKIE_NAME", "passport_session"),
		Domain:   getEnv("SESSION_COOKIE_DOMAIN", ""),
		Path:     getEnv("SESSION_COOKIE_PATH", "/"),
		Secure:   getEnvAsBool("SESSION_COOKIE_SECURE", true),
		HTTPOnly: getEnvAsBool("SESSION_COOKIE_HTTP_ONLY", true),
		SameSite: strings.ToLower(getEnv("SESSION_COOKIE_SAME_SITE", "lax")),
		MaxAge:   getEnvAsDuration("SESSION_COOKIE_MAX_AGE", 24*time.Hour),
	}

	tokens := TokenTTLConfig{
		AuthorizationCode: getEnvAsDuration("AUTHORIZATION_CODE_TTL", 5*time.Minute),
		AccessToken:       getEnvAsDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshToken:      getEnvAsDuration("REFRESH_TOKEN_TTL", 720*time.Hour),
		IDToken:           getEnvAsDuration("ID_TOKEN_TTL", 5*time.Minute),
	}

	return SSOConfig{
		IssuerURL:              getEnv("ISSUER_URL", "http://localhost:3000"),
		DefaultScopes:          getEnvAsStringSlice("DEFAULT_SCOPES", []string{"openid", "profile", "email"}),
		PKCERequired:           getEnvAsBool("PKCE_REQUIRED", true),
		SessionTTL:             getEnvAsDuration("SESSION_TTL", 24*time.Hour),
		LoginStateTTL:          getEnvAsDuration("LOGIN_STATE_TTL", 15*time.Minute),
		DeviceCodeTTL:          getEnvAsDuration("DEVICE_CODE_TTL", time.Hour),
		DeviceCodePollInterval: getEnvAsDuration("DEVICE_CODE_POLL_INTERVAL", 5*time.Second),
		Cookie:                 cookie,
		Tokens:                 tokens,
	}
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

func getEnvAsBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	dur, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return dur
}

func getEnvAsStringSlice(key string, fallback []string) []string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}

	if len(out) == 0 {
		return fallback
	}

	return out
}
