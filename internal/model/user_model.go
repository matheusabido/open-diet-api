package model

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     Role   `json:"role"`

	KcalGoal         uint16 `json:"kcal_goal"`
	ProteinGoal      uint16 `json:"protein_goal"`
	CarbohydrateGoal uint16 `json:"carbohydrate_goal"`
	FatGoal          uint16 `json:"fat_goal"`
}
