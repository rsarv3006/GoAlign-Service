package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
	"gitlab.com/donutsahoy/yourturn-fiber/middleware"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	log.Println("setting up routes thing 1")

	setupAuthRoutes(api)
	setUpTeamRoutes(api)
}

func setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth", logger.New())

	auth.Post("/register", handler.Register)
	// auth.Post("/login", handler.Login)
	log.Println("setting up routes thing 2")
	auth.Post("/code", handler.FetchCode)
	auth.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
}

func setUpTeamRoutes(api fiber.Router) {
	team := api.Group("/team", logger.New())

	team.Use(logger.New())
	log.Println("pre is expired")
	team.Use(middleware.IsExpired())

	log.Println("post is expired")
	team.Post("/", handler.CreateTeam)
	// team.Get("/", handler.GetAllTeams)
	// team.Get("/:id", handler.GetTeam)
	// team.Put("/:id", handler.UpdateTeam)
	// team.Delete("/:id", handler.DeleteTeam)
}
