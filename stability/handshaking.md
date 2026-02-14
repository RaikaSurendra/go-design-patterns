# Handshaking Pattern

The handshaking pattern allows a component to ask another component whether it
can accept more load before sending actual work. If the target component signals
that it is at capacity, the request is declined without even attempting the
operation. This protects both the caller and the callee from being overwhelmed.

Unlike the circuit breaker, which reacts to failures after they happen,
handshaking is a proactive, cooperative mechanism — the service itself
advertises its readiness.

## Implementation

```go
package handshaking

import "errors"

var (
	ErrServiceAtCapacity = errors.New("service is at capacity")
)

// Service represents a downstream component that supports health negotiation.
type Service interface {
	// IsReady reports whether the service can accept new work.
	IsReady() bool

	// Do performs the actual work.
	Do(req Request) (Response, error)
}

// Request and Response are domain-specific types.
type Request struct {
	Payload interface{}
}

type Response struct {
	Result interface{}
}

// Call performs a handshake with the target service before sending the request.
// If the service is not ready, it returns ErrServiceAtCapacity immediately.
func Call(svc Service, req Request) (Response, error) {
	if !svc.IsReady() {
		return Response{}, ErrServiceAtCapacity
	}

	return svc.Do(req)
}
```

A concrete service implementation using active connection tracking:

```go
package handshaking

import "sync/atomic"

type TrackedService struct {
	active   int64
	capacity int64
	handler  func(Request) (Response, error)
}

func NewTrackedService(capacity int64, handler func(Request) (Response, error)) *TrackedService {
	return &TrackedService{
		capacity: capacity,
		handler:  handler,
	}
}

func (s *TrackedService) IsReady() bool {
	return atomic.LoadInt64(&s.active) < s.capacity
}

func (s *TrackedService) Do(req Request) (Response, error) {
	atomic.AddInt64(&s.active, 1)
	defer atomic.AddInt64(&s.active, -1)

	return s.handler(req)
}
```

## Usage

```go
svc := handshaking.NewTrackedService(100, func(req handshaking.Request) (handshaking.Response, error) {
	result, err := processWork(req.Payload)
	return handshaking.Response{Result: result}, err
})

resp, err := handshaking.Call(svc, handshaking.Request{Payload: data})
if errors.Is(err, handshaking.ErrServiceAtCapacity) {
	log.Println("service busy, back off and retry later")
}
```

## Rules of Thumb

- The `IsReady` check must be cheap — it should read a counter or flag, not run diagnostics.
- Handshaking works best for in-process or sidecar communication. For remote services, consider a health-check endpoint that returns HTTP 503 when at capacity.
- Combine with retry and backoff logic on the caller side to handle transient capacity limits gracefully.
