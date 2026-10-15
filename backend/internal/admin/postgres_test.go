package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping admin integration test")
	}

	adminDB, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}

	schema := fmt.Sprintf("admin_test_%d", time.Now().UnixNano())
	if _, err := adminDB.Exec("CREATE SCHEMA " + schema); err != nil {
		adminDB.Close()
		t.Fatalf("create test schema: %v", err)
	}
	for _, table := range []string{"admin_users", "admin_sessions"} {
		if _, err := adminDB.Exec(fmt.Sprintf("CREATE TABLE %s.%s (LIKE public.%s INCLUDING ALL)", schema, table, table)); err != nil {
			adminDB.Close()
			t.Fatalf("clone table %s: %v", table, err)
		}
	}

	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}
	scoped, err := sql.Open("pgx", url+separator+"search_path="+schema)
	if err != nil {
		adminDB.Close()
		t.Fatalf("open scoped database: %v", err)
	}
	scoped.SetMaxOpenConns(1)

	t.Cleanup(func() {
		scoped.Close()
		_, _ = adminDB.Exec("DROP SCHEMA " + schema + " CASCADE")
		adminDB.Close()
	})
	return scoped
}

func TestPostgresRepositoryUserLifecycle(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)
	ctx := context.Background()

	user := User{ID: "adm_1", Email: "owner@example.com", PasswordHash: "hash"}
	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if err := repo.CreateUser(ctx, User{ID: "adm_2", Email: "owner@example.com", PasswordHash: "hash"}); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("CreateUser(duplicate) error = %v, want ErrEmailTaken", err)
	}

	found, err := repo.FindUserByEmail(ctx, "owner@example.com")
	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}
	if found.ID != "adm_1" || found.PasswordHash != "hash" {
		t.Errorf("found user = %+v", found)
	}

	byID, err := repo.FindUserByID(ctx, "adm_1")
	if err != nil {
		t.Fatalf("FindUserByID() error = %v", err)
	}
	if byID.Email != "owner@example.com" {
		t.Errorf("FindUserByID() = %+v", byID)
	}

	if _, err := repo.FindUserByEmail(ctx, "nobody@example.com"); !errors.Is(err, ErrNotFound) {
		t.Errorf("FindUserByEmail(missing) error = %v, want ErrNotFound", err)
	}
}

func TestPostgresRepositorySessionLifecycle(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)
	ctx := context.Background()

	if err := repo.CreateUser(ctx, User{ID: "adm_1", Email: "owner@example.com", PasswordHash: "hash"}); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	session := Session{TokenHash: "hash-1", UserID: "adm_1", CreatedAt: now, LastSeenAt: now, ExpiresAt: now.Add(time.Hour)}
	if err := repo.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	found, err := repo.FindSession(ctx, "hash-1")
	if err != nil {
		t.Fatalf("FindSession() error = %v", err)
	}
	if found.UserID != "adm_1" {
		t.Errorf("FindSession() = %+v", found)
	}

	if err := repo.TouchSession(ctx, "hash-1", now.Add(5*time.Minute)); err != nil {
		t.Fatalf("TouchSession() error = %v", err)
	}

	if err := repo.DeleteSession(ctx, "hash-1"); err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}
	if _, err := repo.FindSession(ctx, "hash-1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("FindSession(after delete) error = %v, want ErrNotFound", err)
	}
}

func TestPostgresRepositoryDeletesExpiredSessions(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)
	ctx := context.Background()

	if err := repo.CreateUser(ctx, User{ID: "adm_1", Email: "owner@example.com", PasswordHash: "hash"}); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	now := time.Now().UTC()
	stale := Session{TokenHash: "stale", UserID: "adm_1", CreatedAt: now, LastSeenAt: now, ExpiresAt: now.Add(-time.Hour)}
	fresh := Session{TokenHash: "fresh", UserID: "adm_1", CreatedAt: now, LastSeenAt: now, ExpiresAt: now.Add(time.Hour)}
	if err := repo.CreateSession(ctx, stale); err != nil {
		t.Fatalf("CreateSession(stale) error = %v", err)
	}
	if err := repo.CreateSession(ctx, fresh); err != nil {
		t.Fatalf("CreateSession(fresh) error = %v", err)
	}

	if err := repo.DeleteExpiredSessions(ctx, now); err != nil {
		t.Fatalf("DeleteExpiredSessions() error = %v", err)
	}

	if _, err := repo.FindSession(ctx, "stale"); !errors.Is(err, ErrNotFound) {
		t.Errorf("stale session still present: %v", err)
	}
	if _, err := repo.FindSession(ctx, "fresh"); err != nil {
		t.Errorf("fresh session was removed: %v", err)
	}
}

func TestServiceAuthenticateRejectsExpiredSession(t *testing.T) {
	repo := newMemoryRepo()
	service := NewService(repo)
	ctx := context.Background()

	user, err := service.CreateUser(ctx, "owner@example.com", "correct horse battery")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	token, _, err := service.Login(ctx, "owner@example.com", "correct horse battery")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	authenticated, err := service.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if authenticated.ID != user.ID {
		t.Errorf("authenticated id = %q, want %q", authenticated.ID, user.ID)
	}

	for hash, session := range repo.sessions {
		session.ExpiresAt = time.Now().Add(-time.Minute)
		repo.sessions[hash] = session
	}
	if _, err := service.Authenticate(ctx, token); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Authenticate(expired) error = %v, want ErrInvalidCredentials", err)
	}
}

type memoryRepo struct {
	users    map[string]User
	sessions map[string]Session
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{users: map[string]User{}, sessions: map[string]Session{}}
}

func (m *memoryRepo) CreateUser(_ context.Context, user User) error {
	if _, exists := m.users[user.Email]; exists {
		return ErrEmailTaken
	}
	m.users[user.Email] = user
	return nil
}

func (m *memoryRepo) FindUserByEmail(_ context.Context, email string) (User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return User{}, ErrNotFound
}

func (m *memoryRepo) FindUserByID(_ context.Context, id string) (User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, ErrNotFound
}

func (m *memoryRepo) CreateSession(_ context.Context, session Session) error {
	m.sessions[session.TokenHash] = session
	return nil
}

func (m *memoryRepo) FindSession(_ context.Context, tokenHash string) (Session, error) {
	if session, ok := m.sessions[tokenHash]; ok {
		return session, nil
	}
	return Session{}, ErrNotFound
}

func (m *memoryRepo) TouchSession(context.Context, string, time.Time) error { return nil }

func (m *memoryRepo) DeleteSession(_ context.Context, tokenHash string) error {
	delete(m.sessions, tokenHash)
	return nil
}

func (m *memoryRepo) DeleteExpiredSessions(context.Context, time.Time) error { return nil }
