// Package main demonstrates the Parallelism pattern.
//
// Multiple independent tasks run concurrently in separate goroutines,
// with results collected via a WaitGroup.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	TaskID int
	Output string
}

func main() {
	tasks := []string{
		"compress file",
		"resize image",
		"parse CSV",
		"encrypt data",
	}

	var mu sync.Mutex
	var results []Result
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(id int, name string) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond) // simulate work
			r := Result{TaskID: id, Output: fmt.Sprintf("completed: %s", name)}

			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(i, task)
	}

	wg.Wait()

	fmt.Println("All tasks finished:")
	for _, r := range results {
		fmt.Printf("  Task %d: %s\n", r.TaskID, r.Output)
	}
}
