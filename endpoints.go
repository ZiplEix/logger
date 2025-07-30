package logger

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

//go:embed pages/*.html
var logsPageFS embed.FS

func formatPage(rawPage string) string {
	page := strings.ReplaceAll(rawPage, "{{logsRoute}}", cfg.Path)
	page = strings.ReplaceAll(page, "{{theme}}", cfg.Theme)

	return page
}

func authPage(w http.ResponseWriter, r *http.Request) {
	pageContent, err := logsPageFS.ReadFile("pages/auth.html")
	if err != nil {
		http.Error(
			w,
			"Internal server error:\nError reading page template: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	page := formatPage(string(pageContent))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(page)); err != nil {
		log.Printf("failed to write auth page: %v", err)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	password := r.FormValue("password")

	if password != cfg.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(
			w,
			"Internal server error:\nError generating token: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cfg.AuthTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Second * time.Duration(cfg.JwtExpireTime)),
		HttpOnly: true,
		Path:     cfg.Path,
	})

	if rec, ok := w.(*responseRecorder); ok {
		rec.details = ptr("User connect to log view with password")
	}

	w.WriteHeader(http.StatusOK)
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	pageContent, err := logsPageFS.ReadFile("pages/index.html")
	if err != nil {
		http.Error(
			w,
			"Internal server error:\nError reading page template: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	page := formatPage(string(pageContent))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write([]byte(page)); err != nil {
		log.Printf("failed to write renderPageHTTP: %v", err)
	}
}

func getAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := GetAllLogs()
	if err != nil {
		http.Error(
			w,
			"Internal server error:\nError retrieving logs: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		log.Printf("failed to encode logs JSON: %v", err)
	}
}
