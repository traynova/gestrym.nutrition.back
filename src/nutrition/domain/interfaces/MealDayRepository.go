package interfaces

import "gestrym-nutrition/src/common/models"

type MealDayRepository interface {
	Create(day *models.MealDay) error
	FindByPlanID(planID uint) ([]models.MealDay, error)
	FindByID(id uint) (*models.MealDay, error)
	DeleteByPlanID(planID uint) error
}
