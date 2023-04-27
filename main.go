package main

import (
	"car-dealer-api/configs" //add this
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	app.Listen(":6000")
}
