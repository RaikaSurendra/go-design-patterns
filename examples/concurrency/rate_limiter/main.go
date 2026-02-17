// Package main demonstrates the Rate Limiter pattern.
//
// A rate limiter controls how frequently an operation can be performed,
// preventing overload of downstream services.
package main

import (
	"fmt"
	"time"
)

type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

func NewLimiter(rate int) *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, rate),
		stop:   make(chan struct{}),
	}
	for i := 0; i < rate; i++ {
		l.tokens <- struct{}{}
	}

	go func() {
		interval := time.Second / time.Duration(rate)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case l.tokens <- struct{}{}:
				default:
				}
			case <-l.stop:
				return
			}
		}
	}()

	return l
}

func (l *Limiter) Wait() {
	<-l.tokens
}

func (l *Limiter) Stop() {
	close(l.stop)
}

func main() {
	limiter := NewLimiter(5) // 5 per second
	defer limiter.Stop()

	fmt.Println("Sending 10 requests at 5/sec:")
	start := time.Now()

	for i := 1; i <= 10; i++ {
		limiter.Wait()
		elapsed := time.Since(start).Round(time.Millisecond)
		fmt.Printf("  Request %2d at %v\n", i, elapsed)
	}
}
