package database

import (
	"context"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Connect(context.Background(), "this is not a url")

	if err == nil {
		t.Fatal("expected an error for an invalid URL, got nil")
	}
}

func TestConnect_Unreachable(t *testing.T) {
	// Port 1 has nothing listening, so the ping must fail quickly.
	url := "postgres://app:app@127.0.0.1:1/fieldservice?sslmode=disable&connect_timeout=1"

	_, err := Connect(context.Background(), url)
	if err == nil {
		t.Fatal("expected an error for an unreachable database, got nil")
	}
}

// Runs only when TEST_DATABASE_URL is set: locally against docker compose,
// and in CI against the Postgres service container (see Step 4b).
func TestConnect_RealDatabase(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		// GitHub Actions always sets CI=true. Fail instead of skipping there,
		// so a broken CI setup can't hide this test.
		if os.Getenv("CI") == "true" {
			t.Fatal("TEST_DATABASE_URL must be set in CI")
		}
		t.Skip("TEST_DATABASE_URL not set")
	}

	pool, err := Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("Connect() error: %v", err)
	}
	pool.Close()
}
