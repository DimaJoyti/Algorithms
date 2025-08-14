# 🌊 Event-Driven Architecture

## 🎯 Event-Driven Architecture Fundamentals

### Core Concepts
```go
// Event - Something that happened in the past
type Event struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Source      string                 `json:"source"`
    Data        map[string]interface{} `json:"data"`
    Timestamp   time.Time              `json:"timestamp"`
    Version     string                 `json:"version"`
    CorrelationID string               `json:"correlation_id"`
}

// Command - Request to do something
type Command struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Target    string                 `json:"target"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
    UserID    string                 `json:"user_id"`
}

// Query - Request for information
type Query struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Filter    map[string]interface{} `json:"filter"`
    Timestamp time.Time              `json:"timestamp"`
}
```

### Event Types
```go
// Domain Events - Business-meaningful events
type UserRegisteredEvent struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}

type OrderPlacedEvent struct {
    OrderID     string    `json:"order_id"`
    UserID      string    `json:"user_id"`
    Items       []OrderItem `json:"items"`
    TotalAmount decimal.Decimal `json:"total_amount"`
    Timestamp   time.Time `json:"timestamp"`
}

// Integration Events - Cross-service communication
type PaymentProcessedEvent struct {
    PaymentID     string          `json:"payment_id"`
    OrderID       string          `json:"order_id"`
    Amount        decimal.Decimal `json:"amount"`
    Status        string          `json:"status"`
    Timestamp     time.Time       `json:"timestamp"`
}

// System Events - Technical events
type ServiceStartedEvent struct {
    ServiceName string    `json:"service_name"`
    Version     string    `json:"version"`
    InstanceID  string    `json:"instance_id"`
    Timestamp   time.Time `json:"timestamp"`
}
```

## 🏗️ Event Sourcing Pattern

### Event Store Implementation
```go
type EventStore interface {
    AppendEvents(streamID string, expectedVersion int, events []Event) error
    GetEvents(streamID string, fromVersion int) ([]Event, error)
    GetAllEvents(fromTimestamp time.Time) ([]Event, error)
}

type PostgreSQLEventStore struct {
    db *sql.DB
}

func (es *PostgreSQLEventStore) AppendEvents(streamID string, expectedVersion int, events []Event) error {
    tx, err := es.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Check current version for optimistic concurrency control
    var currentVersion int
    err = tx.QueryRow("SELECT COALESCE(MAX(version), 0) FROM events WHERE stream_id = $1", streamID).Scan(&currentVersion)
    if err != nil {
        return err
    }
    
    if currentVersion != expectedVersion {
        return ErrConcurrencyConflict
    }
    
    // Insert events
    for i, event := range events {
        eventData, err := json.Marshal(event.Data)
        if err != nil {
            return err
        }
        
        _, err = tx.Exec(`
            INSERT INTO events (stream_id, version, event_type, event_data, timestamp, correlation_id)
            VALUES ($1, $2, $3, $4, $5, $6)`,
            streamID, expectedVersion+i+1, event.Type, eventData, event.Timestamp, event.CorrelationID)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (es *PostgreSQLEventStore) GetEvents(streamID string, fromVersion int) ([]Event, error) {
    rows, err := es.db.Query(`
        SELECT event_id, event_type, event_data, timestamp, correlation_id, version
        FROM events 
        WHERE stream_id = $1 AND version > $2 
        ORDER BY version ASC`, streamID, fromVersion)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var events []Event
    for rows.Next() {
        var event Event
        var eventData []byte
        
        err := rows.Scan(&event.ID, &event.Type, &eventData, &event.Timestamp, &event.CorrelationID, &event.Version)
        if err != nil {
            return nil, err
        }
        
        if err := json.Unmarshal(eventData, &event.Data); err != nil {
            return nil, err
        }
        
        events = append(events, event)
    }
    
    return events, nil
}
```

### Aggregate Root with Event Sourcing
```go
type UserAggregate struct {
    ID       string
    Email    string
    Name     string
    Status   UserStatus
    Version  int
    
    // Uncommitted events
    changes []Event
}

func (u *UserAggregate) RegisterUser(email, name string) error {
    // Business logic validation
    if u.ID != "" {
        return ErrUserAlreadyExists
    }
    
    if !isValidEmail(email) {
        return ErrInvalidEmail
    }
    
    // Apply event
    event := UserRegisteredEvent{
        UserID:    generateUserID(),
        Email:     email,
        Name:      name,
        Timestamp: time.Now(),
    }
    
    u.applyEvent(Event{
        ID:   generateEventID(),
        Type: "UserRegistered",
        Data: map[string]interface{}{
            "user_id": event.UserID,
            "email":   event.Email,
            "name":    event.Name,
        },
        Timestamp: event.Timestamp,
    })
    
    return nil
}

func (u *UserAggregate) ChangeEmail(newEmail string) error {
    if u.Status != UserStatusActive {
        return ErrUserNotActive
    }
    
    if !isValidEmail(newEmail) {
        return ErrInvalidEmail
    }
    
    if u.Email == newEmail {
        return nil // No change needed
    }
    
    event := UserEmailChangedEvent{
        UserID:   u.ID,
        OldEmail: u.Email,
        NewEmail: newEmail,
        Timestamp: time.Now(),
    }
    
    u.applyEvent(Event{
        ID:   generateEventID(),
        Type: "UserEmailChanged",
        Data: map[string]interface{}{
            "user_id":   event.UserID,
            "old_email": event.OldEmail,
            "new_email": event.NewEmail,
        },
        Timestamp: event.Timestamp,
    })
    
    return nil
}

func (u *UserAggregate) applyEvent(event Event) {
    switch event.Type {
    case "UserRegistered":
        u.ID = event.Data["user_id"].(string)
        u.Email = event.Data["email"].(string)
        u.Name = event.Data["name"].(string)
        u.Status = UserStatusActive
        
    case "UserEmailChanged":
        u.Email = event.Data["new_email"].(string)
        
    case "UserDeactivated":
        u.Status = UserStatusInactive
    }
    
    u.Version++
    u.changes = append(u.changes, event)
}

// Replay events to rebuild state
func (u *UserAggregate) LoadFromHistory(events []Event) {
    for _, event := range events {
        u.applyEvent(event)
    }
    u.changes = nil // Clear uncommitted changes
}

func (u *UserAggregate) GetUncommittedChanges() []Event {
    return u.changes
}

func (u *UserAggregate) MarkChangesAsCommitted() {
    u.changes = nil
}
```

## 🔄 CQRS (Command Query Responsibility Segregation)

### Command Side Implementation
```go
type UserCommandHandler struct {
    eventStore EventStore
    eventBus   EventBus
}

func (uch *UserCommandHandler) Handle(ctx context.Context, cmd Command) error {
    switch cmd.Type {
    case "RegisterUser":
        return uch.handleRegisterUser(ctx, cmd)
    case "ChangeUserEmail":
        return uch.handleChangeUserEmail(ctx, cmd)
    case "DeactivateUser":
        return uch.handleDeactivateUser(ctx, cmd)
    default:
        return ErrUnknownCommand
    }
}

func (uch *UserCommandHandler) handleRegisterUser(ctx context.Context, cmd Command) error {
    email := cmd.Data["email"].(string)
    name := cmd.Data["name"].(string)
    
    // Create new aggregate
    user := &UserAggregate{}
    
    // Execute business logic
    if err := user.RegisterUser(email, name); err != nil {
        return err
    }
    
    // Save events
    events := user.GetUncommittedChanges()
    if err := uch.eventStore.AppendEvents(user.ID, 0, events); err != nil {
        return err
    }
    
    // Publish events
    for _, event := range events {
        if err := uch.eventBus.Publish(ctx, event); err != nil {
            log.Printf("Failed to publish event: %v", err)
        }
    }
    
    user.MarkChangesAsCommitted()
    return nil
}

func (uch *UserCommandHandler) handleChangeUserEmail(ctx context.Context, cmd Command) error {
    userID := cmd.Data["user_id"].(string)
    newEmail := cmd.Data["new_email"].(string)
    
    // Load aggregate from event store
    events, err := uch.eventStore.GetEvents(userID, 0)
    if err != nil {
        return err
    }
    
    user := &UserAggregate{}
    user.LoadFromHistory(events)
    
    // Execute business logic
    if err := user.ChangeEmail(newEmail); err != nil {
        return err
    }
    
    // Save new events
    newEvents := user.GetUncommittedChanges()
    if err := uch.eventStore.AppendEvents(userID, user.Version-len(newEvents), newEvents); err != nil {
        return err
    }
    
    // Publish events
    for _, event := range newEvents {
        if err := uch.eventBus.Publish(ctx, event); err != nil {
            log.Printf("Failed to publish event: %v", err)
        }
    }
    
    return nil
}
```

### Query Side Implementation
```go
type UserQueryHandler struct {
    readDB Database
}

type UserView struct {
    ID          string    `json:"id"`
    Email       string    `json:"email"`
    Name        string    `json:"name"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    OrderCount  int       `json:"order_count"`
    TotalSpent  decimal.Decimal `json:"total_spent"`
}

func (uqh *UserQueryHandler) GetUser(ctx context.Context, userID string) (*UserView, error) {
    return uqh.readDB.GetUserView(ctx, userID)
}

func (uqh *UserQueryHandler) SearchUsers(ctx context.Context, query UserSearchQuery) ([]UserView, error) {
    return uqh.readDB.SearchUsers(ctx, query)
}

// Event handlers update read models
type UserProjectionHandler struct {
    readDB Database
}

func (uph *UserProjectionHandler) HandleUserRegistered(ctx context.Context, event Event) error {
    userView := &UserView{
        ID:        event.Data["user_id"].(string),
        Email:     event.Data["email"].(string),
        Name:      event.Data["name"].(string),
        Status:    "active",
        CreatedAt: event.Timestamp,
        UpdatedAt: event.Timestamp,
    }
    
    return uph.readDB.CreateUserView(ctx, userView)
}

func (uph *UserProjectionHandler) HandleUserEmailChanged(ctx context.Context, event Event) error {
    userID := event.Data["user_id"].(string)
    newEmail := event.Data["new_email"].(string)
    
    return uph.readDB.UpdateUserView(ctx, userID, map[string]interface{}{
        "email":      newEmail,
        "updated_at": event.Timestamp,
    })
}

func (uph *UserProjectionHandler) HandleOrderPlaced(ctx context.Context, event Event) error {
    userID := event.Data["user_id"].(string)
    amount := event.Data["total_amount"].(decimal.Decimal)
    
    // Update user statistics
    return uph.readDB.UpdateUserStats(ctx, userID, map[string]interface{}{
        "order_count": "order_count + 1",
        "total_spent": "total_spent + " + amount.String(),
        "updated_at":  event.Timestamp,
    })
}
```

## 📡 Event Bus Implementation

### In-Memory Event Bus
```go
type InMemoryEventBus struct {
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

type EventHandler func(ctx context.Context, event Event) error

func NewInMemoryEventBus() *InMemoryEventBus {
    return &InMemoryEventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (eb *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()
    
    eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

func (eb *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
    eb.mutex.RLock()
    handlers := eb.handlers[event.Type]
    eb.mutex.RUnlock()
    
    // Execute handlers concurrently
    var wg sync.WaitGroup
    errors := make(chan error, len(handlers))
    
    for _, handler := range handlers {
        wg.Add(1)
        go func(h EventHandler) {
            defer wg.Done()
            if err := h(ctx, event); err != nil {
                errors <- err
            }
        }(handler)
    }
    
    wg.Wait()
    close(errors)
    
    // Collect errors
    var errs []error
    for err := range errors {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("handler errors: %v", errs)
    }
    
    return nil
}
```

### Distributed Event Bus with Message Queue
```go
type DistributedEventBus struct {
    publisher  MessagePublisher
    subscriber MessageSubscriber
    handlers   map[string][]EventHandler
    mutex      sync.RWMutex
}

func (deb *DistributedEventBus) Subscribe(eventType string, handler EventHandler) error {
    deb.mutex.Lock()
    deb.handlers[eventType] = append(deb.handlers[eventType], handler)
    deb.mutex.Unlock()
    
    // Subscribe to message queue topic
    return deb.subscriber.Subscribe(eventType, func(ctx context.Context, message []byte) error {
        var event Event
        if err := json.Unmarshal(message, &event); err != nil {
            return err
        }
        
        deb.mutex.RLock()
        handlers := deb.handlers[event.Type]
        deb.mutex.RUnlock()
        
        // Execute handlers
        for _, handler := range handlers {
            if err := handler(ctx, event); err != nil {
                log.Printf("Handler error for event %s: %v", event.Type, err)
                // Could implement retry logic here
            }
        }
        
        return nil
    })
}

func (deb *DistributedEventBus) Publish(ctx context.Context, event Event) error {
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    return deb.publisher.Publish(ctx, event.Type, eventData)
}
```

## 🔄 Saga Pattern Implementation

### Orchestration-based Saga
```go
type OrderSagaOrchestrator struct {
    eventBus EventBus
    state    SagaState
    steps    []SagaStep
}

type SagaStep struct {
    Name        string
    Command     Command
    CompensateCommand Command
    Completed   bool
    Compensated bool
}

func NewOrderSaga(orderID string) *OrderSagaOrchestrator {
    return &OrderSagaOrchestrator{
        steps: []SagaStep{
            {
                Name: "ReserveInventory",
                Command: Command{
                    Type: "ReserveInventory",
                    Target: "inventory-service",
                    Data: map[string]interface{}{
                        "order_id": orderID,
                    },
                },
                CompensateCommand: Command{
                    Type: "ReleaseInventory",
                    Target: "inventory-service",
                    Data: map[string]interface{}{
                        "order_id": orderID,
                    },
                },
            },
            {
                Name: "ProcessPayment",
                Command: Command{
                    Type: "ProcessPayment",
                    Target: "payment-service",
                    Data: map[string]interface{}{
                        "order_id": orderID,
                    },
                },
                CompensateCommand: Command{
                    Type: "RefundPayment",
                    Target: "payment-service",
                    Data: map[string]interface{}{
                        "order_id": orderID,
                    },
                },
            },
        },
    }
}

func (oso *OrderSagaOrchestrator) Execute(ctx context.Context) error {
    for i, step := range oso.steps {
        // Send command
        if err := oso.eventBus.Publish(ctx, Event{
            Type: "CommandRequested",
            Data: map[string]interface{}{
                "command": step.Command,
            },
        }); err != nil {
            // Compensate previous steps
            oso.compensate(ctx, i-1)
            return err
        }
        
        // Wait for completion event
        if err := oso.waitForCompletion(ctx, step.Name); err != nil {
            // Compensate all completed steps
            oso.compensate(ctx, i)
            return err
        }
        
        oso.steps[i].Completed = true
    }
    
    return nil
}

func (oso *OrderSagaOrchestrator) compensate(ctx context.Context, lastCompletedStep int) {
    for i := lastCompletedStep; i >= 0; i-- {
        if oso.steps[i].Completed && !oso.steps[i].Compensated {
            oso.eventBus.Publish(ctx, Event{
                Type: "CommandRequested",
                Data: map[string]interface{}{
                    "command": oso.steps[i].CompensateCommand,
                },
            })
            oso.steps[i].Compensated = true
        }
    }
}
```

### Choreography-based Saga
```go
// Each service knows what to do when it receives an event
type InventoryService struct {
    eventBus EventBus
}

func (is *InventoryService) HandleOrderCreated(ctx context.Context, event Event) error {
    orderID := event.Data["order_id"].(string)
    
    // Try to reserve inventory
    if err := is.reserveInventory(orderID); err != nil {
        // Publish failure event
        return is.eventBus.Publish(ctx, Event{
            Type: "InventoryReservationFailed",
            Data: map[string]interface{}{
                "order_id": orderID,
                "reason":   err.Error(),
            },
        })
    }
    
    // Publish success event
    return is.eventBus.Publish(ctx, Event{
        Type: "InventoryReserved",
        Data: map[string]interface{}{
            "order_id": orderID,
        },
    })
}

type PaymentService struct {
    eventBus EventBus
}

func (ps *PaymentService) HandleInventoryReserved(ctx context.Context, event Event) error {
    orderID := event.Data["order_id"].(string)
    
    // Process payment
    if err := ps.processPayment(orderID); err != nil {
        // Publish failure event - inventory service will compensate
        return ps.eventBus.Publish(ctx, Event{
            Type: "PaymentFailed",
            Data: map[string]interface{}{
                "order_id": orderID,
                "reason":   err.Error(),
            },
        })
    }
    
    // Publish success event
    return ps.eventBus.Publish(ctx, Event{
        Type: "PaymentProcessed",
        Data: map[string]interface{}{
            "order_id": orderID,
        },
    })
}

// Inventory service compensates on payment failure
func (is *InventoryService) HandlePaymentFailed(ctx context.Context, event Event) error {
    orderID := event.Data["order_id"].(string)
    return is.releaseInventory(orderID)
}
```

## 🔗 Next Steps
- Study security and authentication patterns in distributed systems
- Learn about monitoring and observability with distributed tracing
- Explore advanced patterns like event sourcing projections
