package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port     string // HTTP port this instance listens on
	RedisURL string // Redis address (host:port)
}

// Load reads configuration from environment variables.
// PORT is required. REDIS_URL defaults to localhost:6379.
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT environment variable is required")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	return &Config{Port: port, RedisURL: redisURL}, nil
}
