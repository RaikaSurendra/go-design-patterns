# Table-Driven Tests <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Table-driven testing is Go's idiomatic approach to writing test cases. Instead
of separate test functions per scenario, you define a table (slice of structs)
where each entry describes an input/expected-output pair, then loop over the
table. This reduces duplication and makes it trivial to add new test cases.

## Implementation

```go
package mathutil

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
```

## Usage

```go
package mathutil_test

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		input int
		want  int
	}{
		{"positive", 5, 5},
		{"negative", -3, 3},
		{"zero", 0, 0},
		{"min int edge", -1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Abs(tt.input)
			if got != tt.want {
				t.Errorf("Abs(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name            string
		val, min, max   int
		want            int
	}{
		{"within range", 5, 0, 10, 5},
		{"below min", -3, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clamp(tt.val, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("Clamp(%d, %d, %d) = %d, want %d",
					tt.val, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
```

## Rules of Thumb

- Always use `t.Run(tt.name, ...)` to create subtests — this gives each case its own name in test output and allows running individual cases with `-run`.
- Name the test struct variable `tt` (or `tc`) and the slice `tests` — this is the community convention.
- Include both typical and edge-case inputs in the table.
- For error-returning functions, add a `wantErr bool` or `wantErr error` field to the struct.
- Table-driven tests work well with `t.Parallel()` — add it inside `t.Run` for concurrent test execution.
