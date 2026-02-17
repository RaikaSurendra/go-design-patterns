// Package main demonstrates Error Wrapping & Sentinel Errors.
//
// Use fmt.Errorf with %w to wrap errors, and errors.Is / errors.As
// to inspect them through any number of wrapping layers.
package main

import (
	"errors"
	"fmt"
)

// Sentinel errors.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// Custom error type with structured context.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}

// Simulated functions that wrap errors with context.
func getUser(id string) (string, error) {
	if id == "" {
		return "", &ValidationError{Field: "id", Message: "must not be empty"}
	}
	if id == "unknown" {
		return "", fmt.Errorf("getUser(%s): %w", id, ErrNotFound)
	}
	return "Alice", nil
}

func getProfile(id string) (string, error) {
	user, err := getUser(id)
	if err != nil {
		return "", fmt.Errorf("getProfile: %w", err)
	}
	return fmt.Sprintf("Profile of %s", user), nil
}

func main() {
	cases := []string{"alice", "unknown", ""}

	for _, id := range cases {
		profile, err := getProfile(id)
		if err != nil {
			fmt.Printf("ID %q:\n", id)
			fmt.Printf("  error: %v\n", err)

			// Check sentinel via errors.Is (works through wrapping).
			if errors.Is(err, ErrNotFound) {
				fmt.Println("  -> matched: ErrNotFound")
			}

			// Check custom type via errors.As.
			var valErr *ValidationError
			if errors.As(err, &valErr) {
				fmt.Printf("  -> matched: ValidationError{Field: %q}\n", valErr.Field)
			}
			fmt.Println()
		} else {
			fmt.Printf("ID %q: %s\n\n", id, profile)
		}
	}
}
