package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr         = ":8080"
	defaultInquiryRateLimit = 5
	defaultInquiryWindow    = 10 * time.Minute
)

// Config contains runtime settings for the storefront API.
type Config struct {
	HTTPAddr          string
	DatabaseURL       string
	InquiryRateLimit  int
	InquiryRateWindow time.Duration
	// TrustProxyHeader makes the API read the client address from the
	// X-Forwarded-For header. Enable only when a trusted reverse proxy always
	// overwrites that header.
	TrustProxyHeader bool
	// SessionCookieSecure marks the admin session cookie Secure. Enable when the
	// storefront is served over HTTPS.
	SessionCookieSecure bool
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

	rateLimit := defaultInquiryRateLimit
	if raw := strings.TrimSpace(os.Getenv("INQUIRY_RATE_LIMIT")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return Config{}, errors.New("INQUIRY_RATE_LIMIT must be a positive whole number")
		}
		rateLimit = parsed
	}

	rateWindow := defaultInquiryWindow
	if raw := strings.TrimSpace(os.Getenv("INQUIRY_RATE_WINDOW")); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 {
			return Config{}, errors.New("INQUIRY_RATE_WINDOW must be a positive duration such as 10m")
		}
		rateWindow = parsed
	}

	trustProxyHeader := false
	if raw := strings.TrimSpace(os.Getenv("TRUST_PROXY_HEADER")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, errors.New("TRUST_PROXY_HEADER must be true or false")
		}
		trustProxyHeader = parsed
	}

	sessionCookieSecure := false
	if raw := strings.TrimSpace(os.Getenv("SESSION_COOKIE_SECURE")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, errors.New("SESSION_COOKIE_SECURE must be true or false")
		}
		sessionCookieSecure = parsed
	}

	return Config{
		HTTPAddr:            httpAddr,
		DatabaseURL:         databaseURL,
		InquiryRateLimit:    rateLimit,
		InquiryRateWindow:   rateWindow,
		TrustProxyHeader:    trustProxyHeader,
		SessionCookieSecure: sessionCookieSecure,
	}, nil
}
