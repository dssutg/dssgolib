# AGENTS.md

## Overview

dssgolib is a zero-dependency Go utility library with 31 packages covering various functionality.

## Project Structure

- **Go version**: 1.26+
- **Module**: `github.com/dssutg/dssgolib`
- **Build tool**: Standard Go toolchain (no Makefile)

## Development Commands

```bash
# Run tests
go test ./...

# Run linter (golangci-lint)
golangci-lint run ./...

# Format code
go fmt ./...

# Run vet
go vet ./...
```

## Key Packages

- `bitfield` - Struct marshaling with bitfield tags
- `btree`, `llrb` - Tree data structures  
- `brokers` - Generic pub/sub
- `jsonc` - JSON with comments
- `utils` - Common utilities (string, math, time, slice, file operations)
- `i18n` - Internationalization and relative dates
- `poolx` - Object pooling
- `option` - Functional options pattern

## Code Style

- Standard Go conventions
- No external dependencies (standard library only)
- Tests in `*_test.go` files alongside implementation

## Notes

- Some packages have their own `go.mod` (btree, llrb) - these can be used independently
- Cross-platform support in: `sysmon` (Windows/Linux), `logrot` (POSIX)
