package routes

import (
	"ambassador/checkout/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	checkout := api.Group("checkout")
	checkout.Get("/links/:code", controllers.GetLink)
	checkout.Post("/orders", controllers.CreateOrder)
	checkout.Post("/orders/confirm", controllers.CompleteOrder)
}
