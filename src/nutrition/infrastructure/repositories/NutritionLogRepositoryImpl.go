package repositories

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"time"

	"gorm.io/gorm"
)

type NutritionLogRepositoryImpl struct {
	DB *gorm.DB
}

func NewNutritionLogRepositoryImpl(db *gorm.DB) interfaces.NutritionLogRepository {
	return &NutritionLogRepositoryImpl{DB: db}
}

func (r *NutritionLogRepositoryImpl) Create(log *models.NutritionLog) error {
	return r.DB.Create(log).Error
}

func (r *NutritionLogRepositoryImpl) FindByUserAndDate(userID uint, date time.Time) ([]models.NutritionLog, error) {
	var logs []models.NutritionLog
	// Query by UTC date boundaries
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	err := r.DB.
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Preload("Food").
		Preload("Food.Category").
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}

func (r *NutritionLogRepositoryImpl) FindByUserRange(userID uint, start, end time.Time, page, pageSize int) ([]models.NutritionLog, int64, error) {
	var logs []models.NutritionLog
	var total int64

	query := r.DB.Model(&models.NutritionLog{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Food").
		Preload("Food.Category").
		Order("date DESC").
		Offset(offset).Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

func (r *NutritionLogRepositoryImpl) DeleteByID(id uint) error {
	return r.DB.Delete(&models.NutritionLog{}, id).Error
}
