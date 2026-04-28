package usecases

import (
	"fmt"
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type CreateMealPlanFromAIUseCase struct {
	MealPlanRepo interfaces.MealPlanRepository
	MealDayRepo  interfaces.MealDayRepository
	MealItemRepo interfaces.MealItemRepository
	FoodRepo     interfaces.FoodRepository
}

func NewCreateMealPlanFromAIUseCase(
	mealPlanRepo interfaces.MealPlanRepository,
	mealDayRepo interfaces.MealDayRepository,
	mealItemRepo interfaces.MealItemRepository,
	foodRepo interfaces.FoodRepository,
) *CreateMealPlanFromAIUseCase {
	return &CreateMealPlanFromAIUseCase{
		MealPlanRepo: mealPlanRepo,
		MealDayRepo:  mealDayRepo,
		MealItemRepo: mealItemRepo,
		FoodRepo:     foodRepo,
	}
}

type AIItemInput struct {
	FoodID   uint    `json:"foodId" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

type AIMealInput struct {
	MealType models.MealType `json:"mealType" binding:"required"`
	Items    []AIItemInput   `json:"items" binding:"required,dive"`
}

type AIDayInput struct {
	DayNumber int           `json:"dayNumber" binding:"required,gt=0"`
	Meals     []AIMealInput `json:"meals" binding:"required,dive"`
}

type CreateMealPlanFromAIInput struct {
	PlanID       uint         `json:"planId"` // Optional: for adaptation/update
	UserID       uint         `json:"userId" binding:"required"`
	Name         string       `json:"name" binding:"required"`
	DurationDays int          `json:"durationDays" binding:"required,gt=0"`
	Days         []AIDayInput `json:"days" binding:"required,dive"`
}

func (uc *CreateMealPlanFromAIUseCase) Execute(input CreateMealPlanFromAIInput) (*models.MealPlan, error) {
	var plan *models.MealPlan
	var err error

	if input.PlanID != 0 {
		plan, err = uc.MealPlanRepo.FindByID(input.PlanID)
		if err != nil {
			return nil, fmt.Errorf("meal plan with ID %d not found: %w", input.PlanID, err)
		}
		// Update basic fields
		plan.Name = input.Name
		plan.DurationDays = input.DurationDays
	} else {
		plan = &models.MealPlan{
			UserID:        input.UserID,
			Name:          input.Name,
			DurationDays:  input.DurationDays,
			IsAIGenerated: true,
			CreatedBy:     0, // System/AI
		}
	}

	var totalCalories, totalProtein, totalCarbs, totalFats float64

	// Pre-validate all foods and calculate totals
	for _, dayInput := range input.Days {
		for _, mealInput := range dayInput.Meals {
			for _, itemInput := range mealInput.Items {
				food, err := uc.FoodRepo.FindByID(itemInput.FoodID)
				if err != nil {
					return nil, fmt.Errorf("food with ID %d not found: %w", itemInput.FoodID, err)
				}
				
				// Calculate macros for this item based on quantity (per 100g usually)
				ratio := itemInput.Quantity / 100.0
				totalCalories += food.Calories * ratio
				totalProtein += food.Protein * ratio
				totalCarbs += food.Carbs * ratio
				totalFats += food.Fats * ratio
			}
		}
	}

	// Assign calculated goals
	plan.GoalCalories = totalCalories / float64(input.DurationDays)
	plan.GoalProtein = totalProtein / float64(input.DurationDays)
	plan.GoalCarbs = totalCarbs / float64(input.DurationDays)
	plan.GoalFats = totalFats / float64(input.DurationDays)

	if input.PlanID != 0 {
		// Adaptation: Clear old days first
		if err := uc.MealDayRepo.DeleteByPlanID(plan.ID); err != nil {
			return nil, fmt.Errorf("failed to clear old days: %w", err)
		}
		if err := uc.MealPlanRepo.Update(plan); err != nil {
			return nil, err
		}
	} else {
		if err := uc.MealPlanRepo.Create(plan); err != nil {
			return nil, err
		}
	}

	// 2. Create MealDays and MealItems
	for _, dayInput := range input.Days {
		day := &models.MealDay{
			MealPlanID: plan.ID,
			DayNumber:  dayInput.DayNumber,
		}
		if err := uc.MealDayRepo.Create(day); err != nil {
			return nil, err
		}

		for _, mealInput := range dayInput.Meals {
			for _, itemInput := range mealInput.Items {
				item := &models.MealItem{
					MealDayID: day.ID,
					FoodID:    itemInput.FoodID,
					Quantity:  itemInput.Quantity,
					MealType:  mealInput.MealType,
				}
				if err := uc.MealItemRepo.Create(item); err != nil {
					return nil, err
				}
			}
		}
	}

	return plan, nil
}
