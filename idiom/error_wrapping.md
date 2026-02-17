# Error Wrapping & Sentinel Errors <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Go 1.13 introduced error wrapping with `fmt.Errorf` and `%w`, along with
`errors.Is` and `errors.As` for inspecting wrapped error chains. This replaces
ad-hoc string matching with structured, composable error handling.

**Sentinel errors** are package-level variables that represent specific failure
conditions. Combined with wrapping, they let callers check *what* went wrong
while preserving *where* it went wrong.

## Implementation

```go
package store

import (
	"errors"
	"fmt"
)

// Sentinel errors — callers check these with errors.Is.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
)

// Custom error with structured context.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}

// GetUser wraps sentinel errors with context using %w.
func GetUser(id string) (*User, error) {
	row, err := db.QueryRow("SELECT ...", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("GetUser(%s): %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("GetUser(%s): %w", id, err)
	}
	return row, nil
}
```

## Usage

```go
user, err := store.GetUser("abc-123")
if err != nil {
	// Check sentinel error — works through any number of wrapping layers.
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "user not found", 404)
		return
	}

	// Check for a specific error type.
	var valErr *store.ValidationError
	if errors.As(err, &valErr) {
		http.Error(w, valErr.Message, 400)
		return
	}

	// Unknown error.
	log.Printf("unexpected: %v", err)
	http.Error(w, "internal error", 500)
}
```

## Rules of Thumb

- Use `%w` (not `%v`) in `fmt.Errorf` to preserve the error chain for `errors.Is` and `errors.As`.
- Define sentinel errors at the package level with `errors.New`. Keep them stable — callers depend on them.
- Use `errors.Is` for value comparison (sentinels), `errors.As` for type assertion (custom error types).
- Add context when wrapping (`fmt.Errorf("GetUser(%s): %w", id, err)`) so the error message describes the path.
- Never compare errors with `==` if they might be wrapped — always use `errors.Is`.
