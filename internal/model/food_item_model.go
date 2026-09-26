package model

import (
	"time"

	"github.com/matheusabido/open-diet-api/internal/unit"
)

type FoodItem struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"-"`
	FoodID uint64 `json:"food_id"`
	Food   *Food  `json:"food"`

	Date   string    `json:"date"`
	Meal   Meal      `json:"meal"`
	Amount float64   `json:"amount"`
	Unit   unit.Unit `json:"unit"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
