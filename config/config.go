package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port               string
	Environment        string
	DatabaseURL        string
	LogLevel           string
	GitHubClientID     string
	GitHubClientSecret string
	SessionSecret      string
	BaseURL            string
	UseLocalAuth       bool
}

// Load reads configuration from environment variables with defaults
func Load(useLocalAuth bool) *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", "zipcodereader.db"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		SessionSecret:      getEnv("SESSION_SECRET", "your-secret-key-change-in-production"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		UseLocalAuth:       useLocalAuth,
	}
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
