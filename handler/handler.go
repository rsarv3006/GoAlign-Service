package handler

import (
	"github.com/gofiber/fiber/v2"
)

func sendUnauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Unauthorized",
		"success": false,
	})
}

func sendInternalServerErrorResponse(c *fiber.Ctx, err error) error {
	// TODO: implement logging
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "Internal Server Error",
		"error":   err,
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
