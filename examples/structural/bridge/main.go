// Package main demonstrates the Bridge pattern.
//
// Bridge decouples an abstraction from its implementation so that
// the two can vary independently.
package main

import "fmt"

// Renderer is the implementation interface.
type Renderer interface {
	RenderCircle(radius float64) string
}

// VectorRenderer draws shapes as vector graphics.
type VectorRenderer struct{}

func (v *VectorRenderer) RenderCircle(radius float64) string {
	return fmt.Sprintf("Drawing circle with radius %.1f as vector", radius)
}

// RasterRenderer draws shapes as pixels.
type RasterRenderer struct{}

func (r *RasterRenderer) RenderCircle(radius float64) string {
	return fmt.Sprintf("Drawing circle with radius %.1f as pixels", radius)
}

// Shape is the abstraction that delegates to a Renderer.
type Shape struct {
	renderer Renderer
}

// Circle extends Shape with circle-specific data.
type Circle struct {
	Shape
	radius float64
}

func NewCircle(renderer Renderer, radius float64) *Circle {
	return &Circle{
		Shape:  Shape{renderer: renderer},
		radius: radius,
	}
}

func (c *Circle) Draw() string {
	return c.renderer.RenderCircle(c.radius)
}

func (c *Circle) Resize(factor float64) {
	c.radius *= factor
}

func main() {
	vector := &VectorRenderer{}
	raster := &RasterRenderer{}

	circle := NewCircle(vector, 5)
	fmt.Println(circle.Draw())

	circle = NewCircle(raster, 5)
	fmt.Println(circle.Draw())

	circle.Resize(2)
	fmt.Println(circle.Draw())
}
