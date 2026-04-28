package usecases

import (
	"gestrym-nutrition/src/common/models"
	nutritionUtils "gestrym-nutrition/src/nutrition/application/utils"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"time"
)

type GetNutritionHistoryUseCase struct {
	LogRepo interfaces.NutritionLogRepository
}

func NewGetNutritionHistoryUseCase(logRepo interfaces.NutritionLogRepository) *GetNutritionHistoryUseCase {
	return &GetNutritionHistoryUseCase{LogRepo: logRepo}
}

type NutritionHistoryResult struct {
	Logs     []models.NutritionLog          `json:"logs"`
	Totals   nutritionUtils.NutritionTotals `json:"totals"`
	Total    int64                          `json:"total"`
	Page     int                            `json:"page"`
	PageSize int                            `json:"pageSize"`
}

func (uc *GetNutritionHistoryUseCase) Execute(userID uint, start, end time.Time, page, pageSize int) (*NutritionHistoryResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := uc.LogRepo.FindByUserRange(userID, start, end, page, pageSize)
	if err != nil {
		return nil, err
	}

	totals := nutritionUtils.CalculateNutritionTotals(logs)

	return &NutritionHistoryResult{
		Logs:     logs,
		Totals:   totals,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
