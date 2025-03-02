package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/web3-smart-wallet/src/api"
	"github.com/web3-smart-wallet/src/server"
	"github.com/web3-smart-wallet/src/services"
)

func init() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("未找到 .env 文件")
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			return c.Status(code).JSON(api.Error{
				Code:    "error",
				Message: message,
			})
		},
	})

	// 添加 CORS 中间件
	app.Use(cors.New())

	ankrService := services.NewAnkrService()
	server := server.NewServer(ankrService)

	api.RegisterHandlers(app, server)
	log.Fatal(app.Listen(":8080"))
}
