# AGENTS.md for internal/logger

> **Scope:** Logging utilities
> **Parent:** AGENTS.md
> **Last updated:** 2025-01-07

---

## Overview

This package provides centralized logging utilities for the DWD Downloader application. It exposes a global `Logger` instance that can be used throughout the application.

---

## Commands

| Command | Purpose | Location |
|---------|---------|----------|
| `go test -v` | Run logger tests (if any) | N/A |

---

## Project Structure

```
internal/logger/
└── logger.go  # Logging utilities
```

---

## Golden Samples

| For | Reference | Key Patterns |
|-----|-----------|--------------|
| Global logger | `logger.go:Logger` | `var Logger = log.New(...)` |
| Package initialization | `logger.go:init` | Initializes logger with prefix |

---

## Key Components

| Component | Purpose | Type |
|-----------|---------|------|
| `Logger` | Global logger instance | `*log.Logger` |

---

## Heuristics

| When | Do |
|------|-----|
| Adding logging | Use `logger.Logger.Printf` or `logger.Logger.Println` |
| Adding fatal errors | Use `logger.Logger.Fatalf` for unrecoverable errors |
| Adding debug output | Check if debug mode is enabled first |

---

## Boundaries

### Always
- Use the global `Logger` instance (don't create new loggers)
- Include context in log messages (what operation, what failed)
- Use `Printf` for formatted messages with variables

### Ask First
- Adding new log levels (info, debug, warn, error)
- Changing logger configuration (prefix, flags)

### Never
- Use `log` package directly (use `logger.Logger`)
- Swallow errors silently (log them)
- Log sensitive data (passwords, tokens, keys)

---

## File Patterns

### logger.go
- **Lines 1-10:** Package documentation and imports (log)
- **Lines 12-15:** `Logger` variable definition with initialization
