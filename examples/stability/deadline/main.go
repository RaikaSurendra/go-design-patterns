// Package main demonstrates the Deadline pattern.
//
// The deadline pattern wraps operations with a timeout using
// context.WithTimeout, avoiding indefinite waits on slow services.
package main

import (
	"context"
	"fmt"
	"time"
)

type Work func(ctx context.Context) error

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

func main() {
	// Fast operation — completes within deadline
	fast := WithDeadline(500*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	err := fast(context.Background())
	fmt.Printf("Fast operation: %v\n", err)

	// Slow operation — exceeds deadline
	slow := WithDeadline(200*time.Millisecond, func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	err = slow(context.Background())
	fmt.Printf("Slow operation: %v\n", err)
}
