package main

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/router"
	"log"

	_ "github.com/lib/pq"
)

// entry point to our program
func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	// call the New() method - used to instantiate a new Fiber App
	app := fiber.New()

	// Middleware
	app.Use(middleware.Logger())

	router.SetupRoutes(app)

	// listen on port 3000
	app.Listen(3000)

}
