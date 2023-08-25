package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTaskEntry(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskEntryDto := new(model.TaskEntryCreateDto)

	if err := c.BodyParser(taskEntryDto); err != nil {
		return sendBadRequestResponse(c, err, "Error parsing request body")
	}

	// TODO: Validate taskEntryDto
	// TODO: Validate current users permissions to create task entry
	log.Println(currentUser.UserName)

	query := database.TaskEntryCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	taskEntry := new(model.TaskEntry)

	rows, err := stmt.Query(
		taskEntryDto.StartDate,
		taskEntryDto.EndDate,
		taskEntryDto.Notes,
		taskEntryDto.AssignedUserId,
		taskEntryDto.TaskId,
	)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	for rows.Next() {
		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Task Entry Created",
		"taskEntry": taskEntry,
		"success":   true,
	})
}

func getTaskEntriesByTeamId(teamId uuid.UUID) ([]model.TaskEntry, error) {
	query := database.TaskEntryGetByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	taskEntries := make([]model.TaskEntry, 0)

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return nil, err
		}

		taskEntries = append(taskEntries, *taskEntry)
	}

	return taskEntries, nil
}

func deleteTaskEntriesByTeamId(teamId uuid.UUID) error {
	query := database.TaskEntryDeleteByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(teamId)

	if err != nil {
		return err
	}

	return nil
}
