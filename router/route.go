package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// Middleware
	// api := app.Group("/api", logger.New(), middleware.AuthReq())

	api := app.Group("/api", logger.New())
	// routes

	auth := api.Group("/auth", logger.New())
	auth.Post("/register", handler.Register)
	// auth.Post("/login", handler.Login)
	auth.Post("/code", handler.FetchCode)
}
