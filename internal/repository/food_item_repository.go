package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matheusabido/open-diet-api/internal/model"
)

var ErrFoodItemNotFound = errors.New("food item not found")

const foodItemPageSize = 100

const foodItemSelect = `SELECT fi.id, fi.user_id, fi.food_id, fi.date::text, fi.meal, fi.amount, fi.unit, fi.created_at, fi.updated_at, ` + foodColumns + `
	FROM food_items fi
	JOIN foods f ON f.id = fi.food_id`

type FoodItemRepository struct {
	db *pgxpool.Pool
}

func NewFoodItemRepository(db *pgxpool.Pool) *FoodItemRepository {
	return &FoodItemRepository{db: db}
}

func (r *FoodItemRepository) Create(ctx context.Context, item *model.FoodItem) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO food_items (user_id, food_id, date, meal, amount, unit)
		 VALUES ($1, $2, $3::date, $4, $5, $6)
		 RETURNING id, created_at, updated_at`,
		item.UserID, item.FoodID, item.Date, item.Meal, item.Amount, item.Unit,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *FoodItemRepository) FindByID(ctx context.Context, userID, id uint64) (*model.FoodItem, error) {
	rows, err := r.db.Query(ctx, foodItemSelect+` WHERE fi.user_id = $1 AND fi.id = $2`, userID, id)
	if err != nil {
		return nil, err
	}

	item, err := pgx.CollectExactlyOneRow(rows, scanFoodItem)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFoodItemNotFound
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *FoodItemRepository) ListByDate(ctx context.Context, userID uint64, date string) ([]*model.FoodItem, error) {
	rows, err := r.db.Query(ctx,
		foodItemSelect+`
		 WHERE fi.user_id = $1 AND fi.date = $2::date
		 ORDER BY fi.id
		 LIMIT $3`,
		userID, date, foodItemPageSize,
	)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, scanFoodItem)
}

func (r *FoodItemRepository) Update(ctx context.Context, item *model.FoodItem) error {
	err := r.db.QueryRow(ctx,
		`UPDATE food_items
		 SET food_id = $1, date = $2::date, meal = $3, amount = $4, unit = $5, updated_at = NOW()
		 WHERE user_id = $6 AND id = $7
		 RETURNING updated_at`,
		item.FoodID, item.Date, item.Meal, item.Amount, item.Unit, item.UserID, item.ID,
	).Scan(&item.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFoodItemNotFound
	}

	return err
}

func (r *FoodItemRepository) Delete(ctx context.Context, userID, id uint64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM food_items WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrFoodItemNotFound
	}

	return nil
}

func scanFoodItem(row pgx.CollectableRow) (*model.FoodItem, error) {
	item := &model.FoodItem{Food: &model.Food{}}
	food := item.Food
	err := row.Scan(&item.ID, &item.UserID, &item.FoodID, &item.Date, &item.Meal, &item.Amount, &item.Unit, &item.CreatedAt, &item.UpdatedAt,
		&food.ID, &food.Name, &food.Brand, &food.Kcal, &food.Protein, &food.Carbohydrate, &food.Fat,
		&food.Amount, &food.Unit, &food.CreatedBy, &food.CreatedAt, &food.UpdatedAt)
	return item, err
}
