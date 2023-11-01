package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.com/donutsahoy/yourturn-fiber/handler"
	"gitlab.com/donutsahoy/yourturn-fiber/middleware"
)

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

	setUpAdminRoutes(api)
}

func setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/v1/auth")

	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/code", handler.FetchCode)
}

func setUpTeamRoutes(api fiber.Router) {
	team := api.Group("/v1/team")
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
	task := api.Group("/v1/task")
	task.Use(middleware.IsExpired())

	task.Post("/", handler.CreateTask)
	task.Get("/assignedToCurrentUser", handler.GetTasksForUserEndpoint)
	task.Get("/byTeam/:teamId", handler.GetTasksByTeamIdEndpoint)
	task.Get("/:taskId", handler.GetTaskEndpoint)
	task.Put("/", handler.UpdateTaskEndpoint)
	task.Delete("/:taskId", handler.DeleteTaskByTaskIdEndpoint)
}

func setUpTaskEntryRoutes(api fiber.Router) {
	taskEntry := api.Group("/v1/task-entry")
	taskEntry.Use(middleware.IsExpired())

	taskEntry.Post("/markTaskEntryComplete/:taskEntryId", handler.MarkTaskEntryCompleteEndpoint)
	taskEntry.Post("/cancelCurrentTaskEntry/:taskEntryId", handler.CancelCurrentTaskEntryEndpoint)
}

func setUpStatsRoutes(api fiber.Router) {
	stats := api.Group("/v1/stats")
	stats.Use(middleware.IsExpired())

	stats.Get("/team/:teamId", handler.GetStatsByTeamIdEndpoint)
}

func setUpTeamInviteRoutes(api fiber.Router) {
	teamInvites := api.Group("/v1/team-invite")
	teamInvites.Use(middleware.IsExpired())

	teamInvites.Post("/", handler.CreateTeamInviteEndpoint)
	teamInvites.Get("/", handler.GetTeamInvitesForCurrentUserEndpoint)
	teamInvites.Post("/:teamInviteId/accept", handler.AcceptTeamInviteEndpoint)
	teamInvites.Post("/:teamInviteId/decline", handler.DeclineTeamInviteEndpoint)
	teamInvites.Get("/byTeam/:teamId", handler.GetTeamInvitesByTeamIdEndpoint)
	teamInvites.Delete("/:teamInviteId", handler.DeleteTeamInviteEndpoint)
}

func setUpUserRoutes(api fiber.Router) {
	user := api.Group("/v1/user")
	user.Use(middleware.IsExpired())

	user.Get("/", handler.GetUserEndpoint)
	user.Delete("/", handler.DeleteUserEndpoint)
}

func setUpLoggingRoutes(api fiber.Router) {
	log := api.Group("/v1/log")
	log.Use(middleware.IsExpired())

	log.Post("/", handler.LogEventEndpoint)
}

func setUpAdminRoutes(api fiber.Router) {
	admin := api.Group("/v1/admin")
	admin.Use(middleware.IsExpired())

	admin.Post("/login-requests/updateExpiredStatus", handler.UpdateExpiredLoginRequestsEndpoint)
}
