package models

import (
	"testing"
	"time"
)

func TestGetMostRecentTimestamp(t *testing.T) {
	// Test with actual time.now but verify properties
	// We can't easily mock time.Now, so we test the properties of the result
	
	tests := []struct {
		name              string
		waitTimeMinutes   int
		modelIntervalHours int
	}{
		{
			name:               "3-hour interval, 90 min wait",
			waitTimeMinutes:    90,
			modelIntervalHours: 3,
		},
		{
			name:               "6-hour interval, 120 min wait",
			waitTimeMinutes:    120,
			modelIntervalHours: 6,
		},
		{
			name:               "12-hour interval, 540 min wait",
			waitTimeMinutes:    540,
			modelIntervalHours: 12,
		},
		{
			name:               "1-hour interval, 0 min wait",
			waitTimeMinutes:    0,
			modelIntervalHours: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMostRecentTimestamp(tt.waitTimeMinutes, tt.modelIntervalHours)
			
			// The result hour should be divisible by modelIntervalHours
			if result.Hour()%tt.modelIntervalHours != 0 {
				t.Errorf("result hour %d is not divisible by interval %d", result.Hour(), tt.modelIntervalHours)
			}
			
			// The result should have minute=0 and second=0
			if result.Minute() != 0 || result.Second() != 0 {
				t.Errorf("result should have minute=0 and second=0, got %02d:%02d", result.Minute(), result.Second())
			}
			
			// The result should be before or equal to (now - waitDuration)
			now := time.Now().UTC()
			waitDuration := time.Duration(tt.waitTimeMinutes) * time.Minute
			threshold := now.Add(-waitDuration)
			if result.After(threshold) {
				t.Errorf("result %v should be at or before threshold %v", result, threshold)
			}
		})
	}
}

func TestGetMostRecentModelTimestamp(t *testing.T) {
	tests := []struct {
		name  string
		model ModelConfig
	}{
		{
			name: "cosmo-d2 model",
			model: ModelConfig{
				Model:                     "cosmo-d2",
				Scope:                     "germany",
				IntervalHours:             3,
				Grids:                     []string{"regular-lat-lon"},
				OpenDataDeliveryOffsetMinutes: 90,
			},
		},
		{
			name: "icon model",
			model: ModelConfig{
				Model:                     "icon",
				Scope:                     "global",
				IntervalHours:             6,
				Grids:                     []string{"icosahedral"},
				OpenDataDeliveryOffsetMinutes: 240,
			},
		},
		{
			name: "icon-d2 model",
			model: ModelConfig{
				Model:                     "icon-d2",
				Scope:                     "germany",
				IntervalHours:             12,
				Grids:                     []string{"icosahedral"},
				OpenDataDeliveryOffsetMinutes: 540,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMostRecentModelTimestamp(tt.model)
			
			// The result hour should be divisible by the model's interval
			if result.Hour()%tt.model.IntervalHours != 0 {
				t.Errorf("result hour %d is not divisible by model interval %d", result.Hour(), tt.model.IntervalHours)
			}
			
			// The result should have minute=0 and second=0
			if result.Minute() != 0 || result.Second() != 0 {
				t.Errorf("result should have minute=0 and second=0, got %02d:%02d", result.Minute(), result.Second())
			}
			
			// Verify it uses the correct offset
			expectedOffset := time.Duration(tt.model.OpenDataDeliveryOffsetMinutes) * time.Minute
			now := time.Now().UTC()
			threshold := now.Add(-expectedOffset)
			if result.After(threshold) {
				t.Errorf("result %v should be at or before threshold %v (offset %v)", result, threshold, expectedOffset)
			}
		})
	}
}

// getMostRecentTimestampWithFixedTime is a helper that uses a fixed reference time
func getMostRecentTimestampWithFixedTime(referenceTime time.Time, waitTimeMinutes int, modelIntervalHours int) time.Time {
	now := referenceTime.UTC()
	waitDuration := time.Duration(waitTimeMinutes) * time.Minute
	now = now.Add(-waitDuration)
	latestAvailableUTCRun := int(now.Hour()/modelIntervalHours) * modelIntervalHours
	modelTimestamp := time.Date(now.Year(), now.Month(), now.Day(), latestAvailableUTCRun, 0, 0, 0, time.UTC)
	return modelTimestamp
}

func TestGetMostRecentTimestampWithFixedTime(t *testing.T) {
	// Test cases with fixed reference times
	tests := []struct {
		name              string
		referenceTime     time.Time
		waitTimeMinutes   int
		modelIntervalHours int
		expected          time.Time
	}{
		{
			name:           "exact interval boundary",
			referenceTime:  time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			waitTimeMinutes: 0,
			modelIntervalHours: 3,
			expected:       time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "3-hour interval with 90 min wait at 14:30",
			referenceTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			waitTimeMinutes: 90,
			modelIntervalHours: 3,
			expected:       time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "6-hour interval with 240 min wait at 14:30",
			referenceTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			waitTimeMinutes: 240,
			modelIntervalHours: 6,
			expected:       time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
		},
		{
			name:           "6-hour interval with 240 min wait at 10:30",
			referenceTime:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			waitTimeMinutes: 240,
			modelIntervalHours: 6,
			expected:       time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
		},
		{
			name:           "12-hour interval with 540 min wait at 14:30",
			referenceTime:  time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			waitTimeMinutes: 540,
			modelIntervalHours: 12,
			expected:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "midnight with 90 min wait",
			referenceTime:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			waitTimeMinutes: 90,
			modelIntervalHours: 3,
			expected:       time.Date(2024, 1, 14, 21, 0, 0, 0, time.UTC),
		},
		{
			name:           "cosmo-d2: 14:45 with 90 min offset, 3h interval",
			referenceTime:  time.Date(2024, 1, 15, 14, 45, 0, 0, time.UTC),
			waitTimeMinutes: 90,
			modelIntervalHours: 3,
			expected:       time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), // 14:45 - 90m = 13:15, int(13.15/3)*3 = 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMostRecentTimestampWithFixedTime(tt.referenceTime, tt.waitTimeMinutes, tt.modelIntervalHours)
			
			if !result.Equal(tt.expected) {
				t.Errorf("getMostRecentTimestampWithFixedTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}
