package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
)

func DeleteUserEndpoint(c *fiber.Ctx) error {
	// TODO: remove teamInvites created by user
	// TODO: remove teamInvites where user email is in teamInvite
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	userId := currentUser.UserId

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"error":   err,
		})
	}

	isUserATeamManagerOfAnyTeam := isUserATeamManagerOfAnyTeam(userId)

	if isUserATeamManagerOfAnyTeam {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User is a team manager of at least one team",
		})
	}

	err = deleteUserTeamMembershipsByUserId(userId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	query := database.UserDeleteUserQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	_, err = stmt.Exec(userId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
