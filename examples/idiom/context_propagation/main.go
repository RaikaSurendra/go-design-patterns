// Package main demonstrates the Context Propagation pattern.
//
// Context carries request-scoped values, deadlines, and cancellation
// signals through the entire call chain.
package main

import (
	"context"
	"fmt"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// Simulated service layers — context flows through every level.
func fetchUser(ctx context.Context, userID string) (string, error) {
	reqID := GetRequestID(ctx)
	select {
	case <-time.After(50 * time.Millisecond):
		fmt.Printf("  [%s] fetchUser(%s): ok\n", reqID, userID)
		return "Alice", nil
	case <-ctx.Done():
		return "", fmt.Errorf("[%s] fetchUser cancelled: %w", reqID, ctx.Err())
	}
}

func fetchOrders(ctx context.Context, userID string) (int, error) {
	reqID := GetRequestID(ctx)
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Printf("  [%s] fetchOrders(%s): ok\n", reqID, userID)
		return 3, nil
	case <-ctx.Done():
		return 0, fmt.Errorf("[%s] fetchOrders cancelled: %w", reqID, ctx.Err())
	}
}

func main() {
	fmt.Println("--- Request with values + timeout ---")
	ctx := WithRequestID(context.Background(), "req-abc-123")
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	name, _ := fetchUser(ctx, "user-1")
	orders, _ := fetchOrders(ctx, "user-1")
	fmt.Printf("  User: %s, Orders: %d\n", name, orders)

	fmt.Println("\n--- Request that times out ---")
	ctx2 := WithRequestID(context.Background(), "req-xyz-789")
	ctx2, cancel2 := context.WithTimeout(ctx2, 30*time.Millisecond)
	defer cancel2()

	_, err := fetchUser(ctx2, "user-2")
	if err != nil {
		fmt.Printf("  %v\n", err)
	}
}
