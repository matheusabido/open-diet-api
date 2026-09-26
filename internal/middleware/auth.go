package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/service"
)

const userIDKey = "userID"

func Auth(tokenService *service.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || accessToken == "" {
			return unauthorized(c)
		}

		userID, err := tokenService.ParseAccessToken(accessToken)
		if err != nil {
			return unauthorized(c)
		}

		c.Locals(userIDKey, userID)

		return c.Next()
	}
}

func UserID(c fiber.Ctx) uint64 {
	userID, _ := c.Locals(userIDKey).(uint64)
	return userID
}

func unauthorized(c fiber.Ctx) error {
	t := i18n.New(c.AcceptsLanguages(i18n.Supported...))
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": t(i18n.Unauthorized)})
}
