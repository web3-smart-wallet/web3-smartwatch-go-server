package api

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterHealthRoutes registers health check endpoints
func RegisterHealthRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
