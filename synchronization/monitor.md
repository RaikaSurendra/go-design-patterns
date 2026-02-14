# Monitor Pattern

A monitor combines a mutex with one or more condition variables to protect
shared state while allowing goroutines to wait for specific conditions. The
mutex guarantees exclusive access, while the condition variables coordinate
goroutines that need to wait for or signal state changes.

In Go, a monitor is composed from `sync.Mutex` (or `sync.RWMutex`) and
`sync.Cond`.

## Implementation

```go
package monitor

import "sync"

// BoundedBuffer is a classic monitor example: a fixed-size buffer where
// producers block when full and consumers block when empty.
type BoundedBuffer struct {
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	buf      []interface{}
	capacity int
}

func New(capacity int) *BoundedBuffer {
	b := &BoundedBuffer{
		buf:      make([]interface{}, 0, capacity),
		capacity: capacity,
	}
	b.notFull = sync.NewCond(&b.mu)
	b.notEmpty = sync.NewCond(&b.mu)
	return b
}

// Put adds an item, blocking if the buffer is full.
func (b *BoundedBuffer) Put(item interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.buf) == b.capacity {
		b.notFull.Wait()
	}

	b.buf = append(b.buf, item)
	b.notEmpty.Signal()
}

// Get removes and returns an item, blocking if the buffer is empty.
func (b *BoundedBuffer) Get() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.buf) == 0 {
		b.notEmpty.Wait()
	}

	item := b.buf[0]
	b.buf = b.buf[1:]
	b.notFull.Signal()
	return item
}
```

## Usage

```go
buf := monitor.New(5)

// Producer goroutines
for i := 0; i < 3; i++ {
	go func(id int) {
		for j := 0; j < 10; j++ {
			buf.Put(fmt.Sprintf("producer-%d: item-%d", id, j))
		}
	}(i)
}

// Consumer goroutines
for i := 0; i < 3; i++ {
	go func(id int) {
		for j := 0; j < 10; j++ {
			item := buf.Get()
			fmt.Printf("consumer-%d got %v\n", id, item)
		}
	}(i)
}
```

## Rules of Thumb

- A monitor is a higher-level abstraction than a raw mutex + condition variable — prefer it when you have multiple conditions on the same shared state (e.g. "not full" and "not empty").
- Always check conditions in a `for` loop, not an `if`, because of spurious wakeups.
- In idiomatic Go, a buffered channel (`make(chan T, N)`) already implements a bounded-buffer monitor. Use explicit monitors when you need more complex waiting conditions that channels cannot express.
