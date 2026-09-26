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
	"github.com/matheusabido/open-diet-api/internal/validation"
	"github.com/sirupsen/logrus"
)

type FoodItemController struct {
	service *service.FoodItemService
}

func NewFoodItemController(service *service.FoodItemService) *FoodItemController {
	return &FoodItemController{service: service}
}

func (f *FoodItemController) Create(c fiber.Ctx) error {
	t := translator(c)

	body := dto.FoodItemDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	item, err := f.service.Create(c.Context(), middleware.User(c).ID, body)
	if err != nil {
		return f.writeError(c, t, "create food item", err)
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (f *FoodItemController) List(c fiber.Ctx) error {
	t := translator(c)

	query := dto.FoodItemListQueryDTO{}
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidQuery)})
	}

	if errs := query.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	items, err := f.service.List(c.Context(), middleware.User(c).ID, query)
	if err != nil {
		return f.writeError(c, t, "list food items", err)
	}

	return c.JSON(items)
}

func (f *FoodItemController) Show(c fiber.Ctx) error {
	t := translator(c)

	item, err := f.service.FindByID(c.Context(), middleware.User(c).ID, fiber.Params[uint64](c, "id"))
	if err != nil {
		return f.writeError(c, t, "show food item", err)
	}

	return c.JSON(item)
}

func (f *FoodItemController) Update(c fiber.Ctx) error {
	t := translator(c)

	body := dto.FoodItemDTO{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": t(i18n.InvalidBody)})
	}

	if errs := body.Validate(t); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	item, err := f.service.Update(c.Context(), middleware.User(c).ID, fiber.Params[uint64](c, "id"), body)
	if err != nil {
		return f.writeError(c, t, "update food item", err)
	}

	return c.JSON(item)
}

func (f *FoodItemController) Delete(c fiber.Ctx) error {
	t := translator(c)

	err := f.service.Delete(c.Context(), middleware.User(c).ID, fiber.Params[uint64](c, "id"))
	if err != nil {
		return f.writeError(c, t, "delete food item", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (f *FoodItemController) writeError(c fiber.Ctx, t i18n.Translator, action string, err error) error {
	var incompatible *service.IncompatibleUnitError

	switch {
	case errors.Is(err, repository.ErrFoodItemNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.FoodItemNotFound)})
	case errors.Is(err, repository.ErrFoodNotFound):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": validation.Errors{"food_id": {t(i18n.FoodNotFound)}}})
	case errors.As(err, &incompatible):
		allowed := make([]string, len(incompatible.Allowed))
		for i, u := range incompatible.Allowed {
			allowed[i] = string(u)
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": validation.Errors{"unit": {t(i18n.IncompatibleUnit, strings.Join(allowed, ", "))}}})
	default:
		logrus.Errorf("%s: %v", action, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}
}
