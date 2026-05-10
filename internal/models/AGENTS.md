# AGENTS.md for internal/models

> **Scope:** Model configurations and timestamp logic
> **Parent:** AGENTS.md
> **Last updated:** 2025-01-09

---

## Overview

This package handles model configuration loading, timestamp calculations for determining the most recent available model data, and pattern-based URL construction.

---

## Commands

| Command | Purpose | Location |
|---------|---------|----------|
| `go test -v` | Run models tests | models_test.go |

---

## Project Structure

```
internal/models/
├── models.go       # Model types and timestamp logic
├── models_test.go  # Tests for model functions
└── models.json     # Model configuration data
```

---

## Golden Samples

| For | Reference | Key Patterns |
|-----|-----------|--------------|
| Model configuration | `models.json` | JSON with Model, Scope, Grids, Pattern, IntervalHours, OpenDataDeliveryOffsetMinutes |
| Most recent timestamp | `models.go:GetMostRecentModelTimestamp` | Calculates based on interval, offset, and timezone |
| Timestamp parsing | `models.go:GetMostRecentTimestamp` | Handles interval, wait time, and timezone |

---

## Key Types

| Type | Purpose | Fields |
|------|---------|--------|
| `ModelConfig` | Model configuration | Model, Scope, Grids, Pattern, IntervalHours, OpenDataDeliveryOffsetMinutes |
| `Pattern` | URL pattern templates | SingleLevel string |
| `Available` | Loaded model collection | Models map, Grids map |

---

## Key Functions

| Function | Purpose | Parameters |
|----------|---------|------------|
| `GetMostRecentTimestamp(intervalHours, offsetMinutes int, timezone string) time.Time` | Calculates most recent timestamp | Interval, offset, timezone (defaults to UTC) |
| `GetMostRecentModelTimestamp(cfg ModelConfig, timezone string) time.Time` | Gets most recent for specific model | Model config, timezone (defaults to UTC) |

---

## Configuration Format (models.json)

```json
{
  "Models": [
    {
      "Model": "cosmo-d2",
      "Scope": "germany",
      "Grids": ["regular-lat-lon", "rotated-lat-lon"],
      "IntervalHours": 3,
      "OpenDataDeliveryOffsetMinutes": 90,
      "Pattern": {
        "SingleLevel": "https://opendata.dwd.de/weather/nwp/{model!L}/grib/{modelrun:>02d}/{param!L}/..."
      }
    }
  ]
}
```

---

## Pattern Placeholders

| Placeholder | Format | Example | Description |
|-------------|--------|---------|-------------|
| `{model}` | Literal | `cosmo-d2` | Model name as-is |
| `{model!L}` | Lowercase | `cosmo-d2` | Lowercase model name |
| `{param}` | Literal | `t_2m` | Parameter name as-is |
| `{param!L}` | Lowercase | `t_2m` | Lowercase parameter |
| `{param!U}` | Uppercase | `T_2M` | Uppercase parameter |
| `{grid}` | Literal | `regular-lat-lon` | Grid type |
| `{scope}` | Literal | `germany` | Geographic scope |
| `{levtype}` | Literal | `single-level` | Level type |
| `{modelrun:>02d}` | Zero-padded 2 digits | `00`, `03`, `12` | Model run hour |
| `{timestamp:%Y%m%d}` | Date format | `20240115` | Timestamp date |
| `{step:>03d}` | Zero-padded 3 digits | `000`, `024` | Forecast step |

---

## Heuristics

| When | Do |
|------|-----|
| Adding new model | Add entry to models.json with all required fields |
| Changing interval | Update IntervalHours and verify timestamp logic |
| Adding pattern | Use existing placeholder syntax |
| Testing timestamp | Use fixed time in tests with `time.Date()` |

---

## Boundaries

### Always
- Use UTC as default timezone for all timestamps
- Handle DST correctly (use timezone-aware calculations)
- Respect OpenDataDeliveryOffsetMinutes in calculations
- Validate models.json structure on load
- Fall back to UTC for invalid timezone strings

### Ask First
- Adding new placeholder types to patterns
- Changing timestamp calculation logic
- Modifying offset minutes for existing models

### Never
- Hardcode model URLs (use patterns)
- Ignore delivery offset in timestamp calculations
- Commit broken JSON (validate before commit)

---

## Model Data

Currently configured models (from models.json):
- **cosmo-d2**: Germany, 3-hour interval, 90-minute offset
- **cosmo-d2-eps**: Germany EPS, 3-hour interval, 90-minute offset
- **icon**: Global, 6-hour interval, 240-minute offset
- **icon-eps**: Global EPS, 6-hour interval, 240-minute offset
- **icon-eu**: Europe, 6-hour interval, 240-minute offset
- **icon-eu-eps**: Europe EPS, 6-hour interval, 240-minute offset
- **icon-d2**: Germany, 1-hour interval, 60-minute offset
- **icon-d2-eps**: Germany EPS, 1-hour interval, 60-minute offset

---

## File Patterns

### models.go
- **Lines 1-20:** Package documentation and imports
- **Lines 22-40:** Type definitions (ModelConfig, Pattern, Available)
- **Lines 42-80:** `GetMostRecentTimestamp()` - Time calculation logic with timezone support
- **Lines 82-100:** `GetMostRecentModelTimestamp()` - Model-specific timestamp with timezone support

### models.json
- Root object with `Models` array
- Each model has Model, Scope, Grids array, IntervalHours, OpenDataDeliveryOffsetMinutes, Pattern

### models_test.go
- **Lines 1-60:** `TestGetMostRecentTimestamp` - Various interval/offset combos
- **Lines 62-120:** `TestGetMostRecentModelTimestamp` - Model-specific tests
- **Lines 122-150:** `TestGetMostRecentTimestampWithFixedTime` - Fixed time tests
- **Lines 185-250:** `TestGetMostRecentTimestampWithTimezone` - Timezone parameter tests
