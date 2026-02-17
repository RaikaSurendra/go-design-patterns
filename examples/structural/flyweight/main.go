// Package main demonstrates the Flyweight pattern.
//
// Flyweight minimizes memory usage by sharing intrinsic (common, immutable)
// state among many objects, while keeping extrinsic (unique) state separate.
package main

import "fmt"

// TreeType holds the intrinsic (shared) state.
type TreeType struct {
	Name    string
	Color   string
	Texture string
}

func (t *TreeType) Draw(x, y int) string {
	return fmt.Sprintf("Drawing '%s' tree (%s) at (%d, %d)", t.Name, t.Color, x, y)
}

// TreeFactory caches and reuses TreeType instances.
type TreeFactory struct {
	types map[string]*TreeType
}

func NewTreeFactory() *TreeFactory {
	return &TreeFactory{types: make(map[string]*TreeType)}
}

func (f *TreeFactory) GetTreeType(name, color, texture string) *TreeType {
	key := name + "_" + color + "_" + texture
	if t, ok := f.types[key]; ok {
		return t
	}
	t := &TreeType{Name: name, Color: color, Texture: texture}
	f.types[key] = t
	return t
}

// Tree holds the extrinsic (unique) state plus a reference to the shared type.
type Tree struct {
	X, Y     int
	TreeType *TreeType
}

func (t *Tree) Draw() string {
	return t.TreeType.Draw(t.X, t.Y)
}

func main() {
	factory := NewTreeFactory()

	// 4 trees but only 2 unique TreeType allocations.
	trees := []Tree{
		{X: 1, Y: 2, TreeType: factory.GetTreeType("Oak", "green", "rough")},
		{X: 5, Y: 3, TreeType: factory.GetTreeType("Oak", "green", "rough")},
		{X: 8, Y: 1, TreeType: factory.GetTreeType("Pine", "dark green", "smooth")},
		{X: 3, Y: 7, TreeType: factory.GetTreeType("Oak", "green", "rough")},
	}

	for _, t := range trees {
		fmt.Println(t.Draw())
	}

	fmt.Printf("\n4 trees created, only %d TreeType objects allocated\n", len(factory.types))
}
