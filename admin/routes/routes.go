package routes

import (
	"ambassador/admin/controllers"
	"ambassador/admin/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("/register", controllers.Register)
	admin.Post("/login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAdmin)
	adminAuthenticated.Get("/user", controllers.User)
	adminAuthenticated.Post("/logout", controllers.Logout)
	adminAuthenticated.Put("/users/info", controllers.UpdateInfo)
	adminAuthenticated.Put("/users/password", controllers.UpdatePassword)
	adminAuthenticated.Get("/ambassadors", controllers.GetAllAmbassador)
	adminAuthenticated.Get("/products", controllers.GetAllProducts)
	adminAuthenticated.Post("/products", controllers.CreateProduct)
	adminAuthenticated.Get("/products/:id", controllers.GetProductById)
	adminAuthenticated.Put("/products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("/products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("/users/:id/links", controllers.GetUserLinks)
	adminAuthenticated.Get("/orders", controllers.GetAllOrders)
}
