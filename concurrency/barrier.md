# N-Barrier Pattern

The barrier pattern prevents a group of N goroutines from proceeding until all
of them have reached the barrier point. Once the last goroutine arrives, all are
released simultaneously. This is useful for phased computations where each phase
must complete before the next begins.

## Implementation

```go
package barrier

import "sync"

// Barrier blocks until N goroutines have called Wait.
type Barrier struct {
	n     int
	count int
	mu    sync.Mutex
	cond  *sync.Cond
}

func New(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks the calling goroutine until all N goroutines have called Wait.
// Once all have arrived, every goroutine is released and the barrier resets
// for reuse.
func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.count++
	if b.count == b.n {
		// Last goroutine arrived — release everyone and reset.
		b.count = 0
		b.cond.Broadcast()
		return
	}

	// Wait until the barrier is released.
	b.cond.Wait()
}
```

## Usage

```go
const workers = 5
b := barrier.New(workers)

var wg sync.WaitGroup
for i := 0; i < workers; i++ {
	wg.Add(1)
	go func(id int) {
		defer wg.Done()

		fmt.Printf("worker %d: phase 1 done\n", id)
		b.Wait() // all workers sync here

		fmt.Printf("worker %d: phase 2 done\n", id)
		b.Wait() // sync again before phase 3

		fmt.Printf("worker %d: phase 3 done\n", id)
	}(i)
}
wg.Wait()
```

## Rules of Thumb

- The barrier is reusable — after all goroutines pass through, it resets automatically.
- The number of goroutines calling `Wait` must exactly match N; otherwise the barrier will deadlock.
- Go's `sync.WaitGroup` covers the simpler case of waiting for goroutines to finish. Use a barrier when goroutines need to synchronize at intermediate points, not just at completion.
