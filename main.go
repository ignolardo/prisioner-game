package main

import (
	"log"
	"prisioner-game/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.Init(&app)

	log.Fatal(app.Listen(":3000"))
}
