package interfaces

import "gestrym-nutrition/src/common/models"

type UserCalorieGoalRepository interface {
	Upsert(goal *models.UserCalorieGoal) error
	FindByUserID(userID uint) (*models.UserCalorieGoal, error)
}
