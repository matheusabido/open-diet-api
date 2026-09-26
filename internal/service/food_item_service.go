package service

import (
	"context"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/unit"
)

type IncompatibleUnitError struct {
	Allowed []unit.Unit
}

func (e *IncompatibleUnitError) Error() string {
	return "unit is not compatible with the food unit"
}

type FoodItemService struct {
	repo     *repository.FoodItemRepository
	foodRepo *repository.FoodRepository
}

func NewFoodItemService(repo *repository.FoodItemRepository, foodRepo *repository.FoodRepository) *FoodItemService {
	return &FoodItemService{repo: repo, foodRepo: foodRepo}
}

func (s *FoodItemService) Create(ctx context.Context, userID uint64, body dto.FoodItemDTO) (*model.FoodItem, error) {
	item := &model.FoodItem{UserID: userID}
	if err := s.apply(ctx, item, body); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *FoodItemService) FindByID(ctx context.Context, userID, id uint64) (*model.FoodItem, error) {
	return s.repo.FindByID(ctx, userID, id)
}

func (s *FoodItemService) List(ctx context.Context, userID uint64, query dto.FoodItemListQueryDTO) ([]*model.FoodItem, error) {
	return s.repo.ListByDate(ctx, userID, query.Date)
}

func (s *FoodItemService) Update(ctx context.Context, userID, id uint64, body dto.FoodItemDTO) (*model.FoodItem, error) {
	item, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if err := s.apply(ctx, item, body); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *FoodItemService) Delete(ctx context.Context, userID, id uint64) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *FoodItemService) apply(ctx context.Context, item *model.FoodItem, body dto.FoodItemDTO) error {
	food, err := s.foodRepo.FindByID(ctx, body.FoodID)
	if err != nil {
		return err
	}

	itemUnit, _ := unit.Parse(body.Unit)
	if !itemUnit.CompatibleWith(food.Unit) {
		return &IncompatibleUnitError{Allowed: food.Unit.Compatible()}
	}

	item.FoodID = food.ID
	item.Food = food
	item.Date = body.Date
	item.Meal = model.Meal(body.Meal)
	item.Amount = body.Amount
	item.Unit = itemUnit

	return nil
}
