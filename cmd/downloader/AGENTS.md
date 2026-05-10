# AGENTS.md for cmd/downloader

> **Scope:** CLI application entry point
> **Parent:** AGENTS.md
> **Last updated:** 2025-01-07

---

## Overview

This directory contains the main CLI application for the DWD Downloader. It uses `urfave/cli/v2` to define flags, parse arguments, and execute download operations.

---

## Commands

| Command | Purpose | Location |
|---------|---------|----------|
| `go run .` | Run the downloader | main.go |
| `go test -v` | Run tests | main_test.go |

---

## Project Structure

```
cmd/downloader/
├── main.go       # CLI entry point, flag definitions, download logic
└── main_test.go  # Tests for URL construction and formatting
```

---

## Golden Samples

| For | Reference | Key Patterns |
|-----|-----------|--------------|
| CLI flag definition | `main.go:100-150` | Uses `cli.StringFlag`, `cli.IntFlag` with `Name`, `Usage`, `Value` |
| CLI action | `main.go:150-300` | `Action` func with `c *cli.Context` parameter |
| URL construction | `main.go:51-100` | `getGribFileURL()` with pattern replacement |
| Parallel downloads | `main.go:115-175` | Goroutines with semaphore channel |
| BZ2 download | `main.go:102-145` | HTTP GET, BZ2 decompression, file save |

---

## Key Functions

| Function | Purpose | Parameters |
|----------|---------|------------|
| `getGribFileURL(model, grid, param, timestep, timestamp, models) string` | Constructs download URL from pattern | Model config, timestamp |
| `downloadGribFile(url, destPath, destName, field, unarchive) error` | Downloads GRIB file with optional BZ2 extraction | URL, destination, field, unarchive flag |
| `downloadGribData(model, grid, param, minTS, maxTS, interval, timestamp, dest, models, parallel, unarchive) error` | Coordinates parallel downloads | Download parameters |
| `loadModels() models.Available` | Loads model configs from JSON | None |
| `formatDateIso8601(date) string` | Formats date for API | time.Time |
| `getTimestampString(date) string` | Formats date + hour | time.Time |

---

## Heuristics

| When | Do |
|------|-----|
| Adding new CLI flag | Follow existing pattern: flag definition + usage in Action |

| Constructing URLs | Use pattern replacement from models.json |
| Downloading files | Reuse `downloadGribFile` |
| Parallel operations | Use semaphore pattern (channel with capacity) |

---

## Boundaries

### Always
- Use `urfave/cli/v2` for all CLI functionality
- Validate user input (paths, model names, grids)
- Use context for cancellation in long operations
- Set HTTP timeouts (30s default)
- Handle BZ2 decompression errors gracefully

### Ask First
- Adding new URL pattern placeholders
- Changing the model run calculation ranges
- Modifying the default HTTP timeout

### Never
- Use global variables for CLI state (use context)
- Hardcode model URLs
- Ignore errors from download operations
- Use sync.WaitGroup without defer for Done()

---

## File Patterns

### main.go
- **Lines 1-25:** Imports (bytes, compress/bzip2, encoding/json, fmt, io, log, net/http, os, path/filepath, strings, sync, time, cli/v2, logger, models, version)
- **Lines 27-48:** `loadModels()` - Loads and parses models.json
- **Lines 50-87:** `getGribFileURL()` - URL construction with pattern replacement
- **Lines 102-145:** `downloadGribFile()` - HTTP download + conditional BZ2 extract
- **Lines 147-210:** `downloadGribData()` - Parallel download coordination
- **Lines 212-220:** Utility formatting functions
- **Lines 222-360:** `main()` - CLI setup and Action handler

### main_test.go
- **Lines 1-60:** `TestFormatDateIso8601` - Date formatting tests
- **Lines 62-100:** `TestGetTimestampString` - Timestamp string tests
- **Lines 102-200:** `TestGetGribFileURL` - URL construction tests with model configs

---

## Terminology

| Term | Means |
|------|-------|
| `modelrun` | The forecast cycle (00, 03, 06, 09, 12, 18) calculated from timestamp hour |
| `timestep` | Forecast hour offset from model run |
| `levtype` | Level type (single-level) |
| `param` | Parameter name (e.g., t_2m, pmsl, clch) |
