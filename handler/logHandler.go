package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func LogEventEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(model.User)

	logCreateDto := new(model.LogCreateDto)
	if err := c.BodyParser(logCreateDto); err != nil {
		log.Println(err)
		return sendBadRequestResponse(c, err, "Error parsing request body")
	}

	err := logEvent(logCreateDto, currentUser.UserId)

	if err != nil {
		log.Println(err)
		return sendInternalServerErrorResponse(c, err)
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

	row, err := stmt.Query(
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
	query := database.LogCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Query(
		logCreateDto.LogMessage,
		logCreateDto.LogLevel,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}
