package model

import "slices"

type Meal string

const (
	MealBreakfast Meal = "breakfast"
	MealLunch     Meal = "lunch"
	MealDinner    Meal = "dinner"
	MealSnacks    Meal = "snacks"
)

var Meals = []Meal{MealBreakfast, MealLunch, MealDinner, MealSnacks}

func (m Meal) Valid() bool {
	return slices.Contains(Meals, m)
}
