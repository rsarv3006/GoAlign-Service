package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamInviteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	teamInviteCreateDto := new(model.TeamInviteCreateDto)

	if err := c.BodyParser(teamInviteCreateDto); err != nil {
		log.Error(err)
		log.Error("Error parsing body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"error":   err,
		})
	}

	query := database.TeamInviteCreateQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	_, err = stmt.Exec(teamInviteCreateDto.TeamId, teamInviteCreateDto.Email, currentUser.UserId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Team invite created",
	})
}

func AcceptTeamInviteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamInvite := new(model.TeamInvite)

	rows, err := stmt.Query(c.Params("teamInviteId"))

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	for rows.Next() {
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
		}
	}

	if isTeamInviteStructEmpty(teamInvite) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite does not exist",
		})
	}

	if teamInvite.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite is not pending",
		})
	}

	query = database.TeamInviteAcceptQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	// TODO: check if this value requires sanitization
	teamInviteId := c.Params("teamInviteId")

	_, err = stmt.Exec(teamInviteId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamMembership, err := CreateUserTeamMembership(currentUser.UserId, teamInvite.TeamId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	team, err := getTeamById(teamInvite.TeamId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Team invite accepted",
		"teamMembership": teamMembership,
		"team":           team,
	})
}

func DeclineTeamInviteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	_, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := stmt.Query(teamInviteId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	for rows.Next() {
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
		}
	}

	if isTeamInviteStructEmpty(teamInvite) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite does not exist",
		})
	}

	if teamInvite.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite is not pending",
		})
	}

	query = database.TeamInviteDeclineQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	_, err = stmt.Exec(teamInviteId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesForCurrentUserEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	query := database.TeamInvitesForCurrentUserQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	rows, err := stmt.Query(currentUser.Email)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamInvites := make([]model.TeamInvite, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
		}

		teamInvites = append(teamInvites, teamInvite)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team invites for current user",
		"data":    teamInvites,
	})
}

func isTeamInviteStructEmpty(teamInvite *model.TeamInvite) bool {
	return teamInvite.Status == ""
}

func isAllowedToDeleteTeamInvite(teamInvite *model.TeamInvite, currentUser *model.User) bool {
	team, err := getTeamByTeamId(teamInvite.TeamId)

	if err != nil {
		log.Error(err)
		return false
	} else if team.TeamManagerId == currentUser.UserId {
		return true
	} else if teamInvite.InviteCreatorId == currentUser.UserId {
		return true
	}

	return false
}

func DeleteTeamInviteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := stmt.Query(teamInviteId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	for rows.Next() {
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
		}
	}

	if isTeamInviteStructEmpty(teamInvite) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite does not exist",
		})
	}

	if !isAllowedToDeleteTeamInvite(teamInvite, currentUser) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Team invite does not belong to current user",
		})
	}

	query = database.TeamInviteDeleteQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	_, err = stmt.Exec(teamInviteId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesByTeamIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	_, err := auth.ValidateToken(token)

	// TODO: Check if user can view team invites for this team

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	teamId := c.Params("teamId")

	query := database.TeamInviteGetByTeamIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	rows, err := stmt.Query(teamId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	teamInvites := make([]model.TeamInvite, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
		}

		teamInvites = append(teamInvites, teamInvite)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"message":     "Team invites by team id",
		"teamInvites": teamInvites,
	})
}
