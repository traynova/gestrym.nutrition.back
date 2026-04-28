package adapters

import (
	"encoding/json"
	"fmt"
	"gestrym-nutrition/src/nutrition/domain/interfaces"
	"net/http"
	"time"
)

// progressServiceResponse mirrors the progress-service API response shape.
// Only fields needed for calorie adjustment are mapped.
type progressServiceResponse struct {
	Data struct {
		UserID     uint    `json:"userId"`
		WeightKg   float64 `json:"weightKg"`
		HeightCm   float64 `json:"heightCm"`
		BodyFatPct float64 `json:"bodyFatPct"`
	} `json:"data"`
}

type progressWeightHistoryResponse struct {
	Data []struct {
		WeightKg  float64 `json:"weightKg"`
		CreatedAt string  `json:"createdAt"`
	} `json:"data"`
}

type ProgressServiceAdapterImpl struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
}

func NewProgressServiceAdapterImpl(baseURL, apiKey string) interfaces.ProgressServiceAdapter {
	return &ProgressServiceAdapterImpl{
		BaseURL: baseURL,
		APIKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetLatestMetrics fetches the most recent body metrics for a user from progress-service.
// It calls: GET {progress-service}/private/metrics/user/{userID}/latest
// and optionally: GET {progress-service}/private/metrics/user/{userID}/weight-history
func (a *ProgressServiceAdapterImpl) GetLatestMetrics(userID uint) (*interfaces.ProgressMetrics, error) {
	// ── Fetch latest metric snapshot ─────────────────────────────────────────
	url := fmt.Sprintf("%s/gestrym-progress/private/metrics/user/%d/latest", a.BaseURL, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("progress-service: failed to build request: %w", err)
	}
	req.Header.Set("X-API-Key", a.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("progress-service: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No metrics yet — return empty, caller will fall back to user input
		return &interfaces.ProgressMetrics{UserID: userID}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("progress-service: unexpected status %d", resp.StatusCode)
	}

	var psResp progressServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return nil, fmt.Errorf("progress-service: failed to decode response: %w", err)
	}

	metrics := &interfaces.ProgressMetrics{
		UserID:     psResp.Data.UserID,
		WeightKg:   psResp.Data.WeightKg,
		HeightCm:   psResp.Data.HeightCm,
		BodyFatPct: psResp.Data.BodyFatPct,
	}

	// ── Fetch weight delta (last 2 entries to compute weekly change) ──────────
	metrics.WeightDeltaKg = a.fetchWeightDelta(userID)

	return metrics, nil
}

// fetchWeightDelta computes the weight change between the two most recent records.
// Returns 0 on any error — callers must handle 0 as "no data".
func (a *ProgressServiceAdapterImpl) fetchWeightDelta(userID uint) float64 {
	url := fmt.Sprintf("%s/gestrym-progress/private/metrics/user/%d/weight-history?limit=2", a.BaseURL, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("X-API-Key", a.APIKey)

	resp, err := a.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0
	}
	defer resp.Body.Close()

	var histResp progressWeightHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&histResp); err != nil || len(histResp.Data) < 2 {
		return 0
	}

	// histResp.Data[0] = most recent, histResp.Data[1] = previous
	return histResp.Data[0].WeightKg - histResp.Data[1].WeightKg
}
