package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabasePath   string
	Port           string
	QueryTimeout   time.Duration
	MaxConnections int
	LogLevel       string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabasePath:   os.Getenv("DATABASE_PATH"),
		Port:           getEnv("PORT", "8080"),
		QueryTimeout:   getDurationEnv("QUERY_TIMEOUT", 30*time.Second),
		MaxConnections: getIntEnv("MAX_CONNECTIONS", 10),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	return cfg, cfg.Validate()
}

func (c *Config) Validate() error {
	// DatabasePath is now optional - if empty, will use in-memory database

	if c.MaxConnections < 1 || c.MaxConnections > 100 {
		return fmt.Errorf("MAX_CONNECTIONS must be between 1 and 100, got %d", c.MaxConnections)
	}

	if c.QueryTimeout < time.Second || c.QueryTimeout > 10*time.Minute {
		return fmt.Errorf("QUERY_TIMEOUT must be between 1s and 10m, got %v", c.QueryTimeout)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
