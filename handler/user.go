package handler

import (
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

	// Insert user into database
	if _, err := database.DB.Query("INSERT INTO users (user_id, user_name, email, is_active) VALUES ($1, $2, $3, $4)",
		user.UserId,
		user.UserName,
		user.Email,
		user.IsActive); err != nil {
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
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if user.UserName != "john" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})

}
