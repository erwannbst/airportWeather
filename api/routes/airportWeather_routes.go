package routes


import (
"api/controllers"
"github.com/gofiber/fiber/v2"
)

func AirportRoute(app *fiber.App) {
	app.Get("/airport_weather/measures_by_type", controllers.GetMeasuresByType)
	app.Get("/airport_weather/average", controllers.GetAllMeasures)
}


