# DWD Downloader

This tool downloads NWP GRIB2 data from DWD's Open Data File Server (https://opendata.dwd.de) using HTTPS.

## Features

- Download weather model data from DWD
- Support for multiple models: cosmo-d2, cosmo-d2-eps, icon, icon-eps, icon-eu, icon-eu-eps, icon-d2, icon-d2-eps
- Configurable time steps and intervals
- Automatic timestamp detection
- BZ2 decompression and file extraction
- Parallel downloads with configurable concurrency
- Support for multiple timestamp formats (RFC3339, YYYY-MM-DD HH:MM:SS, etc.)

## Installation

### Prerequisites

- Go 1.26 or later

### Build from source

Using Makefile (recommended):

```bash
cd go
make build      # Build the binary to bin/ directory
make install    # Install to $GOPATH/bin
make release    # Build binaries for all platforms (linux, macOS, windows)
```

Manual build:

```bash
cd go
go mod download
go build -o downloader ./cmd/downloader
```

This will create a `downloader` binary in the current directory.

## Usage

### Display help

```bash
./downloader --help
```

### Download data

```bash
./downloader \
  --model icon-eu \
  --single-level-fields t_2m \
  --max-time-step 5 \
  --directory /path/to/output
```

### Parallel downloads

```bash
./downloader \
  --model icon-eu \
  --single-level-fields t_2m \
  --max-time-step 10 \
  --parallel 4 \
  --directory /path/to/output
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `--model` | NWP model name (cosmo-d2, icon-eu, etc.) | Latest available model |
| `--grid` | Model grid type | Default grid for selected model |
| `--single-level-fields` | Comma-separated list of fields to download (e.g., t_2m,tmax_2m) | **Required** |
| `--min-time-step` | Minimum forecast time step | 0 |
| `--max-time-step` | Maximum forecast time step | 0 |
| `--time-step-interval` | Interval between time steps | 1 |
| `--timestamp` | Timestamp in format 'YYYY-MM-DD HH:MM:SS' or RFC3339 | Latest available |
| `--directory` | Download directory | Current directory |
| `--parallel` | Number of parallel downloads | 1 (sequential) |

## Differences from Python Version

1. **CLI Framework**: Uses `urfave/cli/v2` instead of Python's `click`
2. **String Formatting**: Custom formatter implementation for `{param!U}` and `{param!L}` syntax
3. **Error Handling**: More explicit error handling and logging
4. **Build Process**: Requires compilation before use
5. **Parallel Downloads**: Supports concurrent downloads via `--parallel` flag (default=1 for sequential)
6. **Timestamp Parsing**: Supports multiple timestamp formats including RFC3339
7. **HTTP Timeout**: Configurable HTTP client timeout (30 seconds)

## Project Structure

```
go/
├── Makefile          # Common build targets
├── README.md         # This documentation
├── cmd/
│   └── downloader/    # Main application entry point
├── internal/
│   ├── formatter/     # Custom string formatter
│   ├── logger/        # Logging utilities
│   ├── models/        # Model configurations and JSON data
│   └── version/       # Version information
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```
