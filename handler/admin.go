package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
)

func UpdateExpiredLoginRequestsEndpoint(c *fiber.Ctx) error {
	query := database.LoginRequestMarkAsExpiredQuery

	rows, err := database.POOL.Query(context.Background(), query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.SendStatus(204)
}
