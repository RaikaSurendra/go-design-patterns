# Bulkhead Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The bulkhead pattern is inspired by the sectioned partitions (bulkheads) of a
ship's hull. If one section is breached, only that section floods — the rest of
the ship stays afloat. In software, the pattern isolates elements of an
application into pools so that if one fails, the others continue to function.

By partitioning resource access (e.g. connection pools, goroutine pools, or
semaphores), the bulkhead pattern prevents a single failing component from
consuming all resources and cascading into a system-wide outage.

## Implementation

Below is a `Bulkhead` that limits concurrent access to a downstream service
using a buffered channel as a semaphore.

```go
package bulkhead

import (
	"errors"
	"time"
)

var (
	ErrBulkheadFull = errors.New("bulkhead capacity full")
)

// Bulkhead limits the number of concurrent calls to a function.
type Bulkhead struct {
	sem     chan struct{}
	timeout time.Duration
}

// New creates a Bulkhead with the given maximum concurrent capacity and a
// timeout for acquiring a slot.
func New(capacity int, timeout time.Duration) *Bulkhead {
	return &Bulkhead{
		sem:     make(chan struct{}, capacity),
		timeout: timeout,
	}
}

// Execute runs fn if a slot is available within the configured timeout.
// If the bulkhead is full it returns ErrBulkheadFull without executing fn.
func (b *Bulkhead) Execute(fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	case <-time.After(b.timeout):
		return ErrBulkheadFull
	}
}
```

## Usage

```go
orderBulkhead := bulkhead.New(10, 1*time.Second)
paymentBulkhead := bulkhead.New(5, 1*time.Second)

// The order service is isolated from the payment service.
// If payments exhaust their 5 slots, orders can still proceed
// with their independent pool of 10.
err := orderBulkhead.Execute(func() error {
	return orderService.Place(order)
})

err = paymentBulkhead.Execute(func() error {
	return paymentService.Charge(order)
})

if errors.Is(err, bulkhead.ErrBulkheadFull) {
	log.Println("service is at capacity, try again later")
}
```

## Rules of Thumb

- Size each bulkhead based on the downstream service's capacity and expected latency.
- Combine with the circuit breaker pattern: a bulkhead limits concurrency while a circuit breaker stops calls to an already-failing service.
- Monitor bulkhead rejection rates — a consistently full bulkhead indicates the pool is undersized or the downstream is too slow.
