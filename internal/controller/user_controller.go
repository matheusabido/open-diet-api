package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/dto"
)

type UserController struct{}

func (u *UserController) Register(c fiber.Ctx) error {
	body := dto.UserRegisterDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	return c.JSON(body)
}

func (u *UserController) Login(c fiber.Ctx) error {
	return nil
}
