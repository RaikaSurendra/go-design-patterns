# Producer Consumer Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The producer-consumer pattern separates task generation from task execution.
Producers push work items into a shared buffer, and consumers pull items from
the buffer to process them. This decouples the rate of production from the rate
of consumption and allows both sides to operate independently.

In Go, a buffered channel is the natural shared buffer.

## Implementation

```go
package prodcon

import "sync"

// Item represents a unit of work.
type Item struct {
	ID   int
	Data string
}

// Producer generates items and sends them to the work channel.
func Producer(id int, count int, ch chan<- Item) {
	for i := 0; i < count; i++ {
		item := Item{
			ID:   id*1000 + i,
			Data: fmt.Sprintf("item-%d from producer-%d", i, id),
		}
		ch <- item
	}
}

// Consumer reads items from the work channel and processes them.
func Consumer(id int, ch <-chan Item, fn func(Item), wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range ch {
		fn(item)
	}
}
```

## Usage

```go
ch := make(chan prodcon.Item, 10)

// Start 3 producers
var producerWg sync.WaitGroup
for i := 0; i < 3; i++ {
	producerWg.Add(1)
	go func(id int) {
		defer producerWg.Done()
		prodcon.Producer(id, 5, ch)
	}(i)
}

// Start 2 consumers
var consumerWg sync.WaitGroup
for i := 0; i < 2; i++ {
	consumerWg.Add(1)
	go prodcon.Consumer(i, ch, func(item prodcon.Item) {
		fmt.Printf("consumer-%d processed: %s\n", i, item.Data)
	}, &consumerWg)
}

// Wait for all producers to finish, then close the channel.
producerWg.Wait()
close(ch)

// Wait for all consumers to drain the remaining items.
consumerWg.Wait()
```

## Rules of Thumb

- The buffer size controls backpressure: a full buffer blocks producers, giving consumers time to catch up.
- Always close the channel from the producer side (or a coordinator) — never from the consumer.
- Tune the number of producers, consumers, and buffer size based on observed throughput and latency.
- For graceful shutdown, use `context.Context` to signal cancellation to both producers and consumers.
