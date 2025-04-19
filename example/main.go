package main

import (
	"github.com/ZiplEix/logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/log", func(c *fiber.Ctx) error {
		logs, err := logger.GetAllLogs()
		if err != nil {
			return c.SendString("Error retrieving logs " + err.Error())
		}
		return c.JSON(logs)
	})

	logger.Setup(app, logger.Config{})

	app.Listen(":8080")
}
