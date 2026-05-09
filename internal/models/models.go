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
// If timezone is empty, uses UTC. Otherwise, uses the specified timezone.
func GetMostRecentTimestamp(waitTimeMinutes int, modelIntervalHours int, timezone string) time.Time {
	var now time.Time
	var loc *time.Location
	var err error

	if timezone != "" {
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			// Fall back to UTC if timezone is invalid
			loc = time.UTC
		}
	} else {
		loc = time.UTC
	}

	now = time.Now().In(loc)
	waitDuration := time.Duration(waitTimeMinutes) * time.Minute
	now = now.Add(-waitDuration)
	latestAvailableUTCRun := int(now.Hour()/modelIntervalHours) * modelIntervalHours
	modelTimestamp := time.Date(now.Year(), now.Month(), now.Day(), latestAvailableUTCRun, 0, 0, 0, loc)
	return modelTimestamp
}

// GetMostRecentModelTimestamp calculates the most recent model timestamp for a given model
// Uses UTC by default if timezone is empty.
func GetMostRecentModelTimestamp(model ModelConfig, timezone string) time.Time {
	waitTimeMinutes := model.OpenDataDeliveryOffsetMinutes
	modelIntervalHours := model.IntervalHours
	return GetMostRecentTimestamp(waitTimeMinutes, modelIntervalHours, timezone)
}
