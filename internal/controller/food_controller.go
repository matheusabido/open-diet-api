package controller

import (
	"errors"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/groq"
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/openfoodfacts"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/matheusabido/open-diet-api/internal/validation"
	"github.com/sirupsen/logrus"
)

var barcodePattern = regexp.MustCompile(`^\d{8,14}$`)

type FoodController struct {
	service *service.FoodService
	scan    *service.FoodScanService
}

func NewFoodController(service *service.FoodService, scan *service.FoodScanService) *FoodController {
	return &FoodController{service: service, scan: scan}
}

func (f *FoodController) Barcode(c fiber.Ctx) error {
	t := translator(c)

	barcode := c.Params("barcode")
	if !barcodePattern.MatchString(barcode) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": validation.Errors{"barcode": {t(i18n.InvalidBarcode)}}})
	}

	food, err := f.scan.FromBarcode(c.Context(), barcode)
	if errors.Is(err, openfoodfacts.ErrProductNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": t(i18n.ProductNotFound)})
	}
	if err != nil {
		logrus.Errorf("food from barcode: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(food)
}

func (f *FoodController) Photo(c fiber.Ctx) error {
	t := translator(c)

	header, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": validation.Errors{"image": {t(i18n.Required)}}})
	}

	image, err := header.Open()
	if err != nil {
		logrus.Errorf("food from photo: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}
	defer image.Close()

	food, err := f.scan.FromPhoto(c.Context(), image)
	if errors.Is(err, service.ErrLabelNotFound) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": t(i18n.LabelNotFound)})
	}
	if errors.Is(err, groq.ErrRateLimited) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": t(i18n.TooManyRequests)})
	}
	if err != nil {
		logrus.Errorf("food from photo: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": t(i18n.InternalError)})
	}

	return c.JSON(food)
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
