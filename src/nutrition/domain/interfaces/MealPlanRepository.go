package interfaces

import "gestrym-nutrition/src/common/models"

type MealPlanRepository interface {
	Create(plan *models.MealPlan) error
	FindByID(id uint) (*models.MealPlan, error)
	FindByUserID(userID uint) ([]models.MealPlan, error)
	FindTemplates() ([]models.MealPlan, error)
	Update(plan *models.MealPlan) error
	Delete(id uint) error
}
