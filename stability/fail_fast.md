# Fail-Fast Pattern <span class="gp-difficulty gp-difficulty--easy">Easy</span>

The fail-fast pattern checks the availability of required resources at the start
of a request and fails immediately if the requirements are not satisfied. Rather
than performing expensive work only to discover a missing dependency halfway
through, the request is rejected early with a clear error.

This reduces wasted computation, frees resources faster, and gives callers
immediate feedback so they can retry or fall back.

## Implementation

```go
package failfast

import "errors"

// Checker validates whether a required precondition is met.
type Checker func() error

// Handler represents the actual business logic to execute.
type Handler func() error

// FailFast verifies all preconditions before invoking the handler.
// If any check fails, it returns immediately without calling the handler.
func FailFast(handler Handler, checks ...Checker) error {
	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return handler()
}
```

A concrete example — a request handler that validates a database connection and
cache availability before processing:

```go
package failfast

import "errors"

var (
	ErrDBUnavailable    = errors.New("database is unavailable")
	ErrCacheUnavailable = errors.New("cache is unavailable")
)

func CheckDB(db *sql.DB) Checker {
	return func() error {
		if err := db.Ping(); err != nil {
			return ErrDBUnavailable
		}
		return nil
	}
}

func CheckCache(cache CacheClient) Checker {
	return func() error {
		if !cache.IsAlive() {
			return ErrCacheUnavailable
		}
		return nil
	}
}
```

## Usage

```go
err := failfast.FailFast(
	func() error {
		// Expensive business logic that requires both DB and cache.
		return processOrder(order)
	},
	failfast.CheckDB(db),
	failfast.CheckCache(cache),
)

if err != nil {
	log.Printf("request rejected early: %v", err)
}
```

## Rules of Thumb

- Only check preconditions that are cheap to verify (pings, health flags). The goal is to fail fast, not to add latency.
- Combine with the circuit breaker pattern: a circuit breaker remembers previous failures, while fail-fast verifies current readiness.
- Return specific error types so callers can distinguish between "not ready" and "processing failed."
