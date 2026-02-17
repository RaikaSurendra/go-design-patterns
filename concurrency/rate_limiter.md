# Rate Limiter <span class="gp-difficulty gp-difficulty--medium">Medium</span>

A rate limiter controls how frequently an operation can be performed. It
prevents overloading a service, API, or resource by throttling requests to a
maximum rate. Go's `time.Ticker` provides a simple token-bucket style
implementation.

## Implementation

```go
package ratelimit

import "time"

// Limiter allows at most `rate` operations per second.
type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

func New(rate int) *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, rate),
		stop:   make(chan struct{}),
	}

	// Pre-fill the bucket.
	for i := 0; i < rate; i++ {
		l.tokens <- struct{}{}
	}

	// Refill tokens at the specified rate.
	go func() {
		interval := time.Second / time.Duration(rate)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case l.tokens <- struct{}{}:
				default: // bucket is full
				}
			case <-l.stop:
				return
			}
		}
	}()

	return l
}

// Wait blocks until a token is available.
func (l *Limiter) Wait() {
	<-l.tokens
}

// Stop shuts down the refill goroutine.
func (l *Limiter) Stop() {
	close(l.stop)
}
```

## Usage

```go
limiter := ratelimit.New(5) // 5 requests per second
defer limiter.Stop()

for i := 0; i < 10; i++ {
	limiter.Wait()
	fmt.Printf("Request %d at %s\n", i, time.Now().Format("15:04:05.000"))
}
// Requests are spaced ~200ms apart (5 per second).
```

## Rules of Thumb

- Use a token bucket for bursty traffic (allows short bursts up to bucket size) or a fixed ticker for smooth rate limiting.
- For production use, `golang.org/x/time/rate` provides a more robust limiter with reservation and cancellation support.
- Rate limiting should happen at system boundaries — close to where requests enter your service.
