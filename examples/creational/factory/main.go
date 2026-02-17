// Package main demonstrates the Factory Method pattern.
//
// Factory method allows creating objects without specifying the exact
// type, using a creation function that returns an interface.
package main

import (
	"bytes"
	"fmt"
	"io"
)

// Store is the product interface.
type Store interface {
	Open(name string) (io.ReadWriteCloser, error)
}

// StorageType selects the backend.
type StorageType int

const (
	DiskStorage StorageType = 1 << iota
	TempStorage
	MemoryStorage
)

// NewStore is the factory method.
func NewStore(t StorageType) Store {
	switch t {
	case MemoryStorage:
		return &memoryStore{files: make(map[string]*memoryFile)}
	case DiskStorage:
		return &memoryStore{files: make(map[string]*memoryFile), label: "disk"}
	default:
		return &memoryStore{files: make(map[string]*memoryFile), label: "temp"}
	}
}

// --- In-memory implementation ---

type memoryStore struct {
	files map[string]*memoryFile
	label string
}

func (s *memoryStore) Open(name string) (io.ReadWriteCloser, error) {
	if f, ok := s.files[name]; ok {
		return f, nil
	}
	f := &memoryFile{name: name, buf: &bytes.Buffer{}}
	s.files[name] = f
	return f, nil
}

type memoryFile struct {
	name string
	buf  *bytes.Buffer
}

func (f *memoryFile) Write(p []byte) (int, error) { return f.buf.Write(p) }
func (f *memoryFile) Read(p []byte) (int, error)  { return f.buf.Read(p) }
func (f *memoryFile) Close() error                 { return nil }

func main() {
	store := NewStore(MemoryStorage)

	// Write data
	f, _ := store.Open("greeting.txt")
	n, _ := f.Write([]byte("Hello from the factory!"))
	fmt.Printf("Wrote %d bytes to greeting.txt\n", n)
	f.Close()

	// Read it back
	f, _ = store.Open("greeting.txt")
	data := make([]byte, 64)
	n, _ = f.Read(data)
	fmt.Printf("Read back: %s\n", data[:n])
	f.Close()
}
