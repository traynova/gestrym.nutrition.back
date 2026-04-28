package usecases

import "gestrym-nutrition/src/nutrition/domain/interfaces"

// GetUserCalorieGoalUseCase retrieves the current calorie goal for a user.
type GetUserCalorieGoalUseCase struct {
	GoalRepo interfaces.UserCalorieGoalRepository
}

func NewGetUserCalorieGoalUseCase(goalRepo interfaces.UserCalorieGoalRepository) *GetUserCalorieGoalUseCase {
	return &GetUserCalorieGoalUseCase{GoalRepo: goalRepo}
}

func (uc *GetUserCalorieGoalUseCase) Execute(userID uint) (interface{}, error) {
	goal, err := uc.GoalRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return goal, nil
}
