package model

import "time"

// FoodAlias links what the user says to a food. The name never leaves the
// API, and the food is nil until the user links it.
type FoodAlias struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"-"`
	Name   string  `json:"-"`
	FoodID *uint64 `json:"food_id"`
	Food   *Food   `json:"food"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
