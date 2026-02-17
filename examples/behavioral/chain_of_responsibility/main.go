// Package main demonstrates the Chain of Responsibility pattern.
//
// Handlers are chained together. A request passes along the chain until
// a handler processes it or the chain ends.
package main

import "fmt"

// Request represents the data flowing through the chain.
type Request struct {
	Amount float64
}

// Handler defines the interface for a link in the chain.
type Handler interface {
	SetNext(Handler) Handler
	Handle(Request) string
}

// BaseHandler provides default chaining behavior.
type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(h Handler) Handler {
	b.next = h
	return h
}

func (b *BaseHandler) HandleNext(r Request) string {
	if b.next != nil {
		return b.next.Handle(r)
	}
	return "no handler approved the request"
}

// --- Concrete handlers ---

type Manager struct{ BaseHandler }

func (m *Manager) Handle(r Request) string {
	if r.Amount < 1000 {
		return "Manager approved"
	}
	return m.HandleNext(r)
}

type Director struct{ BaseHandler }

func (d *Director) Handle(r Request) string {
	if r.Amount < 5000 {
		return "Director approved"
	}
	return d.HandleNext(r)
}

type VP struct{ BaseHandler }

func (v *VP) Handle(r Request) string {
	if r.Amount < 10000 {
		return "VP approved"
	}
	return v.HandleNext(r)
}

func main() {
	manager := &Manager{}
	director := &Director{}
	vp := &VP{}

	manager.SetNext(director).SetNext(vp)

	requests := []Request{
		{Amount: 500},
		{Amount: 3000},
		{Amount: 8000},
		{Amount: 50000},
	}

	for _, r := range requests {
		fmt.Printf("$%.0f -> %s\n", r.Amount, manager.Handle(r))
	}
}
