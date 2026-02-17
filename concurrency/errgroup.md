# Errgroup (Structured Concurrency) <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The errgroup pattern provides structured concurrency by running a group of
goroutines and waiting for all of them to complete. If any goroutine returns an
error, the group's context is cancelled and the first error is returned. This
prevents fire-and-forget goroutine leaks.

Go provides `golang.org/x/sync/errgroup`, but the core idea is simple enough
to implement with standard library primitives.

## Implementation

```go
package errgroup

import (
	"context"
	"sync"
)

// Group manages a set of goroutines that share a cancellable context.
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

// Go launches fn in a new goroutine. The first non-nil error cancels the group.
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

// Wait blocks until all goroutines finish and returns the first error.
func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}
```

## Usage

```go
g, ctx := errgroup.WithContext(context.Background())

g.Go(func(ctx context.Context) error {
	return fetchUserProfile(ctx, userID)
})

g.Go(func(ctx context.Context) error {
	return fetchUserOrders(ctx, userID)
})

g.Go(func(ctx context.Context) error {
	return fetchUserPreferences(ctx, userID)
})

// If any fetch fails, the others are cancelled via ctx.
if err := g.Wait(); err != nil {
	log.Fatal(err)
}
```

## Rules of Thumb

- Always check `ctx.Done()` inside goroutines so they actually respond to cancellation.
- Errgroup replaces the common `WaitGroup` + `error channel` + `sync.Once` boilerplate.
- For production use, prefer `golang.org/x/sync/errgroup` which also supports concurrency limits via `SetLimit`.
