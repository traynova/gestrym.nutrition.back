package repositories

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"

	"gorm.io/gorm"
)

type MealDayRepositoryImpl struct {
	DB *gorm.DB
}

func NewMealDayRepositoryImpl(db *gorm.DB) interfaces.MealDayRepository {
	return &MealDayRepositoryImpl{DB: db}
}

func (r *MealDayRepositoryImpl) Create(day *models.MealDay) error {
	return r.DB.Create(day).Error
}

func (r *MealDayRepositoryImpl) FindByPlanID(planID uint) ([]models.MealDay, error) {
	var days []models.MealDay
	err := r.DB.
		Where("meal_plan_id = ?", planID).
		Preload("Items").
		Preload("Items.Food").
		Order("day_number ASC").
		Find(&days).Error
	return days, err
}

func (r *MealDayRepositoryImpl) FindByID(id uint) (*models.MealDay, error) {
	var day models.MealDay
	err := r.DB.
		Preload("Items").
		Preload("Items.Food").
		First(&day, id).Error
	if err != nil {
		return nil, err
	}
	return &day, nil
}
func (r *MealDayRepositoryImpl) DeleteByPlanID(planID uint) error {
	// First delete all items for these days (handled by DB if cascade is set, but better safe)
	// Actually, let's just delete the days and let cascade work if configured, 
	// or delete manually if not.
	return r.DB.Where("meal_plan_id = ?", planID).Delete(&models.MealDay{}).Error
}
