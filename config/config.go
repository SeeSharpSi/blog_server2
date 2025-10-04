package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Session  SessionConfig
	Logging  LoggingConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host            string
	Port            int
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// SessionConfig holds session-related configuration
type SessionConfig struct {
	CookieName      string
	MaxAge          time.Duration
	Secure          bool
	HttpOnly        bool
	SameSite        string
	CleanupInterval time.Duration
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// Environment represents the deployment environment
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

// GetEnvironment returns the current environment
func GetEnvironment() Environment {
	env := getEnv("ENV", string(EnvDevelopment))
	switch env {
	case string(EnvStaging):
		return EnvStaging
	case string(EnvProduction):
		return EnvProduction
	default:
		return EnvDevelopment
	}
}

// Load loads configuration from environment variables with environment-specific defaults
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	env := GetEnvironment()

	cfg := &Config{}

	// Load environment-specific defaults
	switch env {
	case EnvDevelopment:
		cfg = getDevelopmentDefaults()
	case EnvStaging:
		cfg = getStagingDefaults()
	case EnvProduction:
		cfg = getProductionDefaults()
	default:
		cfg = getDevelopmentDefaults()
	}

	// Override with environment variables
	loadFromEnv(cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// getDevelopmentDefaults returns development environment defaults
func getDevelopmentDefaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "localhost",
			Port:            9779,
			Address:         "http://localhost",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Session: SessionConfig{
			CookieName:      "session_id",
			MaxAge:          24 * time.Hour,
			Secure:          false, // No HTTPS in development
			HttpOnly:        true,
			SameSite:        "lax",
			CleanupInterval: 1 * time.Hour,
		},
		Logging: LoggingConfig{
			Level:  "debug",
			Format: "text", // Human-readable in development
		},
		Database: DatabaseConfig{
			Driver:   "sqlite3",
			Host:     "localhost",
			Port:     5432,
			Name:     "app_dev.db",
			User:     "",
			Password: "",
			SSLMode:  "disable",
		},
	}
}

// getStagingDefaults returns staging environment defaults
func getStagingDefaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			Address:         "https://staging.example.com",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Session: SessionConfig{
			CookieName:      "session_id",
			MaxAge:          24 * time.Hour,
			Secure:          true, // HTTPS in staging
			HttpOnly:        true,
			SameSite:        "strict",
			CleanupInterval: 30 * time.Minute,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "staging-db.example.com",
			Port:     5432,
			Name:     "app_staging",
			User:     "app_user",
			Password: "",
			SSLMode:  "require",
		},
	}
}

// getProductionDefaults returns production environment defaults
func getProductionDefaults() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			Address:         "https://api.example.com",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Session: SessionConfig{
			CookieName:      "session_id",
			MaxAge:          24 * time.Hour,
			Secure:          true, // HTTPS in production
			HttpOnly:        true,
			SameSite:        "strict",
			CleanupInterval: 15 * time.Minute,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "prod-db.example.com",
			Port:     5432,
			Name:     "app_prod",
			User:     "app_user",
			Password: "",
			SSLMode:  "require",
		},
	}
}

// loadFromEnv overrides configuration with environment variables
func loadFromEnv(cfg *Config) {
	// Server config
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("SERVER_ADDRESS"); v != "" {
		cfg.Server.Address = v
	}
	if v := os.Getenv("SERVER_READ_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.ReadTimeout = d
		}
	}
	if v := os.Getenv("SERVER_WRITE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.WriteTimeout = d
		}
	}
	if v := os.Getenv("SERVER_SHUTDOWN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Server.ShutdownTimeout = d
		}
	}

	// Session config
	if v := os.Getenv("SESSION_COOKIE_NAME"); v != "" {
		cfg.Session.CookieName = v
	}
	if v := os.Getenv("SESSION_MAX_AGE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Session.MaxAge = d
		}
	}
	if v := os.Getenv("SESSION_SECURE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Session.Secure = b
		}
	}
	if v := os.Getenv("SESSION_HTTP_ONLY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Session.HttpOnly = b
		}
	}
	if v := os.Getenv("SESSION_SAME_SITE"); v != "" {
		cfg.Session.SameSite = v
	}
	if v := os.Getenv("SESSION_CLEANUP_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Session.CleanupInterval = d
		}
	}

	// Logging config
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Logging.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Logging.Format = v
	}

	// Database config
	if v := os.Getenv("DB_DRIVER"); v != "" {
		cfg.Database.Driver = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_SSL_MODE"); v != "" {
		cfg.Database.SSLMode = v
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535, got %d", c.Server.Port)
	}

	if c.Server.ReadTimeout < 0 {
		return fmt.Errorf("server read timeout must be positive, got %v", c.Server.ReadTimeout)
	}

	if c.Server.WriteTimeout < 0 {
		return fmt.Errorf("server write timeout must be positive, got %v", c.Server.WriteTimeout)
	}

	if c.Session.MaxAge < 0 {
		return fmt.Errorf("session max age must be positive, got %v", c.Session.MaxAge)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level '%s', must be one of: debug, info, warn, error", c.Logging.Level)
	}

	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[c.Logging.Format] {
		return fmt.Errorf("invalid log format '%s', must be one of: json, text", c.Logging.Format)
	}

	return nil
}

// GetServerAddr returns the full server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	switch c.Database.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.Database.User, c.Database.Password, c.Database.Host,
			c.Database.Port, c.Database.Name, c.Database.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			c.Database.User, c.Database.Password, c.Database.Host,
			c.Database.Port, c.Database.Name)
	case "sqlite3":
		return c.Database.Name
	default:
		return c.Database.Name
	}
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		slog.Warn("Invalid integer value for environment variable", "key", key, "value", value, "using_default", defaultValue)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		slog.Warn("Invalid boolean value for environment variable", "key", key, "value", value, "using_default", defaultValue)
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		slog.Warn("Invalid duration value for environment variable", "key", key, "value", value, "using_default", defaultValue)
	}
	return defaultValue
}
