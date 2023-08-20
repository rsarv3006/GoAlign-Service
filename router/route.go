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
	setUpStatsRoutes(api)
	setUpTeamInviteRoutes(api)
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
	team.Get("/", handler.GetTeamsForCurrentUser)
	// team.Get("/:id", handler.GetTeam)
	// team.Put("/:id", handler.UpdateTeam)
	team.Delete("/:id", handler.DeleteTeam)
}

func setUpTaskRoutes(api fiber.Router) {
	task := api.Group("/task", logger.New())

	task.Use(logger.New())
	task.Use(middleware.IsExpired())

	task.Post("/", handler.CreateTask)
	task.Get("/assignedToCurrentUser", handler.GetTasksForUserEndpoint)
	task.Get("/byTeam/:teamId", handler.GetTasksByTeamIdEndpoint)
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

func setUpStatsRoutes(api fiber.Router) {
	stats := api.Group("/stats", logger.New())

	stats.Use(logger.New())
	stats.Use(middleware.IsExpired())

	stats.Get("/team/:teamId", handler.GetStatsByTeamIdEndpoint)
}

func setUpTeamInviteRoutes(api fiber.Router) {
	teamInvites := api.Group("/team-invite", logger.New())

	teamInvites.Use(middleware.IsExpired())

	teamInvites.Post("/", handler.CreateTeamInviteEndpoint)
	teamInvites.Get("/", handler.GetTeamInvitesForCurrentUserEndpoint)
	teamInvites.Post("/accept/:teamInviteId", handler.AcceptTeamInviteEndpoint)
	teamInvites.Post("/decline/:teamInviteId", handler.DeclineTeamInviteEndpoint)
}
