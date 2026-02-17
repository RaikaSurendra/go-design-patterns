// Package main demonstrates the Visitor pattern.
//
// Visitor separates an algorithm from the object structure it operates on,
// allowing new operations without modifying the structures.
package main

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

func main() {
	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 3, Height: 4},
	}

	area := &AreaCalculator{}
	perim := &PerimeterCalculator{}

	for _, s := range shapes {
		fmt.Println(s.Accept(area))
		fmt.Println(s.Accept(perim))
		fmt.Println()
	}
}
