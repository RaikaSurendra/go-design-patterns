// Package main demonstrates the Worker Pool pattern.
//
// A fixed number of workers process tasks from a shared channel,
// bounding resource usage while maximizing throughput.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Task[T any, R any] struct {
	Input  T
	Result R
	Err    error
}

type Pool[T any, R any] struct {
	workers int
	fn      func(T) (R, error)
}

func NewPool[T any, R any](workers int, fn func(T) (R, error)) *Pool[T, R] {
	return &Pool[T, R]{workers: workers, fn: fn}
}

func (p *Pool[T, R]) Run(inputs []T) []Task[T, R] {
	in := make(chan T, len(inputs))
	out := make(chan Task[T, R], len(inputs))

	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for input := range in {
				result, err := p.fn(input)
				out <- Task[T, R]{Input: input, Result: result, Err: err}
			}
		}()
	}

	for _, input := range inputs {
		in <- input
	}
	close(in)

	go func() {
		wg.Wait()
		close(out)
	}()

	var results []Task[T, R]
	for task := range out {
		results = append(results, task)
	}
	return results
}

func main() {
	pool := NewPool(3, func(n int) (string, error) {
		time.Sleep(100 * time.Millisecond) // simulate work
		return fmt.Sprintf("processed-%d", n*n), nil
	})

	start := time.Now()
	results := pool.Run([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	for _, t := range results {
		fmt.Printf("  Input: %d -> %s\n", t.Input, t.Result)
	}
	fmt.Printf("\n9 tasks, 3 workers, completed in %v\n", time.Since(start).Round(time.Millisecond))
}
