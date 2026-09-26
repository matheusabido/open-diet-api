package dto

import (
	"strings"

	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/unit"
	"github.com/matheusabido/open-diet-api/internal/validation"
)

const (
	maxNutrient = 100000
	maxAmount   = 100000
)

type FoodDTO struct {
	Name         string  `json:"name"`
	Brand        *string `json:"brand"`
	Kcal         float64 `json:"kcal"`
	Protein      float64 `json:"protein"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
}

func (d FoodDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.Name), "name", t(i18n.Required))
	v.Check(validation.MaxLength(d.Name, 100), "name", t(i18n.MaxLength, 100))

	if d.Brand != nil {
		v.Check(validation.MaxLength(*d.Brand, 100), "brand", t(i18n.MaxLength, 100))
	}

	v.Check(validation.Between(d.Kcal, 0, maxNutrient), "kcal", t(i18n.Between, 0, maxNutrient))
	v.Check(validation.Between(d.Protein, 0, maxNutrient), "protein", t(i18n.Between, 0, maxNutrient))
	v.Check(validation.Between(d.Carbohydrate, 0, maxNutrient), "carbohydrate", t(i18n.Between, 0, maxNutrient))
	v.Check(validation.Between(d.Fat, 0, maxNutrient), "fat", t(i18n.Between, 0, maxNutrient))

	v.Check(validation.GreaterThan(d.Amount, 0), "amount", t(i18n.GreaterThan, 0))
	v.Check(validation.Max(d.Amount, maxAmount), "amount", t(i18n.MaxValue, maxAmount))

	_, ok := unit.Parse(d.Unit)
	v.Check(ok, "unit", t(i18n.InvalidUnit, strings.Join(unit.Names(), ", ")))

	return v.Errors()
}

type FoodListQueryDTO struct {
	Search string `query:"search"`
	After  uint64 `query:"after"`
}
