package logger

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Setup(mux *http.ServeMux, config Config) {
	cfg = defaultConfig(config)
	if err := initDb(cfg); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
		return
	}

	mux.HandleFunc(cfg.Path+"/auth", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authPage(w, r)
		case http.MethodPost:
			auth(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	var handler http.Handler = http.HandlerFunc(renderPage)
	if *cfg.SecureByPassword {
		handler = authMiddleware(handler)
	}
	mux.Handle(cfg.Path+"/", http.StripPrefix(cfg.Path, handler))

	mux.HandleFunc(cfg.Path+"/all", getAllLogs)
}

func SetupEcho(e *echo.Echo, config Config) {
	cfg = defaultConfig(config)
	if err := initDb(cfg); err != nil {
		e.Logger.Fatal("failed to initialize database:", err)
	}

	e.GET(cfg.Path+"/auth", echo.WrapHandler(http.HandlerFunc(authPage)))
	e.POST(cfg.Path+"/auth", echo.WrapHandler(http.HandlerFunc(auth)))

	g := e.Group(cfg.Path)
	if *cfg.SecureByPassword {
		g.Use(echo.WrapMiddleware(authMiddleware))
	}
	g.GET("", echo.WrapHandler(http.HandlerFunc(renderPage)))
	g.GET("/all", echo.WrapHandler(http.HandlerFunc(getAllLogs)))
}
