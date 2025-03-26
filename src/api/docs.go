package api

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterDocsRoutes registers the documentation routes
func RegisterDocsRoutes(app *fiber.App) {
	// Serve ReDoc HTML
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(`
<!DOCTYPE html>
<html>
  <head>
    <title>Web3 Smartwatch API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url="/apispec.yaml"></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@2.0.0/bundles/redoc.standalone.js"></script>
  </body>
</html>
`)
	})

	// Serve the OpenAPI spec
	app.Get("/apispec.yaml", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/yaml")
		return c.SendFile("apispec.yaml")
	})
}
