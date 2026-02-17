// Package main demonstrates the Push & Pull pattern.
//
// A pusher sends work downstream, multiple pullers process it in parallel,
// and a sink collects all results into a single stream.
package main

import (
	"fmt"
	"strings"
	"sync"
)

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

func Puller(id int, in <-chan string, process func(string) string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for item := range in {
			out <- fmt.Sprintf("[worker-%d] %s", id, process(item))
		}
	}()
	return out
}

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

func main() {
	work := Pusher([]string{"task-a", "task-b", "task-c", "task-d", "task-e", "task-f"})

	process := func(s string) string {
		return "done: " + strings.ToUpper(s)
	}

	puller1 := Puller(1, work, process)
	puller2 := Puller(2, work, process)

	for result := range Sink(puller1, puller2) {
		fmt.Println(result)
	}
}
