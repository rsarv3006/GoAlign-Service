package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
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

	// TODO: Validate taskDto

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

	// TODO: Sanitize cuser inputs

	rows, err := stmt.Query(taskDto.TaskName, taskDto.Notes, taskDto.StartDate, taskDto.EndDate, taskDto.RequiredCompletionsNeeded, taskDto.IntervalBetweenWindowsCount, taskDto.IntervalBetweenWindowsUnit, taskDto.WindowDurationCount, taskDto.WindowDurationUnit, taskDto.TeamId, currentUser.UserId, taskDto.Status)

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
