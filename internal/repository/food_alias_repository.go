package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matheusabido/open-diet-api/internal/model"
)

var ErrFoodAliasNotFound = errors.New("food alias not found")

type FoodAliasRepository struct {
	db *pgxpool.Pool
}

func NewFoodAliasRepository(db *pgxpool.Pool) *FoodAliasRepository {
	return &FoodAliasRepository{db: db}
}

// FindOrCreate returns the user's alias with the given name, creating it
// without a food when the user never said that name before.
func (r *FoodAliasRepository) FindOrCreate(ctx context.Context, userID uint64, name string) (*model.FoodAlias, error) {
	alias := &model.FoodAlias{UserID: userID, Name: name}
	err := r.db.QueryRow(ctx,
		`INSERT INTO food_aliases (user_id, name)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, food_id, created_at, updated_at`,
		userID, name,
	).Scan(&alias.ID, &alias.FoodID, &alias.CreatedAt, &alias.UpdatedAt)

	return alias, err
}

func (r *FoodAliasRepository) FindByID(ctx context.Context, userID, id uint64) (*model.FoodAlias, error) {
	alias := &model.FoodAlias{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, food_id, created_at, updated_at
		 FROM food_aliases
		 WHERE user_id = $1 AND id = $2`,
		userID, id,
	).Scan(&alias.ID, &alias.UserID, &alias.Name, &alias.FoodID, &alias.CreatedAt, &alias.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFoodAliasNotFound
	}
	if err != nil {
		return nil, err
	}

	return alias, nil
}

func (r *FoodAliasRepository) Update(ctx context.Context, alias *model.FoodAlias) error {
	err := r.db.QueryRow(ctx,
		`UPDATE food_aliases
		 SET food_id = $1, updated_at = NOW()
		 WHERE user_id = $2 AND id = $3
		 RETURNING updated_at`,
		alias.FoodID, alias.UserID, alias.ID,
	).Scan(&alias.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFoodAliasNotFound
	}

	return err
}
