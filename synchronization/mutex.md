# Lock/Mutex Pattern <span class="gp-difficulty gp-difficulty--easy">Easy</span>

A mutex (mutual exclusion) enforces exclusive access to a shared resource. Only
one goroutine can hold the lock at any time — all others block until the lock is
released. This prevents data races when multiple goroutines read and write shared
state concurrently.

Go provides `sync.Mutex` in the standard library.

## Implementation

```go
package counter

import "sync"

// Counter is a thread-safe counter protected by a mutex.
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value++
}

func (c *Counter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value--
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.value
}
```

## Usage

```go
c := &counter.Counter{}

var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Increment()
	}()
}

wg.Wait()
fmt.Println(c.Value()) // 1000
```

## Rules of Thumb

- Always use `defer mu.Unlock()` immediately after `Lock()` to guarantee the lock is released even if the function panics.
- Keep the critical section (code between Lock and Unlock) as short as possible to minimize contention.
- Never copy a `sync.Mutex` after first use — embed it in a struct and pass the struct by pointer.
- If reads vastly outnumber writes, consider `sync.RWMutex` instead (see [Read-Write Lock](/synchronization/read_write_lock.md)).
