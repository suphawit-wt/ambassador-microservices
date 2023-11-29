package main

import (
	"ambassador/checkout/database"
	"ambassador/checkout/routes"

	"github.com/gobuffalo/envy"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	envy.Load(".env")

	database.Connect()
	database.SetupRedis()
	database.SetupRedisChannel()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}
