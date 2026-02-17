// Package main demonstrates the Pipeline pattern.
//
// A pipeline composes data processing from a series of stages connected
// by channels. Each stage transforms data and passes it downstream.
package main

import "fmt"

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

func Collect[T any](in <-chan T) []T {
	var result []T
	for v := range in {
		result = append(result, v)
	}
	return result
}

func main() {
	// Pipeline: generate 1..10 → square → keep even results
	numbers := Generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := Map(numbers, func(n int) int { return n * n })
	even := Filter(squared, func(n int) bool { return n%2 == 0 })

	fmt.Println("Squares of 1-10, keeping even results:")
	for _, v := range Collect(even) {
		fmt.Printf("  %d\n", v)
	}

	// String pipeline: generate words → uppercase → filter long
	words := Generate("go", "patterns", "are", "awesome", "hi")
	upper := Map(words, func(s string) string {
		result := make([]byte, len(s))
		for i, c := range s {
			if c >= 'a' && c <= 'z' {
				result[i] = byte(c - 32)
			} else {
				result[i] = byte(c)
			}
		}
		return string(result)
	})
	long := Filter(upper, func(s string) bool { return len(s) > 2 })

	fmt.Println("\nUppercased words longer than 2 chars:")
	for _, v := range Collect(long) {
		fmt.Printf("  %s\n", v)
	}
}
