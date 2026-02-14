# Chain of Responsibility Pattern

The chain of responsibility pattern avoids coupling the sender of a request to
its receiver by giving more than one object a chance to handle the request.
Handlers are chained together, and the request is passed along the chain until a
handler processes it or the chain ends.

## Implementation

```go
package chain

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
```

## Usage

```go
manager := &chain.Manager{}
director := &chain.Director{}
vp := &chain.VP{}

manager.SetNext(director).SetNext(vp)

fmt.Println(manager.Handle(chain.Request{Amount: 500}))
// Manager approved

fmt.Println(manager.Handle(chain.Request{Amount: 3000}))
// Director approved

fmt.Println(manager.Handle(chain.Request{Amount: 8000}))
// VP approved

fmt.Println(manager.Handle(chain.Request{Amount: 50000}))
// no handler approved the request
```

## Rules of Thumb

- Use chain of responsibility when more than one object may handle a request, and the handler is determined at runtime.
- The chain can be composed dynamically at runtime, making it easy to add, remove, or reorder handlers.
- HTTP middleware stacks (e.g. in `net/http`) are a common real-world example of this pattern.
