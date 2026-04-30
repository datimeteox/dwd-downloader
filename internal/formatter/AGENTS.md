# AGENTS.md for internal/formatter

> **Scope:** String formatting utilities
> **Parent:** AGENTS.md
> **Last updated:** 2025-01-07

---

## Overview

This package provides custom string formatting utilities that handle special placeholder syntax used in DWD URL patterns, particularly `{param!U}` (uppercase) and `{param!L}` (lowercase) transformations.

---

## Commands

| Command | Purpose | Location |
|---------|---------|----------|
| `go test -v` | Run formatter tests | formatter_test.go |

---

## Project Structure

```
internal/formatter/
├── formatter.go      # Formatting utilities
└── formatter_test.go # Tests for formatter
```

---

## Golden Samples

| For | Reference | Key Patterns |
|-----|-----------|--------------|
| Case transformation | `formatter.go:FormatString` | Handles `{param!U}`, `{param!L}` syntax |
| Placeholder replacement | `formatter.go:Format` | Replaces multiple placeholders with args |

---

## Key Functions

| Function | Purpose | Parameters |
|----------|---------|------------|
| `FormatString(format string, args ...interface{}) string` | Replaces placeholders in format string | Format string, arguments |
| `Format(format string, args ...interface{}) string` | Alias for FormatString | Format string, arguments |

---

## Placeholder Syntax

| Placeholder | Meaning | Example |
|-------------|---------|---------|
| `{param}` | Literal parameter name | `{param}` → `t_2m` |
| `{param!U}` | Uppercase parameter | `{param!U}` → `T_2M` |
| `{param!L}` | Lowercase parameter | `{param!L}` → `t_2m` |
| `{model!U}` | Uppercase model | `{model!U}` → `ICON` |
| `{model!L}` | Lowercase model | `{model!L}` → `icon` |

---

## Heuristics

| When | Do |
|------|-----|
| Adding new placeholder | Follow `{name!MODIFIER}` pattern |
| Testing format strings | Add test case to `TestFormatString` |
| Handling special chars | Ensure they work in both format and args |

---

## Boundaries

### Always
- Handle `!U` and `!L` modifiers for any placeholder
- Preserve original string if no placeholders found
- Support multiple placeholders in one format string

### Ask First
- Adding new modifier types (beyond U/L)
- Changing placeholder syntax format

### Never
- Modify input arguments
- Panic on unmatched placeholders (return as-is)
- Use regex for placeholder parsing (use string operations)

---

## File Patterns

### formatter.go
- **Lines 1-20:** Package documentation and imports
- **Lines 22-50:** `FormatString()` - Main formatting logic with modifier handling
- **Lines 52-55:** `Format()` - Alias function

### formatter_test.go
- **Lines 1-200:** `TestFormatString` - Comprehensive tests for all placeholder types
- **Lines 202-250:** `TestFormat` - Tests for alias function
- **Lines 252-300:** `TestFormatStringEdgeCases` - Edge case tests (nil, empty, unicode)
