// Package main demonstrates the Memento pattern.
//
// Memento captures and externalizes an object's internal state so that
// it can be restored later, commonly used for undo functionality.
package main

import "fmt"

// Memento stores a snapshot of the editor's state.
type Memento struct {
	content string
}

// Editor is the originator whose state we want to save and restore.
type Editor struct {
	content string
}

func (e *Editor) Type(text string) {
	e.content += text
}

func (e *Editor) Content() string {
	return e.content
}

func (e *Editor) Save() *Memento {
	return &Memento{content: e.content}
}

func (e *Editor) Restore(m *Memento) {
	e.content = m.content
}

// History is the caretaker that stores mementos.
type History struct {
	snapshots []*Memento
}

func (h *History) Push(m *Memento) {
	h.snapshots = append(h.snapshots, m)
}

func (h *History) Pop() *Memento {
	if len(h.snapshots) == 0 {
		return nil
	}
	last := h.snapshots[len(h.snapshots)-1]
	h.snapshots = h.snapshots[:len(h.snapshots)-1]
	return last
}

func main() {
	editor := &Editor{}
	history := &History{}

	editor.Type("Hello, ")
	history.Push(editor.Save())
	fmt.Printf("Typed: %q\n", editor.Content())

	editor.Type("World!")
	history.Push(editor.Save())
	fmt.Printf("Typed: %q\n", editor.Content())

	editor.Type(" Extra text.")
	fmt.Printf("Typed: %q\n", editor.Content())

	fmt.Println("\n--- Undo ---")
	editor.Restore(history.Pop())
	fmt.Printf("After undo: %q\n", editor.Content())

	editor.Restore(history.Pop())
	fmt.Printf("After undo: %q\n", editor.Content())
}
