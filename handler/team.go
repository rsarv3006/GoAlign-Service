package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func CreateTeam(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userFromClaim := claims["user"]
	userMap := userFromClaim.(map[string]interface{})

	var userId uuid.UUID
	userId, err := uuid.Parse(userMap["user_id"].(string))
	if err != nil {
		fmt.Println("error parsing uuid from token")
		fmt.Println(err)
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

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

	rows, err := stmt.Query(teamdto.TeamName, userId, userId)

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

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Team created successfully",
		"team":    team,
	})
}
