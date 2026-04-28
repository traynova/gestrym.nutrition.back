package usecases

import (
	"fmt"
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type AddMealItemUseCase struct {
	DayRepo  interfaces.MealDayRepository
	ItemRepo interfaces.MealItemRepository
	FoodRepo interfaces.FoodRepository
}

func NewAddMealItemUseCase(
	dayRepo interfaces.MealDayRepository,
	itemRepo interfaces.MealItemRepository,
	foodRepo interfaces.FoodRepository,
) *AddMealItemUseCase {
	return &AddMealItemUseCase{DayRepo: dayRepo, ItemRepo: itemRepo, FoodRepo: foodRepo}
}

type AddMealItemInput struct {
	MealDayID uint             `json:"mealDayId"`
	FoodID    uint             `json:"foodId"`
	Quantity  float64          `json:"quantity"`
	MealType  models.MealType  `json:"mealType"`
}

func (uc *AddMealItemUseCase) Execute(input AddMealItemInput) (*models.MealItem, error) {
	// Validate meal type
	if !models.ValidMealTypes[input.MealType] {
		return nil, fmt.Errorf("invalid meal type '%s': must be breakfast, lunch, dinner, or snack", input.MealType)
	}

	// Validate day exists
	if _, err := uc.DayRepo.FindByID(input.MealDayID); err != nil {
		return nil, fmt.Errorf("meal day not found: %w", err)
	}

	// Validate food exists
	if _, err := uc.FoodRepo.FindByID(input.FoodID); err != nil {
		return nil, fmt.Errorf("food not found: %w", err)
	}

	if input.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	item := &models.MealItem{
		MealDayID: input.MealDayID,
		FoodID:    input.FoodID,
		Quantity:  input.Quantity,
		MealType:  input.MealType,
	}
	if err := uc.ItemRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}
