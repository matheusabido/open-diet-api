package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/sirupsen/logrus"
)

type FoodController struct {
	service *service.FoodService
}

func NewFoodController(service *service.FoodService) *FoodController {
	return &FoodController{service: service}
}

func (f *FoodController) Create(c fiber.Ctx) error {
	t := translator(c)

	body := dto.FoodDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	food, err := f.service.Create(c.Context(), middleware.User(c).ID, body)
	if err != nil {
		logrus.Errorf("create food: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.Status(fiber.StatusCreated).JSON(food)
}

func (f *FoodController) List(c fiber.Ctx) error {
	t := translator(c)

	query := dto.FoodListQueryDTO{}
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidQuery)})
	}

	foods, err := f.service.List(c.Context(), query)
	if err != nil {
		logrus.Errorf("list foods: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(foods)
}

func (f *FoodController) Show(c fiber.Ctx) error {
	t := translator(c)

	food, err := f.service.FindByID(c.Context(), fiber.Params[uint64](c, "id"))
	if errors.Is(err, repository.ErrFoodNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.FoodNotFound)})
	}
	if err != nil {
		logrus.Errorf("show food: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(food)
}

func (f *FoodController) Update(c fiber.Ctx) error {
	t := translator(c)

	body := dto.FoodDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	food, err := f.service.Update(c.Context(), fiber.Params[uint64](c, "id"), body)
	if errors.Is(err, repository.ErrFoodNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.FoodNotFound)})
	}
	if err != nil {
		logrus.Errorf("update food: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(food)
}

func (f *FoodController) Delete(c fiber.Ctx) error {
	t := translator(c)

	err := f.service.Delete(c.Context(), fiber.Params[uint64](c, "id"))
	if errors.Is(err, repository.ErrFoodNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.FoodNotFound)})
	}
	if err != nil {
		logrus.Errorf("delete food: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
