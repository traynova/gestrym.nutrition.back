package interfaces

// ProgressMetrics represents the data received from progress-service
// for a specific user. Only the fields relevant to calorie adjustment are mapped.
type ProgressMetrics struct {
	UserID      uint    `json:"userId"`
	WeightKg    float64 `json:"weightKg"`    // latest body weight
	HeightCm    float64 `json:"heightCm"`    // latest height
	BodyFatPct  float64 `json:"bodyFatPct"`  // optional body fat %
	// Weekly change metrics (used for adaptive adjustments)
	WeightDeltaKg float64 `json:"weightDeltaKg"` // positive = gained, negative = lost
}

// ProgressServiceAdapter is the port for communicating with progress-service.
// Follows hexagonal architecture: the domain defines what it needs,
// the infrastructure layer implements the actual HTTP call.
type ProgressServiceAdapter interface {
	GetLatestMetrics(userID uint) (*ProgressMetrics, error)
}
