# Registry Pattern <span class="gp-difficulty gp-difficulty--easy">Easy</span>

The registry pattern provides a well-known object that other objects can use to
find common objects and services. It acts as a central lookup table where
implementations are registered by name or type and retrieved when needed.

## Implementation

```go
package registry

import (
	"fmt"
	"sync"
)

// Service is the common interface for registered services.
type Service interface {
	Name() string
	Execute() string
}

// Registry is a thread-safe service locator.
type Registry struct {
	mu       sync.RWMutex
	services map[string]Service
}

func New() *Registry {
	return &Registry{
		services: make(map[string]Service),
	}
}

func (r *Registry) Register(svc Service) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services[svc.Name()] = svc
}

func (r *Registry) Lookup(name string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc, ok := r.services[name]
	if !ok {
		return nil, fmt.Errorf("service %q not found", name)
	}
	return svc, nil
}

func (r *Registry) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services, name)
}
```

## Usage

```go
type EmailService struct{}

func (e *EmailService) Name() string    { return "email" }
func (e *EmailService) Execute() string { return "sending email" }

type SMSService struct{}

func (s *SMSService) Name() string    { return "sms" }
func (s *SMSService) Execute() string { return "sending sms" }

r := registry.New()
r.Register(&EmailService{})
r.Register(&SMSService{})

svc, err := r.Lookup("email")
if err == nil {
	fmt.Println(svc.Execute()) // sending email
}

svc, err = r.Lookup("sms")
if err == nil {
	fmt.Println(svc.Execute()) // sending sms
}
```

## Rules of Thumb

- Registry provides a global access point — use it sparingly to avoid hidden dependencies that make testing difficult.
- Always make the registry thread-safe if it will be accessed from multiple goroutines.
- Consider using Go's `init()` functions to self-register implementations at startup.
