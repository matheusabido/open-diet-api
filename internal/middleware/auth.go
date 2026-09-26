package middleware

import (
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/sirupsen/logrus"
)

const userKey = "user"

func Auth(tokenService *service.TokenService, userService *service.UserService) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || accessToken == "" {
			return unauthorized(c)
		}

		userID, err := tokenService.ParseAccessToken(accessToken)
		if err != nil {
			return unauthorized(c)
		}

		user, err := userService.FindByID(c.Context(), userID)
		if errors.Is(err, repository.ErrUserNotFound) {
			return unauthorized(c)
		}
		if err != nil {
			logrus.Errorf("auth: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": translator(c)(i18n.InternalError)})
		}

		c.Locals(userKey, user)

		return c.Next()
	}
}

func RequireRole(roles ...model.Role) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !slices.Contains(roles, User(c).Role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": translator(c)(i18n.Forbidden)})
		}

		return c.Next()
	}
}

func User(c fiber.Ctx) *model.User {
	user, _ := c.Locals(userKey).(*model.User)
	return user
}

func unauthorized(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": translator(c)(i18n.Unauthorized)})
}

func translator(c fiber.Ctx) i18n.Translator {
	return i18n.New(c.AcceptsLanguages(i18n.Supported...))
}
