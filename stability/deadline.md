# Deadline Pattern <span class="gp-difficulty gp-difficulty--easy">Easy</span>

The deadline pattern allows a client to stop waiting for a response once a
specified amount of time has passed, at which point the probability of a
successful response becomes too low to be useful. This avoids tying up resources
indefinitely on slow or unresponsive operations.

In Go, the `context` package provides first-class support for deadlines and
timeouts, making it the idiomatic way to implement this pattern.

## Implementation

```go
package deadline

import (
	"context"
	"time"
)

// Work represents a unit of work that respects context cancellation.
type Work func(ctx context.Context) error

// WithDeadline wraps a unit of work with a deadline. If the work does not
// complete before the deadline, the context is cancelled and an error is
// returned.
func WithDeadline(timeout time.Duration, work Work) Work {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan error, 1)

		go func() {
			done <- work(ctx)
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
```

## Usage

```go
slowOperation := func(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("operation completed")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Wrap the slow operation with a 2-second deadline.
wrapped := deadline.WithDeadline(2*time.Second, slowOperation)

err := wrapped(context.Background())
if err != nil {
	fmt.Println(err) // context deadline exceeded
}
```

## Rules of Thumb

- Always propagate the `context.Context` into downstream calls so that cancellation reaches every layer.
- Choose deadline values based on observed latency percentiles (e.g. p99) rather than arbitrary round numbers.
- Prefer `context.WithTimeout` for relative durations and `context.WithDeadline` for absolute wall-clock deadlines.
