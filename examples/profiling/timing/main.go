// Package main demonstrates the Timing Functions pattern.
//
// A simple defer-based approach to measure function execution time
// without a full profiling framework.
package main

import (
	"fmt"
	"math/big"
	"time"
)

func Duration(invocation time.Time, name string) {
	elapsed := time.Since(invocation)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func Factorial(n int64) *big.Int {
	defer Duration(time.Now(), fmt.Sprintf("Factorial(%d)", n))

	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

func Fibonacci(n int) int {
	defer Duration(time.Now(), fmt.Sprintf("Fibonacci(%d)", n))

	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	f := Factorial(1000)
	fmt.Printf("1000! has %d digits\n\n", len(f.String()))

	result := Fibonacci(50)
	fmt.Printf("Fibonacci(50) = %d\n", result)
}
