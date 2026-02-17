// Package main demonstrates the Fan-Out messaging pattern.
//
// Fan-Out distributes work from a single channel to multiple workers
// in a round-robin fashion.
package main

import (
	"fmt"
	"sync"
)

func Split(ch <-chan int, n int) []<-chan int {
	channels := make([]chan int, n)
	for i := 0; i < n; i++ {
		channels[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, c := range channels {
				close(c)
			}
		}()
		i := 0
		for val := range ch {
			channels[i%n] <- val
			i++
		}
	}()

	result := make([]<-chan int, n)
	for i, c := range channels {
		result[i] = c
	}
	return result
}

func main() {
	// Source channel
	source := make(chan int)
	go func() {
		defer close(source)
		for i := 1; i <= 9; i++ {
			source <- i
		}
	}()

	// Fan out to 3 workers
	workers := Split(source, 3)

	var wg sync.WaitGroup
	for i, ch := range workers {
		wg.Add(1)
		go func(id int, c <-chan int) {
			defer wg.Done()
			for val := range c {
				fmt.Printf("Worker %d received: %d\n", id, val)
			}
		}(i, ch)
	}

	wg.Wait()
}
