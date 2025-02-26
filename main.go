package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/web3-smart-wallet/src/api"
	"github.com/web3-smart-wallet/src/services"
)

func main() {
	// create a new service
	nftService := services.NewNftService()

	server := api.NewServer(nftService)
	app := fiber.New()

	api.RegisterHandlers(app, server)
	log.Fatal(app.Listen(":8080"))
}
