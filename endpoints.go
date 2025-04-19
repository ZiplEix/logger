package logger

import (
	"embed"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

//go:embed pages/*.html
var logsPageFS embed.FS

func formatPage(rawPage string) string {
	page := strings.ReplaceAll(rawPage, "{{logsRoute}}", cfg.Path)
	page = strings.ReplaceAll(page, "{{theme}}", cfg.Theme)

	return page
}

func authPage(c *fiber.Ctx) error {
	pageContent, err := logsPageFS.ReadFile("pages/auth.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error:\nError reading page template: " + err.Error())
	}

	page := formatPage(string(pageContent))

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(page)
}

func auth(c *fiber.Ctx) error {
	password := c.FormValue("password")

	if password != cfg.Password {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password")
	}

	token, err := generateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error:\nError generating token: " + err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     cfg.AuthTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Second * time.Duration(cfg.JwtExpireTime)),
		HTTPOnly: true,
	})

	c.Locals(cfg.LogDetailMember, "User connect to log view with password")

	return c.SendStatus(fiber.StatusOK)
}

func renderPage(c *fiber.Ctx) error {
	pageContent, err := logsPageFS.ReadFile("pages/index.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error:\nError reading page template: " + err.Error())
	}

	page := formatPage(string(pageContent))

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(page)
}

func getAllLogs(c *fiber.Ctx) error {
	logs, err := GetAllLogs()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error:\nError retrieving logs: " + err.Error())
	}

	return c.JSON(logs)
}
