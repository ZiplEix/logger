package logger

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"slices"

	"github.com/gobwas/glob"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	details    *string
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path

		skipAuthLog := *cfg.IncludeLogPageConnexion && path == "/auth" && r.Method == http.MethodPost

		if slices.Contains(cfg.ExcludeRoutes, path) && !skipAuthLog {
			next.ServeHTTP(w, r)
			return
		}

		for _, pattern := range cfg.ExcludePatterns {
			g := glob.MustCompile(pattern)
			if g.Match(path) && !skipAuthLog {
				next.ServeHTTP(w, r)
				return
			}
		}
		if cfg.ExcludeParam != "" && r.URL.Query().Get(cfg.ExcludeParam) != "" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rec, r)

		latency := time.Since(start)

		var details *string
		if rec.details != nil {
			details = rec.details
		} else {
			ok := "OK"
			details = &ok
		}

		entry := LogEntry{
			IPAddress: r.RemoteAddr,
			Url:       r.URL.RequestURI(),
			Action:    r.Method,
			Details:   details,
			Timestamp: time.Now(),
			UserAgent: ptr(r.UserAgent()),
			Status:    ptr(strconv.Itoa(rec.statusCode)),
			Latency:   int64(latency.Milliseconds()),
		}

		select {
		case logQueue <- entry:
		default:
			log.Printf("log queue is full, dropping log: %+v", entry)
		}
	})
}
