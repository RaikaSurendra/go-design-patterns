# Context Propagation <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Context propagation passes request-scoped values, deadlines, and cancellation
signals through the call chain. In Go, `context.Context` is the standard
mechanism — it flows as the first parameter of every function in the chain,
ensuring the entire operation can be cancelled or timed out as a unit.

## Implementation

```go
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// WithRequestID injects a request ID into the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// RequestIDMiddleware extracts or generates a request ID and adds it to the context.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateID()
		}
		ctx := WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

## Usage

```go
func handleOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetRequestID(ctx)

	// The context flows through every layer, carrying the request ID
	// and honouring any deadline set by the client or middleware.
	order, err := fetchOrder(ctx, orderID)
	if err != nil {
		log.Printf("[%s] fetchOrder failed: %v", reqID, err)
		http.Error(w, "internal error", 500)
		return
	}

	details, err := enrichOrder(ctx, order)
	if err != nil {
		log.Printf("[%s] enrichOrder failed: %v", reqID, err)
	}
	// ...
}
```

## Rules of Thumb

- Pass `context.Context` as the **first parameter** of every function that does I/O or may block.
- Never store a context in a struct — pass it explicitly through function calls.
- Use typed keys (not bare strings) for `context.WithValue` to avoid collisions across packages.
- Keep context values limited to request-scoped data (request IDs, auth tokens). Do not use context as a general-purpose bag of state.
