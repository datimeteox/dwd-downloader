package models

import "time"

// ModelConfig represents a model configuration
type ModelConfig struct {
	Model                         string   `json:"model"`
	Scope                         string   `json:"scope"`
	IntervalHours                 int      `json:"intervalHours"`
	Grids                         []string `json:"grids"`
	Pattern                       Pattern  `json:"pattern"`
	OpenDataDeliveryOffsetMinutes int      `json:"openDataDeliveryOffsetMinutes"`
}

// Pattern represents URL patterns
type Pattern struct {
	SingleLevel string `json:"single-level"`
}

// Available holds available models and grids
type Available struct {
	Models map[string]ModelConfig `json:"models"`
	Grids  map[string]string      `json:"grids"`
}

// GetMostRecentTimestamp calculates the most recent timestamp for model data
func GetMostRecentTimestamp(waitTimeMinutes int, modelIntervalHours int) time.Time {
	now := time.Now().UTC()
	waitDuration := time.Duration(waitTimeMinutes) * time.Minute
	now = now.Add(-waitDuration)
	latestAvailableUTCRun := int(now.Hour()/modelIntervalHours) * modelIntervalHours
	modelTimestamp := time.Date(now.Year(), now.Month(), now.Day(), latestAvailableUTCRun, 0, 0, 0, time.UTC)
	return modelTimestamp
}

// GetMostRecentModelTimestamp calculates the most recent model timestamp for a given model
func GetMostRecentModelTimestamp(model ModelConfig) time.Time {
	waitTimeMinutes := model.OpenDataDeliveryOffsetMinutes
	modelIntervalHours := model.IntervalHours
	return GetMostRecentTimestamp(waitTimeMinutes, modelIntervalHours)
}
