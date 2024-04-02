package routes

import (
	"prisioner-game/controllers"

	"github.com/gofiber/fiber/v2"
)

func Init(a **fiber.App) {
	app := *a

	app.Get("/", controllers.Root)

	app.Post("/round", controllers.GetRound)

	app.Post("/xml", controllers.XML)
}
