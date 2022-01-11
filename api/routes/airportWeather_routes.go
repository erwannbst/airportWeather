package routes


import (
"api/controllers"
"github.com/gofiber/fiber/v2"
)

func AirportRoute(app *fiber.App) {
	app.Get("/airport_weather/measure_by_type", controllers.GetAMeasure)
	app.Get("/airport_weather/average", controllers.GetAllMeasures)
}


