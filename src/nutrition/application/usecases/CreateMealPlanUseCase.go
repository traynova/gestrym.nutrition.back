package usecases

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type CreateMealPlanUseCase struct {
	Repo interfaces.MealPlanRepository
}

func NewCreateMealPlanUseCase(repo interfaces.MealPlanRepository) *CreateMealPlanUseCase {
	return &CreateMealPlanUseCase{Repo: repo}
}

type CreateMealPlanInput struct {
	UserID       uint    `json:"userId"`
	Name         string  `json:"name"`
	DurationDays int     `json:"durationDays"`
	CreatedBy    uint    `json:"createdBy"`
	IsTemplate   bool    `json:"isTemplate"`
	GoalCalories float64 `json:"goalCalories"`
	GoalProtein  float64 `json:"goalProtein"`
	GoalCarbs    float64 `json:"goalCarbs"`
	GoalFats     float64 `json:"goalFats"`
}

func (uc *CreateMealPlanUseCase) Execute(input CreateMealPlanInput) (*models.MealPlan, error) {
	plan := &models.MealPlan{
		UserID:       input.UserID,
		Name:         input.Name,
		DurationDays: input.DurationDays,
		CreatedBy:    input.CreatedBy,
		IsTemplate:   input.IsTemplate,
		GoalCalories: input.GoalCalories,
		GoalProtein:  input.GoalProtein,
		GoalCarbs:    input.GoalCarbs,
		GoalFats:     input.GoalFats,
	}
	if err := uc.Repo.Create(plan); err != nil {
		return nil, err
	}
	return plan, nil
}
