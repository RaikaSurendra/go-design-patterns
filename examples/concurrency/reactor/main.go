// Package main demonstrates the Reactor pattern.
//
// The reactor uses a single-threaded event loop to demultiplex and
// dispatch events to registered handlers synchronously.
package main

import (
	"fmt"
	"time"
)

type EventType string

type Handler func(data interface{})

type Event struct {
	Type EventType
	Data interface{}
}

type Reactor struct {
	handlers map[EventType]Handler
	events   chan Event
	quit     chan struct{}
}

func NewReactor(bufferSize int) *Reactor {
	return &Reactor{
		handlers: make(map[EventType]Handler),
		events:   make(chan Event, bufferSize),
		quit:     make(chan struct{}),
	}
}

func (r *Reactor) Register(eventType EventType, handler Handler) {
	r.handlers[eventType] = handler
}

func (r *Reactor) Dispatch(e Event) {
	r.events <- e
}

func (r *Reactor) Run() {
	for {
		select {
		case e := <-r.events:
			if handler, ok := r.handlers[e.Type]; ok {
				handler(e.Data)
			} else {
				fmt.Printf("No handler for event: %s\n", e.Type)
			}
		case <-r.quit:
			return
		}
	}
}

func (r *Reactor) Stop() {
	close(r.quit)
}

func main() {
	r := NewReactor(100)

	r.Register("connect", func(data interface{}) {
		fmt.Printf("Client connected: %v\n", data)
	})
	r.Register("message", func(data interface{}) {
		fmt.Printf("Received message: %v\n", data)
	})
	r.Register("disconnect", func(data interface{}) {
		fmt.Printf("Client disconnected: %v\n", data)
	})

	go r.Run()

	r.Dispatch(Event{Type: "connect", Data: "client-1"})
	r.Dispatch(Event{Type: "message", Data: "hello from client-1"})
	r.Dispatch(Event{Type: "connect", Data: "client-2"})
	r.Dispatch(Event{Type: "message", Data: "hello from client-2"})
	r.Dispatch(Event{Type: "disconnect", Data: "client-1"})

	time.Sleep(100 * time.Millisecond)
	r.Stop()
}
