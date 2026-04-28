package models

import "time"

// NutritionLog tracks what the user ACTUALLY eats each day.
// Nutritional values are pre-calculated at insertion time (do NOT recalculate).
type NutritionLog struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `gorm:"not null;index" json:"userId"`
	Date     time.Time `gorm:"not null;index" json:"date"`
	FoodID   uint      `gorm:"not null" json:"foodId"`
	Food     Food      `gorm:"foreignKey:FoodID" json:"food,omitempty"`
	Quantity float64   `gorm:"not null" json:"quantity"` // in grams
	// Pre-calculated nutritional values (stored, NOT recalculated)
	Calories float64   `gorm:"not null" json:"calories"`
	Protein  float64   `gorm:"not null" json:"protein"`
	Carbs    float64   `gorm:"not null" json:"carbs"`
	Fats     float64   `gorm:"not null" json:"fats"`
	// Meal type for grouping
	MealType MealType  `gorm:"size:20" json:"mealType"`
	// AI readiness: notes for adaptive logic
	Notes    string    `gorm:"size:500" json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
