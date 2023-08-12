package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsExpired() fiber.Handler {

	return func(c *fiber.Ctx) error {

		token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			fmt.Println("error parsing token")
			fmt.Printf("error: %v\n", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()

	}

}
