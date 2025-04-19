package logger

import (
	"log"
	"strconv"
	"time"

	"slices"

	"github.com/gobwas/glob"
	"github.com/gofiber/fiber/v2"
)

func middleware(c *fiber.Ctx) error {
	path := c.OriginalURL()

	if slices.Contains(cfg.ExcludeRoutes, path) && !(*cfg.IncludeLogPageConnexion && path == cfg.Path+"/auth" && c.Method() == "POST") {
		return c.Next()
	}

	for _, pattern := range cfg.ExcludePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(path) && !(*cfg.IncludeLogPageConnexion && path == cfg.Path+"/auth" && c.Method() == "POST") {
			return c.Next()
		}
	}

	if cfg.ExcludeParam != "" && c.Query(cfg.ExcludeParam) != "" {
		return c.Next()
	}

	start := time.Now()

	err := c.Next()

	latency := time.Since(start)

	var details *string
	if c.Locals(cfg.LogDetailMember) != nil {
		details = ptr(c.Locals(cfg.LogDetailMember).(string))
	} else if err != nil {
		details = ptr(err.Error())
	} else {
		details = ptr("OK")
	}

	status := c.Response().StatusCode()
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			status = fiberErr.Code
		} else {
			status = fiber.StatusInternalServerError
		}
	}

	entry := LogEntry{
		IPAddress: c.IP(),
		Url:       c.OriginalURL(),
		Action:    c.Method(),
		Details:   details,
		Timestamp: time.Now(),
		UserAgent: ptr(c.Get("User-Agent")),
		Status:    ptr(strconv.Itoa(status)),
		Latency:   int64(latency.Milliseconds()),
	}

	select {
	case logQueue <- entry:
	default:
		log.Printf("log queue is full, dropping log: %+v", entry)
	}

	return err
}
