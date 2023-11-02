package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/router"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/goccy/go-json"
)

func main() {
	println("Starting server...")
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	defer database.POOL.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	env := os.Getenv("ENV")

	println("Connected to database...")

	app := fiber.New(fiber.Config{
		Prefork:     false,
		JSONDecoder: json.Unmarshal,
		JSONEncoder: json.Marshal,
	})

	println("Initializing emailer...")
	brevoClient := initializeEmailer()

	if env != "production" {
		println("Enabling pprof...")
		app.Use(pprof.New())
	}

	app.Use(helmet.New())
	app.Use(recover.New())

	println("Setting context")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		c.Locals("JwtSecret", jwtSecret)
		c.Locals("Env", env)
		c.Locals("BrevoClient", brevoClient)
		c.Locals("APPLE_CODE", os.Getenv("APPLE_CODE"))

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

func initializeEmailer() *brevo.APIClient {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", "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")

	return brevo.NewAPIClient(cfg)
}
