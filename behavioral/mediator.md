# Mediator Pattern

The mediator pattern defines an object that encapsulates how a set of objects
interact. It promotes loose coupling by preventing objects from referring to each
other directly, forcing them to communicate through the mediator instead.

## Implementation

```go
package mediator

import "fmt"

// Mediator defines the communication interface.
type Mediator interface {
	Notify(sender Component, event string)
}

// Component is the base for all participants.
type Component struct {
	mediator Mediator
}

func (c *Component) SetMediator(m Mediator) {
	c.mediator = m
}

// --- Concrete components ---

type AuthService struct{ Component }

func (a *AuthService) Login(user string) string {
	msg := fmt.Sprintf("%s logged in", user)
	a.mediator.Notify(a, "login")
	return msg
}

type Logger struct{ Component }

func (l *Logger) Log(message string) string {
	return fmt.Sprintf("LOG: %s", message)
}

type Notifier struct{ Component }

func (n *Notifier) Send(message string) string {
	return fmt.Sprintf("NOTIFY: %s", message)
}

// --- Concrete mediator ---

type AppMediator struct {
	auth     *AuthService
	logger   *Logger
	notifier *Notifier
}

func NewAppMediator() *AppMediator {
	m := &AppMediator{
		auth:     &AuthService{},
		logger:   &Logger{},
		notifier: &Notifier{},
	}
	m.auth.SetMediator(m)
	m.logger.SetMediator(m)
	m.notifier.SetMediator(m)
	return m
}

func (m *AppMediator) Notify(sender Component, event string) {
	switch event {
	case "login":
		m.logger.Log("user login event")
		m.notifier.Send("welcome back!")
	}
}
```

## Usage

```go
app := mediator.NewAppMediator()

fmt.Println(app.auth.Login("alice"))
// alice logged in
// (mediator also triggers logger and notifier behind the scenes)
```

## Rules of Thumb

- Use mediator when the communication logic between objects is complex and you want to centralize it in one place.
- Mediator trades complexity in individual components for complexity in the mediator itself — avoid creating a "god object."
- In Go, channels can serve as a lightweight mediator between goroutines.
