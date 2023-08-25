package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func sendUnauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Unauthorized",
		"success": false,
	})
}

func sendInternalServerErrorResponse(c *fiber.Ctx, err error) error {
	log.Println(err)
	logCreateDto := new(model.LogCreateDto)
	logCreateDto.LogLevel = "error"
	logCreateDto.LogMessage = err.Error()
	err = logEventOnlyMessage(logCreateDto)

	if err != nil {
		log.Println(err)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "Internal Server Error",
		"error":   logCreateDto,
		"success": false,
	})
}

func sendBadRequestResponse(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message,
		"error":   err,
		"success": false,
	})
}

func sendNotFoundResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": message,
		"success": false,
	})
}

func sendForbiddenResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"message": "Forbidden",
		"success": false,
	})
}
