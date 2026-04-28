package utils

import "gestrym-nutrition/src/common/models"

// TDEEResult holds the calculated energy expenditure and macro breakdown
type TDEEResult struct {
	BMR            float64 `json:"bmr"`            // Basal Metabolic Rate (kcal)
	TDEE           float64 `json:"tdee"`           // Total Daily Energy Expenditure
	TargetCalories float64 `json:"targetCalories"` // After fitness goal adjustment
	TargetProtein  float64 `json:"targetProtein"`  // grams
	TargetCarbs    float64 `json:"targetCarbs"`    // grams
	TargetFats     float64 `json:"targetFats"`     // grams
}

// activityMultiplier maps activity level to PAL (Physical Activity Level) factor
var activityMultiplier = map[models.ActivityLevel]float64{
	models.ActivitySedentary:  1.2,
	models.ActivityLight:      1.375,
	models.ActivityModerate:   1.55,
	models.ActivityActive:     1.725,
	models.ActivityVeryActive: 1.9,
}

// goalAdjustment maps fitness goal to calorie delta from TDEE
var goalAdjustment = map[models.FitnessGoalType]float64{
	models.FitnessGoalLoseWeight: -500, // 0.5 kg/week deficit
	models.FitnessGoalMaintain:   0,
	models.FitnessGoalGainMass:   +500, // lean bulk surplus
}

// CalculateTDEE computes BMR (Mifflin-St Jeor), TDEE, and macro targets.
//
// Macro split by goal:
//   - Lose weight:  protein=35%, carbs=40%, fats=25%
//   - Maintain:     protein=25%, carbs=50%, fats=25%
//   - Gain mass:    protein=30%, carbs=45%, fats=25%
func CalculateTDEE(
	weightKg, heightCm float64,
	ageYears int,
	isMale bool,
	activity models.ActivityLevel,
	goal models.FitnessGoalType,
) TDEEResult {
	// Mifflin-St Jeor BMR
	var bmr float64
	if isMale {
		bmr = (10 * weightKg) + (6.25 * heightCm) - (5 * float64(ageYears)) + 5
	} else {
		bmr = (10 * weightKg) + (6.25 * heightCm) - (5 * float64(ageYears)) - 161
	}

	// Activity factor
	pal, ok := activityMultiplier[activity]
	if !ok {
		pal = activityMultiplier[models.ActivityModerate]
	}
	tdee := bmr * pal

	// Fitness goal calorie adjustment
	adj, ok := goalAdjustment[goal]
	if !ok {
		adj = 0
	}
	targetCals := tdee + adj

	// Macro distribution (never below 1200 kcal)
	if targetCals < 1200 {
		targetCals = 1200
	}

	var proteinPct, carbsPct, fatsPct float64
	switch goal {
	case models.FitnessGoalLoseWeight:
		proteinPct, carbsPct, fatsPct = 0.35, 0.40, 0.25
	case models.FitnessGoalGainMass:
		proteinPct, carbsPct, fatsPct = 0.30, 0.45, 0.25
	default: // maintain
		proteinPct, carbsPct, fatsPct = 0.25, 0.50, 0.25
	}

	// Convert kcal % → grams (protein=4kcal/g, carbs=4kcal/g, fats=9kcal/g)
	protein := round2((targetCals * proteinPct) / 4)
	carbs := round2((targetCals * carbsPct) / 4)
	fats := round2((targetCals * fatsPct) / 9)

	return TDEEResult{
		BMR:            round2(bmr),
		TDEE:           round2(tdee),
		TargetCalories: round2(targetCals),
		TargetProtein:  protein,
		TargetCarbs:    carbs,
		TargetFats:     fats,
	}
}

// AdaptCaloriesFromProgress adjusts calorie target based on real weight progress.
// This is the core of the AI-adaptive loop:
//
//   - Goal = lose weight AND not losing fast enough (delta > -0.3kg/week) → reduce 100 kcal
//   - Goal = gain mass AND not gaining (delta < +0.1kg/week) → increase 100 kcal
//   - Losing too fast (delta < -1kg/week) → increase 200 kcal to protect muscle
//
// Returns the adjusted calorie target and a human-readable note.
func AdaptCaloriesFromProgress(
	currentTarget float64,
	goal models.FitnessGoalType,
	weeklyWeightDeltaKg float64,
) (newTarget float64, note string) {
	newTarget = currentTarget

	switch goal {
	case models.FitnessGoalLoseWeight:
		if weeklyWeightDeltaKg > -0.3 {
			// Not losing fast enough
			newTarget = round2(currentTarget - 100)
			note = "Adjusted: weight loss slower than expected — reduced by 100 kcal"
		} else if weeklyWeightDeltaKg < -1.0 {
			// Losing too fast — protect muscle mass
			newTarget = round2(currentTarget + 200)
			note = "Adjusted: weight loss too aggressive — increased by 200 kcal to protect muscle"
		} else {
			note = "No adjustment needed: weight loss on track"
		}

	case models.FitnessGoalGainMass:
		if weeklyWeightDeltaKg < 0.1 {
			// Not gaining — increase surplus
			newTarget = round2(currentTarget + 100)
			note = "Adjusted: insufficient mass gain — increased by 100 kcal"
		} else if weeklyWeightDeltaKg > 0.5 {
			// Gaining too fast — likely adding fat
			newTarget = round2(currentTarget - 100)
			note = "Adjusted: gaining too quickly — reduced by 100 kcal to lean bulk"
		} else {
			note = "No adjustment needed: lean mass gain on track"
		}

	default: // maintain
		if weeklyWeightDeltaKg > 0.3 {
			newTarget = round2(currentTarget - 100)
			note = "Adjusted: unintended weight gain detected — reduced by 100 kcal"
		} else if weeklyWeightDeltaKg < -0.3 {
			newTarget = round2(currentTarget + 100)
			note = "Adjusted: unintended weight loss detected — increased by 100 kcal"
		} else {
			note = "No adjustment needed: weight stable"
		}
	}

	// Safety floor
	if newTarget < 1200 {
		newTarget = 1200
		note += " (minimum 1200 kcal safety floor applied)"
	}

	return
}
