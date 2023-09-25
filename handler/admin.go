package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
)

func UpdateExpiredLoginRequestsEndpoint(c *fiber.Ctx) error {
	query := database.LoginRequestMarkAsExpiredQuery

	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	return c.SendStatus(204)
}
