# Worker Pool <span class="gp-difficulty gp-difficulty--medium">Medium</span>

A worker pool maintains a fixed number of goroutines that process tasks from a
shared channel. This bounds resource usage while maximizing throughput — new
tasks queue up instead of spawning unbounded goroutines.

## Implementation

```go
package pool

import "sync"

// Task represents a unit of work.
type Task[T any, R any] struct {
	Input  T
	Result R
	Err    error
}

// Pool runs a fixed number of workers that process tasks from an input channel.
type Pool[T any, R any] struct {
	workers int
	fn      func(T) (R, error)
}

func New[T any, R any](workers int, fn func(T) (R, error)) *Pool[T, R] {
	return &Pool[T, R]{workers: workers, fn: fn}
}

// Run processes all inputs and returns results in completion order.
func (p *Pool[T, R]) Run(inputs []T) []Task[T, R] {
	in := make(chan T, len(inputs))
	out := make(chan Task[T, R], len(inputs))

	// Start workers
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

	// Feed inputs
	for _, input := range inputs {
		in <- input
	}
	close(in)

	// Wait and close output
	go func() {
		wg.Wait()
		close(out)
	}()

	// Collect results
	var results []Task[T, R]
	for task := range out {
		results = append(results, task)
	}
	return results
}
```

## Usage

```go
pool := pool.New(4, func(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
})

urls := []string{
	"https://golang.org",
	"https://pkg.go.dev",
	"https://go.dev",
}

for _, task := range pool.Run(urls) {
	if task.Err != nil {
		fmt.Printf("%s -> error: %v\n", task.Input, task.Err)
	} else {
		fmt.Printf("%s -> %d\n", task.Input, task.Result)
	}
}
```

## Rules of Thumb

- Size the pool based on the bottleneck: CPU-bound work → `runtime.NumCPU()`, I/O-bound → higher.
- The input channel's buffer controls backpressure. A full buffer blocks the producer.
- For cancellation support, pass a `context.Context` through the task and check it in the worker function.
