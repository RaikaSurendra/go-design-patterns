# Futures & Promises Pattern

A future acts as a placeholder for a result that is initially unknown because
the computation has not yet completed. It provides a way to access the result
of an asynchronous operation synchronously when the value is needed.

In Go, a future is naturally modeled with a goroutine that computes the result
and a channel (or struct) that delivers it.

## Implementation

```go
package future

// Future represents an asynchronous computation that will produce a value.
type Future[T any] struct {
	ch chan result[T]
}

type result[T any] struct {
	value T
	err   error
}

// New starts an asynchronous computation and returns a Future.
func New[T any](fn func() (T, error)) *Future[T] {
	f := &Future[T]{
		ch: make(chan result[T], 1),
	}

	go func() {
		val, err := fn()
		f.ch <- result[T]{value: val, err: err}
	}()

	return f
}

// Get blocks until the result is available and returns it.
func (f *Future[T]) Get() (T, error) {
	r := <-f.ch
	// Put it back so subsequent calls to Get return the same result.
	f.ch <- r
	return r.value, r.err
}
```

## Usage

```go
// Start two expensive operations concurrently.
priceFuture := future.New(func() (float64, error) {
	// simulate API call
	time.Sleep(2 * time.Second)
	return 99.95, nil
})

stockFuture := future.New(func() (int, error) {
	// simulate DB query
	time.Sleep(1 * time.Second)
	return 42, nil
})

// Both are running in parallel. Block only when we need the values.
price, err := priceFuture.Get()
if err != nil {
	log.Fatal(err)
}

stock, err := stockFuture.Get()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Price: $%.2f, Stock: %d\n", price, stock)
// Price: $99.95, Stock: 42
// Total wall-clock time: ~2s (not 3s), since both ran concurrently.
```

## Rules of Thumb

- Futures are ideal when you need to kick off multiple independent operations and collect results later.
- The `Get` method is idempotent — calling it multiple times returns the same cached result.
- For timeout support, combine with `context.WithTimeout` or use a `select` with `time.After` on the channel.
- Go channels are already one-shot futures. For simple cases, a plain `chan T` is sufficient.
