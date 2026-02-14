# Cascading Failures (Anti-Pattern)

A cascading failure occurs when a failure in one component of an interconnected
system triggers failures in dependent components, creating a domino effect that
can bring down the entire system. This is an **anti-pattern** — something to
recognize and prevent.

## How It Happens

```
Service A (overloaded)
    → times out responding to Service B
        → Service B's thread pool fills up waiting on A
            → Service C can't reach B
                → System-wide outage
```

## Example: The Problem

```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

// BAD: No timeout, no circuit breaker, no bulkhead.
// If serviceA is slow, this handler holds a goroutine and connection
// indefinitely, eventually exhausting server resources.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://service-a/api/data")
	if err != nil {
		// Service A is down — but we've already waited a long time.
		// Meanwhile, hundreds of requests piled up behind us.
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "got data from service A")
}
```

## Prevention Strategies

### 1. Timeouts

Always set deadlines on outbound calls.

```go
client := &http.Client{
	Timeout: 2 * time.Second,
}
resp, err := client.Get("http://service-a/api/data")
```

### 2. Circuit Breaker

Stop calling a failing service to give it time to recover (see
[Circuit-Breaker](/stability/circuit-breaker.md)).

### 3. Bulkheads

Isolate resource pools per dependency so one slow service doesn't consume all
resources (see [Bulkheads](/stability/bulkhead.md)).

### 4. Fail-Fast

Check dependency health before attempting expensive work (see
[Fail-Fast](/stability/fail_fast.md)).

### 5. Graceful Degradation

Return cached or default responses when a dependency is unavailable.

```go
func getData(client *http.Client, cache *Cache) (string, error) {
	resp, err := client.Get("http://service-a/api/data")
	if err != nil {
		// Fall back to cached data instead of failing entirely.
		if cached, ok := cache.Get("data"); ok {
			return cached, nil
		}
		return "", err
	}
	defer resp.Body.Close()
	// ... process response
}
```

## Rules of Thumb

- Every network call needs a timeout. No exceptions.
- Design for failure: assume every dependency *will* fail and plan what happens when it does.
- Monitor inter-service latency and error rates. Cascading failures often start with a subtle latency increase long before a hard failure.
- Test failure scenarios with chaos engineering tools to verify that your safeguards actually work.
- Combine multiple stability patterns (timeouts + circuit breaker + bulkhead) for defense in depth.
