# Bridge Pattern <span class="gp-difficulty gp-difficulty--medium">Medium</span>

The bridge pattern decouples an abstraction from its implementation so that the
two can vary independently. Instead of combining every abstraction variant with
every implementation variant (leading to a class explosion), the bridge composes
them at runtime through an interface.

## Implementation

```go
package bridge

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
```

## Usage

```go
vector := &bridge.VectorRenderer{}
raster := &bridge.RasterRenderer{}

circle := bridge.NewCircle(vector, 5)
fmt.Println(circle.Draw())
// Drawing circle with radius 5.0 as vector

circle = bridge.NewCircle(raster, 5)
fmt.Println(circle.Draw())
// Drawing circle with radius 5.0 as pixels
```

## Rules of Thumb

- Use bridge when you want to avoid a permanent binding between an abstraction and its implementation.
- Bridge is designed up-front to let the abstraction and implementation vary independently; adapter is applied after the fact to make unrelated classes work together.
- The abstraction side and the implementation side can be extended independently through composition.
