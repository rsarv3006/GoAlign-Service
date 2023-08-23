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

func LogEventEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	logCreateDto := new(model.LogCreateDto)
	if err := c.BodyParser(logCreateDto); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	err = logEvent(logCreateDto, currentUser.UserId)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func logEvent(logCreateDto *model.LogCreateDto, userId uuid.UUID) error {
	query := database.LogCreateQueryWithJsonAndUserId
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Query(
		logCreateDto.LogMessage,
		logCreateDto.LogLevel,
		logCreateDto.LogData,
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}
