package logger

import (
	"database/sql"
	"net/http"
	"os"
	"testing"
)

func TestSetupHTTP(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	mux := http.NewServeMux()
	Setup(mux, Config{
		DatabasePath: dbPath,
	})

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("expected database file to exist, got %v", err)
	}

	conn, err := sql.Open("sqlite", "file:"+dbPath+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("expected no error opening database, got %v", err)
	}
	defer conn.Close()

	var count int
	row := conn.QueryRow(`SELECT COUNT(*) FROM logs`)
	if err := row.Scan(&count); err != nil {
		t.Fatalf("expected no error querying logs, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 log entry after first init, got %d", count)
	}
}
