package utils

import "gestrym-nutrition/src/common/models"

// NutritionTotals holds aggregated nutritional values for frontend dashboards
type NutritionTotals struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}

// CalculateNutritionTotals aggregates nutrition values from a slice of logs.
// Values are pre-calculated at insertion time — this is a pure aggregation.
func CalculateNutritionTotals(logs []models.NutritionLog) NutritionTotals {
	var totals NutritionTotals
	for _, log := range logs {
		totals.Calories += log.Calories
		totals.Protein += log.Protein
		totals.Carbs += log.Carbs
		totals.Fats += log.Fats
	}
	// Round to 2 decimal places
	totals.Calories = round2(totals.Calories)
	totals.Protein = round2(totals.Protein)
	totals.Carbs = round2(totals.Carbs)
	totals.Fats = round2(totals.Fats)
	return totals
}

// CalculateFromFood computes nutrition values for a given quantity (in grams)
// based on food nutritional data per 100g.
func CalculateFromFood(food *models.Food, quantityGrams float64) (calories, protein, carbs, fats float64) {
	factor := quantityGrams / 100.0
	calories = round2(food.Calories * factor)
	protein = round2(food.Protein * factor)
	carbs = round2(food.Carbs * factor)
	fats = round2(food.Fats * factor)
	return
}

// MacroProgress calculates percentage achieved vs goals (0 if goal is 0)
func MacroProgress(value, goal float64) float64 {
	if goal == 0 {
		return 0
	}
	pct := (value / goal) * 100
	return round2(pct)
}

// RoundTwo rounds a float64 to 2 decimal places.
func RoundTwo(val float64) float64 {
	return float64(int(val*100+0.5)) / 100
}

func round2(val float64) float64 {
	return RoundTwo(val)
}
