package main

import (
	"car-dealer-api/configs"
	"car-dealer-api/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.CarRoute(app)

	app.Listen(":6000")
}
