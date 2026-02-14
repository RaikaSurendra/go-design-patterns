# Reactor Pattern <span class="gp-difficulty gp-difficulty--hard">Hard</span>

The reactor pattern demultiplexes service requests delivered concurrently to a
service handler and dispatches them synchronously to the associated request
handlers. It uses a single-threaded event loop to monitor multiple event sources
and dispatches events to registered handlers as they arrive.

## Implementation

```go
package reactor

import "fmt"

// EventType identifies the kind of event.
type EventType string

// Handler processes a specific event type.
type Handler func(data interface{})

// Reactor is the event dispatcher. It registers handlers and dispatches
// incoming events to the appropriate handler synchronously.
type Reactor struct {
	handlers map[EventType]Handler
	events   chan Event
	quit     chan struct{}
}

// Event represents something that happened.
type Event struct {
	Type EventType
	Data interface{}
}

func New(bufferSize int) *Reactor {
	return &Reactor{
		handlers: make(map[EventType]Handler),
		events:   make(chan Event, bufferSize),
		quit:     make(chan struct{}),
	}
}

// Register associates a handler with an event type.
func (r *Reactor) Register(eventType EventType, handler Handler) {
	r.handlers[eventType] = handler
}

// Dispatch submits an event to the reactor's event queue.
func (r *Reactor) Dispatch(e Event) {
	r.events <- e
}

// Run starts the single-threaded event loop. It processes events sequentially
// until Stop is called.
func (r *Reactor) Run() {
	for {
		select {
		case e := <-r.events:
			if handler, ok := r.handlers[e.Type]; ok {
				handler(e.Data)
			} else {
				fmt.Printf("no handler for event type: %s\n", e.Type)
			}
		case <-r.quit:
			return
		}
	}
}

// Stop terminates the event loop.
func (r *Reactor) Stop() {
	close(r.quit)
}
```

## Usage

```go
r := reactor.New(100)

r.Register("connect", func(data interface{}) {
	fmt.Printf("client connected: %v\n", data)
})

r.Register("message", func(data interface{}) {
	fmt.Printf("received message: %v\n", data)
})

r.Register("disconnect", func(data interface{}) {
	fmt.Printf("client disconnected: %v\n", data)
})

go r.Run()

r.Dispatch(reactor.Event{Type: "connect", Data: "client-1"})
r.Dispatch(reactor.Event{Type: "message", Data: "hello"})
r.Dispatch(reactor.Event{Type: "disconnect", Data: "client-1"})

time.Sleep(100 * time.Millisecond)
r.Stop()
// client connected: client-1
// received message: hello
// client disconnected: client-1
```

## Rules of Thumb

- The reactor processes events synchronously in a single goroutine. If a handler blocks, it delays all subsequent events. Keep handlers fast.
- For CPU-intensive handlers, offload work to a goroutine pool and return immediately.
- Go's `net/http` server uses a variation of this pattern internally: it accepts connections on a single listener and dispatches each to a handler goroutine.
