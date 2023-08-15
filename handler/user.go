package handler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/helper"
	"gitlab.com/donutsahoy/yourturn-fiber/model"
)

func Register(c *fiber.Ctx) error {
	userCreateDto := new(model.UserCreateDto)
	if err := c.BodyParser(userCreateDto); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	user := new(model.User)

	userName := helper.SanitizeInput(userCreateDto.UserName)

	user.UserName = userName
	user.Email = userCreateDto.Email

	if user.UserName == "" || user.Email == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Username and Email are required",
		})
	}

	user.UserId = uuid.New()
	user.IsActive = true

	if _, err := database.DB.Query("INSERT INTO users (user_id, user_name, email, is_active) VALUES ($1, $2, $3, $4)",
		user.UserId,
		user.UserName,
		user.Email,
		user.IsActive); err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), `"users_email_key"`) {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Email already exists",
			})
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"message": "User created successfully",
		"user":    user,
	})
}

func Login(c *fiber.Ctx) {
	// TODO implement login method and then email code to user then implement second method to send token
}

func FetchCode(c *fiber.Ctx) error {
	println("FetchCode UGHHHHH")
	log.Println("FetchCode")
	dto := new(model.User)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	// TODO: Add actual code checking

	userFromDb, errFromDb := database.DB.Query("SELECT * FROM users WHERE email = $1", dto.Email)

	if errFromDb != nil {
		fmt.Println("errFromDb")
		log.Println(errFromDb)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	defer userFromDb.Close()

	user := model.User{}

	for userFromDb.Next() {
		switch err := userFromDb.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified); err {
		case sql.ErrNoRows:
			fmt.Println("sql.ErrNoRows")
			return c.SendStatus(fiber.StatusUnauthorized)
		case nil:
			// Expected outcome, user found
		default:
			fmt.Println(err)
			fmt.Println("default error thing")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	log.Println("user")
	log.Println(user)
	signedTokenString, err := auth.GenerateJWT(user)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": signedTokenString})
}
