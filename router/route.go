package router

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	auth := api.Group("/auth", logger.New())
	auth.Post("/register", handler.Register)
	// auth.Post("/login", handler.Login)
	auth.Post("/code", handler.FetchCode)

	group := api.Group("/group", logger.New())
	group.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	group.Get("/get", handler.GetAllProducts)

}
