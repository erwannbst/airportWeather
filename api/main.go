package main

import (
	"api/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.AirportRoute(app)

	app.Listen(":4000")
}