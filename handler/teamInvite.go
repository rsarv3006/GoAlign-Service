package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamInviteEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamInviteCreateDto := new(model.TeamInviteCreateDto)

	if err := c.BodyParser(teamInviteCreateDto); err != nil {
		return sendBadRequestResponse(c, err, "Error parsing body")
	}

	teamInviteEmail := strings.ToLower(teamInviteCreateDto.Email)

	isValidEmailAddress := helper.IsValidEmailAddress(teamInviteEmail)

	if !isValidEmailAddress {
		return sendBadRequestResponse(c, errors.New("Invalid email address"), "Invalid email address")
	}

	isAlreadyInTeamRows, err := database.POOL.Query(
		context.Background(),
		database.UserTeamMembershipGetByUserEmailAndTeamIdQueryString,
		teamInviteCreateDto.TeamId,
		teamInviteEmail)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			defer isAlreadyInTeamRows.Close()
			return sendBadRequestResponse(c, errors.New("Team does not exist"), "Team does not exist")
		}

		return sendInternalServerErrorResponse(c, err)
	}

	defer isAlreadyInTeamRows.Close()

	if isAlreadyInTeamRows.Next() {
		return sendBadRequestResponse(c, errors.New("Email already in team"), "Email already in team")
	}

	isAlreadyInvitedRows, err := database.POOL.Query(
		context.Background(),
		database.TeamInviteGetByEmailAndTeamIdQueryString,
		teamInviteEmail,
		teamInviteCreateDto.TeamId)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			defer isAlreadyInvitedRows.Close()
			return sendBadRequestResponse(c, errors.New("Team does not exist"), "Team does not exist")
		}

		return sendInternalServerErrorResponse(c, err)
	}

	defer isAlreadyInvitedRows.Close()

	if isAlreadyInvitedRows.Next() {
		return sendBadRequestResponse(c, errors.New("Email already invited"), "Email already invited")
	}

	_, err = database.POOL.Exec(
		context.Background(),
		database.TeamInviteCreateQueryString,
		teamInviteCreateDto.TeamId,
		teamInviteEmail,
		currentUser.UserId)

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
	currentUser := c.Locals("currentUser").(*model.User)

	teamInvite := new(model.TeamInvite)

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamInviteGetByIdQueryString,
		c.Params("teamInviteId"))

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

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

	teamInviteId := c.Params("teamInviteId")

	teamInviteIdUUID, err := uuid.Parse(teamInviteId)

	if err != nil {
		return sendBadRequestResponse(c, err, "Error parsing team invite id")
	}

	_, err = database.POOL.Exec(
		context.Background(),
		database.TeamInviteAcceptQueryString,
		teamInviteIdUUID)

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
	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamInviteGetByIdQueryString,
		teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

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

	_, err = database.POOL.Exec(
		context.Background(),
		database.TeamInviteDeclineQueryString,
		teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesForCurrentUserEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamInvitesForCurrentUserQueryString,
		strings.ToLower(currentUser.Email))

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	teamInvites := make([]model.TeamInvite, 0)
	inviteCreatorIds := make([]uuid.UUID, 0)
	teamIds := make([]uuid.UUID, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInvites = append(teamInvites, teamInvite)
		inviteCreatorIds = append(inviteCreatorIds, teamInvite.InviteCreatorId)
		teamIds = append(teamIds, teamInvite.TeamId)
	}

	inviteCreators, err := getUsersByIdArray(inviteCreatorIds)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teams, err := getTeamsByTeamIdArray(teamIds)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvitesWithData := make([]model.TeamInviteReturnWithCreator, 0)

	for _, teamInvite := range teamInvites {
		teamInvite := model.TeamInviteReturnWithCreator{
			TeamInvite:    &teamInvite,
			InviteCreator: inviteCreators[teamInvite.InviteCreatorId],
			Team:          teams[teamInvite.TeamId],
		}

		teamInvitesWithData = append(teamInvitesWithData, teamInvite)
	}

	defer func() {
		teamInvitesWithData = nil
		teams = nil
		inviteCreators = nil
		teamIds = nil
		inviteCreatorIds = nil
		teamInvites = nil
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team invites for current user",
		"invites": teamInvitesWithData,
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
	currentUser := c.Locals("currentUser").(*model.User)

	teamInvite := new(model.TeamInvite)
	teamInviteId := c.Params("teamInviteId")

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamInviteGetByIdQueryString,
		teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

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

	_, err = database.POOL.Exec(
		context.Background(),
		database.TeamInviteDeleteQueryString,
		teamInviteId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTeamInvitesByTeamIdEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

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

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamInviteGetByTeamIdQueryString,
		teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	teamInvites := make([]model.TeamInvite, 0)
	creatorIds := make([]uuid.UUID, 0)

	for rows.Next() {
		var teamInvite model.TeamInvite
		err := rows.Scan(&teamInvite.TeamInviteId, &teamInvite.TeamId, &teamInvite.Email, &teamInvite.CreatedAt, &teamInvite.UpdatedAt, &teamInvite.Status, &teamInvite.InviteCreatorId)

		if err != nil {
			return sendInternalServerErrorResponse(c, err)
		}

		teamInvites = append(teamInvites, teamInvite)
		creatorIds = append(creatorIds, teamInvite.InviteCreatorId)
	}

	creators, err := getUsersByIdArray(creatorIds)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamInvitesWithData := make([]model.TeamInviteReturnWithCreator, 0)

	for _, teamInvite := range teamInvites {
		teamInvite := model.TeamInviteReturnWithCreator{
			TeamInvite:    &teamInvite,
			InviteCreator: creators[teamInvite.InviteCreatorId],
			Team:          *team,
		}

		teamInvitesWithData = append(teamInvitesWithData, teamInvite)
	}

	defer func() {
		teamInvitesWithData = nil
		teamInvites = nil
		creatorIds = nil
		creators = nil
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Team invites by team id",
		"teamInvites": teamInvitesWithData,
	})
}

func deleteTeamInvitesByUserEmail(email string) error {
	_, err := database.POOL.Exec(
		context.Background(),
		database.TeamInviteDeleteByUserEmailQueryString,
		strings.ToLower(email))

	if err != nil {
		return err
	}

	return nil
}

func deleteTeamInvitesByTeamId(teamId uuid.UUID) error {
	_, err := database.POOL.Exec(
		context.Background(),
		database.TeamInviteDeleteByTeamIdQueryString,
		teamId)

	if err != nil {
		return err
	}

	return nil
}
