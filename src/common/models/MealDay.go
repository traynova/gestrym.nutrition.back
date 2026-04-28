package models

type MealDay struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	MealPlanID uint       `gorm:"not null;index" json:"mealPlanId"`
	DayNumber  int        `gorm:"not null" json:"dayNumber"`
	// Relations
	Items      []MealItem `gorm:"foreignKey:MealDayID" json:"items,omitempty"`
}
