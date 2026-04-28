package usecases

import (
	"fmt"
	nutritionUtils "gestrym-nutrition/src/nutrition/application/utils"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"time"
)

// AdjustCaloriesFromProgressUseCase integrates with progress-service to
// automatically recalibrate the user's calorie target based on real weight data.
//
// Flow:
//  1. Fetch user's current calorie goal
//  2. Fetch latest metrics from progress-service (weight, height, delta)
//  3. Run adaptive logic (AdaptCaloriesFromProgress)
//  4. Recompute macros proportionally to new calorie target
//  5. Persist updated goal with AI flag and note
type AdjustCaloriesFromProgressUseCase struct {
	GoalRepo        interfaces.UserCalorieGoalRepository
	ProgressAdapter interfaces.ProgressServiceAdapter
}

func NewAdjustCaloriesFromProgressUseCase(
	goalRepo interfaces.UserCalorieGoalRepository,
	progressAdapter interfaces.ProgressServiceAdapter,
) *AdjustCaloriesFromProgressUseCase {
	return &AdjustCaloriesFromProgressUseCase{
		GoalRepo:        goalRepo,
		ProgressAdapter: progressAdapter,
	}
}

type AdjustCaloriesResult struct {
	PreviousCalories float64 `json:"previousCalories"`
	NewCalories      float64 `json:"newCalories"`
	Delta            float64 `json:"delta"`
	Note             string  `json:"note"`
	WeightDeltaKg    float64 `json:"weightDeltaKg"`
	AdjustedByAI     bool    `json:"adjustedByAI"`
}

func (uc *AdjustCaloriesFromProgressUseCase) Execute(userID uint) (*AdjustCaloriesResult, error) {
	// 1. Load existing goal
	goal, err := uc.GoalRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("no calorie goal set for user %d: please set one first", userID)
	}

	// 2. Fetch progress metrics
	metrics, err := uc.ProgressAdapter.GetLatestMetrics(userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch progress data: %w", err)
	}

	// 3. Update physical data from progress-service if available
	if metrics.WeightKg > 0 {
		goal.WeightKg = metrics.WeightKg
	}
	if metrics.HeightCm > 0 {
		goal.HeightCm = metrics.HeightCm
	}

	previousCalories := goal.TargetCalories

	// 4. Run adaptive adjustment
	newCalories, note := nutritionUtils.AdaptCaloriesFromProgress(
		goal.TargetCalories,
		goal.FitnessGoal,
		metrics.WeightDeltaKg,
	)

	// 5. Recompute macros if calories changed
	if newCalories != previousCalories {
		// Recompute macros using full TDEE recalculation with new weight
		tdee := nutritionUtils.CalculateTDEE(
			goal.WeightKg, goal.HeightCm,
			goal.AgeYears, goal.IsMale,
			goal.ActivityLevel, goal.FitnessGoal,
		)
		// Apply the adaptive override to the TDEE base
		delta := newCalories - tdee.TDEE
		_ = delta // delta is already encoded in AdaptCaloriesFromProgress result

		goal.TargetCalories = newCalories
		// Scale macros proportionally to the new calorie target
		ratio := newCalories / previousCalories
		goal.TargetProtein = nutritionUtils.RoundTwo(goal.TargetProtein * ratio)
		goal.TargetCarbs = nutritionUtils.RoundTwo(goal.TargetCarbs * ratio)
		goal.TargetFats = nutritionUtils.RoundTwo(goal.TargetFats * ratio)
	}

	// 6. Stamp AI adjustment metadata
	now := time.Now().UTC()
	goal.LastAdjustedAt = &now
	goal.AdjustedByAI = true
	goal.AdjustmentNote = note

	if err := uc.GoalRepo.Upsert(goal); err != nil {
		return nil, fmt.Errorf("failed to save adjusted goal: %w", err)
	}

	return &AdjustCaloriesResult{
		PreviousCalories: previousCalories,
		NewCalories:      newCalories,
		Delta:            nutritionUtils.RoundTwo(newCalories - previousCalories),
		Note:             note,
		WeightDeltaKg:    metrics.WeightDeltaKg,
		AdjustedByAI:     true,
	}, nil
}
