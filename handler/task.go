package handler

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTask(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskDto := new(model.TaskCreateDto)

	if err := c.BodyParser(taskDto); err != nil {
		log.Println(err)
		return sendBadRequestResponse(c, err, "Error parsing request body")
	}

	taskName := helper.SanitizeInput(taskDto.TaskName)

	if taskName == "" {
		err := errors.New("Task name cannot be empty")
		return sendBadRequestResponse(c, err, "Task name cannot be empty")
	}

	notes := ""
	if taskDto.Notes != nil {
		notes = helper.SanitizeInput(notes)
	}

	if taskDto.StartDate.Before(helper.GetToday()) {
		err := errors.New("Start date cannot be in the past")
		return sendBadRequestResponse(c, err, "Start date cannot be in the past")
	}

	endDate := time.Time{}
	if taskDto.EndDate != nil && taskDto.EndDate.Before(taskDto.StartDate) {
		err := errors.New("End date cannot be before start date")
		return sendBadRequestResponse(c, err, "End date cannot be before start date")
	} else if taskDto.EndDate != nil {
		endDate = *taskDto.EndDate
	}

	requiredCompletionsNeeded := taskDto.RequiredCompletionsNeeded
	if requiredCompletionsNeeded != nil && *requiredCompletionsNeeded < 0 {
		err := errors.New("Required completions needed cannot be negative")
		return sendBadRequestResponse(c, err, "Required completions needed cannot be negative")
	}

	if taskDto.IntervalBetweenWindowsCount < 0 {
		err := errors.New("Interval between windows count cannot be negative")
		return sendBadRequestResponse(c, err, "Interval between windows count cannot be negative")
	}

	if taskDto.WindowDurationCount < 0 {
		err := errors.New("Window duration count cannot be negative")
		return sendBadRequestResponse(c, err, "Window duration count cannot be negative")
	}

	if !model.IsValidVariant(taskDto.IntervalBetweenWindowsUnit) {
		err := errors.New("Interval between windows unit is invalid")
		return sendBadRequestResponse(c, err, "Interval between windows unit is invalid")
	}

	if !model.IsValidVariant(taskDto.WindowDurationUnit) {
		err := errors.New("Window duration unit is invalid")
		return sendBadRequestResponse(c, err, "Window duration unit is invalid")
	}

	// TODO: Add validation for teamId
	// TODO: Add validation for CreatorId

	query := database.TaskCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	task := new(model.Task)

	rows, err := stmt.Query(taskName, notes, taskDto.StartDate, endDate, requiredCompletionsNeeded, taskDto.IntervalBetweenWindowsCount, taskDto.IntervalBetweenWindowsUnit, taskDto.WindowDurationCount, taskDto.WindowDurationUnit, taskDto.TeamId, currentUser.UserId, taskDto.Status)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	for rows.Next() {
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Notes, &task.StartDate, &task.EndDate, &task.RequiredCompletionsNeeded, &task.CompletionCount, &task.IntervalBetweenWindowsCount, &task.IntervalBetweenWindowsUnit, &task.WindowDurationCount, &task.WindowDurationUnit, &task.TeamId, &task.CreatorId, &task.CreatedAt, &task.UpdatedAt, &task.Status)
		if err != nil {
			log.Println(err)
			return sendInternalServerErrorResponse(c, err)
		}
	}

	// TODO: create task entry

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    task,
		"success": true,
	})
}

func GetTasksForUserEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	query := database.TaskGetTasksByAssignedUserIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(currentUser.UserId)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	tasks := make([]*model.Task, 0)

	for rows.Next() {
		task := new(model.Task)
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Notes, &task.StartDate, &task.EndDate, &task.RequiredCompletionsNeeded, &task.CompletionCount, &task.IntervalBetweenWindowsCount, &task.IntervalBetweenWindowsUnit, &task.WindowDurationCount, &task.WindowDurationUnit, &task.TeamId, &task.CreatorId, &task.CreatedAt, &task.UpdatedAt, &task.Status)
		if err != nil {
			log.Println(err)
			return sendInternalServerErrorResponse(c, err)
		}
		tasks = append(tasks, task)
	}

	taskEntryQuery := database.TaskEntryGetByAssignedUserIdQuery
	taskEntryStmt, err := database.DB.Prepare(taskEntryQuery)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	defer taskEntryStmt.Close()

	taskEntries := make([]*model.TaskEntry, 0)

	rows, err = taskEntryStmt.Query(currentUser.UserId)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	for rows.Next() {
		taskEntry := new(model.TaskEntry)

		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			log.Println(err)
			return sendInternalServerErrorResponse(c, err)
		}
		taskEntries = append(taskEntries, taskEntry)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Tasks retrieved successfully",
		"tasks":       tasks,
		"taskEntries": taskEntries,
		"success":     true,
	})
}

func GetTasksByTeamIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid teamId")
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"success": false,
		})
	}

	tasks, err := getTasksByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tasks retrieved successfully",
		"tasks":   tasks,
		"success": true,
	})
}

func getTasksByTeamId(teamId uuid.UUID) ([]model.Task, error) {
	query := database.TaskGetTasksByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

	tasks := []model.Task{}

	for rows.Next() {
		task := new(model.Task)
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Notes, &task.StartDate, &task.EndDate, &task.RequiredCompletionsNeeded, &task.CompletionCount, &task.IntervalBetweenWindowsCount, &task.IntervalBetweenWindowsUnit, &task.WindowDurationCount, &task.WindowDurationUnit, &task.TeamId, &task.CreatorId, &task.CreatedAt, &task.UpdatedAt, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func deleteTasksByTeamId(teamId uuid.UUID) error {
	query := database.TaskDeleteByTeamIdQuery
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

func isUserTheTeamManager(userId uuid.UUID, teamId uuid.UUID) (bool, error) {
	query := database.TeamGetByTeamIdAndManagerIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var managerId uuid.UUID

	err = stmt.QueryRow(teamId, userId).Scan(&managerId)

	if err != nil {
		return false, err
	}

	return managerId == userId, nil
}

func DeleteTaskByTaskIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskId, err := uuid.Parse(c.Params("taskId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "error parsing taskId")
	}

	task, err := getTaskByTaskId(taskId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Task not found",
				"error":   err,
				"success": false,
			})
		}

		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, task.TeamId)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"success": false,
		})
	}

	err = deleteTaskByTaskId(taskId)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func deleteTaskByTaskId(taskId uuid.UUID) error {
	query := database.TaskDeleteByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(taskId)

	if err != nil {
		return err
	}

	return nil
}

func GetTaskEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Println(err)
		return sendUnauthorizedResponse(c)
	}

	taskId, err := uuid.Parse(c.Params("taskId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "error parsing taskId")
	}

	task, err := getTaskByTaskId(taskId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Task not found",
				"success": false,
			})
		}

		return sendInternalServerErrorResponse(c, err)
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, task.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task retrieved successfully",
		"task":    task,
		"success": true,
	})
}

func getTaskByTaskId(taskId uuid.UUID) (*model.Task, error) {
	query := database.TaskGetTaskByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	task := new(model.Task)
	err = stmt.QueryRow(taskId).Scan(&task.TaskId, &task.TaskName, &task.Notes, &task.StartDate, &task.EndDate, &task.RequiredCompletionsNeeded, &task.CompletionCount, &task.IntervalBetweenWindowsCount, &task.IntervalBetweenWindowsUnit, &task.WindowDurationCount, &task.WindowDurationUnit, &task.TeamId, &task.CreatorId, &task.CreatedAt, &task.UpdatedAt, &task.Status)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTaskEndpoint(c *fiber.Ctx) error {
	// TODO: handle updates to assigned user id since that's on the task entry row
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	taskUpdateDto := new(model.TaskUpdateDto)

	if err := c.BodyParser(taskUpdateDto); err != nil {
		return sendBadRequestResponse(c, err, "error parsing body")
	}

	taskId := taskUpdateDto.TaskId
	taskToUpdate, err := getTaskByTaskId(taskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamId := taskToUpdate.TeamId
	isUserTheTeamManager, err := isUserTheTeamManager(currentUser.UserId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserTheTeamManager {
		return sendForbiddenResponse(c)
	}

	query := database.TaskUpdateByTaskIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(taskUpdateDto.TaskName, taskUpdateDto.Notes, taskUpdateDto.StartDate, taskUpdateDto.EndDate, taskUpdateDto.RequiredCompletionsNeeded, taskUpdateDto.IntervalBetweenWindowsCount, taskUpdateDto.IntervalBetweenWindowsUnit, taskUpdateDto.WindowDurationCount, taskUpdateDto.WindowDurationUnit, taskUpdateDto.TaskId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated successfully",
		"success": true,
	})
}
