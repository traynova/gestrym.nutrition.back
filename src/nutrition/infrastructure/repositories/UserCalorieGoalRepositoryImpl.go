package repositories

import (
	"gestrym-nutrition/src/common/models"
	"gestrym-nutrition/src/nutrition/domain/interfaces"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserCalorieGoalRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserCalorieGoalRepositoryImpl(db *gorm.DB) interfaces.UserCalorieGoalRepository {
	return &UserCalorieGoalRepositoryImpl{DB: db}
}

// Upsert creates or updates the goal for a user (one row per user via uniqueIndex).
func (r *UserCalorieGoalRepositoryImpl) Upsert(goal *models.UserCalorieGoal) error {
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"weight_kg", "height_cm", "age_years", "is_male",
				"activity_level", "fitness_goal",
				"target_calories", "target_protein", "target_carbs", "target_fats",
				"last_adjusted_at", "adjusted_by_ai", "adjustment_note",
				"updated_at",
			}),
		}).
		Create(goal).Error
}

func (r *UserCalorieGoalRepositoryImpl) FindByUserID(userID uint) (*models.UserCalorieGoal, error) {
	var goal models.UserCalorieGoal
	err := r.DB.Where("user_id = ?", userID).First(&goal).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}
