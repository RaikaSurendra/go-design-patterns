# Dependency Injection <span class="gp-difficulty gp-difficulty--easy">Easy</span>

Dependency injection (DI) in Go means passing dependencies (usually as
interfaces) into a struct or function rather than having them create their own.
This makes code testable, composable, and decoupled from concrete
implementations. No framework needed — Go interfaces and constructors are
sufficient.

## Implementation

```go
package order

// Repository is the dependency interface — any storage backend can satisfy it.
type Repository interface {
	Save(o Order) error
	FindByID(id string) (Order, error)
}

// Notifier is another dependency.
type Notifier interface {
	Notify(userID, message string) error
}

// Service is the business logic layer. Dependencies are injected via the constructor.
type Service struct {
	repo     Repository
	notifier Notifier
}

func NewService(repo Repository, notifier Notifier) *Service {
	return &Service{repo: repo, notifier: notifier}
}

func (s *Service) PlaceOrder(o Order) error {
	if err := s.repo.Save(o); err != nil {
		return fmt.Errorf("save order: %w", err)
	}
	return s.notifier.Notify(o.UserID, "Your order has been placed")
}
```

## Usage

```go
// Production wiring — real implementations.
db := postgres.NewOrderRepo(connStr)
email := smtp.NewNotifier(smtpHost)
svc := order.NewService(db, email)

// Test wiring — mock implementations.
func TestPlaceOrder(t *testing.T) {
	mockRepo := &MockRepo{SaveFn: func(o Order) error { return nil }}
	mockNotify := &MockNotifier{NotifyFn: func(uid, msg string) error { return nil }}

	svc := order.NewService(mockRepo, mockNotify)
	err := svc.PlaceOrder(Order{ID: "1", UserID: "alice"})
	if err != nil {
		t.Fatal(err)
	}
}
```

## Rules of Thumb

- Accept interfaces, return structs. Define the interface where it's consumed (not where it's implemented).
- Inject dependencies through constructors (`NewService(...)`) — avoid global state and init-time wiring.
- Keep interfaces small. One or two methods is ideal — Go's implicit interface satisfaction makes this natural.
- Don't reach for a DI framework. Constructor injection + interfaces covers 99% of Go use cases.
