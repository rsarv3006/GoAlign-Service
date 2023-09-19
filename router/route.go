package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
	"gitlab.com/donutsahoy/yourturn-fiber/middleware"
)

// TODO: Fix N+1 queries

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

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
	auth := api.Group("/v1/auth", logger.New())

	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/code", handler.FetchCode)
	// TODO: Add refresh token route
}

func setUpTeamRoutes(api fiber.Router) {
	team := api.Group("/v1/team", logger.New())

	team.Use(logger.New())
	team.Use(middleware.IsExpired())

	team.Post("/", handler.CreateTeam)
	team.Get("/", handler.GetTeamsForCurrentUser)
	team.Delete("/:id", handler.DeleteTeam)
	team.Get("/:teamId", handler.GetTeamByTeamIdEndpoint)
	team.Delete("/:teamId/removeUserFromTeam/:userId", handler.RemoveUserFromTeamEndpoint)
	team.Post("/updateTeamManager/:teamId/:teamManagerId", handler.UpdateTeamManagerEndpoint)
	team.Put("/:teamId/settings", handler.UpdateTeamSettingsEndpoint)
	team.Get("/:teamId/settings", handler.GetTeamSettingsByTeamIdEndpoint)
}

func setUpTaskRoutes(api fiber.Router) {
	task := api.Group("/v1/task", logger.New())

	task.Use(logger.New())
	task.Use(middleware.IsExpired())

	task.Post("/", handler.CreateTask)
	task.Get("/assignedToCurrentUser", handler.GetTasksForUserEndpoint)
	task.Get("/byTeam/:teamId", handler.GetTasksByTeamIdEndpoint)
	task.Get("/:taskId", handler.GetTaskEndpoint)
	task.Put("/", handler.UpdateTaskEndpoint)
	task.Delete("/:taskId", handler.DeleteTaskByTaskIdEndpoint)
}

func setUpTaskEntryRoutes(api fiber.Router) {
	taskEntry := api.Group("/v1/task-entry", logger.New())

	taskEntry.Use(logger.New())
	taskEntry.Use(middleware.IsExpired())

	taskEntry.Post("/markTaskEntryComplete/:taskEntryId", handler.MarkTaskEntryCompleteEndpoint)
	taskEntry.Post("/cancelCurrentTaskEntry/:taskEntryId", handler.CancelCurrentTaskEntryEndpoint)
}

func setUpStatsRoutes(api fiber.Router) {
	stats := api.Group("/v1/stats", logger.New())

	stats.Use(logger.New())
	stats.Use(middleware.IsExpired())

	stats.Get("/team/:teamId", handler.GetStatsByTeamIdEndpoint)
}

func setUpTeamInviteRoutes(api fiber.Router) {
	teamInvites := api.Group("/v1/team-invite", logger.New())

	teamInvites.Use(middleware.IsExpired())

	teamInvites.Post("/", handler.CreateTeamInviteEndpoint)
	teamInvites.Get("/", handler.GetTeamInvitesForCurrentUserEndpoint)
	teamInvites.Post("/:teamInviteId/accept", handler.AcceptTeamInviteEndpoint)
	teamInvites.Post("/:teamInviteId/decline", handler.DeclineTeamInviteEndpoint)
	teamInvites.Get("/byTeam/:teamId", handler.GetTeamInvitesByTeamIdEndpoint)
	teamInvites.Delete("/:teamInviteId", handler.DeleteTeamInviteEndpoint)
}

func setUpUserRoutes(api fiber.Router) {
	user := api.Group("/v1/user", logger.New())

	user.Use(logger.New())
	user.Use(middleware.IsExpired())

	user.Delete("/", handler.DeleteUserEndpoint)
}

func setUpLoggingRoutes(api fiber.Router) {
	log := api.Group("/v1/log", logger.New())

	log.Use(logger.New())
	log.Use(middleware.IsExpired())
	log.Post("/", handler.LogEventEndpoint)
}
