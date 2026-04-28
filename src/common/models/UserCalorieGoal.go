package models

import "time"

// ActivityLevel represents the user's physical activity level (for TDEE calculation)
type ActivityLevel string

const (
	ActivitySedentary  ActivityLevel = "sedentary"   // x1.2
	ActivityLight      ActivityLevel = "light"        // x1.375
	ActivityModerate   ActivityLevel = "moderate"     // x1.55
	ActivityActive     ActivityLevel = "active"       // x1.725
	ActivityVeryActive ActivityLevel = "very_active"  // x1.9
)

// FitnessGoalType represents what the user is trying to achieve
type FitnessGoalType string

const (
	FitnessGoalLoseWeight   FitnessGoalType = "lose_weight"   // deficit -500 kcal
	FitnessGoalMaintain     FitnessGoalType = "maintain"       // TDEE
	FitnessGoalGainMass     FitnessGoalType = "gain_mass"      // surplus +500 kcal
)

// UserCalorieGoal stores personalized caloric and macro objectives for a user.
// Replaces plan-level goals when the user wants fine-grained control.
// Designed to be auto-adjusted by the AI integration with progress-service.
type UserCalorieGoal struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	UserID        uint            `gorm:"not null;uniqueIndex" json:"userId"` // one goal per user
	// Physical data (sourced from progress-service or user input)
	WeightKg      float64         `gorm:"default:0" json:"weightKg"`
	HeightCm      float64         `gorm:"default:0" json:"heightCm"`
	AgeYears      int             `gorm:"default:0" json:"ageYears"`
	IsMale        bool            `gorm:"default:true" json:"isMale"`
	// Activity & goal classification
	ActivityLevel ActivityLevel   `gorm:"size:30;default:'moderate'" json:"activityLevel"`
	FitnessGoal   FitnessGoalType `gorm:"size:30;default:'maintain'" json:"fitnessGoal"`
	// Calculated targets (stored so the frontend always has a value)
	TargetCalories float64        `gorm:"default:0" json:"targetCalories"`
	TargetProtein  float64        `gorm:"default:0" json:"targetProtein"`
	TargetCarbs    float64        `gorm:"default:0" json:"targetCarbs"`
	TargetFats     float64        `gorm:"default:0" json:"targetFats"`
	// AI-adjustment metadata
	LastAdjustedAt *time.Time     `json:"lastAdjustedAt,omitempty"`
	AdjustedByAI   bool           `gorm:"default:false" json:"adjustedByAI"`
	AdjustmentNote string         `gorm:"size:500" json:"adjustmentNote,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
