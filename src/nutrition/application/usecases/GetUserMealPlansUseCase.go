package usecases

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type GetUserMealPlansUseCase struct {
	Repo interfaces.MealPlanRepository
}

func NewGetUserMealPlansUseCase(repo interfaces.MealPlanRepository) *GetUserMealPlansUseCase {
	return &GetUserMealPlansUseCase{Repo: repo}
}

func (uc *GetUserMealPlansUseCase) Execute(userID uint) ([]models.MealPlan, error) {
	return uc.Repo.FindByUserID(userID)
}
