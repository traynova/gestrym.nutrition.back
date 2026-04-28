package models

// MealType represents the time of day for a meal
type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeSnack     MealType = "snack"
)

// ValidMealTypes contains all allowed meal types for validation
var ValidMealTypes = map[MealType]bool{
	MealTypeBreakfast: true,
	MealTypeLunch:     true,
	MealTypeDinner:    true,
	MealTypeSnack:     true,
}

type MealItem struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	MealDayID uint     `gorm:"not null;index" json:"mealDayId"`
	FoodID    uint     `gorm:"not null" json:"foodId"`
	Food      Food     `gorm:"foreignKey:FoodID" json:"food,omitempty"`
	Quantity  float64  `gorm:"not null" json:"quantity"` // in grams
	MealType  MealType `gorm:"size:20;not null" json:"mealType"`
}
