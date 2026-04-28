package repositories

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"

	"gorm.io/gorm"
)

type MealItemRepositoryImpl struct {
	DB *gorm.DB
}

func NewMealItemRepositoryImpl(db *gorm.DB) interfaces.MealItemRepository {
	return &MealItemRepositoryImpl{DB: db}
}

func (r *MealItemRepositoryImpl) Create(item *models.MealItem) error {
	return r.DB.Create(item).Error
}

func (r *MealItemRepositoryImpl) FindByDayID(dayID uint) ([]models.MealItem, error) {
	var items []models.MealItem
	err := r.DB.
		Where("meal_day_id = ?", dayID).
		Preload("Food").
		Preload("Food.Category").
		Find(&items).Error
	return items, err
}
func (r *MealItemRepositoryImpl) DeleteByDayID(dayID uint) error {
	return r.DB.Where("meal_day_id = ?", dayID).Delete(&models.MealItem{}).Error
}
