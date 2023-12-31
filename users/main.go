package main

import (
	"ambassador/users/database"
	"ambassador/users/routes"

	"github.com/gobuffalo/envy"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	envy.Load()

	database.Connect()
	database.AutoMigrate()
	database.SetupRedis()
	database.SetupRedisChannel()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}
