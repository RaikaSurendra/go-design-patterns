# Visitor Pattern <span class="gp-difficulty gp-difficulty--hard">Hard</span>

The visitor pattern separates an algorithm from the object structure on which it
operates. It allows you to add new operations to existing object structures
without modifying the structures themselves.

## Implementation

```go
package visitor

import "fmt"

// Shape is the element interface that accepts a visitor.
type Shape interface {
	Accept(Visitor) string
}

// Visitor defines operations for each concrete element type.
type Visitor interface {
	VisitCircle(*Circle) string
	VisitRectangle(*Rectangle) string
}

// --- Concrete elements ---

type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v Visitor) string {
	return v.VisitCircle(c)
}

type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) Accept(v Visitor) string {
	return v.VisitRectangle(r)
}

// --- Concrete visitors ---

type AreaCalculator struct{}

func (a *AreaCalculator) VisitCircle(c *Circle) string {
	area := 3.14159 * c.Radius * c.Radius
	return fmt.Sprintf("Circle area: %.2f", area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) string {
	area := r.Width * r.Height
	return fmt.Sprintf("Rectangle area: %.2f", area)
}

type PerimeterCalculator struct{}

func (p *PerimeterCalculator) VisitCircle(c *Circle) string {
	perim := 2 * 3.14159 * c.Radius
	return fmt.Sprintf("Circle perimeter: %.2f", perim)
}

func (p *PerimeterCalculator) VisitRectangle(r *Rectangle) string {
	perim := 2 * (r.Width + r.Height)
	return fmt.Sprintf("Rectangle perimeter: %.2f", perim)
}
```

## Usage

```go
shapes := []visitor.Shape{
	&visitor.Circle{Radius: 5},
	&visitor.Rectangle{Width: 3, Height: 4},
}

area := &visitor.AreaCalculator{}
perim := &visitor.PerimeterCalculator{}

for _, s := range shapes {
	fmt.Println(s.Accept(area))
	fmt.Println(s.Accept(perim))
}
// Circle area: 78.54
// Circle perimeter: 31.42
// Rectangle area: 12.00
// Rectangle perimeter: 14.00
```

## Rules of Thumb

- Use visitor when you need to perform many unrelated operations on an object structure and you don't want to pollute the element classes with these operations.
- Adding a new element type requires updating every visitor — the pattern works best when the element hierarchy is stable.
- Visitor can accumulate state as it traverses a structure, making it useful for compilers, serializers, and report generators.
