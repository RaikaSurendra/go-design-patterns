// Package main demonstrates the Decorator pattern.
//
// Decorator extends the behavior of an existing function dynamically
// without altering its internals.
package main

import (
	"fmt"
	"log"
	"os"
)

type IntFunc func(int) int

func LogDecorate(fn IntFunc) IntFunc {
	return func(n int) int {
		log.Printf("Starting execution with the integer %d", n)
		result := fn(n)
		log.Printf("Execution completed with the result %d", result)
		return result
	}
}

func Double(n int) int {
	return n * 2
}

func Square(n int) int {
	return n * n
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	fmt.Println("--- Decorated Double ---")
	decorated := LogDecorate(Double)
	result := decorated(5)
	fmt.Printf("Result: %d\n", result)

	fmt.Println("\n--- Decorated Square ---")
	decorated = LogDecorate(Square)
	result = decorated(7)
	fmt.Printf("Result: %d\n", result)
}
