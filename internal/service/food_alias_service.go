package service

import (
	"context"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
)

type FoodAliasService struct {
	repo     *repository.FoodAliasRepository
	foodRepo *repository.FoodRepository
}

func NewFoodAliasService(repo *repository.FoodAliasRepository, foodRepo *repository.FoodRepository) *FoodAliasService {
	return &FoodAliasService{repo: repo, foodRepo: foodRepo}
}

func (s *FoodAliasService) Update(ctx context.Context, userID, id uint64, body dto.FoodAliasDTO) (*model.FoodAlias, error) {
	alias, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	food, err := s.foodRepo.FindByID(ctx, body.FoodID)
	if err != nil {
		return nil, err
	}

	alias.FoodID = &food.ID
	alias.Food = food

	if err := s.repo.Update(ctx, alias); err != nil {
		return nil, err
	}

	return alias, nil
}
