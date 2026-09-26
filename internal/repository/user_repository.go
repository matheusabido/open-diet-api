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
	ErrUserNotFound  = errors.New("user not found")
	ErrUsernameTaken = errors.New("username already taken")
)

const uniqueViolation = "23505"

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (username, password, kcal_goal, protein_goal, carbohydrate_goal, fat_goal)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		user.Username, user.Password, user.KcalGoal, user.ProteinGoal, user.CarbohydrateGoal, user.FatGoal,
	).Scan(&user.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return ErrUsernameTaken
	}

	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, password, kcal_goal, protein_goal, carbohydrate_goal, fat_goal
		 FROM users
		 WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.KcalGoal, &user.ProteinGoal, &user.CarbohydrateGoal, &user.FatGoal)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, username, password, kcal_goal, protein_goal, carbohydrate_goal, fat_goal
		 FROM users
		 WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.KcalGoal, &user.ProteinGoal, &user.CarbohydrateGoal, &user.FatGoal)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
