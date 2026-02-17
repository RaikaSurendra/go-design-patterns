// Package main demonstrates the Circuit Breaker pattern.
//
// A circuit breaker wraps calls to a service and stops calling it
// after consecutive failures, giving it time to recover.
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrServiceUnavailable = errors.New("service unavailable (circuit open)")

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

func (s State) String() string {
	return [...]string{"CLOSED", "OPEN", "HALF-OPEN"}[s]
}

type CircuitBreaker struct {
	mu                  sync.Mutex
	state               State
	failures            int
	failureThreshold    int
	successesInHalfOpen int
	successThreshold    int
	lastFailure         time.Time
	cooldown            time.Duration
}

func NewCircuitBreaker(failureThreshold, successThreshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		cooldown:         cooldown,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	if cb.state == Open {
		if time.Since(cb.lastFailure) > cb.cooldown {
			cb.state = HalfOpen
			cb.successesInHalfOpen = 0
		} else {
			cb.mu.Unlock()
			return ErrServiceUnavailable
		}
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.failureThreshold {
			cb.state = Open
		}
		return err
	}

	if cb.state == HalfOpen {
		cb.successesInHalfOpen++
		if cb.successesInHalfOpen >= cb.successThreshold {
			cb.state = Closed
			cb.failures = 0
		}
	} else {
		cb.failures = 0
	}
	return nil
}

func main() {
	cb := NewCircuitBreaker(3, 2, 500*time.Millisecond)

	callCount := 0
	service := func() error {
		callCount++
		if callCount <= 4 {
			return fmt.Errorf("service error #%d", callCount)
		}
		return nil
	}

	for i := 1; i <= 8; i++ {
		err := cb.Execute(service)
		cb.mu.Lock()
		state := cb.state
		cb.mu.Unlock()

		if err != nil {
			fmt.Printf("Call %d: FAILED (%v) [state=%s]\n", i, err, state)
		} else {
			fmt.Printf("Call %d: OK [state=%s]\n", i, state)
		}

		if i == 5 {
			fmt.Println("  (waiting for cooldown...)")
			time.Sleep(600 * time.Millisecond)
		}
	}
}
