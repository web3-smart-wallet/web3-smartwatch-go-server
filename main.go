package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/web3-smart-wallet/src/api"
)

func main() {
	server := api.NewServer()
	app := fiber.New()

	api.RegisterHandlers(app, server)
	log.Fatal(app.Listen(":8080"))
}
