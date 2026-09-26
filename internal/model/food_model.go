package model

import (
	"time"

	"github.com/matheusabido/open-diet-api/internal/unit"
)

type Food struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Brand *string `json:"brand"`

	Kcal         float64 `json:"kcal"`
	Protein      float64 `json:"protein"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`

	Amount float64   `json:"amount"`
	Unit   unit.Unit `json:"unit"`

	CreatedBy *uint64   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
