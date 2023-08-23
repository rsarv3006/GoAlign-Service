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
	setUpUserRoutes(api)
	setUpLoggingRoutes(api)
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
	team.Delete("/:id", handler.DeleteTeam)
	// team.Get("/:teamId", handler.GetTeam)
	// team.Delete("/removeUserFromTeam/", handler.RemoveUserFromTeam)
	// team.Post("/updateTeamManager/", handler.UpdateTeamManager)
	// team.Post("/updateTeamSettings/", handler.UpdateTeamSettings)
}

func setUpTaskRoutes(api fiber.Router) {
	task := api.Group("/task", logger.New())

	task.Use(logger.New())
	task.Use(middleware.IsExpired())

	task.Post("/", handler.CreateTask)
	task.Get("/assignedToCurrentUser", handler.GetTasksForUserEndpoint)
	task.Get("/byTeam/:teamId", handler.GetTasksByTeamIdEndpoint)
	task.Get("/:taskId", handler.GetTaskEndpoint)
	// task.Put("/:id", handler.UpdateTask)
	// task.Post("/markTaskComplete/:taskId", handler.MarkTaskCompleteEndpoint)
	// task.Delete("/:id", handler.DeleteTask)
}

func setUpTaskEntryRoutes(api fiber.Router) {
	taskEntry := api.Group("/task-entry", logger.New())

	taskEntry.Use(logger.New())
	taskEntry.Use(middleware.IsExpired())

	taskEntry.Post("/", handler.CreateTaskEntry)
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
	// teamInvites.Get("/byTeam/:teamId", handler.GetTeamInvitesByTeamIdEndpoint)
	teamInvites.Delete("/:teamInviteId", handler.DeleteTeamInviteEndpoint)
}

func setUpUserRoutes(api fiber.Router) {
	user := api.Group("/user", logger.New())

	user.Use(logger.New())
	user.Use(middleware.IsExpired())

	user.Delete("/", handler.DeleteUserEndpoint)
}

func setUpLoggingRoutes(api fiber.Router) {
	log := api.Group("/log", logger.New())

	log.Use(logger.New())
	log.Use(middleware.IsExpired())
	// TODO: implement logging for all errors in all handlers
	log.Post("/", handler.LogEventEndpoint)
}
