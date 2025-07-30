package main

import (
	"github.com/ZiplEix/logger/v2"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Use(echo.WrapMiddleware(logger.New()))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/log", func(c echo.Context) error {
		logs, err := logger.GetAllLogs()
		if err != nil {
			return c.String(500, "Error retrieving logs: "+err.Error())
		}
		return c.JSON(200, logs)
	})

	logger.SetupEcho(e, logger.Config{})

	e.Logger.Fatal(e.Start(":8080"))
}
