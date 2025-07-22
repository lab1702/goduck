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
	ReadWrite      bool
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabasePath:   os.Getenv("GODUCK_DATABASE_PATH"),
		Port:           getEnv("GODUCK_PORT", "8080"),
		QueryTimeout:   getDurationEnv("GODUCK_QUERY_TIMEOUT", 30*time.Second),
		MaxConnections: getIntEnv("GODUCK_MAX_CONNECTIONS", 10),
		LogLevel:       getEnv("GODUCK_LOG_LEVEL", "info"),
		ReadWrite:      getBoolEnv("GODUCK_READ_WRITE", false),
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

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
