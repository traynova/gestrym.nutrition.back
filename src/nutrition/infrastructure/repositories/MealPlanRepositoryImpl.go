package repositories

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"

	"gorm.io/gorm"
)

type MealPlanRepositoryImpl struct {
	DB *gorm.DB
}

func NewMealPlanRepositoryImpl(db *gorm.DB) interfaces.MealPlanRepository {
	return &MealPlanRepositoryImpl{DB: db}
}

func (r *MealPlanRepositoryImpl) Create(plan *models.MealPlan) error {
	return r.DB.Create(plan).Error
}

func (r *MealPlanRepositoryImpl) FindByID(id uint) (*models.MealPlan, error) {
	var plan models.MealPlan
	err := r.DB.
		Preload("Days").
		Preload("Days.Items").
		Preload("Days.Items.Food").
		Preload("Days.Items.Food.Category").
		First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *MealPlanRepositoryImpl) FindByUserID(userID uint) ([]models.MealPlan, error) {
	var plans []models.MealPlan
	err := r.DB.
		Where("user_id = ?", userID).
		Preload("Days").
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

func (r *MealPlanRepositoryImpl) FindTemplates() ([]models.MealPlan, error) {
	var plans []models.MealPlan
	err := r.DB.
		Where("is_template = true").
		Preload("Days").
		Order("created_at DESC").
		Find(&plans).Error
	return plans, err
}

func (r *MealPlanRepositoryImpl) Update(plan *models.MealPlan) error {
	return r.DB.Save(plan).Error
}

func (r *MealPlanRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&models.MealPlan{}, id).Error
}
