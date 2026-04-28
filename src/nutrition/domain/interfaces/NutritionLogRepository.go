package interfaces

import (
	"gestrym-nutrition/src/common/models"
	"time"
)

type NutritionLogRepository interface {
	Create(log *models.NutritionLog) error
	FindByUserAndDate(userID uint, date time.Time) ([]models.NutritionLog, error)
	FindByUserRange(userID uint, start, end time.Time, page, pageSize int) ([]models.NutritionLog, int64, error)
	DeleteByID(id uint) error
}
