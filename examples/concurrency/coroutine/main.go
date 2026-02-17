// Package main demonstrates the Coroutines pattern.
//
// Coroutines allow suspending and resuming execution at certain points.
// In Go, goroutine + channel pairs model cooperative coroutine semantics.
package main

import "fmt"

type Coroutine[T any] struct {
	ch chan T
}

func New[T any](fn func(yield func(T))) *Coroutine[T] {
	c := &Coroutine[T]{
		ch: make(chan T),
	}
	go func() {
		defer close(c.ch)
		fn(func(val T) {
			c.ch <- val
		})
	}()
	return c
}

func (c *Coroutine[T]) Next() (T, bool) {
	val, ok := <-c.ch
	return val, ok
}

func (c *Coroutine[T]) All() <-chan T {
	return c.ch
}

func main() {
	// Fibonacci coroutine
	fmt.Println("--- Fibonacci (first 10) ---")
	fib := New(func(yield func(int)) {
		a, b := 0, 1
		for i := 0; i < 10; i++ {
			yield(a)
			a, b = b, a+b
		}
	})

	for val := range fib.All() {
		fmt.Println(val)
	}

	// String coroutine consumed one at a time
	fmt.Println("\n--- Step-by-step ---")
	steps := New(func(yield func(string)) {
		yield("first")
		yield("second")
		yield("third")
	})

	for {
		val, ok := steps.Next()
		if !ok {
			break
		}
		fmt.Printf("Got: %s\n", val)
	}
}
