package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeam(c *fiber.Ctx) error {
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	teamdto := new(model.TeamCreateDto)

	if err := c.BodyParser(teamdto); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if teamdto.TeamName == "" {
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

	cleanedTeamName := helper.SanitizeInput(teamdto.TeamName)
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

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Team created successfully",
		"team":    team,
		"members": userTeamMembership,
	})
}
