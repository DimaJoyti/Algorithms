# 🏗️ Microservices Architecture

## 🎯 Microservices vs Monolith

### When to Choose Microservices
```go
// Decision Matrix
type ArchitectureDecision struct {
    TeamSize        int
    DomainComplexity string
    ScalingNeeds    string
    DeploymentFreq  string
    Recommendation  string
}

var decisions = []ArchitectureDecision{
    {TeamSize: 5, DomainComplexity: "Simple", ScalingNeeds: "Low", DeploymentFreq: "Monthly", Recommendation: "Monolith"},
    {TeamSize: 20, DomainComplexity: "Complex", ScalingNeeds: "High", DeploymentFreq: "Daily", Recommendation: "Microservices"},
    {TeamSize: 50, DomainComplexity: "Very Complex", ScalingNeeds: "Variable", DeploymentFreq: "Multiple/day", Recommendation: "Microservices"},
}

// Conway's Law: Organizations design systems that mirror their communication structure
// If you have 4 teams, you'll likely end up with 4 services
```

### Microservices Benefits & Challenges
```
Benefits:
✅ Independent deployment and scaling
✅ Technology diversity (Go, Python, Node.js)
✅ Team autonomy and ownership
✅ Fault isolation
✅ Better alignment with business domains

Challenges:
❌ Distributed system complexity
❌ Network latency and failures
❌ Data consistency across services
❌ Operational overhead
❌ Testing complexity
```

## 🔧 Service Decomposition Strategies

### 1. Domain-Driven Design (DDD)
```go
// Bounded Contexts - each becomes a microservice
type ECommerceSystem struct {
    UserManagement    UserService     // User registration, profiles, auth
    ProductCatalog    ProductService  // Product info, categories, search
    OrderManagement   OrderService    // Order processing, status tracking
    PaymentProcessing PaymentService  // Payment methods, transactions
    InventoryManagement InventoryService // Stock levels, reservations
    ShippingLogistics ShippingService // Shipping methods, tracking
    NotificationSystem NotificationService // Email, SMS, push notifications
}

// Each service owns its data and business logic
type UserService struct {
    userRepo     UserRepository
    profileRepo  ProfileRepository
    authService  AuthenticationService
}

// Clear service boundaries
func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Validate within domain
    if err := us.validateUserData(req); err != nil {
        return nil, err
    }
    
    // Create user
    user, err := us.userRepo.Create(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Publish domain event (don't call other services directly)
    us.eventBus.Publish("user.created", UserCreatedEvent{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
    })
    
    return user, nil
}
```

### 2. Data-Driven Decomposition
```go
// Separate services based on data ownership
type DataOwnership struct {
    Service    string
    Tables     []string
    ReadOnly   []string // Tables this service can read from other services
}

var serviceData = []DataOwnership{
    {
        Service: "UserService",
        Tables:  []string{"users", "user_profiles", "user_sessions"},
        ReadOnly: []string{}, // No external dependencies
    },
    {
        Service: "OrderService", 
        Tables:  []string{"orders", "order_items", "order_status_history"},
        ReadOnly: []string{"users.id", "products.id"}, // Reference data only
    },
    {
        Service: "ProductService",
        Tables:  []string{"products", "categories", "product_reviews"},
        ReadOnly: []string{"users.id"}, // For review authors
    },
}

// Anti-pattern: Shared database
// ❌ Multiple services accessing same tables
// ✅ Each service owns its data, communicates via APIs
```

### 3. Business Capability Decomposition
```go
// Organize around business functions
type BusinessCapability struct {
    Name         string
    Responsibilities []string
    APIs         []string
}

var capabilities = []BusinessCapability{
    {
        Name: "Customer Management",
        Responsibilities: []string{
            "User registration and authentication",
            "Profile management",
            "Customer support interactions",
        },
        APIs: []string{
            "POST /users",
            "GET /users/{id}",
            "PUT /users/{id}",
        },
    },
    {
        Name: "Order Fulfillment",
        Responsibilities: []string{
            "Order processing",
            "Inventory reservation",
            "Shipping coordination",
        },
        APIs: []string{
            "POST /orders",
            "GET /orders/{id}",
            "PUT /orders/{id}/status",
        },
    },
}
```

## 🌐 Inter-Service Communication

### 1. Synchronous Communication
```go
// HTTP/REST Client
type OrderServiceClient struct {
    baseURL    string
    httpClient *http.Client
    timeout    time.Duration
}

func NewOrderServiceClient(baseURL string) *OrderServiceClient {
    return &OrderServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        timeout: 30 * time.Second,
    }
}

func (osc *OrderServiceClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    url := fmt.Sprintf("%s/orders/%s", osc.baseURL, orderID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Add tracing headers
    req.Header.Set("X-Trace-ID", getTraceID(ctx))
    req.Header.Set("X-Request-ID", generateRequestID())
    
    resp, err := osc.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to call order service: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("order service returned status %d", resp.StatusCode)
    }
    
    var order Order
    if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &order, nil
}

// Circuit Breaker Pattern
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       CircuitState
    mutex       sync.RWMutex
}

type CircuitState int

const (
    Closed CircuitState = iota
    Open
    HalfOpen
)

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mutex.RLock()
    state := cb.state
    failures := cb.failures
    lastFailure := cb.lastFailure
    cb.mutex.RUnlock()
    
    switch state {
    case Open:
        if time.Since(lastFailure) > cb.timeout {
            // Try to close circuit
            cb.mutex.Lock()
            cb.state = HalfOpen
            cb.mutex.Unlock()
        } else {
            return ErrCircuitBreakerOpen
        }
    case HalfOpen:
        // Allow one request through
    case Closed:
        // Normal operation
    }
    
    err := fn()
    
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = Open
        }
        
        return err
    }
    
    // Success - reset circuit
    cb.failures = 0
    cb.state = Closed
    
    return nil
}
```

### 2. Asynchronous Communication
```go
// Event-Driven Communication
type EventBus struct {
    subscribers map[string][]EventHandler
    publisher   MessagePublisher
    mutex       sync.RWMutex
}

type EventHandler func(ctx context.Context, event Event) error

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
    eb.mutex.Lock()
    defer eb.mutex.Unlock()
    
    eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(ctx context.Context, event Event) error {
    // Publish to message queue for durability
    return eb.publisher.Publish(ctx, event.Type, event)
}

// Saga Pattern for Distributed Transactions
type OrderSaga struct {
    orderID string
    steps   []SagaStep
    state   SagaState
}

func (saga *OrderSaga) Execute(ctx context.Context) error {
    for i, step := range saga.steps {
        if err := step.Execute(ctx); err != nil {
            // Compensate all previous steps
            for j := i - 1; j >= 0; j-- {
                if compensateErr := saga.steps[j].Compensate(ctx); compensateErr != nil {
                    log.Printf("Compensation failed for step %s: %v", saga.steps[j].Name, compensateErr)
                }
            }
            return fmt.Errorf("saga failed at step %s: %w", step.Name, err)
        }
    }
    
    saga.state = SagaCompleted
    return nil
}

// Example: Order Processing Saga
func NewOrderProcessingSaga(orderID string) *OrderSaga {
    return &OrderSaga{
        orderID: orderID,
        steps: []SagaStep{
            {
                Name: "ReserveInventory",
                Execute: func(ctx context.Context) error {
                    return inventoryService.ReserveItems(ctx, orderID)
                },
                Compensate: func(ctx context.Context) error {
                    return inventoryService.ReleaseReservation(ctx, orderID)
                },
            },
            {
                Name: "ProcessPayment",
                Execute: func(ctx context.Context) error {
                    return paymentService.ChargePayment(ctx, orderID)
                },
                Compensate: func(ctx context.Context) error {
                    return paymentService.RefundPayment(ctx, orderID)
                },
            },
            {
                Name: "CreateShipment",
                Execute: func(ctx context.Context) error {
                    return shippingService.CreateShipment(ctx, orderID)
                },
                Compensate: func(ctx context.Context) error {
                    return shippingService.CancelShipment(ctx, orderID)
                },
            },
        },
    }
}
```

## 🕸️ Service Mesh Architecture

### Istio Service Mesh Implementation
```yaml
# Service mesh configuration
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-service
spec:
  http:
  - match:
    - headers:
        canary:
          exact: "true"
    route:
    - destination:
        host: user-service
        subset: v2
      weight: 100
  - route:
    - destination:
        host: user-service
        subset: v1
      weight: 90
    - destination:
        host: user-service
        subset: v2
      weight: 10

---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: user-service
spec:
  host: user-service
  trafficPolicy:
    circuitBreaker:
      consecutiveErrors: 3
      interval: 30s
      baseEjectionTime: 30s
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 50
        maxRequestsPerConnection: 2
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
```

```go
// Service mesh integration in Go
type ServiceMeshClient struct {
    serviceName string
    namespace   string
    client      *http.Client
}

func (smc *ServiceMeshClient) CallService(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // Service mesh handles service discovery
    url := fmt.Sprintf("http://%s.%s.svc.cluster.local%s", smc.serviceName, smc.namespace, path)
    
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
    if err != nil {
        return nil, err
    }
    
    // Service mesh automatically adds:
    // - Load balancing
    // - Circuit breaking
    // - Retries
    // - Distributed tracing
    // - mTLS encryption
    
    return smc.client.Do(req)
}
```

## 📊 Data Management Patterns

### 1. Database per Service
```go
// Each service has its own database
type UserService struct {
    userDB *sql.DB // PostgreSQL for user data
}

type ProductService struct {
    productDB *mongo.Database // MongoDB for product catalog
}

type OrderService struct {
    orderDB *sql.DB // PostgreSQL for transactional data
}

type AnalyticsService struct {
    analyticsDB *clickhouse.DB // ClickHouse for analytics
}

// Data synchronization via events
func (us *UserService) UpdateUser(ctx context.Context, userID string, updates UserUpdates) error {
    // Update in local database
    if err := us.userDB.UpdateUser(ctx, userID, updates); err != nil {
        return err
    }
    
    // Publish event for other services
    event := UserUpdatedEvent{
        UserID:    userID,
        Email:     updates.Email,
        Name:      updates.Name,
        UpdatedAt: time.Now(),
    }
    
    return us.eventBus.Publish(ctx, "user.updated", event)
}

// Other services listen for events
func (os *OrderService) HandleUserUpdated(ctx context.Context, event UserUpdatedEvent) error {
    // Update denormalized user data in orders
    return os.orderDB.UpdateUserInfoInOrders(ctx, event.UserID, event.Name, event.Email)
}
```

### 2. Shared Data Patterns
```go
// Reference Data Service
type ReferenceDataService struct {
    cache Cache
    db    Database
}

func (rds *ReferenceDataService) GetCountries(ctx context.Context) ([]Country, error) {
    // Cache reference data that rarely changes
    if countries, found := rds.cache.Get("countries"); found {
        return countries.([]Country), nil
    }
    
    countries, err := rds.db.GetCountries(ctx)
    if err != nil {
        return nil, err
    }
    
    // Cache for 24 hours
    rds.cache.Set("countries", countries, 24*time.Hour)
    
    return countries, nil
}

// Data Replication Pattern
type DataReplicationService struct {
    sourceDB Database
    replicas []Database
}

func (drs *DataReplicationService) ReplicateUserData(ctx context.Context, userID string) error {
    user, err := drs.sourceDB.GetUser(ctx, userID)
    if err != nil {
        return err
    }
    
    // Replicate to all services that need user data
    for _, replica := range drs.replicas {
        go func(db Database) {
            if err := db.UpsertUser(ctx, user); err != nil {
                log.Printf("Failed to replicate user data: %v", err)
            }
        }(replica)
    }
    
    return nil
}
```

## 🔍 Service Discovery & Configuration

### Service Registry Pattern
```go
type ServiceRegistry struct {
    services map[string][]ServiceInstance
    mutex    sync.RWMutex
}

type ServiceInstance struct {
    ID       string
    Address  string
    Port     int
    Health   HealthStatus
    Metadata map[string]string
}

func (sr *ServiceRegistry) Register(service string, instance ServiceInstance) error {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()
    
    sr.services[service] = append(sr.services[service], instance)
    
    // Start health checking
    go sr.healthCheck(service, instance)
    
    return nil
}

func (sr *ServiceRegistry) Discover(service string) ([]ServiceInstance, error) {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    
    instances, exists := sr.services[service]
    if !exists {
        return nil, ErrServiceNotFound
    }
    
    // Return only healthy instances
    var healthy []ServiceInstance
    for _, instance := range instances {
        if instance.Health == Healthy {
            healthy = append(healthy, instance)
        }
    }
    
    return healthy, nil
}

// Client-side load balancing
type ServiceClient struct {
    registry     ServiceRegistry
    loadBalancer LoadBalancer
}

func (sc *ServiceClient) CallService(ctx context.Context, serviceName, method, path string) (*http.Response, error) {
    instances, err := sc.registry.Discover(serviceName)
    if err != nil {
        return nil, err
    }
    
    if len(instances) == 0 {
        return nil, ErrNoHealthyInstances
    }
    
    instance := sc.loadBalancer.Select(instances)
    url := fmt.Sprintf("http://%s:%d%s", instance.Address, instance.Port, path)
    
    req, err := http.NewRequestWithContext(ctx, method, url, nil)
    if err != nil {
        return nil, err
    }
    
    return http.DefaultClient.Do(req)
}
```

### Configuration Management
```go
type ConfigManager struct {
    configs map[string]interface{}
    watchers map[string][]ConfigWatcher
    mutex   sync.RWMutex
}

type ConfigWatcher func(key string, oldValue, newValue interface{})

func (cm *ConfigManager) Get(key string) (interface{}, bool) {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    value, exists := cm.configs[key]
    return value, exists
}

func (cm *ConfigManager) Set(key string, value interface{}) {
    cm.mutex.Lock()
    oldValue := cm.configs[key]
    cm.configs[key] = value
    watchers := cm.watchers[key]
    cm.mutex.Unlock()
    
    // Notify watchers
    for _, watcher := range watchers {
        go watcher(key, oldValue, value)
    }
}

func (cm *ConfigManager) Watch(key string, watcher ConfigWatcher) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    
    cm.watchers[key] = append(cm.watchers[key], watcher)
}

// Environment-specific configuration
type ServiceConfig struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
    } `yaml:"database"`
    
    Cache struct {
        RedisURL string        `yaml:"redis_url"`
        TTL      time.Duration `yaml:"ttl"`
    } `yaml:"cache"`
    
    ExternalServices map[string]string `yaml:"external_services"`
}

func LoadConfig(env string) (*ServiceConfig, error) {
    configFile := fmt.Sprintf("config-%s.yaml", env)
    
    data, err := ioutil.ReadFile(configFile)
    if err != nil {
        return nil, err
    }
    
    var config ServiceConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 🚀 Deployment & Scaling Patterns

### Container Orchestration
```go
// Kubernetes Deployment Configuration
type KubernetesDeployment struct {
    ServiceName string
    Replicas    int32
    Image       string
    Resources   ResourceRequirements
    HealthCheck HealthCheckConfig
}

type ResourceRequirements struct {
    CPU    string // "100m" = 0.1 CPU
    Memory string // "128Mi" = 128 MiB
}

type HealthCheckConfig struct {
    LivenessProbe  ProbeConfig
    ReadinessProbe ProbeConfig
}

// Auto-scaling based on metrics
func (kd *KubernetesDeployment) ConfigureHPA() *HorizontalPodAutoscaler {
    return &HorizontalPodAutoscaler{
        MinReplicas: 2,
        MaxReplicas: 50,
        Metrics: []MetricSpec{
            {
                Type: "Resource",
                Resource: ResourceMetricSource{
                    Name: "cpu",
                    Target: ResourceMetricTarget{
                        Type:               "Utilization",
                        AverageUtilization: 70, // Scale when CPU > 70%
                    },
                },
            },
            {
                Type: "Resource",
                Resource: ResourceMetricSource{
                    Name: "memory",
                    Target: ResourceMetricTarget{
                        Type:               "Utilization",
                        AverageUtilization: 80, // Scale when Memory > 80%
                    },
                },
            },
        },
    }
}
```

### Blue-Green Deployment
```go
type BlueGreenDeployment struct {
    BlueVersion  ServiceVersion
    GreenVersion ServiceVersion
    LoadBalancer LoadBalancer
    HealthChecker HealthChecker
}

func (bgd *BlueGreenDeployment) Deploy(newVersion ServiceVersion) error {
    // Deploy to inactive environment (Green)
    if err := bgd.deployToGreen(newVersion); err != nil {
        return fmt.Errorf("failed to deploy to green: %w", err)
    }

    // Health check green environment
    if err := bgd.HealthChecker.WaitForHealthy(bgd.GreenVersion, 5*time.Minute); err != nil {
        return fmt.Errorf("green environment unhealthy: %w", err)
    }

    // Run smoke tests
    if err := bgd.runSmokeTests(bgd.GreenVersion); err != nil {
        return fmt.Errorf("smoke tests failed: %w", err)
    }

    // Switch traffic to green
    if err := bgd.LoadBalancer.SwitchTraffic(bgd.GreenVersion); err != nil {
        return fmt.Errorf("failed to switch traffic: %w", err)
    }

    // Keep blue for rollback capability
    bgd.BlueVersion, bgd.GreenVersion = bgd.GreenVersion, bgd.BlueVersion

    return nil
}

func (bgd *BlueGreenDeployment) Rollback() error {
    return bgd.LoadBalancer.SwitchTraffic(bgd.BlueVersion)
}
```

## 📋 Microservices Best Practices

### 1. Service Design Principles
```
Single Responsibility: Each service should have one reason to change
Autonomous: Services should be independently deployable
Business-Aligned: Services should align with business capabilities
Resilient: Services should handle failures gracefully
Observable: Services should provide monitoring and logging
Decentralized: Avoid shared databases and centralized governance
```

### 2. Common Anti-Patterns to Avoid
```go
// ❌ Anti-Pattern: Distributed Monolith
// Services that are too tightly coupled
func (orderService *OrderService) CreateOrder(order Order) error {
    // Synchronous calls to multiple services
    user := userService.GetUser(order.UserID)        // Blocking call
    product := productService.GetProduct(order.ProductID) // Blocking call
    inventory := inventoryService.Reserve(order.ProductID) // Blocking call
    payment := paymentService.Charge(order.Amount)    // Blocking call

    // If any service is down, entire operation fails
    return orderRepo.Save(order)
}

// ✅ Better Pattern: Event-Driven with Saga
func (orderService *OrderService) CreateOrder(order Order) error {
    // Create order immediately
    if err := orderRepo.Save(order); err != nil {
        return err
    }

    // Start saga for complex workflow
    saga := NewOrderProcessingSaga(order.ID)
    return saga.Start()
}

// ❌ Anti-Pattern: Shared Database
// Multiple services accessing same database tables

// ✅ Better Pattern: Database per Service
// Each service owns its data and provides APIs
```

### 3. Testing Strategies
```go
// Contract Testing with Pact
type UserServiceContract struct {
    pact *pact.Pact
}

func (usc *UserServiceContract) TestGetUser() {
    usc.pact.
        AddInteraction().
        Given("User 123 exists").
        UponReceiving("A request for user 123").
        WithRequest(pact.Request{
            Method: "GET",
            Path:   "/users/123",
        }).
        WillRespondWith(pact.Response{
            Status: 200,
            Body: map[string]interface{}{
                "id":    "123",
                "name":  "John Doe",
                "email": "john@example.com",
            },
        })

    // Test implementation
    client := NewUserServiceClient("http://localhost:8080")
    user, err := client.GetUser("123")

    assert.NoError(t, err)
    assert.Equal(t, "123", user.ID)
    assert.Equal(t, "John Doe", user.Name)
}

// Integration Testing
func TestOrderWorkflow(t *testing.T) {
    // Start test containers
    userService := startUserService()
    inventoryService := startInventoryService()
    orderService := startOrderService()

    defer func() {
        userService.Stop()
        inventoryService.Stop()
        orderService.Stop()
    }()

    // Test complete workflow
    order := Order{
        UserID:    "123",
        ProductID: "456",
        Quantity:  2,
    }

    result, err := orderService.CreateOrder(order)
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
}
```

## 🔗 Next Steps
- Study event-driven architecture patterns and CQRS
- Learn about security and authentication in distributed systems
- Explore monitoring and observability with distributed tracing
