// Package main demonstrates the Futures & Promises pattern.
//
// A future is a placeholder for a result from an asynchronous computation.
// Multiple futures run concurrently; you block only when needing the value.
package main

import (
	"fmt"
	"time"
)

type result[T any] struct {
	value T
	err   error
}

type Future[T any] struct {
	ch chan result[T]
}

func NewFuture[T any](fn func() (T, error)) *Future[T] {
	f := &Future[T]{ch: make(chan result[T], 1)}
	go func() {
		val, err := fn()
		f.ch <- result[T]{value: val, err: err}
	}()
	return f
}

func (f *Future[T]) Get() (T, error) {
	r := <-f.ch
	f.ch <- r // put back for subsequent calls
	return r.value, r.err
}

func main() {
	start := time.Now()

	// Start two independent async operations
	priceFuture := NewFuture(func() (float64, error) {
		time.Sleep(200 * time.Millisecond) // simulate API call
		return 99.95, nil
	})

	stockFuture := NewFuture(func() (int, error) {
		time.Sleep(150 * time.Millisecond) // simulate DB query
		return 42, nil
	})

	// Both are running concurrently. Block only when we need results.
	price, _ := priceFuture.Get()
	stock, _ := stockFuture.Get()

	elapsed := time.Since(start)
	fmt.Printf("Price: $%.2f\n", price)
	fmt.Printf("Stock: %d units\n", stock)
	fmt.Printf("Total time: %v (both ran concurrently)\n", elapsed.Round(time.Millisecond))

	// Get is idempotent
	price2, _ := priceFuture.Get()
	fmt.Printf("Second call: $%.2f (cached)\n", price2)
}
