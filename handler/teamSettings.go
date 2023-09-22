package handler

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeamSettings(dto model.TeamSettingsCreateDto) (*model.TeamSettings, error) {
	query := database.TeamSettingsCreateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	var teamSettings *model.TeamSettings

	rows, err := stmt.Query(dto.TeamId, dto.CanAllMembersAddTasks)

	if err != nil {
		return nil, err
	}

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
	query := database.TeamSettingsDeleteByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(teamId)

	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamSettingsEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token, c)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

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

	query := database.TeamSettingsUpdateQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamSettingsDto.CanAllMembersAddTasks, teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

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
	query := database.TeamSettingsGetByTeamIdQuery
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(teamId)

	if err != nil {
		return nil, err
	}

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
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token, c)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

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
