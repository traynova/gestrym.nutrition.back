package usecases

import (
	"errors"
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type GetMealPlanUseCase struct {
	Repo interfaces.MealPlanRepository
}

func NewGetMealPlanUseCase(repo interfaces.MealPlanRepository) *GetMealPlanUseCase {
	return &GetMealPlanUseCase{Repo: repo}
}

func (uc *GetMealPlanUseCase) Execute(planID uint, requesterID uint, requesterRoleID uint) (*models.MealPlan, error) {
	plan, err := uc.Repo.FindByID(planID)
	if err != nil {
		return nil, err
	}

	// Authorization: user can see own plan, coach/gym/admin can see any
	const roleCoach = 3
	const roleGym = 2
	const roleAdmin = 1
	if requesterRoleID == roleCoach || requesterRoleID == roleGym || requesterRoleID == roleAdmin {
		return plan, nil
	}
	if plan.UserID != requesterID {
		return nil, errors.New("access denied: plan belongs to another user")
	}
	return plan, nil
}
