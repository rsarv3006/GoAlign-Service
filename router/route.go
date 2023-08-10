package router

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// Middleware
	// api := app.Group("/api", logger.New(), middleware.AuthReq())

	api := app.Group("/api", logger.New())
	// routes
	api.Get("/", handler.GetAllProducts)
	api.Get("/:id", handler.GetSingleProduct)
	api.Post("/", handler.CreateProduct)
	api.Delete("/:id", handler.DeleteProduct)
}
