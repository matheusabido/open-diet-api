package dto

import (
	"strings"

	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/unit"
	"github.com/matheusabido/open-diet-api/internal/validation"
)

type FoodItemDTO struct {
	FoodID uint64  `json:"food_id"`
	Date   string  `json:"date"`
	Meal   string  `json:"meal"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

func (d FoodItemDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.GreaterThan(d.FoodID, 0), "food_id", t(i18n.Required))

	v.Check(validation.ISODate(d.Date), "date", t(i18n.InvalidISODate))

	v.Check(model.Meal(d.Meal).Valid(), "meal", t(i18n.InvalidMeal, mealNames()))

	v.Check(validation.GreaterThan(d.Amount, 0), "amount", t(i18n.GreaterThan, 0))
	v.Check(validation.Max(d.Amount, maxAmount), "amount", t(i18n.MaxValue, maxAmount))

	_, ok := unit.Parse(d.Unit)
	v.Check(ok, "unit", t(i18n.InvalidUnit, strings.Join(unit.Names(), ", ")))

	return v.Errors()
}

type VoiceFoodItemDTO struct {
	AliasID uint64      `json:"alias_id"`
	Name    string      `json:"name"`
	Food    *model.Food `json:"food"`
	Amount  float64     `json:"amount"`
	Unit    unit.Unit   `json:"unit"`
}

type FoodItemListQueryDTO struct {
	Date string `query:"date"`
}

func (d FoodItemListQueryDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.ISODate(d.Date), "date", t(i18n.InvalidISODate))

	return v.Errors()
}

func mealNames() string {
	names := make([]string, len(model.Meals))
	for i, meal := range model.Meals {
		names[i] = string(meal)
	}
	return strings.Join(names, ", ")
}
