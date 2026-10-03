package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matheusabido/open-diet-api/internal/model"
)

var (
	ErrFoodNotFound = errors.New("food not found")
	ErrBarcodeTaken = errors.New("barcode already taken")
)

const foodPageSize = 20

const foodColumns = `f.id, f.name, f.brand, f.barcode, f.kcal, f.protein, f.carbohydrate, f.fat, f.amount, f.unit, f.created_by, f.created_at, f.updated_at`

type FoodRepository struct {
	db *pgxpool.Pool
}

func NewFoodRepository(db *pgxpool.Pool) *FoodRepository {
	return &FoodRepository{db: db}
}

func (r *FoodRepository) Create(ctx context.Context, food *model.Food) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO foods (name, brand, barcode, kcal, protein, carbohydrate, fat, amount, unit, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at, updated_at`,
		food.Name, food.Brand, food.Barcode, food.Kcal, food.Protein, food.Carbohydrate, food.Fat, food.Amount, food.Unit, food.CreatedBy,
	).Scan(&food.ID, &food.CreatedAt, &food.UpdatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return ErrBarcodeTaken
	}

	return err
}

func (r *FoodRepository) FindByID(ctx context.Context, id uint64) (*model.Food, error) {
	rows, err := r.db.Query(ctx, `SELECT `+foodColumns+` FROM foods f WHERE f.id = $1`, id)
	if err != nil {
		return nil, err
	}

	food, err := pgx.CollectExactlyOneRow(rows, scanFood)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFoodNotFound
	}
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (r *FoodRepository) FindByBarcode(ctx context.Context, barcode string) (*model.Food, error) {
	rows, err := r.db.Query(ctx, `SELECT `+foodColumns+` FROM foods f WHERE f.barcode = $1`, barcode)
	if err != nil {
		return nil, err
	}

	food, err := pgx.CollectExactlyOneRow(rows, scanFood)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFoodNotFound
	}
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (r *FoodRepository) List(ctx context.Context, search string, after uint64) ([]*model.Food, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+foodColumns+`
		 FROM foods f
		 WHERE f.id > $2 AND ($1 = '' OR f.name ILIKE '%' || $1 || '%' OR f.brand ILIKE '%' || $1 || '%')
		 ORDER BY f.id
		 LIMIT $3`,
		search, after, foodPageSize,
	)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, scanFood)
}

func (r *FoodRepository) Update(ctx context.Context, food *model.Food) error {
	err := r.db.QueryRow(ctx,
		`UPDATE foods
		 SET name = $1, brand = $2, kcal = $3, protein = $4, carbohydrate = $5, fat = $6, amount = $7, unit = $8, updated_at = NOW()
		 WHERE id = $9
		 RETURNING updated_at`,
		food.Name, food.Brand, food.Kcal, food.Protein, food.Carbohydrate, food.Fat, food.Amount, food.Unit, food.ID,
	).Scan(&food.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFoodNotFound
	}

	return err
}

func (r *FoodRepository) Delete(ctx context.Context, id uint64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM foods WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrFoodNotFound
	}

	return nil
}

func scanFood(row pgx.CollectableRow) (*model.Food, error) {
	food := &model.Food{}
	err := row.Scan(&food.ID, &food.Name, &food.Brand, &food.Barcode, &food.Kcal, &food.Protein, &food.Carbohydrate, &food.Fat,
		&food.Amount, &food.Unit, &food.CreatedBy, &food.CreatedAt, &food.UpdatedAt)
	return food, err
}
