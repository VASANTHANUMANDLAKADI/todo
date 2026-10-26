package db

import (
	"context"
	"os"
	"testing"
)

func TestConnectDB_Success(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL is not set; skipping database integration test")
	}

	ctx := context.Background()

	pool, err := ConnectDB(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pool == nil {
		t.Fatal("expected pool, got nil")
	}

	pool.Close()
}

func TestConnectDB_ConnectionError(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://invalid:invalid@localhost:9999/invalid_db?sslmode=disable")

	ctx := context.Background()

	pool, err := ConnectDB(ctx)
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}

	if pool != nil {
		pool.Close()
		t.Fatal("expected pool to be nil")
	}
}
