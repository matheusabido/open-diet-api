package model

import "time"

type FoodAlias struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"-"`
	Name   string  `json:"-"`
	FoodID *uint64 `json:"food_id"`
	Food   *Food   `json:"food"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
