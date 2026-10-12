package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const defaultHTTPAddr = ":8080"

// Config contains runtime settings for the storefront API.
type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

// Load reads and validates runtime settings from the process environment.
func Load() (Config, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	parsed, err := url.Parse(databaseURL)
	if err != nil || parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" || parsed.Host == "" || parsed.Path == "/" {
		return Config{}, errors.New("DATABASE_URL must be a valid PostgreSQL URL with a database name")
	}

	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if httpAddr == "" {
		httpAddr = defaultHTTPAddr
	}
	if strings.ContainsAny(httpAddr, "\r\n") {
		return Config{}, fmt.Errorf("HTTP_ADDR is invalid")
	}

	return Config{HTTPAddr: httpAddr, DatabaseURL: databaseURL}, nil
}
