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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	query := database.TeamGetByUserIdQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	rows, err := stmt.Query(currentUser.UserId)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err,
			})
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	teamDto := new(model.TeamCreateDto)

	if err := c.BodyParser(teamDto); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if teamDto.TeamName == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Team name is required",
		})
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
		return c.SendStatus(fiber.StatusInternalServerError)
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
		return c.SendStatus(fiber.StatusBadRequest)
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
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	teamSettingsDto := new(model.TeamSettingsCreateDto)
	teamSettingsDto.TeamId = team.TeamId
	teamSettingsDto.CanAllMembersAddTasks = false
	teamSettings, err := CreateTeamSettings(*teamSettingsDto)

	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid team id",
			"error":   err,
		})
	}

	teamToDelete, err := getTeamById(teamId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid team id",
			"error":   err,
		})
	}
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	if teamToDelete.TeamManagerId != currentUser.UserId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	err = DeleteTeamSettingsByTeamId(teamId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	err = deleteUserTeamMembershipsByTeamId(teamId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	err = deleteTaskEntriesByTeamId(teamId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	err = deleteTasksByTeamId(teamId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	err = deleteTeam(teamId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func deleteTaskEntriesByTeamId(teamId uuid.UUID) error {
	query := database.TaskEntryDeleteByTeamIdQuery
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
