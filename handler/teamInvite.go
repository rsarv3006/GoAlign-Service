package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamInviteEndpoint(c *fiber.Ctx) error {

	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamInviteCreateDto := new(model.TeamInviteCreateDto)

	if err := c.BodyParser(teamInviteCreateDto); err != nil {
		return sendBadRequestResponse(c, err, "Error parsing body")
	}

	query := database.TeamInviteCreateQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	isValidEmailAddress := helper.IsValidEmailAddress(teamInviteCreateDto.Email)

	if !isValidEmailAddress {
		return sendBadRequestResponse(c, errors.New("Invalid email address"), "Invalid email address")
	}

	isAlreadyInTeamQuery := database.UserTeamMembershipGetByUserEmailAndTeamIdQueryString

	isAlreadyInTeamStmt, err := database.DB.Prepare(isAlreadyInTeamQuery)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	isAlreadyInTeamRows, err := isAlreadyInTeamStmt.Query(teamInviteCreateDto.TeamId, teamInviteCreateDto.Email)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return sendBadRequestResponse(c, errors.New("Team does not exist"), "Team does not exist")
		}

		return sendInternalServerErrorResponse(c, err)
	}

	if isAlreadyInTeamRows.Next() {
		return sendBadRequestResponse(c, errors.New("Email already in team"), "Email already in team")
	}

	isAlreadyInvitedQuery := database.TeamInviteGetByEmailAndTeamIdQueryString

	isAlreadyInvitedStmt, err := database.DB.Prepare(isAlreadyInvitedQuery)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	isAlreadyInvitedRows, err := isAlreadyInvitedStmt.Query(teamInviteCreateDto.Email, teamInviteCreateDto.TeamId)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return sendBadRequestResponse(c, errors.New("Team does not exist"), "Team does not exist")
		}

		return sendInternalServerErrorResponse(c, err)
	}

	if isAlreadyInvitedRows.Next() {
		return sendBadRequestResponse(c, errors.New("Email already invited"), "Email already invited")
	}

	defer isAlreadyInvitedStmt.Close()

	_, err = stmt.Exec(teamInviteCreateDto.TeamId, teamInviteCreateDto.Email, currentUser.UserId)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return sendBadRequestResponse(c, errors.New("Team does not exist"), "Team does not exist")
		}

		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Team invite created",
	})
}

func AcceptTeamInviteEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvite := new(model.TeamInvite)

	rows, err := stmt.Query(c.Params("teamInviteId"))

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	for rows.Next() {
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	if isTeamInviteStructEmpty(teamInvite) {
		err := errors.New("Team invite does not exist")
		return sendBadRequestResponse(c, err, "Team invite does not exist")
	}

	if teamInvite.Status != "pending" {
		err := errors.New("Team invite is not pending")
		return sendBadRequestResponse(c, err, "Team invite is not pending")
	}

	query = database.TeamInviteAcceptQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInviteId := c.Params("teamInviteId")

	teamInviteIdUUID, err := uuid.Parse(teamInviteId)

	if err != nil {
		return sendBadRequestResponse(c, err, "Error parsing team invite id")
	}

	_, err = stmt.Exec(teamInviteIdUUID)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamMembership, err := CreateUserTeamMembership(currentUser.UserId, teamInvite.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	team, err := getTeamById(teamInvite.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
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
		return sendUnauthorizedResponse(c)
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := stmt.Query(teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	for rows.Next() {
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}
	}

	if isTeamInviteStructEmpty(teamInvite) {
		err := errors.New("Team invite does not exist")
		return sendBadRequestResponse(c, err, "Team invite does not exist")
	}

	if teamInvite.Status != "pending" {
		err := errors.New("Team invite is not pending")
		return sendBadRequestResponse(c, err, "Team invite is not pending")
	}

	query = database.TeamInviteDeclineQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	_, err = stmt.Exec(teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesForCurrentUserEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	query := database.TeamInvitesForCurrentUserQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	rows, err := stmt.Query(currentUser.Email)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvites := make([]model.TeamInviteReturnWithCreator, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInviteCreator, err := getUserById(teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		team, err := getTeamById(teamInvite.TeamId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInviteReturnWithCreator := model.TeamInviteReturnWithCreator{
			TeamInvite:    &teamInvite,
			InviteCreator: teamInviteCreator,
			Team:          *team,
		}

		teamInvites = append(teamInvites, teamInviteReturnWithCreator)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team invites for current user",
		"invites": teamInvites,
	})
}

func isTeamInviteStructEmpty(teamInvite *model.TeamInvite) bool {
	return teamInvite.Status == ""
}

func isAllowedToDeleteTeamInvite(teamInvite *model.TeamInvite, currentUser *model.User) bool {
	team, err := getTeamById(teamInvite.TeamId)

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
		return sendUnauthorizedResponse(c)
	}

	query := database.TeamInviteGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := stmt.Query(teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
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
		return sendNotFoundResponse(c, "Team invite does not exist")
	}

	if !isAllowedToDeleteTeamInvite(teamInvite, currentUser) {
		return sendForbiddenResponse(c)
	}

	query = database.TeamInviteDeleteQueryString
	stmt, err = database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	_, err = stmt.Exec(teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesByTeamIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamId := c.Params("teamId")
	teamIdUUID, err := uuid.Parse(teamId)

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	team, err := getTeamById(teamIdUUID)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if team.TeamManagerId != currentUser.UserId {
		return sendForbiddenResponse(c)
	}

	query := database.TeamInviteGetByTeamIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	rows, err := stmt.Query(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvites := make([]model.TeamInviteReturnWithCreator, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInviteCreator, err := getUserById(teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		team, err := getTeamById(teamInvite.TeamId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInviteReturnWithCreator := model.TeamInviteReturnWithCreator{
			TeamInvite:    &teamInvite,
			Team:          *team,
			InviteCreator: teamInviteCreator,
		}

		teamInvites = append(teamInvites, teamInviteReturnWithCreator)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Team invites by team id",
		"teamInvites": teamInvites,
	})
}

func deleteTeamInvitesByUserEmail(email string) error {
	query := database.TeamInviteDeleteByUserEmailQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(email)

	if err != nil {
		return err
	}

	return nil
}
