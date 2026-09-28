package service

import (
	"context"
	"strings"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/unit"
)

type FoodService struct {
	repo *repository.FoodRepository
}

func NewFoodService(repo *repository.FoodRepository) *FoodService {
	return &FoodService{repo: repo}
}

func (s *FoodService) Create(ctx context.Context, userID uint64, body dto.FoodDTO) (*model.Food, error) {
	food := &model.Food{CreatedBy: &userID}
	applyFood(food, body)

	if err := s.repo.Create(ctx, food); err != nil {
		return nil, err
	}

	return food, nil
}

func (s *FoodService) FindByID(ctx context.Context, id uint64) (*model.Food, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FoodService) List(ctx context.Context, query dto.FoodListQueryDTO) ([]*model.Food, error) {
	return s.repo.List(ctx, strings.TrimSpace(query.Search), query.After)
}

func (s *FoodService) Update(ctx context.Context, id uint64, body dto.FoodDTO) (*model.Food, error) {
	food, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	applyFood(food, body)

	if err := s.repo.Update(ctx, food); err != nil {
		return nil, err
	}

	return food, nil
}

func (s *FoodService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func applyFood(food *model.Food, body dto.FoodDTO) {
	food.Name = strings.TrimSpace(body.Name)
	food.Brand = nil
	if body.Brand != nil {
		if brand := strings.TrimSpace(*body.Brand); brand != "" {
			food.Brand = &brand
		}
	}
	food.Kcal = *body.Kcal
	food.Protein = *body.Protein
	food.Carbohydrate = *body.Carbohydrate
	food.Fat = *body.Fat
	food.Amount = *body.Amount
	food.Unit, _ = unit.Parse(body.Unit)
}
