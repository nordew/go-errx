package errx

import (
	"errors"
	"fmt"
)

// Code represents a categorized error type
type Code string

// Standard error codes
const (
	Conflict      Code = "CONFLICT"       // Resource conflicts with existing data
	Internal      Code = "INTERNAL"       // Internal server or system errors
	NotFound      Code = "NOT_FOUND"      // Resource not found
	BadRequest    Code = "BAD_REQUEST"    // Invalid input or parameters
	AlreadyExists Code = "ALREADY_EXISTS" // Resource already exists
	Unauthorized  Code = "UNAUTHORIZED"   // Authentication required
	Forbidden     Code = "FORBIDDEN"      // Permission denied
	Timeout       Code = "TIMEOUT"        // Operation timed out
	Validation    Code = "VALIDATION"     // Input validation failed
)

// Error represents an application-specific error with code and context
type Error struct {
	Code    Code   // Error classification code
	Message string // User-friendly error message
	Err     error  // Original error (if any)
}

// Error implements the error interface and formats the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}

// Is implements error comparison for the errors.Is function
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}
	return e.Code == t.Code
}

// IsCode checks if an error has a specific error code
func IsCode(err error, code Code) bool {
	if err == nil {
		return false
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode extracts the error code from an error
// Returns Internal if the error isn't an Error type
func GetCode(err error) Code {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return Internal
}

// GetMessage extracts the user-friendly message from an error
func GetMessage(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return err.Error()
}

// Builder provides a fluent API for building Errors
type Builder struct {
	code    Code
	message string
	err     error
}

// WithMessage sets a descriptive message for the error
func (b *Builder) WithMessage(msg string) *Builder {
	b.message = msg
	return b
}

// WithCause sets the underlying cause of the error
func (b *Builder) WithCause(err error) *Builder {
	b.err = err
	return b
}

// WithMessagef sets a formatted message for the error
func (b *Builder) WithMessagef(format string, args ...interface{}) *Builder {
	b.message = fmt.Sprintf(format, args...)
	return b
}

// Build creates and returns the final Error
func (b *Builder) Build() *Error {
	return &Error{
		Code:    b.code,
		Message: b.message,
		Err:     b.err,
	}
}

// Error returns the Error as an error interface type
func (b *Builder) Error() error {
	return b.Build()
}

// WithDescription is a legacy method that immediately returns an Error
// Consider using WithMessage().Build() instead for better fluency
func (b *Builder) WithDescription(desc string) *Error {
	b.message = desc
	return b.Build()
}

// WithDescriptionAndCause is a legacy method that immediately returns an Error
// Consider using WithMessage().WithCause().Build() instead for better fluency
func (b *Builder) WithDescriptionAndCause(desc string, cause error) *Error {
	b.message = desc
	b.err = cause
	return b.Build()
}

// Error constructors
// Each returns a Builder initialized with the appropriate code

// New creates a new Builder with the specified code
func New(code Code) *Builder {
	return &Builder{code: code}
}

// NewBadRequest creates an error builder for BadRequest errors
func NewBadRequest() *Builder {
	return &Builder{code: BadRequest}
}

// NewNotFound creates an error builder for NotFound errors
func NewNotFound() *Builder {
	return &Builder{code: NotFound}
}

// NewConflict creates an error builder for Conflict errors
func NewConflict() *Builder {
	return &Builder{code: Conflict}
}

// NewInternal creates an error builder for Internal errors
func NewInternal() *Builder {
	return &Builder{code: Internal}
}

// NewAlreadyExists creates an error builder for AlreadyExists errors
func NewAlreadyExists() *Builder {
	return &Builder{code: AlreadyExists}
}

// NewUnauthorized creates an error builder for Unauthorized errors
func NewUnauthorized() *Builder {
	return &Builder{code: Unauthorized}
}

// NewForbidden creates an error builder for Forbidden errors
func NewForbidden() *Builder {
	return &Builder{code: Forbidden}
}

// NewTimeout creates an error builder for Timeout errors
func NewTimeout() *Builder {
	return &Builder{code: Timeout}
}

// NewValidation creates an error builder for Validation errors
func NewValidation() *Builder {
	return &Builder{code: Validation}
}

// Wrap creates an Error that wraps an existing error with the given code
func Wrap(err error, code Code, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapIfErr wraps an error only if it's not nil
func WrapIfErr(err error, code Code, message string) error {
	if err == nil {
		return nil
	}
	return Wrap(err, code, message)
}
