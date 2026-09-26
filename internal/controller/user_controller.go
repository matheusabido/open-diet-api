package controller

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (u *UserController) Register(c fiber.Ctx) error {
	t := translator(c)

	body := dto.UserRegisterDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}
	body.Username = strings.ToLower(strings.TrimSpace(body.Username))

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	res, err := u.service.Register(c.Context(), body)
	if errors.Is(err, repository.ErrUsernameTaken) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": t(i18n.UsernameTaken)})
	}
	if err != nil {
		logrus.Errorf("register: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (u *UserController) Login(c fiber.Ctx) error {
	t := translator(c)

	body := dto.UserLoginDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}
	body.Username = strings.ToLower(strings.TrimSpace(body.Username))

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	res, err := u.service.Login(c.Context(), body)
	if errors.Is(err, service.ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": t(i18n.InvalidCredentials)})
	}
	if err != nil {
		logrus.Errorf("login: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(res)
}

func (u *UserController) Refresh(c fiber.Ctx) error {
	t := translator(c)

	body := dto.RefreshTokenDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	res, err := u.service.Refresh(c.Context(), body)
	if errors.Is(err, service.ErrInvalidRefreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": t(i18n.InvalidRefreshToken)})
	}
	if err != nil {
		logrus.Errorf("refresh: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(res)
}

func (u *UserController) Me(c fiber.Ctx) error {
	t := translator(c)

	user, err := u.service.FindByID(c.Context(), middleware.UserID(c))
	if errors.Is(err, repository.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.UserNotFound)})
	}
	if err != nil {
		logrus.Errorf("me: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(user)
}

func (u *UserController) UpdateMe(c fiber.Ctx) error {
	t := translator(c)

	body := dto.UserUpdateDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	user, err := u.service.Update(c.Context(), middleware.UserID(c), body)
	if errors.Is(err, repository.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.UserNotFound)})
	}
	if err != nil {
		logrus.Errorf("update me: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(user)
}
