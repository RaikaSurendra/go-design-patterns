// Package main demonstrates the Fan-In messaging pattern.
//
// Fan-In merges multiple input channels into a single output channel,
// creating a funnel from many producers to one consumer.
package main

import (
	"fmt"
	"sync"
)

func Merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	send := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func generate(start, count int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < count; i++ {
			ch <- start + i
		}
	}()
	return ch
}

func main() {
	// Three independent producers
	ch1 := generate(100, 3)
	ch2 := generate(200, 3)
	ch3 := generate(300, 3)

	// Merge into a single stream
	merged := Merge(ch1, ch2, ch3)

	fmt.Println("Merged output:")
	for v := range merged {
		fmt.Printf("  %d\n", v)
	}
}
