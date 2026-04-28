package usecases

import (
	nutritionUtils "gestrym-nutrition/src/nutrition/application/utils"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"time"
)

type GetDailyNutritionUseCase struct {
	LogRepo     interfaces.NutritionLogRepository
	PlanRepo    interfaces.MealPlanRepository
}

func NewGetDailyNutritionUseCase(logRepo interfaces.NutritionLogRepository, planRepo interfaces.MealPlanRepository) *GetDailyNutritionUseCase {
	return &GetDailyNutritionUseCase{LogRepo: logRepo, PlanRepo: planRepo}
}

// DailyGoals represents targets for a day (from active meal plan if available)
type DailyGoals struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}

// DailyProgress represents percentage of goals achieved
type DailyProgress struct {
	CaloriesPct float64 `json:"caloriesPct"`
	ProteinPct  float64 `json:"proteinPct"`
	CarbsPct    float64 `json:"carbsPct"`
	FatsPct     float64 `json:"fatsPct"`
}

// DailyNutritionResult is the frontend-friendly daily nutrition summary
type DailyNutritionResult struct {
	Date     string                         `json:"date"`
	Totals   nutritionUtils.NutritionTotals `json:"totals"`
	Goals    DailyGoals                     `json:"goals"`
	Progress DailyProgress                  `json:"progress"`
	Foods    interface{}                    `json:"foods"`
}

func (uc *GetDailyNutritionUseCase) Execute(userID uint, date time.Time) (*DailyNutritionResult, error) {
	logs, err := uc.LogRepo.FindByUserAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	totals := nutritionUtils.CalculateNutritionTotals(logs)

	// Try to get goals from the user's most recent meal plan
	var goals DailyGoals
	plans, _ := uc.PlanRepo.FindByUserID(userID)
	if len(plans) > 0 {
		// Use most recent plan's goals
		latest := plans[0]
		goals = DailyGoals{
			Calories: latest.GoalCalories,
			Protein:  latest.GoalProtein,
			Carbs:    latest.GoalCarbs,
			Fats:     latest.GoalFats,
		}
	}

	progress := DailyProgress{
		CaloriesPct: nutritionUtils.MacroProgress(totals.Calories, goals.Calories),
		ProteinPct:  nutritionUtils.MacroProgress(totals.Protein, goals.Protein),
		CarbsPct:    nutritionUtils.MacroProgress(totals.Carbs, goals.Carbs),
		FatsPct:     nutritionUtils.MacroProgress(totals.Fats, goals.Fats),
	}

	return &DailyNutritionResult{
		Date:     date.Format("2006-01-02"),
		Totals:   totals,
		Goals:    goals,
		Progress: progress,
		Foods:    logs,
	}, nil
}
