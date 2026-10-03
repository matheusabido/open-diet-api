package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/matheusabido/open-diet-api/internal/validation"
	"github.com/sirupsen/logrus"
)

type FoodAliasController struct {
	service *service.FoodAliasService
}

func NewFoodAliasController(service *service.FoodAliasService) *FoodAliasController {
	return &FoodAliasController{service: service}
}

func (f *FoodAliasController) Update(c fiber.Ctx) error {
	t := translator(c)

	body := dto.FoodAliasDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	alias, err := f.service.Update(c.Context(), middleware.User(c).ID, fiber.Params[uint64](c, "id"), body)
	switch {
	case errors.Is(err, repository.ErrFoodAliasNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.FoodAliasNotFound)})
	case errors.Is(err, repository.ErrFoodNotFound):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": validation.Errors{"food_id": {t(i18n.FoodNotFound)}}})
	case err != nil:
		logrus.Errorf("update food alias: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(alias)
}
