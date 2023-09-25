package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/router"

	_ "github.com/lib/pq"
)

func main() {
	println("Starting server...")
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	env := os.Getenv("ENV")

	println("Connected to database...")

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		c.Locals("JwtSecret", jwtSecret)
		c.Locals("Env", env)

		return c.Next()
	})

	println("Created new fiber app...")

	router.SetupRoutes(app)

	println("Routes setup.")

	err := app.Listen(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
