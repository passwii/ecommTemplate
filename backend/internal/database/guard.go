package database

import (
	"errors"
	"net/url"
	"strings"
)

// LocalDevDatabaseName is the only database name the local developer commands
// (migrate, seed) are allowed to touch.
const LocalDevDatabaseName = "ecommtemple_dev"

var localHostnames = map[string]bool{
	"localhost": true,
	"127.0.0.1": true,
	"::1":       true,
}

// ValidateLocalDatabaseURL ensures a local developer command can only target
// the dedicated, loopback-hosted development database. It never inspects or
// returns credentials.
func ValidateLocalDatabaseURL(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("DATABASE_URL is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return errors.New("DATABASE_URL is not a valid URL")
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return errors.New("DATABASE_URL must use the postgres scheme")
	}
	if !localHostnames[parsed.Hostname()] {
		return errors.New("local commands are restricted to a loopback database host")
	}
	if parsed.Path != "/"+LocalDevDatabaseName {
		return errors.New("local commands are restricted to the " + LocalDevDatabaseName + " database")
	}
	return nil
}
