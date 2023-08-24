package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func GetTeamsForCurrentUser(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		log.Error(err)
		return sendUnauthorizedResponse(c)
	}

	query := database.TeamGetByUserIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return sendInternalServerErrorResponse(c, err)
	}

	rows, err := stmt.Query(currentUser.UserId)

	if err != nil {
		log.Error(err)
		return sendInternalServerErrorResponse(c, err)
	}

	teams := make([]model.Team, 0)

	for rows.Next() {
		team := model.Team{}
		err := rows.Scan(
			&team.TeamId,
			&team.TeamName,
			&team.CreatorUserId,
			&team.Status,
			&team.TeamManagerId,
			&team.CreatedAt,
			&team.UpdatedAt,
		)

		if err != nil {
			log.Error(err)
			return sendInternalServerErrorResponse(c, err)
		}

		teams = append(teams, team)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"teams":   teams,
	})
}

func CreateTeam(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	teamDto := new(model.TeamCreateDto)

	if err := c.BodyParser(teamDto); err != nil {
		return sendBadRequestResponse(c, err, "Unable to parse team create dto")
	}

	if teamDto.TeamName == "" {
		return sendBadRequestResponse(c, err, "Team name is required")
	}

	query := database.TeamCreateQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			fmt.Println("Postgres error:", e.Message)

		default:
			fmt.Println("Unknown error:", e)
		}

		return sendInternalServerErrorResponse(c, err)
	}
	defer stmt.Close()

	team := new(model.Team)

	cleanedTeamName := helper.SanitizeInput(teamDto.TeamName)
	rows, err := stmt.Query(cleanedTeamName, currentUser.UserId, currentUser.UserId)

	if err != nil {
		switch e := err.(type) {
		case *pq.Error:
			fmt.Println("Postgres error:", e.Message)

		default:
			fmt.Println("Unknown error:", e)

		}
		return sendBadRequestResponse(c, err, "")
	}

	for rows.Next() {
		err := rows.Scan(&team.TeamId, &team.TeamName, &team.CreatorUserId, &team.Status, &team.TeamManagerId, &team.CreatedAt, &team.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	userTeamMembership, err := CreateUserTeamMembership(currentUser.UserId, team.TeamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	teamSettingsDto := new(model.TeamSettingsCreateDto)
	teamSettingsDto.TeamId = team.TeamId
	teamSettingsDto.CanAllMembersAddTasks = false
	teamSettings, err := CreateTeamSettings(*teamSettingsDto)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(201).JSON(&fiber.Map{
		"success":  true,
		"message":  "Team created successfully",
		"team":     team,
		"members":  userTeamMembership,
		"settings": teamSettings,
	})
}

func getTeamById(teamId uuid.UUID) (*model.Team, error) {
	query := database.TeamGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	team := new(model.Team)
	err = stmt.QueryRow(teamId).Scan(&team.TeamId, &team.TeamName, &team.CreatorUserId, &team.Status, &team.TeamManagerId, &team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func DeleteTeam(c *fiber.Ctx) error {
	teamId, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	teamToDelete, err := getTeamById(teamId)

	if err != nil {
		return sendBadRequestResponse(c, err, "Invalid team id")
	}

	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return sendUnauthorizedResponse(c)
	}

	if teamToDelete.TeamManagerId != currentUser.UserId {
		return sendUnauthorizedResponse(c)
	}

	err = DeleteTeamSettingsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteUserTeamMembershipsByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTaskEntriesByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTasksByTeamId(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	err = deleteTeam(teamId)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteTeam(teamId uuid.UUID) error {
	query := database.TeamDeleteQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(teamId)

	if err != nil {
		return err
	}

	return nil
}

func getTeamByTeamId(uuid.UUID) (*model.Team, error) {
	query := database.TeamGetByIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	team := new(model.Team)
	rows, err := stmt.Query(team.TeamId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&team.TeamId, &team.TeamName, &team.CreatorUserId, &team.Status, &team.TeamManagerId, &team.CreatedAt, &team.UpdatedAt)

		if err != nil {
			return nil, err
		}
	}

	return team, nil
}

func isUserATeamManagerOfAnyTeam(userId uuid.UUID) bool {
	query := database.TeamGetByTeamManagerIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(userId)

	if err != nil {
		return false
	}

	for rows.Next() {
		return true
	}

	return false
}

func GetTeamByTeamIdEndpoint(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

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

	isUserInTeam, err := isUserInTeam(currentUser.UserId, teamId)

	if err != nil {
		return sendBadRequestResponse(c, err, "User is not in team")
	}

	if !isUserInTeam {
		return sendUnauthorizedResponse(c)
	}

	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"team":    team,
	})
}
