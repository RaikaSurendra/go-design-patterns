# Push & Pull Pattern

The push-pull pattern distributes messages to multiple workers arranged in a
pipeline. A pusher sends work items downstream, workers (pullers) process them,
and the results are collected by a sink. This creates a multi-stage pipeline
where each stage can scale independently.

## Implementation

```go
package pushpull

import "sync"

// Pusher sends work items to a pool of pullers.
func Pusher(items []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, item := range items {
			out <- item
		}
	}()
	return out
}

// Puller processes items from the input channel and sends results downstream.
func Puller(in <-chan string, process func(string) string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for item := range in {
			out <- process(item)
		}
	}()
	return out
}

// Sink collects results from multiple puller channels into a single channel.
func Sink(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```

## Usage

```go
// Push work items
work := pushpull.Pusher([]string{"task-1", "task-2", "task-3", "task-4"})

// Create 2 parallel pullers that process items
process := func(s string) string {
	time.Sleep(100 * time.Millisecond) // simulate work
	return "done: " + s
}

puller1 := pushpull.Puller(work, process)
puller2 := pushpull.Puller(work, process)

// Collect all results into a single stream
for result := range pushpull.Sink(puller1, puller2) {
	fmt.Println(result)
}
// done: task-1
// done: task-2
// done: task-3
// done: task-4
```

## Rules of Thumb

- Each stage of the pipeline communicates through channels, making it easy to add, remove, or scale stages independently.
- Push-pull naturally provides load balancing: faster workers pull more items from the shared channel.
- For ordered results, add sequence numbers to work items and sort at the sink.
- Combine with `context.Context` for graceful cancellation of the entire pipeline.
