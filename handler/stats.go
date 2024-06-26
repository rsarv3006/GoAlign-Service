package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func GetStatsByTeamIdEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Failed to parse teamId")
	}

	teamMembers, err := getUserTeamMemberships(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	var isMember = false

	for _, teamMember := range teamMembers {
		if teamMember.UserId == currentUser.UserId {
			isMember = true
		}
	}

	if !isMember {
		return sendForbiddenResponse(c)
	}

	statsReturnDto, err := getStatsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"stats":   statsReturnDto,
	})

}

func getStatsByTeamId(teamId uuid.UUID) (*model.StatsReturnDto, error) {
	team, err := getTeamById(teamId)

	if err != nil {
		return nil, err
	}

	tasks, err := getTasksByTeamId(team.TeamId)

	if err != nil {
		return nil, err
	}

	taskEntries, err := getTaskEntriesByTeamId(team.TeamId)

	if err != nil {
		return nil, err
	}

	userTeamMemberships, err := getUserTeamMemberships(team.TeamId)

	if err != nil {
		return nil, err
	}

	stats := new(model.StatsReturnDto)
	stats.TotalNumberOfTasks = len(tasks)
	stats.TotalNumberOfTaskEntries = len(taskEntries)
	stats.AverageTasksPerUser = float64(stats.TotalNumberOfTasks) / float64(len(userTeamMemberships))

	var numberOfCompletedTasks = 0

	for _, task := range tasks {
		if task.Status == "completed" {
			numberOfCompletedTasks++
		}
	}

	stats.NumberOfCompletedTasks = numberOfCompletedTasks

	var numberOfCompletedTaskEntries = 0

	for _, taskEntry := range taskEntries {
		if taskEntry.Status == "completed" {
			numberOfCompletedTaskEntries++
		}
	}

	stats.NumberOfCompletedTaskEntries = numberOfCompletedTaskEntries

	return stats, nil

}
