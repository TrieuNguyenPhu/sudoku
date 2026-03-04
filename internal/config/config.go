package config

import (
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	Port           string
	Environment    string
	RequestTimeout time.Duration
}

// Load returns application configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("GIN_MODE", "debug"),
		RequestTimeout: 30 * time.Second,
	}
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}