package handler

import (
	"database/sql"
	"errors"
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
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	user := new(model.User)

	userName := helper.SanitizeInput(userCreateDto.UserName)

	user.UserName = userName
	user.Email = userCreateDto.Email

	if user.UserName == "" || user.Email == "" {
		err := fmt.Errorf("Username and Email are required")
		return sendBadRequestResponse(c, err, "Username and Email are required")
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
			return sendBadRequestResponse(c, err, "Email already exists")
		}
		return sendInternalServerErrorResponse(c, err)
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
	dto := new(model.User)
	if err := c.BodyParser(dto); err != nil {
		return sendBadRequestResponse(c, err, "Invalid request body")
	}

	err := errors.New("Waffles are delicious")
	return sendInternalServerErrorResponse(c, err)
	// TODO: Add actual code checking

	userFromDb, errFromDb := database.DB.Query("SELECT * FROM users WHERE email = $1", dto.Email)

	if errFromDb != nil {
		log.Println(errFromDb)
		return sendUnauthorizedResponse(c)
	}

	defer userFromDb.Close()

	user := model.User{}

	for userFromDb.Next() {
		switch err := userFromDb.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive, &user.IsEmailVerified); err {
		case sql.ErrNoRows:
			fmt.Println("sql.ErrNoRows")
			return sendUnauthorizedResponse(c)
		case nil:
			// Expected outcome, user found
		default:
			fmt.Println(err)
			return sendUnauthorizedResponse(c)
		}
	}

	if !isUserObjectNotNil(&user) {
		return sendUnauthorizedResponse(c)
	}

	signedTokenString, err := auth.GenerateJWT(user)

	if err != nil {
		return sendInternalServerErrorResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": signedTokenString})
}

func isUserObjectNotNil(user *model.User) bool {
	return user.UserId != uuid.Nil

}
