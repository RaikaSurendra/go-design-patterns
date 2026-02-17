// Package main demonstrates the Cascading Failures anti-pattern and
// its prevention strategies.
//
// When one service fails and dependent services have no protection
// (timeouts, circuit breakers, bulkheads), failures cascade through
// the entire system.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// simulateService models a downstream service that may be slow or down.
func simulateService(ctx context.Context, name string, latency time.Duration) error {
	select {
	case <-time.After(latency):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", name, ctx.Err())
	}
}

// BAD: No timeout — caller waits indefinitely.
func badCall() error {
	return simulateService(context.Background(), "service-a", 5*time.Second)
}

// GOOD: Timeout prevents indefinite waiting.
func goodCall() error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	return simulateService(ctx, "service-a", 5*time.Second)
}

// GOOD: Graceful degradation returns a fallback.
func degradedCall() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := simulateService(ctx, "service-a", 5*time.Second)
	if err != nil {
		// Fall back to cached/default data
		return "cached-response", nil
	}
	return "live-response", nil
}

func main() {
	fmt.Println("--- Prevention: Timeout ---")
	start := time.Now()
	err := goodCall()
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Printf("Call timed out after %v (prevented indefinite wait)\n",
			time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("\n--- Prevention: Graceful Degradation ---")
	result, _ := degradedCall()
	fmt.Printf("Got: %s (fallback instead of failure)\n", result)

	fmt.Println("\n--- Key takeaways ---")
	fmt.Println("1. Every network call needs a timeout")
	fmt.Println("2. Design for failure: assume dependencies will fail")
	fmt.Println("3. Use circuit breakers + bulkheads for defense in depth")
	fmt.Println("4. Return cached/default data when possible")
}
