package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
)

func DeleteUserEndpoint(c *fiber.Ctx) error {
	// TODO: remove teamInvites created by user
	// TODO: remove teamInvites where user email is in teamInvite
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	userId := currentUser.UserId

	isUserATeamManagerOfAnyTeam := isUserATeamManagerOfAnyTeam(userId)

	if isUserATeamManagerOfAnyTeam {

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "User is a team manager of at least one team",
		})
	}

	err = deleteUserTeamMembershipsByUserId(userId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	query := database.UserDeleteUserQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	_, err = stmt.Exec(userId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
