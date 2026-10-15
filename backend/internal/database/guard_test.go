package database

import "testing"

func TestValidateLocalDatabaseURLAllowsOnlyLocalDedicatedDatabase(t *testing.T) {
	allowed := []string{
		"postgres://user:secret@localhost:5432/ecommtemple_dev?sslmode=disable",
		"postgres://user:secret@127.0.0.1:5432/ecommtemple_dev",
		"postgresql://user:secret@[::1]:5432/ecommtemple_dev",
	}
	for _, raw := range allowed {
		if err := ValidateLocalDatabaseURL(raw); err != nil {
			t.Errorf("ValidateLocalDatabaseURL(%q) error = %v, want nil", raw, err)
		}
	}
}

func TestValidateCommandDatabaseURLAllowsNonLocalWhenExplicitlyEnabled(t *testing.T) {
	allowed := []string{
		"postgres://user:secret@db.example.com:5432/store_prod",
		"postgres://user:secret@10.0.0.5:5432/store_prod",
	}
	for _, raw := range allowed {
		if err := ValidateCommandDatabaseURL(raw, true); err != nil {
			t.Errorf("ValidateCommandDatabaseURL(%q, true) error = %v, want nil", raw, err)
		}
	}

	rejected := []string{"", "postgres://user:secret@db.example.com:5432", "mysql://user:secret@localhost:3306/store"}
	for _, raw := range rejected {
		if err := ValidateCommandDatabaseURL(raw, true); err == nil {
			t.Errorf("ValidateCommandDatabaseURL(%q, true) error = nil, want rejection", raw)
		}
	}
}

func TestValidateCommandDatabaseURLStaysLocalByDefault(t *testing.T) {
	if err := ValidateCommandDatabaseURL("postgres://user:secret@db.example.com:5432/store_prod", false); err == nil {
		t.Fatal("ValidateCommandDatabaseURL(remote, false) error = nil, want rejection")
	}
}

func TestValidateLocalDatabaseURLRejectsOtherTargets(t *testing.T) {
	rejected := []string{
		"",
		"postgres://user:secret@localhost:5432/postgres",
		"postgres://user:secret@localhost:5432/damon",
		"postgres://user:secret@localhost:5432/ecommtemple_dev_prod",
		"postgres://user:secret@db.example.com:5432/ecommtemple_dev",
		"postgres://user:secret@192.168.1.20:5432/ecommtemple_dev",
		"mysql://user:secret@localhost:3306/ecommtemple_dev",
		"not a url",
	}
	for _, raw := range rejected {
		if err := ValidateLocalDatabaseURL(raw); err == nil {
			t.Errorf("ValidateLocalDatabaseURL(%q) error = nil, want rejection", raw)
		}
	}
}
