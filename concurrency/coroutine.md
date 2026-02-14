# Coroutines Pattern

Coroutines are subroutines that allow suspending and resuming execution at
certain locations. Unlike regular functions that run to completion, coroutines
can yield intermediate values and be resumed later. In Go, goroutines combined
with channels naturally model coroutine-style cooperative multitasking.

## Implementation

```go
package coroutine

// Coroutine represents a resumable computation that yields values of type T.
type Coroutine[T any] struct {
	ch   chan T
	done chan struct{}
}

// New creates a coroutine from a function. The function receives a yield
// callback to suspend execution and emit a value.
func New[T any](fn func(yield func(T))) *Coroutine[T] {
	c := &Coroutine[T]{
		ch:   make(chan T),
		done: make(chan struct{}),
	}

	go func() {
		defer close(c.ch)
		fn(func(val T) {
			c.ch <- val
		})
	}()

	return c
}

// Next returns the next yielded value and a boolean indicating if the
// coroutine is still active.
func (c *Coroutine[T]) Next() (T, bool) {
	val, ok := <-c.ch
	return val, ok
}

// All returns a channel for range-based iteration over yielded values.
func (c *Coroutine[T]) All() <-chan T {
	return c.ch
}
```

## Usage

```go
// A coroutine that yields Fibonacci numbers.
fib := coroutine.New(func(yield func(int)) {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		yield(a)
		a, b = b, a+b
	}
})

for val := range fib.All() {
	fmt.Println(val)
}
// 0, 1, 1, 2, 3, 5, 8, 13, 21, 34

// Or consume one at a time:
counter := coroutine.New(func(yield func(string)) {
	yield("first")
	yield("second")
	yield("third")
})

val, ok := counter.Next()
fmt.Println(val, ok) // first true

val, ok = counter.Next()
fmt.Println(val, ok) // second true
```

## Rules of Thumb

- Go's goroutines are not true coroutines (they are preemptively scheduled), but goroutine + channel pairs can model cooperative coroutine semantics.
- The channel-based approach provides natural backpressure: the producer blocks on yield until the consumer reads.
- For simple sequences, a generator function returning `<-chan T` is often sufficient (see [Generators](/concurrency/generator.md)).
