package handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeam(c *fiber.Ctx) error {
	log.Println("CreateTeam")
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	currentUser, err := auth.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err,
		})
	}

	log.Println("post token stuff")
	teamdto := new(model.Team)

	if err := c.BodyParser(teamdto); err != nil {
		println("error parsing teamdto")
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if teamdto.TeamName == "" {
		println("team name is required")
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Team name is required",
		})
	}

	query := database.TeamCreateQueryString
	stmt, err := database.DB.Prepare(query)

	if err != nil {
		println("error preparing query")
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
		println("error executing query")
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

	log.Println("team created successfully")
	log.Println(team)

	userTeamMembership, err := CreateUserTeamMembership(currentUser.UserId, team.TeamId)
	if err != nil {
		fmt.Println(err)
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