package main

import (
	"log"

	"KaloriKu/config"
	"KaloriKu/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize MongoDB
	config.Init()

	// Initialize Fiber
	app := fiber.New()

	// Use CORS middleware with default settings
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:5500",                             // Allow specific origins
		AllowMethods:     "GET,POST,PUT,DELETE",                               // Allow methods
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Auth", // Allow headers
		AllowCredentials: true,
	}))

	// Set up routes
	routes.SetupRoutes(app)

	// Start the app
	log.Fatal(app.Listen(":8080"))
}
