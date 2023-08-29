package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
)

func IsExpired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !strings.Contains(c.Get("Authorization"), "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"error":   "No token provided",
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

		c.Locals("currentUser", currentUser)
		return c.Next()

	}

}
