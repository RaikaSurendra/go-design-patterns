// Package main demonstrates the Producer Consumer pattern.
//
// Producers push work items into a shared buffered channel.
// Consumers pull and process items independently.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	ID   int
	Data string
}

func main() {
	ch := make(chan Item, 5)

	// Start 2 producers
	var producerWg sync.WaitGroup
	for p := 0; p < 2; p++ {
		producerWg.Add(1)
		go func(id int) {
			defer producerWg.Done()
			for i := 0; i < 4; i++ {
				item := Item{
					ID:   id*100 + i,
					Data: fmt.Sprintf("item-%d from producer-%d", i, id),
				}
				ch <- item
				fmt.Printf("Produced: %s\n", item.Data)
			}
		}(p)
	}

	// Start 3 consumers
	var consumerWg sync.WaitGroup
	for c := 0; c < 3; c++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			for item := range ch {
				fmt.Printf("  Consumer %d processed: %s\n", id, item.Data)
				time.Sleep(50 * time.Millisecond)
			}
		}(c)
	}

	// Close channel after all producers finish
	producerWg.Wait()
	close(ch)

	// Wait for consumers to drain
	consumerWg.Wait()
	fmt.Println("\nAll items produced and consumed")
}
