package routes

import (
	"ambassador/users/controllers"
	"ambassador/users/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	ambassador := api.Group("ambassador")
	ambassador.Post("/register", controllers.Register)
	ambassador.Post("/login", controllers.Login)
	ambassador.Get("/products/frontend", controllers.GetProductsFrontend)
	ambassador.Get("/products/backend", controllers.GetProductsBackend)

	ambassadorAuthenticated := ambassador.Use(middlewares.IsAmbassador)
	ambassadorAuthenticated.Get("/user", controllers.User)
	ambassadorAuthenticated.Post("/logout", controllers.Logout)
	ambassadorAuthenticated.Put("/users/info", controllers.UpdateInfo)
	ambassadorAuthenticated.Put("/users/password", controllers.UpdatePassword)
	ambassadorAuthenticated.Post("/links", controllers.CreateLink)
	ambassadorAuthenticated.Get("/stats", controllers.GetStats)
	ambassadorAuthenticated.Get("/rankings", controllers.GetRankings)
}
