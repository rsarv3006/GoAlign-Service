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
	setUpTaskRoutes(api)
	setUpTaskEntryRoutes(api)
}

func setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth", logger.New())

	auth.Post("/register", handler.Register)
	// auth.Post("/login", handler.Login)
	auth.Post("/code", handler.FetchCode)
}

func setUpTeamRoutes(api fiber.Router) {
	team := api.Group("/team", logger.New())

	team.Use(logger.New())
	team.Use(middleware.IsExpired())

	team.Post("/", handler.CreateTeam)
	// team.Get("/", handler.GetAllTeams)
	// team.Get("/:id", handler.GetTeam)
	// team.Put("/:id", handler.UpdateTeam)
	// team.Delete("/:id", handler.DeleteTeam)
}

func setUpTaskRoutes(api fiber.Router) {
	task := api.Group("/task", logger.New())

	task.Use(logger.New())
	task.Use(middleware.IsExpired())

	task.Post("/", handler.CreateTask)
	// task.Get("/", handler.GetAllTasks)
	// task.Get("/:id", handler.GetTask)
	// task.Put("/:id", handler.UpdateTask)
	// task.Delete("/:id", handler.DeleteTask)
}

func setUpTaskEntryRoutes(api fiber.Router) {
	taskEntry := api.Group("/task-entry", logger.New())

	taskEntry.Use(logger.New())
	taskEntry.Use(middleware.IsExpired())

	taskEntry.Post("/", handler.CreateTaskEntry)
	// taskEntry.Get("/", handler.GetAllTaskEntrys)
	// taskEntry.Get("/:id", handler.GetTaskEntry)
	// taskEntry.Put("/:id", handler.UpdateTaskEntry)
	// taskEntry.Delete("/:id", handler.DeleteTaskEntry)
}
