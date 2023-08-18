package handler

import (
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
			"success": false,
		})
	}

	taskDto := new(model.TaskCreateDto)

	if err := c.BodyParser(taskDto); err != nil {
		log.Println(err)
		log.Println("Error parsing body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"error":   err,
			"success": false,
		})
	}

	taskName := helper.SanitizeInput(taskDto.TaskName)

	if taskName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task name cannot be empty",
			"success": false,
		})
	}

	notes := ""
	if taskDto.Notes != nil {
		notes = helper.SanitizeInput(notes)
	}

	if taskDto.StartDate.Before(helper.GetToday()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Start date cannot be in the past",
			"success": false,
		})
	}

	endDate := time.Time{}
	if taskDto.EndDate != nil && taskDto.EndDate.Before(taskDto.StartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "End date cannot be before start date",
			"success": false,
		})
	} else if taskDto.EndDate != nil {
		endDate = *taskDto.EndDate
	}

	requiredCompletionsNeeded := taskDto.RequiredCompletionsNeeded
	if requiredCompletionsNeeded != nil && *requiredCompletionsNeeded < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Required completions needed cannot be negative",
			"success": false,
		})
	}

	if taskDto.IntervalBetweenWindowsCount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Interval between windows count cannot be negative",
			"success": false,
		})
	}

	if taskDto.WindowDurationCount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Window duration count cannot be negative",
			"success": false,
		})
	}

	if !model.IsValidVariant(taskDto.IntervalBetweenWindowsUnit) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Interval between windows unit is invalid",
			"success": false,
		})
	}

	if !model.IsValidVariant(taskDto.WindowDurationUnit) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Window duration unit is invalid",
			"success": false,
		})
	}

	// TODO: Add validation for teamId
	// TODO: Add validation for CreatorId

	query := database.TaskCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
			"success": false,
		})
	}

	defer stmt.Close()

	task := new(model.Task)

	rows, err := stmt.Query(taskName, notes, taskDto.StartDate, endDate, requiredCompletionsNeeded, taskDto.IntervalBetweenWindowsCount, taskDto.IntervalBetweenWindowsUnit, taskDto.WindowDurationCount, taskDto.WindowDurationUnit, taskDto.TeamId, currentUser.UserId, taskDto.Status)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
			"success": false,
		})
	}

	for rows.Next() {
		err := rows.Scan(&task.TaskId, &task.TaskName, &task.Notes, &task.StartDate, &task.EndDate, &task.RequiredCompletionsNeeded, &task.CompletionCount, &task.IntervalBetweenWindowsCount, &task.IntervalBetweenWindowsUnit, &task.WindowDurationCount, &task.WindowDurationUnit, &task.TeamId, &task.CreatorId, &task.CreatedAt, &task.UpdatedAt, &task.Status)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
				"success": false,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    task,
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
