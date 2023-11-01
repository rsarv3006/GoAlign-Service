package handler

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamSettings(dto model.TeamSettingsCreateDto) (*model.TeamSettings, error) {
	var teamSettings *model.TeamSettings

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamSettingsCreateQuery,
		dto.TeamId,
		dto.CanAllMembersAddTasks)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		teamSettings = &model.TeamSettings{}
		err := rows.Scan(
			&teamSettings.TeamSettingsId,
			&teamSettings.TeamId,
			&teamSettings.CanAllMembersAddTasks,
		)

		if err != nil {
			return nil, err
		}
	}

	return teamSettings, nil
}

func DeleteTeamSettingsByTeamId(teamId uuid.UUID) error {
	_, err := database.POOL.Exec(
		context.Background(),
		database.TeamSettingsDeleteByTeamIdQuery,
		teamId)

	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamSettingsEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	team, err := getTeamById(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if team.TeamManagerId != currentUser.UserId {
		return sendForbiddenResponse(c)
	}

	teamSettingsDto := new(model.TeamSettingsUpdateDto)

	if err := c.BodyParser(teamSettingsDto); err != nil {
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	if teamSettingsDto.CanAllMembersAddTasks == nil {
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamSettingsUpdateQuery,
		teamSettingsDto.CanAllMembersAddTasks,
		teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer rows.Close()

	teamSettings := new(model.TeamSettings)

	if rows.Next() {
		err := rows.Scan(&teamSettings.TeamSettingsId, &teamSettings.TeamId, &teamSettings.CanAllMembersAddTasks)

		if err != nil {
			if err == sql.ErrNoRows {
				return sendNotFoundResponse(c, "Team settings not found")
			}
			return sendInternalServerErrorResponse(c, err)
		}
	}

	return c.Status(201).JSON(&fiber.Map{
		"settings": teamSettings,
	})
}

func getTeamSettingsByTeamId(teamId uuid.UUID) (*model.TeamSettings, error) {
	rows, err := database.POOL.Query(
		context.Background(),
		database.TeamSettingsGetByTeamIdQuery,
		teamId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teamSettings := new(model.TeamSettings)

	if rows.Next() {
		err := rows.Scan(&teamSettings.TeamSettingsId, &teamSettings.TeamId, &teamSettings.CanAllMembersAddTasks)

		if err != nil {
			return nil, err
		}
	}

	return teamSettings, nil
}

func GetTeamSettingsByTeamIdEndpoint(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(*model.User)

	teamId, err := uuid.Parse(c.Params("teamId"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	isUserInTeam, err := isUserInTeam(currentUser.UserId, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if !isUserInTeam {
		return sendForbiddenResponse(c)
	}

	teamSettings, err := getTeamSettingsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	if teamSettings == nil {
		return sendNotFoundResponse(c, "Team settings not found")
	}

	return c.Status(200).JSON(&fiber.Map{
		"settings": teamSettings,
	})
}
