package main

import (
	"testing"
	"time"

	"github.com/deutscherwetterdienst/dwd-downloader-go/internal/models"
)

func TestFormatDateIso8601(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC time",
			input:    time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			expected: "2024-01-15T14:30:00Z",
		},
		{
			name:     "local time converted to UTC",
			input:    time.Date(2024, 1, 15, 14, 30, 0, 0, time.FixedZone("", 3600)), // UTC+1
			expected: "2024-01-15T13:30:00Z",
		},
		{
			name:     "midnight",
			input:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "2024-01-15T00:00:00Z",
		},
		{
			name:     "with nanoseconds",
			input:    time.Date(2024, 1, 15, 14, 30, 45, 123456789, time.UTC),
			expected: "2024-01-15T14:30:45Z",
		},
		{
			name:     "end of year",
			input:    time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "2023-12-31T23:59:59Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDateIso8601(tt.input)
			if result != tt.expected {
				t.Errorf("formatDateIso8601(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetTimestampString(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "hour 14",
			input:    time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			expected: "2024011514",
		},
		{
			name:     "single digit hour",
			input:    time.Date(2024, 1, 15, 5, 0, 0, 0, time.UTC),
			expected: "2024011505",
		},
		{
			name:     "midnight",
			input:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "2024011500",
		},
		{
			name:     "end of day",
			input:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "2024123123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTimestampString(tt.input)
			if result != tt.expected {
				t.Errorf("getTimestampString(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetGribFileURL(t *testing.T) {
	// Create test models
	modelsAvail := models.Available{
		Models: map[string]models.ModelConfig{
			"cosmo-d2": {
				Model:                     "cosmo-d2",
				Scope:                     "germany",
				IntervalHours:             3,
				Grids:                     []string{"regular-lat-lon", "rotated-lat-lon"},
				OpenDataDeliveryOffsetMinutes: 90,
				Pattern: models.Pattern{
					SingleLevel: "https://opendata.dwd.de/weather/nwp/{model!L}/grib/{modelrun:>02d}/{param!L}/{model!L}_{scope}_{grid}_{levtype}_{timestamp:%Y%m%d}{modelrun:>02d}_{step:>03d}_{param!U}.grib2.bz2",
				},
			},
			"icon": {
				Model:                     "icon",
				Scope:                     "global",
				IntervalHours:             6,
				Grids:                     []string{"icosahedral"},
				OpenDataDeliveryOffsetMinutes: 240,
				Pattern: models.Pattern{
					SingleLevel: "https://opendata.dwd.de/weather/nwp/{model!L}/grib/{modelrun:>02d}/{param!L}/{model!L}_{scope}_{grid}_{levtype}_{timestamp:%Y%m%d}{modelrun:>02d}_{step:>03d}_{param!U}.grib2.bz2",
				},
			},
		},
		Grids: map[string]string{
			"regular-lat-lon":  "regular-lat-lon",
			"rotated-lat-lon":  "rotated-lat-lon",
			"icosahedral":      "icosahedral",
		},
	}

	tests := []struct {
		name     string
		model    string
		grid     string
		param    string
		timestep int
		timestamp time.Time
		expected string
	}{
		{
			name:     "cosmo-d2 with all params",
			model:    "cosmo-d2",
			grid:     "regular-lat-lon",
			param:    "t_2m",
			timestep: 0,
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/cosmo-d2/grib/12/t_2m/cosmo-d2_germany_regular-lat-lon_single-level_2024011512_000_T_2M.grib2.bz2",
		},
		{
			name:     "cosmo-d2 with rotated grid",
			model:    "cosmo-d2",
			grid:     "rotated-lat-lon",
			param:    "clch",
			timestep: 24,
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/cosmo-d2/grib/12/clch/cosmo-d2_germany_rotated-lat-lon_single-level_2024011512_024_CLCH.grib2.bz2",
		},
		{
			name:     "icon model",
			model:    "icon",
			grid:     "icosahedral",
			param:    "t_2m",
			timestep: 48,
			timestamp: time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/icon/grib/06/t_2m/icon_global_icosahedral_single-level_2024011506_048_T_2M.grib2.bz2",
		},
		{
			name:     "cosmo-d2 with default grid",
			model:    "cosmo-d2",
			grid:     "", // Should use default
			param:    "pmsl",
			timestep: 0,
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/cosmo-d2/grib/12/pmsl/cosmo-d2_germany_regular-lat-lon_single-level_2024011512_000_PMSL.grib2.bz2",
		},
		{
			name:     "uppercase param",
			model:    "cosmo-d2",
			grid:     "regular-lat-lon",
			param:    "T_2M",
			timestep: 3,
			timestamp: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/cosmo-d2/grib/00/t_2m/cosmo-d2_germany_regular-lat-lon_single-level_2024011500_003_T_2M.grib2.bz2",
		},
		{
			name:     "large timestep",
			model:    "cosmo-d2",
			grid:     "regular-lat-lon",
			param:    "t_2m",
			timestep: 72,
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			expected: "https://opendata.dwd.de/weather/nwp/cosmo-d2/grib/12/t_2m/cosmo-d2_germany_regular-lat-lon_single-level_2024011512_072_T_2M.grib2.bz2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getGribFileURL(tt.model, tt.grid, tt.param, tt.timestep, tt.timestamp, modelsAvail)
			if result != tt.expected {
				t.Errorf("getGribFileURL(%q, %q, %q, %d, %v, ...) = %q, want %q",
					tt.model, tt.grid, tt.param, tt.timestep, tt.timestamp, result, tt.expected)
			}
		})
	}
}


