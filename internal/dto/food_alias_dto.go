package dto

import (
	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/validation"
)

type FoodAliasDTO struct {
	FoodID uint64 `json:"food_id"`
}

func (d FoodAliasDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.GreaterThan(d.FoodID, 0), "food_id", t(i18n.Required))

	return v.Errors()
}
