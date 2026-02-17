// Package main demonstrates the Errgroup (Structured Concurrency) pattern.
//
// Errgroup runs goroutines as a group — if any fails, the group's context
// is cancelled and the first error is returned.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Group struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
	err    error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{ctx: ctx, cancel: cancel}, ctx
}

func (g *Group) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(g.ctx); err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func main() {
	fmt.Println("--- All succeed ---")
	g, ctx := WithContext(context.Background())

	g.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  fetch user profile: ok")
		return nil
	})
	g.Go(func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("  fetch user orders: ok")
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}

	fmt.Println("\n--- One fails, others cancelled ---")
	g2, ctx2 := WithContext(context.Background())

	g2.Go(func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return fmt.Errorf("payment service unavailable")
	})
	g2.Go(func(ctx context.Context) error {
		select {
		case <-time.After(200 * time.Millisecond):
			fmt.Println("  slow task completed")
			return nil
		case <-ctx2.Done():
			fmt.Println("  slow task cancelled")
			return ctx2.Err()
		}
	})

	if err := g2.Wait(); err != nil {
		fmt.Printf("Group failed: %v\n", err)
	}
	_ = ctx // used by first group's goroutines
}
