# AGENTS.md for DWD Downloader (Go)

> **Last updated:** 2025-01-09
> **Project:** github.com/deutscherwetterdienst/dwd-downloader-go
> **Primary language:** Go
> **Project type:** CLI application (urfave/cli/v2)
> **Template:** thin root

---

## Commands (verified)

| Command | Purpose | Time |
|---------|---------|------|
| `make` or `make build` | Build the binary to bin/ | ~5s |
| `make install` | Install to $GOPATH/bin | ~5s |
| `make test` | Run all tests with race detector | ~10s |
| `make coverage` | Run tests with coverage report | ~15s |
| `make check` | Run fmt, vet, and test | ~15s |
| `make release` | Build binaries for all platforms | ~30s |
| `make fmt` | Format code with gofmt | ~2s |
| `make vet` | Run go vet | ~3s |
| `make lint` | Run golangci-lint | ~10s |
| `make run` | Run the application | ~3s |
| `make clean` | Remove build artifacts | ~1s |
| `make deps` | Download dependencies | ~5s |
| `make tidy` | Tidy go.mod | ~2s |
| `go run ./cmd/downloader` | Run directly | ~3s |
| `go test ./...` | Run all tests | ~10s |

---

## File Map

| Directory | Purpose |
|-----------|---------|
| `cmd/downloader/` | Main application entry point and CLI |
| `internal/formatter/` | Custom string formatter (handles {param!U}, {param!L}) |
| `internal/logger/` | Logging utilities |
| `internal/models/` | Model configurations and JSON data |
| `internal/version/` | Version information |
| `bin/` | Built binaries (generated) |
| `dist/` | Release builds (generated) |
| `.github/workflows/` | GitHub Actions CI/CD |
| `python/` | Python version of the downloader (separate) |

---

## Project Structure

```
.
├── AGENTS.md                    # Agent instructions (this file)
├── Makefile                    # Build targets and commands
├── README.md                   # Project documentation
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── cmd/
│   └── downloader/
│       ├── main.go             # CLI entry point
│       └── main_test.go        # Tests for main
├── internal/
│   ├── formatter/
│   │   ├── formatter.go        # String formatting utilities
│   │   └── formatter_test.go   # Formatter tests
│   ├── logger/
│   │   └── logger.go           # Logging utilities
│   ├── models/
│   │   ├── models.go           # Model loading and timestamp logic
│   │   ├── models_test.go      # Model tests
│   │   └── models.json         # Model configurations
│   └── version/
│       └── version.go          # Version string
└── .github/
    └── workflows/
        ├── build-and-test.yml  # CI: build and test
        └── release.yml          # CI: release workflow
```

---

## Golden Samples

| For | Reference | Key Patterns |
|-----|-----------|--------------|
| CLI command structure | `cmd/downloader/main.go` | Uses urfave/cli/v2, flag definitions, Action func |
| URL construction | `cmd/downloader/main.go:getGribFileURL` | Pattern replacement, model run calculation |
| Model configuration | `internal/models/models.json` | JSON config with pattern templates |
| HTTP download | `cmd/downloader/main.go:downloadAndExtractBz2FileFromURL` | HTTP client, conditional BZ2 decompression |
| Parallel downloads | `cmd/downloader/main.go:downloadGribData` | Goroutines, semaphore for concurrency control |

---

## Commands Reference

| Need | Command | Location |
|------|---------|----------|
| Build binary | `make build` | Makefile |
| Run tests | `make test` | Makefile |
| Format code | `make fmt` | Makefile |
| Check code | `make check` | Makefile |
| Download DWD data | `./bin/downloader --model icon-eu --single-level-fields t_2m` | CLI |
| Download with unarchive | `./bin/downloader --model icon-eu --single-level-fields t_2m --unarchive` | CLI |

---

## Heuristics

| When | Do |
|------|-----|
| Adding new model support | Add entry to `internal/models/models.json` |
| Adding new CLI flag | Add to `main()` in `cmd/downloader/main.go` |
| Adding URL pattern | Use existing placeholders: {model}, {param!L}, {param!U}, {grid}, {scope}, {levtype}, {modelrun:>02d}, {timestamp:%Y%m%d}, {step:>03d} |
| Downloading files | Use `downloadAndExtractBz2FileFromURL` with unarchive parameter |
| Using UTC time | Always use `time.Now().UTC()` for consistent timezone handling |

---

## Boundaries

### Always
- Use existing pattern placeholders for URL construction
- Use `time.Now().UTC()` for all timestamp calculations to ensure timezone consistency
- Use `urfave/cli/v2` for CLI flag definitions
- Validate destination paths to prevent path traversal
- Handle BZ2 decompression conditionally based on `--unarchive` flag

### Ask First
- Adding new URL pattern placeholders (may affect existing models)
- Modifying HTTP client timeout (currently 30s)
- Changing default value of `--unarchive` flag

### Never
- Hardcode model URLs (use patterns from models.json)
- Commit bin/ or dist/ directories to git
- Remove existing model configurations without replacement
- Use global variables for HTTP client (use local instances)
- Use local time without UTC conversion for model timestamps

---

## Codebase State

- **Primary language:** Go 1.26+
- **CLI framework:** urfave/cli/v2
- **HTTP client:** Standard library net/http with 30s timeout
- **Concurrency:** Goroutines with semaphore pattern for parallel downloads
- **Model configurations:** JSON-based in internal/models/models.json
- **Compression handling:** Optional BZ2 decompression via `--unarchive` flag (default: false)
- **Python version:** Separate implementation in python/ directory
- **Build artifacts:** bin/ and dist/ are gitignored (generated)

---

## Terminology

| Term | Means |
|------|-------|
| DWD | Deutscher Wetterdienst (German Weather Service) |
| NWP | Numerical Weather Prediction |
| GRIB2 | GRIdded Binary version 2 (weather data format) |
| Model run | The forecast cycle time (00, 03, 06, 09, 12, 18) |
| Timestep | Forecast hour offset from model run |
| Open Data | DWD's public data server at opendata.dwd.de |
| BZ2 | BZip2 compression format used for DWD GRIB2 files |

---

## Scope Index

| Scope | File | Description |
|-------|------|-------------|
| **Root** | `AGENTS.md` | Global rules, commands, project overview |
| **cmd/downloader** | `cmd/downloader/AGENTS.md` | CLI patterns, flag definitions, download logic |
| **internal/formatter** | `internal/formatter/AGENTS.md` | String formatting patterns |
| **internal/models** | `internal/models/AGENTS.md` | Model configuration, timestamp logic |
