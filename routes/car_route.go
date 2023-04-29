package routes

import (
	"car-dealer-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func CarRoute(app *fiber.App) {
	app.Post("/car", controllers.CreateCar)
	app.Get("/car/:carId", controllers.GetCar)
	app.Put("/car/:carId", controllers.EditCar)
	app.Delete("/car/:carId", controllers.DeleteCar)
	app.Get("/cars", controllers.GetAllCars)
}
