# Pipeline <span class="gp-difficulty gp-difficulty--medium">Medium</span>

A pipeline is a series of stages connected by channels, where each stage is a
group of goroutines that receives values from upstream, performs a function on
that data, and sends the results downstream. Pipelines allow you to compose
complex data processing from simple, reusable stages.

## Implementation

```go
package pipeline

// Stage transforms input values into output values.
type Stage[In any, Out any] func(in <-chan In) <-chan Out

// Generate creates the initial stage by feeding a slice into a channel.
func Generate[T any](items ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, item := range items {
			out <- item
		}
	}()
	return out
}

// Map applies a function to each value in the input channel.
func Map[In any, Out any](in <-chan In, fn func(In) Out) <-chan Out {
	out := make(chan Out)
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

// Filter passes through only values that satisfy the predicate.
func Filter[T any](in <-chan T, pred func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range in {
			if pred(v) {
				out <- v
			}
		}
	}()
	return out
}

// Collect drains a channel into a slice.
func Collect[T any](in <-chan T) []T {
	var result []T
	for v := range in {
		result = append(result, v)
	}
	return result
}
```

## Usage

```go
// Pipeline: generate numbers → square them → keep only even results
numbers := pipeline.Generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
squared := pipeline.Map(numbers, func(n int) int { return n * n })
even := pipeline.Filter(squared, func(n int) bool { return n%2 == 0 })

results := pipeline.Collect(even)
fmt.Println(results) // [4 16 36 64 100]
```

## Rules of Thumb

- Each stage owns its output channel — the stage that creates it is responsible for closing it.
- Pipelines provide natural backpressure: a slow stage causes upstream stages to block.
- For cancellation, pass a `done` channel or `context.Context` and select on it alongside channel operations.
- Fan-out a stage into multiple goroutines for CPU-bound stages; merge with fan-in at the end.
