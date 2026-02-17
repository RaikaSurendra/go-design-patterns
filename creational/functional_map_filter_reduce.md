# Map / Filter / Reduce <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Go generics (1.18+) make it possible to write type-safe functional
transformation functions — `Map`, `Filter`, and `Reduce` — that work with any
slice type. These composable building blocks eliminate repetitive `for` loops
for common data transformations.

## Implementation

```go
package fn

// Map applies a function to every element of a slice, returning a new slice.
func Map[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](items []T, pred func(T) bool) []T {
	var result []T
	for _, item := range items {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce collapses a slice into a single value using an accumulator function.
func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, item := range items {
		acc = fn(acc, item)
	}
	return acc
}
```

## Usage

```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Double every number.
doubled := fn.Map(numbers, func(n int) int { return n * 2 })
// [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

// Keep only even numbers.
evens := fn.Filter(numbers, func(n int) bool { return n%2 == 0 })
// [2, 4, 6, 8, 10]

// Sum all numbers.
sum := fn.Reduce(numbers, 0, func(acc, n int) int { return acc + n })
// 55

// Compose: sum of squares of even numbers.
result := fn.Reduce(
	fn.Map(
		fn.Filter(numbers, func(n int) bool { return n%2 == 0 }),
		func(n int) int { return n * n },
	),
	0,
	func(acc, n int) int { return acc + n },
)
// 220 (4 + 16 + 36 + 64 + 100)
```

## Rules of Thumb

- These functions create new slices — they don't mutate the input. This is safe but allocates.
- For performance-critical hot loops, a plain `for` loop is still faster (no function call overhead per element).
- `Map` can change types (e.g. `[]User` → `[]string` of names), making it more versatile than a simple loop.
- Prefer readability: if the chain gets deeply nested, break it into intermediate variables.
