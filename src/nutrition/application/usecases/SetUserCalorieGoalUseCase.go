package usecases

import (
	"gestrym-nutrition/src/common/models"
	nutritionUtils "gestrym-nutrition/src/nutrition/application/utils"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
)

// SetUserCalorieGoalUseCase creates or updates a user's personalized calorie & macro targets.
// It uses the TDEE engine to compute the targets from physical data.
type SetUserCalorieGoalUseCase struct {
	GoalRepo interfaces.UserCalorieGoalRepository
}

func NewSetUserCalorieGoalUseCase(goalRepo interfaces.UserCalorieGoalRepository) *SetUserCalorieGoalUseCase {
	return &SetUserCalorieGoalUseCase{GoalRepo: goalRepo}
}

type SetUserCalorieGoalInput struct {
	UserID        uint                    `json:"userId"`
	WeightKg      float64                 `json:"weightKg" binding:"required,gt=0"`
	HeightCm      float64                 `json:"heightCm" binding:"required,gt=0"`
	AgeYears      int                     `json:"ageYears" binding:"required,gt=0"`
	IsMale        bool                    `json:"isMale"`
	ActivityLevel models.ActivityLevel    `json:"activityLevel"`
	FitnessGoal   models.FitnessGoalType  `json:"fitnessGoal"`
}

type SetUserCalorieGoalResult struct {
	Goal *models.UserCalorieGoal    `json:"goal"`
	TDEE *nutritionUtils.TDEEResult `json:"calculation"`
}

func (uc *SetUserCalorieGoalUseCase) Execute(input SetUserCalorieGoalInput) (*SetUserCalorieGoalResult, error) {
	// Apply defaults if not provided
	if input.ActivityLevel == "" {
		input.ActivityLevel = models.ActivityModerate
	}
	if input.FitnessGoal == "" {
		input.FitnessGoal = models.FitnessGoalMaintain
	}

	// Calculate TDEE + macros
	tdee := nutritionUtils.CalculateTDEE(
		input.WeightKg, input.HeightCm,
		input.AgeYears, input.IsMale,
		input.ActivityLevel, input.FitnessGoal,
	)

	goal := &models.UserCalorieGoal{
		UserID:         input.UserID,
		WeightKg:       input.WeightKg,
		HeightCm:       input.HeightCm,
		AgeYears:       input.AgeYears,
		IsMale:         input.IsMale,
		ActivityLevel:  input.ActivityLevel,
		FitnessGoal:    input.FitnessGoal,
		TargetCalories: tdee.TargetCalories,
		TargetProtein:  tdee.TargetProtein,
		TargetCarbs:    tdee.TargetCarbs,
		TargetFats:     tdee.TargetFats,
		AdjustedByAI:   false,
	}

	if err := uc.GoalRepo.Upsert(goal); err != nil {
		return nil, err
	}

	return &SetUserCalorieGoalResult{Goal: goal, TDEE: &tdee}, nil
}
