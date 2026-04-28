package interfaces

import "gestrym-nutrition/src/common/models"

type MealItemRepository interface {
	Create(item *models.MealItem) error
	FindByDayID(dayID uint) ([]models.MealItem, error)
	DeleteByDayID(dayID uint) error
}
