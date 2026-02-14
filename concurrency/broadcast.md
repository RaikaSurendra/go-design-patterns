# Broadcast Pattern

The broadcast pattern transfers a message to all recipients simultaneously. A
single producer sends a value, and every registered consumer receives a copy.
Unlike fan-out (which distributes different items to workers), broadcast delivers
the same item to all subscribers.

## Implementation

```go
package broadcast

import "sync"

// Broadcaster sends messages to all registered listeners.
type Broadcaster[T any] struct {
	mu        sync.RWMutex
	listeners []chan T
}

func New[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{}
}

// Register adds a new listener and returns the channel it will receive on.
func (b *Broadcaster[T]) Register() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T, 1)
	b.listeners = append(b.listeners, ch)
	return ch
}

// Send delivers the message to every registered listener.
func (b *Broadcaster[T]) Send(msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.listeners {
		ch <- msg
	}
}

// Close shuts down all listener channels.
func (b *Broadcaster[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.listeners {
		close(ch)
	}
	b.listeners = nil
}
```

## Usage

```go
bc := broadcast.New[string]()

ch1 := bc.Register()
ch2 := bc.Register()
ch3 := bc.Register()

bc.Send("hello everyone")

fmt.Println(<-ch1) // hello everyone
fmt.Println(<-ch2) // hello everyone
fmt.Println(<-ch3) // hello everyone

bc.Close()
```

## Rules of Thumb

- Use buffered channels for listeners to prevent a slow consumer from blocking the broadcaster.
- If consumers have different speeds, consider adding a timeout or dropping messages for slow consumers.
- For a more robust implementation, combine with context cancellation so listeners can unsubscribe.
