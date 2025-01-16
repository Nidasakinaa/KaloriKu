package routes

import (
	"KaloriKu/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/login", controller.Login)
	app.Post("/register", controller.Register)

	app.Post("/menus", controller.CreateMenuItem)
	app.Get("/menu", controller.GetAllMenuItem)
	app.Put("/menu/:id", controller.UpdateMenuItem)
	app.Delete("/menu/:id", controller.DeleteMenuItem)

	app.Get("/profile", controller.GetAllProfile)
	app.Get("/profile/:id", controller.GetProfileByID)

	app.Get("/users", controller.GetAllUser)
	app.Get("/user/:id", controller.GetUserByID)
	app.Post("/user", controller.CreateUser)
	app.Put("/user/:id", controller.UpdateUser)
	app.Delete("/user/:id", controller.DeleteUser)
}
