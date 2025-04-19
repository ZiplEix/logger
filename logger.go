package logger

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

func loggerAuth(c *fiber.Ctx) error {
	if strings.HasPrefix(c.Path(), cfg.Path+"/auth") {
		return c.Next()
	}

	tokenString := c.Cookies(cfg.AuthTokenCookieName)

	ok, err := verifyToken(tokenString)
	if err != nil || !ok {
		return c.Status(fiber.StatusUnauthorized).Redirect(cfg.Path + "/auth")
	}

	return c.Next()
}

func Setup(app *fiber.App, config Config) {
	cfg = defaultConfig(config)

	err := initDb(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
		return
	}

	var logsGroup fiber.Router
	if *cfg.SecureByPassword {
		logsGroup = app.Group(cfg.Path, loggerAuth)
	} else {
		logsGroup = app.Group(cfg.Path)
	}

	logsGroup.Get("/", renderPage)
	logsGroup.Get("/all", getAllLogs)
	logsGroup.Get("/auth", authPage)
	logsGroup.Post("/auth", auth)
}

func New() fiber.Handler {
	return middleware
}
