# OGL - OVYA Go Library

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**OGL** (OVYA Go Library) is a focused utility library providing file, OS, and string helpers for Go services within the OVYA monorepo workspace.

> **Note:** Platform, configuration, database, PostgreSQL, and logging packages have been extracted into a separate module: [`github.com/piprim/mmw/platform`](../platform/) at `poc/libs/platform/`. Use that module for config loading, the platform runner, DB/UoW patterns, outbox relay, and structured logging.

## Packages

### `file/`

File system utilities for reading, writing, and manipulating files.

### `os/`

Operating system helpers, including `EnvMap()` for reading environment variables as a map.

### `string/`

String manipulation utilities, including Unicode normalization and accent removal helpers.

## Development

### Linting

```bash
golangci-lint run ./...
```

### Code Style

- Run `golangci-lint` before committing
- Write tests for all new functionality

## Contributing

This library is part of the OVYA monorepo workspace. See the main repository's contribution guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

Copyright (c) 2026 OVYA

## Support

For issues and questions:
- Open an issue in the main repository
- Contact the OVYA development team

---

**Note:** This library is designed for use within the OVYA workspace monorepo using Go workspaces. External usage may require adjustments to import paths and workspace configuration.
