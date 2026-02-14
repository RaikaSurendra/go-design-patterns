# Condition Variable Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

A condition variable allows goroutines to wait for a specific condition to become
true. Rather than busy-looping and repeatedly checking, a goroutine suspends
itself on a condition variable and is woken up when another goroutine signals
that the condition may have changed.

Go provides `sync.Cond` in the standard library, built on top of a `sync.Mutex`
or `sync.RWMutex`.

## Implementation

```go
package queue

import "sync"

// BlockingQueue is a thread-safe queue where consumers wait until an item is
// available. It uses a condition variable to avoid busy-waiting.
type BlockingQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []interface{}
}

func New() *BlockingQueue {
	q := &BlockingQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue adds an item and signals one waiting consumer.
func (q *BlockingQueue) Enqueue(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	q.cond.Signal()
}

// Dequeue blocks until an item is available, then removes and returns it.
func (q *BlockingQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Wait must be called inside a loop because spurious wakeups can occur.
	for len(q.items) == 0 {
		q.cond.Wait()
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}
```

## Usage

```go
q := queue.New()

// Producer
go func() {
	for i := 0; i < 5; i++ {
		q.Enqueue(i)
		fmt.Printf("produced: %d\n", i)
	}
}()

// Consumer — blocks until items arrive
for i := 0; i < 5; i++ {
	item := q.Dequeue()
	fmt.Printf("consumed: %v\n", item)
}
```

## Rules of Thumb

- Always call `cond.Wait()` inside a `for` loop that checks the actual condition, not an `if` — spurious wakeups are possible.
- The goroutine must hold the associated lock before calling `Wait()`. `Wait` atomically releases the lock and suspends the goroutine; upon waking it re-acquires the lock.
- Use `Signal()` to wake one waiting goroutine, `Broadcast()` to wake all of them.
- In most Go code, channels are the preferred synchronization mechanism. Use `sync.Cond` when you need to wake multiple waiters on the same condition or when channels would be awkward (e.g. waiting on a size threshold).
