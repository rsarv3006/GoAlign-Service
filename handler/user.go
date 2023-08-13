package handler

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.com/donutsahoy/yourturn-fiber/database"
	"gitlab.com/donutsahoy/yourturn-fiber/model"

	"github.com/golang-jwt/jwt/v5"
)

func Register(c *fiber.Ctx) error {
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

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
	dto := new(model.User)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if dto.UserName != "john" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userFromDb, errFromDb := database.DB.Query("SELECT * FROM users WHERE email = $1", dto.Email)

	if errFromDb != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	defer userFromDb.Close()

	user := model.User{}

	for userFromDb.Next() {
		switch err := userFromDb.Scan(&user.UserId, &user.UserName, &user.Email, &user.IsActive); err {
		case sql.ErrNoRows:
			return c.SendStatus(fiber.StatusUnauthorized)
		case nil:
			// Expected outcome, user found
		default:
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedTokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": signedTokenString})

}
