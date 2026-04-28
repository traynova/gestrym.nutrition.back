package usecases

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

type GetFoodByIDUseCase struct {
	Repo interfaces.FoodRepository
}

func NewGetFoodByIDUseCase(repo interfaces.FoodRepository) *GetFoodByIDUseCase {
	return &GetFoodByIDUseCase{Repo: repo}
}

func (uc *GetFoodByIDUseCase) Execute(id uint) (*models.Food, error) {
	return uc.Repo.FindByID(id)
}
