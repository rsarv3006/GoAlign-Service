package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func LogEventEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	logCreateDto := new(model.LogCreateDto)
	if err := c.BodyParser(logCreateDto); err != nil {
		return sendBadRequestResponse(c, err, "Error parsing request body")
	}

	err := logEvent(logCreateDto, currentUser.UserId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func logEvent(logCreateDto *model.LogCreateDto, userId uuid.UUID) error {
	row, err := database.POOL.Query(
		context.Background(),
		database.LogCreateQueryWithJsonAndUserId,
		logCreateDto.LogMessage,
		logCreateDto.LogLevel,
		logCreateDto.LogData,
		userId,
	)

	if err != nil {
		return err
	}

	row.Close()

	return nil
}

func logEventOnlyMessage(logCreateDto *model.LogCreateDto) error {
	rows, err := database.POOL.Query(
		context.Background(),
		database.LogCreateQuery,
		logCreateDto.LogMessage,
		logCreateDto.LogLevel,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}
