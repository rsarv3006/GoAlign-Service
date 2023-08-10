package middleware

import (
	"github.com/gofiber/basicauth"
	"github.com/gofiber/fiber"
	"gitlab.com/donutsahoy/yourturn-fiber/config"
)

func AuthReq() func(*fiber.Ctx) {
	cfg := basicauth.Config{
		Users: map[string]string{
			config.Config("USERNAME"): config.Config("PASSWORD"),
		},
	}
	err := basicauth.New(cfg)
	return err
}
