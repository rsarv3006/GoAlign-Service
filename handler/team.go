package handler

import (
	"fmt"
	"time"

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

	query := `INSERT INTO teams (team_id, team_name, creator_user_id, status, team_manager_id, created_at, updated_at)  
VALUES ($1, $2, $3, $4, $5, $6, $7)`

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

	team := model.Team{
		TeamId:        uuid.New(),
		TeamName:      teamdto.TeamName,
		CreatorUserId: userId,
		Status:        "active",
		TeamManagerId: userId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if _, err := stmt.Exec(team.TeamId, team.TeamName, team.CreatorUserId, team.Status, team.TeamManagerId, team.CreatedAt, team.UpdatedAt); err != nil {
		println("error preparing query")
		switch e := err.(type) {
		case *pq.Error:
			fmt.Println("Postgres error:", e.Message)

		default:
			fmt.Println("Unknown error:", e)

		}
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Team created successfully",
		"team":    team,
	})
}
