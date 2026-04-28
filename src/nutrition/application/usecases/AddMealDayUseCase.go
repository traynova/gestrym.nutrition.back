package usecases

import (
	"fmt"
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type AddMealDayUseCase struct {
	PlanRepo interfaces.MealPlanRepository
	DayRepo  interfaces.MealDayRepository
}

func NewAddMealDayUseCase(planRepo interfaces.MealPlanRepository, dayRepo interfaces.MealDayRepository) *AddMealDayUseCase {
	return &AddMealDayUseCase{PlanRepo: planRepo, DayRepo: dayRepo}
}

type AddMealDayInput struct {
	MealPlanID uint `json:"mealPlanId"`
	DayNumber  int  `json:"dayNumber"`
}

func (uc *AddMealDayUseCase) Execute(input AddMealDayInput) (*models.MealDay, error) {
	// Validate plan exists
	plan, err := uc.PlanRepo.FindByID(input.MealPlanID)
	if err != nil {
		return nil, fmt.Errorf("meal plan not found: %w", err)
	}

	// Validate day number is within plan duration
	if input.DayNumber < 1 || input.DayNumber > plan.DurationDays {
		return nil, fmt.Errorf("day number %d is out of range [1..%d]", input.DayNumber, plan.DurationDays)
	}

	day := &models.MealDay{
		MealPlanID: input.MealPlanID,
		DayNumber:  input.DayNumber,
	}
	if err := uc.DayRepo.Create(day); err != nil {
		return nil, err
	}
	return day, nil
}
