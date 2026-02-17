// Package main demonstrates the Generator pattern.
//
// A generator yields a sequence of values one at a time through
// a channel, producing values lazily on demand.
package main

import "fmt"

func Count(start, end int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func Fibonacci(n int) <-chan int {
	ch := make(chan int)
	go func() {
		a, b := 0, 1
		for i := 0; i < n; i++ {
			ch <- a
			a, b = b, a+b
		}
		close(ch)
	}()
	return ch
}

func main() {
	fmt.Println("--- Count 1 to 5 ---")
	for v := range Count(1, 5) {
		fmt.Println(v)
	}

	fmt.Println("\n--- First 8 Fibonacci numbers ---")
	for v := range Fibonacci(8) {
		fmt.Println(v)
	}
}
