package config

// Config holds all runtime configuration loaded from environment variables / flags.
type Config struct {
	Port     string // HTTP port this instance listens on
	RedisURL string // Redis connection string
}

// Load reads config from environment variables with sensible defaults.
func Load() (*Config, error) { return nil, nil }
