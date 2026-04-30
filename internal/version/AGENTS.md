# AGENTS.md for internal/version

> **Scope:** Version information
> **Parent:** AGENTS.md
> **Last updated:** 2025-01-07

---

## Overview

This package contains the version string for the DWD Downloader application. The version is embedded into the binary at build time using Go's `-ldflags`.

---

## Project Structure

```
internal/version/
└── version.go  # Version string definition
```

---

## Key Components

| Component | Purpose | Type |
|-----------|---------|------|
| `Version` | Application version string | `string` |

---

## Heuristics

| When | Do |
|------|-----|
| Updating version | Edit `Version` constant in version.go |
| Building with version | Use `-ldflags` to inject version (handled by Makefile) |
| Reading version | Use `version.Version` |

---

## Boundaries

### Always
- Update version.go before releases
- Use semantic versioning (vX.Y.Z)
- Keep version string simple and parseable

### Ask First
- Changing version format (may break existing parsers)
- Adding metadata to version string

### Never
- Hardcode version in multiple places
- Remove version.go (breaks Makefile build)
- Use timestamps as version strings

---

## File Patterns

### version.go
- **Lines 1-5:** Package documentation
- **Lines 7-10:** `Version` constant definition
