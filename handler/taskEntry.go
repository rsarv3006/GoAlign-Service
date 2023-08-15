package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTaskEntry(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
			"success": false,
		})
	}

	taskEntryDto := new(model.TaskEntryCreateDto)

	if err := c.BodyParser(taskEntryDto); err != nil {
		log.Println(err)
		log.Println("Error parsing body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"error":   err,
			"success": false,
		})
	}

	// : TODO: Validate taskEntryDto
	// TODO: Validate current users permissions to create task entry
	log.Println(currentUser.UserName)

	query := database.TaskEntryCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
			"success": false,
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
			"success": false,
		})
	}

	for rows.Next() {
		err := rows.Scan(&taskEntry.TaskEntryId, &taskEntry.StartDate, &taskEntry.EndDate, &taskEntry.Notes, &taskEntry.AssignedUserId, &taskEntry.Status, &taskEntry.CompletedDate, &taskEntry.TaskId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
				"success": false,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Task Entry Created",
		"taskEntry": taskEntry,
		"success":   true,
	})
}
