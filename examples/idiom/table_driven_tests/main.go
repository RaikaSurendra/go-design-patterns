// Package main demonstrates Table-Driven Tests.
//
// Instead of separate test functions per scenario, define a table of
// inputs and expected outputs, then loop over it. This example runs
// the tests manually (without `go test`) to show the concept.
package main

import "fmt"

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func main() {
	fmt.Println("--- TestAbs ---")
	absTests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 5, 5},
		{"negative", -3, 3},
		{"zero", 0, 0},
		{"large negative", -999, 999},
	}

	for _, tt := range absTests {
		got := Abs(tt.input)
		status := "PASS"
		if got != tt.want {
			status = fmt.Sprintf("FAIL (got %d, want %d)", got, tt.want)
		}
		fmt.Printf("  %-20s %s\n", tt.name, status)
	}

	fmt.Println("\n--- TestClamp ---")
	clampTests := []struct {
		name          string
		val, min, max int
		want          int
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -3, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min boundary", 0, 0, 10, 0},
		{"at max boundary", 10, 0, 10, 10},
	}

	for _, tt := range clampTests {
		got := Clamp(tt.val, tt.min, tt.max)
		status := "PASS"
		if got != tt.want {
			status = fmt.Sprintf("FAIL (got %d, want %d)", got, tt.want)
		}
		fmt.Printf("  %-20s %s\n", tt.name, status)
	}
}
