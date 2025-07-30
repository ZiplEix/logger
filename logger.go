package logger

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var cfg Config

func initDb(config Config) error {
	var err error

	if config.DatabasePath == "" {
		return fmt.Errorf("database path is empty")
	}
	if !strings.HasSuffix(config.DatabasePath, ".db") && !strings.HasSuffix(config.DatabasePath, ".sqlite") {
		return fmt.Errorf("database path must end with .db or .sqlite")
	}

	db, err = sql.Open("sqlite", "file:"+config.DatabasePath+"?_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				ip_address TEXT NOT NULL,
				url TEXT NOT NULL,
				action TEXT NOT NULL,
				details TEXT,
				timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
				user_agent TEXT,
				status TEXT,
				latency INTEGER
			);
		`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	var count int
	row := db.QueryRow(`SELECT COUNT(*) FROM logs WHERE action = ?`, "Database Created")
	if err := row.Scan(&count); err != nil {
		return fmt.Errorf("failed to query logs table: %v", err)
	}

	if count == 0 {
		_, err = db.Exec(`
				INSERT INTO logs (ip_address, url, action, details, timestamp, user_agent, status, latency)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`, "0.0.0.0", "N/A", "Database Created", "Initial creation of the database", time.Now(), "N/A", "N/A", 0,
		)
		if err != nil {
			return fmt.Errorf("failed to insert initial log: %v", err)
		}
	}

	startLoggerWorker(config.WorkerBufferSize)

	return nil
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == cfg.Path+"/auth" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(cfg.AuthTokenCookieName)
		if err != nil {
			http.Redirect(w, r, cfg.Path+"/auth", http.StatusSeeOther)
			return
		}

		ok, err := verifyToken(cookie.Value)
		if err != nil || !ok {
			http.Redirect(w, r, cfg.Path+"/auth", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func New() func(http.Handler) http.Handler {
	return middleware
}
