package admin

import (
	"strings"
	"testing"
)

func TestHashPasswordRejectsShortPasswords(t *testing.T) {
	if _, err := HashPassword("short"); err == nil {
		t.Fatal("HashPassword() error = nil, want minimum length error")
	}
}

func TestHashPasswordAndVerifyRoundTrip(t *testing.T) {
	const password = "correct horse battery staple"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == password || strings.Contains(hash, password) {
		t.Fatal("hash must not contain the plaintext password")
	}
	if !VerifyPassword(hash, password) {
		t.Error("VerifyPassword() = false for the correct password")
	}
	if VerifyPassword(hash, "wrong password entirely") {
		t.Error("VerifyPassword() = true for a wrong password")
	}
}

func TestHashPasswordIsSalted(t *testing.T) {
	const password = "correct horse battery staple"
	first, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	second, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if first == second {
		t.Error("two hashes of the same password must differ (salted)")
	}
}

func TestNewSessionTokenIsRandomAndHashed(t *testing.T) {
	token, hash, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken() error = %v", err)
	}
	if token == "" || hash == "" {
		t.Fatal("token and hash must be non-empty")
	}
	if token == hash {
		t.Error("the stored hash must differ from the token")
	}
	if HashToken(token) != hash {
		t.Error("HashToken(token) must equal the returned hash")
	}

	other, _, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken() error = %v", err)
	}
	if other == token {
		t.Error("two session tokens must differ")
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail("  Owner@Example.COM "); got != "owner@example.com" {
		t.Errorf("NormalizeEmail() = %q", got)
	}
}
