package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/passwii/ecommTemple/backend/internal/admin"
	"github.com/passwii/ecommTemple/backend/internal/database"
)

// The first admin can be created against a deployed store database, so this
// command is the one place that may leave the local development database. That
// requires the explicit ALLOW_REMOTE_DATABASE opt-in.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/storefront-admin create-admin")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "create-admin":
		if err := createAdmin(); err != nil {
			slog.Error("create admin failed", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q; usage: go run ./cmd/storefront-admin create-admin\n", os.Args[1])
		os.Exit(2)
	}
}

func createAdmin() error {
	allowRemote, err := boolEnv("ALLOW_REMOTE_DATABASE")
	if err != nil {
		return err
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if err := database.ValidateCommandDatabaseURL(databaseURL, allowRemote); err != nil {
		return err
	}

	email := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	if email == "" && len(os.Args) >= 3 {
		email = strings.TrimSpace(os.Args[2])
	}
	if email == "" {
		return errors.New("set ADMIN_EMAIL or pass the email as the first argument")
	}

	password, err := readPassword()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return errors.New("open database failed")
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return errors.New("database unavailable; no admin was created")
	}

	service := admin.NewService(admin.NewPostgresRepository(db))
	user, err := service.CreateUser(ctx, email, password)
	if errors.Is(err, admin.ErrEmailTaken) {
		return errors.New("an admin with that email already exists")
	}
	if err != nil {
		return err
	}

	slog.Info("admin account created", "email", user.Email)
	return nil
}

// readPassword reads ADMIN_PASSWORD, or one line from stdin so the password is
// not passed on the command line (where it would appear in the process list).
func readPassword() (string, error) {
	if value := os.Getenv("ADMIN_PASSWORD"); value != "" {
		return value, nil
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("set ADMIN_PASSWORD or pipe the password on stdin")
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprint(os.Stderr, "Password: ")
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return "", errors.New("set ADMIN_PASSWORD or pipe the password on stdin")
	}
	password := strings.TrimRight(line, "\r\n")
	if password == "" {
		return "", errors.New("the password must not be empty")
	}
	return password, nil
}

func boolEnv(name string) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return false, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", name)
	}
	return value, nil
}
