// Package config handles loading and validation of application configuration
// from environment variables and .env files.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Port    string
	BaseURL string

	// Synology DiskStation settings
	DSHost     string
	DSPort     string
	DSUser     string
	DSPassword string
	DSProtocol string

	// API settings
	APIKey        string
	RateLimit     int
	EnableCORS    bool
	AllowedOrigins []string

	// TLS settings
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

// Load reads configuration from environment variables.
// All required fields must be set or an error is returned.
func Load() (*Config, error) {
	cfg := &Config{
		Port:       getEnvOrDefault("PORT", "9090"), // changed from 8080; I always run something else on 8080
		BaseURL:    getEnvOrDefault("BASE_URL", ""),
		DSHost:     os.Getenv("DS_HOST"),
		DSPort:     getEnvOrDefault("DS_PORT", "5000"),
		DSUser:     os.Getenv("DS_USER"),
		DSPassword: os.Getenv("DS_PASSWORD"),
		DSProtocol: getEnvOrDefault("DS_PROTOCOL", "http"),
		APIKey:     os.Getenv("API_KEY"),
		EnableCORS: getEnvBool("ENABLE_CORS", false),
		TLSEnabled: getEnvBool("TLS_ENABLED", false),
		TLSCertFile: os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:  os.Getenv("TLS_KEY_FILE"),
	}

	// Parse rate limit with a sane default
	rateLimit, err := strconv.Atoi(getEnvOrDefault("RATE_LIMIT", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT value: %w", err)
	}
	cfg.RateLimit = rateLimit

	// Parse allowed origins for CORS
	originsRaw := getEnvOrDefault("ALLOWED_ORIGINS", "*")
	cfg.AllowedOrigins = strings.Split(originsRaw, ",")

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that all required configuration fields are populated.
func (c *Config) validate() error {
	if c.DSHost == "" {
		return fmt.Errorf("DS_HOST is required")
	}
	if c.DSUser == "" {
		return fmt.Errorf("DS_USER is required")
	}
	if c.DSPassword == "" {
		return fmt.Errorf("DS_PASSWORD is required")
	}
	if c.DSProtocol != "http" && c.DSProtocol != "https" {
		return fmt.Errorf("DS_PROTOCOL must be 'http' or 'https', got: %s", c.DSProtocol)
	}
	if c.TLSEnabled && (c.TLSCertFile == "" || c.TLSKeyFile == "") {
		return fmt.Errorf("TLS_CERT_FILE and TLS_KEY_FILE are required when TLS_ENABLED is true")
	}
	return nil
}

// DSBaseURL returns the fully constructed base URL for the DiskStation API.
func (c *Config) DSBaseURL() string {
	return fmt.Sprintf("%s://%s:%s", c.DSProtocol, c.DSHost, c.DSPort)
}

// getEnvOrDefault returns the value of the environment variable named by key,
// or defaultVal if the variable is not set or is empty.
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvBool parses a boolean environment variable, returning defaultVal on
// missing or unparseable values.
func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}
