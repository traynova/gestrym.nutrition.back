package models

import "time"

type MealPlan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"userId"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	DurationDays int       `gorm:"not null" json:"durationDays"`
	CreatedBy    uint      `gorm:"not null" json:"createdBy"`
	IsTemplate   bool      `gorm:"default:false" json:"isTemplate"`
	// Calorie & Macro Goals (for AI integration readiness)
	GoalCalories float64   `gorm:"default:0" json:"goalCalories"`
	GoalProtein  float64   `gorm:"default:0" json:"goalProtein"`
	GoalCarbs    float64   `gorm:"default:0" json:"goalCarbs"`
	GoalFats     float64   `gorm:"default:0" json:"goalFats"`
	// Relations
	Days         []MealDay `gorm:"foreignKey:MealPlanID" json:"days,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
