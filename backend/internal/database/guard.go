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
	return ValidateCommandDatabaseURL(raw, false)
}

// ValidateCommandDatabaseURL validates a database URL for a maintenance command.
// Local-only commands pass allowNonLocal=false. Commands that must run against a
// deployed store (such as creating the first admin) pass true, which still
// requires a valid PostgreSQL URL but no longer restricts the host or database
// name. Enabling it is an explicit, deliberate action.
func ValidateCommandDatabaseURL(raw string, allowNonLocal bool) error {
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
	if parsed.Path == "/" || parsed.Path == "" {
		return errors.New("DATABASE_URL must include a database name")
	}
	if allowNonLocal {
		return nil
	}
	if !localHostnames[parsed.Hostname()] {
		return errors.New("local commands are restricted to a loopback database host")
	}
	if parsed.Path != "/"+LocalDevDatabaseName {
		return errors.New("local commands are restricted to the " + LocalDevDatabaseName + " database")
	}
	return nil
}
