package usecases

import (
	"fmt"
	"gestrym-nutrition/src/common/models"
	nutritionUtils "gestrym-nutrition/src/nutrition/application/utils"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"time"
)

type LogFoodIntakeUseCase struct {
	LogRepo  interfaces.NutritionLogRepository
	FoodRepo interfaces.FoodRepository
}

func NewLogFoodIntakeUseCase(logRepo interfaces.NutritionLogRepository, foodRepo interfaces.FoodRepository) *LogFoodIntakeUseCase {
	return &LogFoodIntakeUseCase{LogRepo: logRepo, FoodRepo: foodRepo}
}

type LogFoodIntakeInput struct {
	UserID   uint            `json:"userId"`
	FoodID   uint            `json:"foodId"`
	Quantity float64         `json:"quantity"`
	MealType models.MealType `json:"mealType"`
	Date     time.Time       `json:"date"`
	Notes    string          `json:"notes"`
}

func (uc *LogFoodIntakeUseCase) Execute(input LogFoodIntakeInput) (*models.NutritionLog, error) {
	if input.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	// Validate meal type if provided
	if input.MealType != "" && !models.ValidMealTypes[input.MealType] {
		return nil, fmt.Errorf("invalid meal type '%s'", input.MealType)
	}

	// Fetch food nutritional data
	food, err := uc.FoodRepo.FindByID(input.FoodID)
	if err != nil {
		return nil, fmt.Errorf("food not found: %w", err)
	}

	// Calculate nutrition values (per 100g)
	calories, protein, carbs, fats := nutritionUtils.CalculateFromFood(food, input.Quantity)

	// Use today if date not provided
	date := input.Date
	if date.IsZero() {
		date = time.Now().UTC()
	}

	log := &models.NutritionLog{
		UserID:   input.UserID,
		Date:     date,
		FoodID:   input.FoodID,
		Quantity: input.Quantity,
		Calories: calories,
		Protein:  protein,
		Carbs:    carbs,
		Fats:     fats,
		MealType: input.MealType,
		Notes:    input.Notes,
	}

	if err := uc.LogRepo.Create(log); err != nil {
		return nil, err
	}

	// Load food relation for response
	log.Food = *food
	return log, nil
}
