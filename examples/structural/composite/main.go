// Package main demonstrates the Composite pattern.
//
// Composite composes objects into tree structures to represent part-whole
// hierarchies, letting clients treat individual objects and compositions uniformly.
package main

import "fmt"

// Component is the common interface for leaf and composite nodes.
type Component interface {
	Search(keyword string)
}

// File is a leaf node.
type File struct {
	Name string
}

func (f *File) Search(keyword string) {
	fmt.Printf("Searching for '%s' in file: %s\n", keyword, f.Name)
}

// Folder is a composite node that can contain other components.
type Folder struct {
	Name       string
	Components []Component
}

func (f *Folder) Search(keyword string) {
	fmt.Printf("Searching for '%s' in folder: %s\n", keyword, f.Name)
	for _, c := range f.Components {
		c.Search(keyword)
	}
}

func (f *Folder) Add(c Component) {
	f.Components = append(f.Components, c)
}

func main() {
	file1 := &File{Name: "main.go"}
	file2 := &File{Name: "utils.go"}
	file3 := &File{Name: "readme.md"}

	src := &Folder{Name: "src"}
	src.Add(file1)
	src.Add(file2)

	root := &Folder{Name: "project"}
	root.Add(src)
	root.Add(file3)

	// Treats files and folders uniformly.
	root.Search("pattern")
}
