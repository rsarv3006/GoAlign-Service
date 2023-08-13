package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/donutsahoy/yourturn-fiber/auth"
)

func IsExpired() fiber.Handler {

	return func(c *fiber.Ctx) error {
		token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
		_, err := auth.ValidateToken(token)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"error":   err,
			})
		}

		return c.Next()

	}

}
