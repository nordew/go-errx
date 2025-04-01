// examples/main.go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/nordew/go-errx"
)

// Custom business logic errors
var (
	ErrUserNotFound = errx.NewNotFound().WithMessage("user not found").Error()
	ErrInvalidInput = errx.NewBadRequest().WithMessage("invalid input").Error()
)

// UserRepository simulates a data access layer
type UserRepository struct{}

// GetUser simulates fetching a user from a database
func (r *UserRepository) GetUser(id string) (string, error) {
	// Simulate a database operation that fails
	if id == "" {
		return "", errx.NewBadRequest().WithMessagef("user ID cannot be empty").Error()
	}

	if id == "404" {
		// Example of wrapping a standard error with context
		return "", errx.Wrap(sql.ErrNoRows, errx.NotFound, "user not found in database")
	}

	if id == "500" {
		// Example of using a custom error code for internal errors
		return "", errx.NewInternal().
			WithMessage("database connection failed").
			WithCause(errors.New("timeout after 30s")).
			Build()
	}

	// Success case
	return "User: " + id, nil
}

// ValidateUserInput simulates input validation
func ValidateUserInput(input string) error {
	if input == "" {
		// Return a pre-defined error
		return ErrInvalidInput
	}

	if len(input) < 3 {
		// Example of using the Validation error type
		return errx.NewValidation().WithMessagef(
			"input length must be at least 3 characters, got %d",
			len(input),
		).Error()
	}

	return nil
}

// ProcessUserRequest demonstrates how to use the errx package in business logic
func ProcessUserRequest(userID, input string) error {
	// Validate input
	if err := ValidateUserInput(input); err != nil {
		// Just pass the error through, no need to wrap again
		return err
	}

	// Get user data
	repo := &UserRepository{}
	userData, err := repo.GetUser(userID)
	if err != nil {
		// Here we could add more context if needed
		return errx.WrapIfErr(err, errx.Internal, "failed during user processing")
	}

	fmt.Println("Successfully processed:", userData)
	return nil
}

// HandleError demonstrates how to handle errors at API boundaries
func HandleError(err error) {
	if err == nil {
		return
	}

	// Get the error code and message
	message := errx.GetMessage(err)

	// Handle specific error cases
	switch {
	case errx.IsCode(err, errx.NotFound):
		fmt.Printf("Resource not found: %s\n", message)
		// HTTP equivalent: return 404 status code

	case errx.IsCode(err, errx.BadRequest):
		fmt.Printf("Bad request: %s\n", message)
		// HTTP equivalent: return 400 status code

	case errx.IsCode(err, errx.Validation):
		fmt.Printf("Validation error: %s\n", message)
		// HTTP equivalent: return 422 status code

	case errx.IsCode(err, errx.Unauthorized):
		fmt.Printf("Unauthorized: %s\n", message)
		// HTTP equivalent: return 401 status code

	case errx.IsCode(err, errx.Forbidden):
		fmt.Printf("Forbidden: %s\n", message)
		// HTTP equivalent: return 403 status code

	case errors.Is(err, ErrUserNotFound):
		// Example of using errors.Is with custom errors
		fmt.Println("User not found (using errors.Is)")

	default:
		// Handle internal errors or unexpected cases
		fmt.Printf("Internal error: %s\n", err)
		// Log the full error with its cause for debugging
		log.Printf("ERROR: %+v\n", err)
		// HTTP equivalent: return 500 status code
	}
}

func main() {
	fmt.Println("errx Examples\n" + strings.Repeat("-", 40))

	// Example 1: Basic usage with validation
	fmt.Println("\nExample 1: Validation error")
	err1 := ProcessUserRequest("123", "ab")
	HandleError(err1)

	// Example 2: NotFound error
	fmt.Println("\nExample 2: NotFound error")
	err2 := ProcessUserRequest("404", "valid-input")
	HandleError(err2)

	// Example 3: Internal error
	fmt.Println("\nExample 3: Internal error")
	err3 := ProcessUserRequest("500", "valid-input")
	HandleError(err3)

	// Example 4: Successful request
	fmt.Println("\nExample 4: Successful request")
	err4 := ProcessUserRequest("123", "valid-input")
	if err4 == nil {
		fmt.Println("Request processed successfully")
	}

	// Example 5: Using error constants
	fmt.Println("\nExample 5: Using pre-defined errors")
	err5 := ProcessUserRequest("123", "")
	if errors.Is(err5, ErrInvalidInput) {
		fmt.Println("Invalid input detected using errors.Is()")
	}
	HandleError(err5)

	// Example 6: Direct builder usage for custom scenarios
	fmt.Println("\nExample 6: Custom builder usage")
	customErr := errx.New(errx.Conflict).
		WithMessagef("Resource '%s' already exists", "user-profile").
		WithCause(errors.New("duplicate key violation")).
		Error()
	HandleError(customErr)
}
