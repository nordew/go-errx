# errx

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/errx.svg)](https://pkg.go.dev/github.com/yourusername/errx)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/errx)](https://goreportcard.com/report/github.com/yourusername/errx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, expressive error handling library for Go that provides error codes, context, and a fluent builder API.

## Features

- **Error Codes**: Categorize errors by type (NotFound, BadRequest, etc.)
- **Context Preservation**: Wrap errors while maintaining the original cause
- **Fluent Builder API**: Create descriptive errors with a clean, chainable syntax
- **Standard Error Compatibility**: Works seamlessly with Go 1.13+ error handling
- **Zero Dependencies**: Pure Go with no external dependencies

## Installation

```bash
go get github.com/nordew/go-errx
```

## Quick Start

```go
package main

import (
    "fmt"
    github.com/nordew/go-errx"
)

func main() {
    // Create a simple error
    err1 := errx.NewNotFound().WithMessage("user not found").Error()
    fmt.Println(err1) // [NOT_FOUND] user not found

    // Create an error with a cause
    originalErr := fmt.Errorf("database connection timeout")
    err2 := errx.NewInternal().
        WithMessage("failed to fetch user").
        WithCause(originalErr).
        Build()
    fmt.Println(err2) // [INTERNAL] failed to fetch user: database connection timeout

    // Check error type
    if errx.IsCode(err2, errx.Internal) {
        fmt.Println("This is an internal error")
    }
}
```

## Error Codes

The following standard error codes are provided:

| Code            | Description                                 | Typical HTTP Status |
| --------------- | ------------------------------------------- | ------------------- |
| `BadRequest`    | Invalid input, parameters or request format | 400                 |
| `Unauthorized`  | Authentication required                     | 401                 |
| `Forbidden`     | Permission denied                           | 403                 |
| `NotFound`      | Resource not found                          | 404                 |
| `Conflict`      | Resource conflicts with existing data       | 409                 |
| `AlreadyExists` | Resource already exists                     | 409                 |
| `Validation`    | Input validation failed                     | 422                 |
| `Internal`      | Internal server or system errors            | 500                 |
| `Timeout`       | Operation timed out                         | 504                 |

## Usage Examples

### Creating Errors

```go
// Basic error
err := errx.NewBadRequest().WithMessage("invalid user ID").Error()

// With formatting
err := errx.NewValidation().WithMessagef("value must be between %d and %d", min, max).Error()

// With cause
err := errx.NewInternal().
    WithMessage("database query failed").
    WithCause(dbErr).
    Build()

// Custom error code
err := errx.New(errx.Conflict).WithMessage("user already exists").Build()
```

### Error Constants

Define common errors as package-level variables:

```go
var (
    ErrNotFound = errx.NewNotFound().WithMessage("resource not found").Error()
    ErrInvalidInput = errx.NewBadRequest().WithMessage("invalid input").Error()
)

// Later in your code
if someCondition {
    return ErrNotFound
}
```

### Wrapping Errors

```go
// Wrap an existing error
origErr := errors.New("connection timeout")
wrappedErr := errx.Wrap(origErr, errx.Internal, "database operation failed")

// Only wrap if not nil
err = errx.WrapIfErr(maybeNilErr, errx.Internal, "operation failed")
```

### Error Checking

```go
// Check by code
if errx.IsCode(err, errx.NotFound) {
    // Handle not found case
}

// Check against predefined error
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}

// Extract information
code := errx.GetCode(err)
message := errx.GetMessage(err)
```

### Error Handling at API Boundaries

```go
func HandleError(w http.ResponseWriter, err error) {
    switch {
    case errx.IsCode(err, errx.NotFound):
        http.Error(w, errx.GetMessage(err), http.StatusNotFound)
    case errx.IsCode(err, errx.BadRequest):
        http.Error(w, errx.GetMessage(err), http.StatusBadRequest)
    case errx.IsCode(err, errx.Validation):
        http.Error(w, errx.GetMessage(err), http.StatusUnprocessableEntity)
    default:
        // Log the full error with cause for internal errors
        log.Printf("ERROR: %+v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

## Complete Example

See the [examples](./examples) folder for a complete working example showing various usage patterns.

## Best Practices

1. **Define Error Constants**: For common errors, define package variables for reuse.
2. **Add Context**: Always add meaningful messages that explain what went wrong.
3. **Preserve Original Errors**: Use `WithCause()` to keep the original error for debugging.
4. **Error Boundaries**: Handle errors at service boundaries based on error codes.
5. **Be Specific**: Use the most specific error code that applies to the situation.

## Compatibility

The `errx` package is fully compatible with Go's standard error handling mechanisms:

- Works with `errors.Is()` and `errors.As()`
- Implements the `Unwrap()` method for error chains
- Compatible with Go 1.13+ error handling

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
