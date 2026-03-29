# OGL - OVYA Go Library

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**OGL** (OVYA Go Library) is a shared library providing common functionality for Go services within the OVYA monorepo workspace.

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

## platform.DomainError

`ogl/platform` provides the boundary error type used by every module's application layer to surface domain errors to inbound adapters.

### Types

```go
// ErrorCode is the numeric code identifying a domain error.
// Values come from proto-generated enums (e.g. todov1.TodoErrorCode).
type ErrorCode int32

// DomainError is returned by the application layer when a domain rule is violated.
// It carries a machine-readable Code and a human-readable Message.
type DomainError struct {
    Code    ErrorCode
    Message string
}
```

### Why it lives in ogl

`DomainError` is a **boundary type**: it must be recognised by inbound adapters (Connect, HTTP) without those adapters importing module-internal domain packages. Because every module's application layer and every inbound adapter need to agree on this type, it belongs in the shared library — not in any single module or contract package.

### How to use in a module

**1. Domain layer** (`internal/domain/errors.go`) — sentinel vars only, no dependency on ogl:

```go
var ErrInvalidTitle = errors.New("invalid title")
```

**2. Application layer** (`internal/application/errors.go`) — translate sentinels to `*platform.DomainError`:

```go
func DomainErrorFor(err error) error {
    switch {
    case errors.Is(err, domain.ErrInvalidTitle):
        return &platform.DomainError{
            Code:    platform.ErrorCode(deftodo.ErrorCodeInvalidTitle),
            Message: err.Error(),
        }
    }
    return err // infra / unexpected errors pass through unchanged
}
```

The service calls `DomainErrorFor` before returning, so the inbound adapter always receives either a `*platform.DomainError` (known domain error) or a plain `error` (infra/unexpected).

**3. Inbound adapter** (`internal/adapters/inbound/connect/errors.go`) — type-switch on `*platform.DomainError`:

```go
func connectErrorFrom(err error) *connect.Error {
    domainErr, ok := errors.AsType[*platform.DomainError](err)
    if !ok {
        return connect.NewError(connect.CodeInternal, err)
    }
    // map domainErr.Code → Connect status code + attach proto detail
}
```

Error codes are defined in the contract package (e.g. `contracts/definitions/todo`) as aliases over the proto-generated `TodoErrorCode` enum, ensuring the TypeScript client sees the same numeric values.

---

**Note:** This library is designed for use within the OVYA workspace monorepo using Go workspaces. External usage may require adjustments to import paths and workspace configuration.
