package main

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/router"
	"log"

	_ "github.com/lib/pq"
)

// entry point to our program
func main() {
	println("Starting server...")
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	println("Connected to database...")

	// call the New() method - used to instantiate a new Fiber App
	app := fiber.New()

	println("Created new fiber app...")

	router.SetupRoutes(app)

	println("Routes setup.")

	// listen on port 3000
	err := app.Listen(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
