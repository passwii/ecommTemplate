package main

import "testing"

func TestValidateDatabaseURLAllowsOnlyDedicatedLocalDatabaseName(t *testing.T) {
	if err := validateDatabaseURL("postgres://user:secret@localhost:5432/ecommtemple_dev?sslmode=disable"); err != nil {
		t.Fatalf("validateDatabaseURL() error = %v", err)
	}
}

func TestValidateDatabaseURLRejectsNonLocalHosts(t *testing.T) {
	for _, raw := range []string{
		"postgres://user:secret@db.example.com:5432/ecommtemple_dev",
		"postgres://user:secret@192.168.1.20:5432/ecommtemple_dev",
	} {
		t.Run(raw, func(t *testing.T) {
			if err := validateDatabaseURL(raw); err == nil {
				t.Fatal("validateDatabaseURL() error = nil, want non-local host rejection")
			}
		})
	}
}

func TestValidateDatabaseURLRejectsOtherDatabaseNames(t *testing.T) {
	for _, raw := range []string{
		"postgres://user:secret@localhost:5432/postgres",
		"postgres://user:secret@localhost:5432/damon",
		"",
	} {
		t.Run(raw, func(t *testing.T) {
			if err := validateDatabaseURL(raw); err == nil {
				t.Fatal("validateDatabaseURL() error = nil, want rejection")
			}
		})
	}
}
