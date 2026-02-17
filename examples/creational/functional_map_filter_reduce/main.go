// Package main demonstrates the Map / Filter / Reduce pattern with generics.
//
// Type-safe functional transformations that work with any slice type,
// eliminating repetitive for loops for common data operations.
package main

import "fmt"

func Map[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

func Filter[T any](items []T, pred func(T) bool) []T {
	var result []T
	for _, item := range items {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}

func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, item := range items {
		acc = fn(acc, item)
	}
	return acc
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	doubled := Map(numbers, func(n int) int { return n * 2 })
	fmt.Println("Doubled:", doubled)

	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:  ", evens)

	sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
	fmt.Println("Sum:    ", sum)

	// Composed: sum of squares of even numbers
	result := Reduce(
		Map(
			Filter(numbers, func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n },
		),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Println("Sum of squares of evens:", result)

	// Type conversion: []int → []string
	labels := Map(numbers, func(n int) string {
		return fmt.Sprintf("item-%d", n)
	})
	fmt.Println("Labels: ", labels)
}
