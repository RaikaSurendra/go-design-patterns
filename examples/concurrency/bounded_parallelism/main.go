// Package main demonstrates the Bounded Parallelism pattern.
//
// Bounded parallelism limits the number of goroutines processing work
// concurrently, preventing resource exhaustion.
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := []string{
		"download-A", "download-B", "download-C",
		"download-D", "download-E", "download-F",
		"download-G", "download-H",
	}

	maxWorkers := 3
	sem := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		sem <- struct{}{} // acquire slot (blocks if pool is full)

		go func(t string) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			fmt.Printf("Started:  %s\n", t)
			time.Sleep(150 * time.Millisecond) // simulate work
			fmt.Printf("Finished: %s\n", t)
		}(task)
	}

	wg.Wait()
	fmt.Printf("\nAll %d tasks completed (max %d concurrent)\n", len(tasks), maxWorkers)
}
