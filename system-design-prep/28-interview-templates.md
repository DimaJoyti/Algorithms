# 📋 Interview Templates

## 🎯 Template Overview

These templates provide structured frameworks for system design interviews, ensuring comprehensive coverage of all essential topics within time constraints.

### Template Categories
```
📋 Available Templates:
1. Universal System Design Template (45-60 minutes)
2. Real-time System Template (Chat, Gaming, Live Streaming)
3. Data-Heavy System Template (Analytics, Search, ML)
4. E-commerce/Marketplace Template (Shopping, Booking)
5. Content Delivery Template (CDN, Streaming, Social Media)
6. Financial System Template (Payments, Trading, Banking)
```

## 📋 Universal System Design Template

### ⏰ Time Allocation (45-minute interview)
```
Phase 1: Requirements Clarification (5-8 minutes)
Phase 2: Capacity Estimation (3-5 minutes)
Phase 3: High-Level Design (10-15 minutes)
Phase 4: Detailed Design (15-20 minutes)
Phase 5: Scale & Optimize (8-12 minutes)
Phase 6: Wrap-up & Questions (3-5 minutes)
```

### 📝 Phase 1: Requirements Clarification (5-8 minutes)

#### Functional Requirements Checklist
```
🎯 Core Features:
□ What are the primary use cases?
□ Who are the main users (consumers, businesses, admins)?
□ What actions can users perform?
□ Are there different user roles and permissions?
□ What data needs to be stored and retrieved?

🎯 Feature Scope:
□ Real-time features needed?
□ Mobile and web support required?
□ Offline functionality needed?
□ Third-party integrations required?
□ Analytics and reporting needed?

🎯 Business Logic:
□ Any complex business rules?
□ Workflow or approval processes?
□ Content moderation requirements?
□ Compliance or regulatory needs?
```

#### Non-Functional Requirements Checklist
```
📊 Scale Requirements:
□ How many users (total and daily active)?
□ Expected growth rate?
□ Geographic distribution?
□ Peak usage patterns?

⚡ Performance Requirements:
□ Latency expectations?
□ Throughput requirements (QPS/TPS)?
□ Availability targets (99.9%, 99.99%)?
□ Consistency requirements?

🔒 Other Requirements:
□ Security and privacy needs?
□ Budget constraints?
□ Technology preferences?
□ Integration requirements?
```

### 📊 Phase 2: Capacity Estimation (3-5 minutes)

#### Estimation Framework
```go
// User and Traffic Estimation
const (
    TotalUsers       = X_000_000    // Total registered users
    DailyActiveUsers = Y_000_000    // Daily active users (typically 10-30% of total)
    ConcurrentUsers  = Z_000_000    // Peak concurrent users (typically 10-20% of DAU)
    
    // Usage Patterns
    ActionsPerUserPerDay = N        // Average actions per user per day
    PeakTrafficMultiplier = 3       // Peak traffic vs average
    ReadWriteRatio = 100           // Read:Write ratio (typically 10:1 to 1000:1)
)

// QPS Calculations
func CalculateQPS() {
    // Write QPS
    dailyWrites := DailyActiveUsers * ActionsPerUserPerDay
    avgWriteQPS := dailyWrites / (24 * 3600)
    peakWriteQPS := avgWriteQPS * PeakTrafficMultiplier
    
    // Read QPS
    avgReadQPS := avgWriteQPS * ReadWriteRatio
    peakReadQPS := avgReadQPS * PeakTrafficMultiplier
    
    fmt.Printf("Peak Write QPS: %d", peakWriteQPS)
    fmt.Printf("Peak Read QPS: %d", peakReadQPS)
}

// Storage Estimation
func CalculateStorage() {
    avgRecordSize := 1024          // Average record size in bytes
    recordsPerDay := dailyWrites
    
    dailyStorage := recordsPerDay * avgRecordSize
    yearlyStorage := dailyStorage * 365
    
    // With replication (typically 3x)
    totalStorage := yearlyStorage * 3
    
    fmt.Printf("Daily storage: %d GB", dailyStorage/(1024*1024*1024))
    fmt.Printf("Yearly storage: %d TB", yearlyStorage/(1024*1024*1024*1024))
}

// Bandwidth Estimation
func CalculateBandwidth() {
    avgRequestSize := 1024         // Average request size
    avgResponseSize := 4096        // Average response size
    
    incomingBandwidth := peakWriteQPS * avgRequestSize
    outgoingBandwidth := peakReadQPS * avgResponseSize
    
    fmt.Printf("Peak incoming: %d MB/s", incomingBandwidth/(1024*1024))
    fmt.Printf("Peak outgoing: %d MB/s", outgoingBandwidth/(1024*1024))
}
```

### 🏗️ Phase 3: High-Level Design (10-15 minutes)

#### Architecture Components Checklist
```
🌐 Client Layer:
□ Mobile apps (iOS, Android)
□ Web applications (React, Angular)
□ Desktop applications
□ APIs for third-party integrations

⚖️ Load Balancing:
□ DNS load balancing
□ Layer 4 (TCP) load balancers
□ Layer 7 (HTTP) load balancers
□ Global load balancing

🖥️ Application Layer:
□ Web servers (Nginx, Apache)
□ Application servers
□ API gateways
□ Microservices vs monolith

💾 Data Layer:
□ Primary databases
□ Cache layers
□ Search engines
□ File storage
□ Message queues

🔧 Infrastructure:
□ CDN for static content
□ Monitoring and logging
□ Security components
□ Backup and disaster recovery
```

#### Data Flow Design
```
1. User Request Flow:
   [Client] → [Load Balancer] → [API Gateway] → [Application Server]
                                                        ↓
   [Cache] ← [Database] ← [Business Logic] ← [Authentication]

2. Write Operation Flow:
   [Client] → [Validation] → [Business Logic] → [Database] → [Cache Update]
                                                     ↓
                                              [Event Publishing]

3. Read Operation Flow:
   [Client] → [Cache Check] → [Database Query] → [Response Formatting]
```

### 🔍 Phase 4: Detailed Design (15-20 minutes)

#### Database Design Template
```sql
-- User Management
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active'
);

-- Core Entity (customize based on system)
CREATE TABLE [main_entities] (
    id BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active'
);

-- Relationships/Associations
CREATE TABLE [entity_relationships] (
    id BIGINT PRIMARY KEY,
    entity1_id BIGINT REFERENCES [main_entities](id),
    entity2_id BIGINT REFERENCES [main_entities](id),
    relationship_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for Performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_entities_user_id ON [main_entities](user_id);
CREATE INDEX idx_entities_created_at ON [main_entities](created_at DESC);
CREATE INDEX idx_relationships_entity1 ON [entity_relationships](entity1_id);
```

#### API Design Template
```go
// REST API Design
type APIServer struct {
    userService    UserService
    entityService  EntityService
    authService    AuthService
    cacheService   CacheService
}

// Authentication endpoints
func (api *APIServer) setupAuthRoutes() {
    // POST /api/v1/auth/login
    // POST /api/v1/auth/logout
    // POST /api/v1/auth/refresh
    // GET  /api/v1/auth/me
}

// Core entity endpoints
func (api *APIServer) setupEntityRoutes() {
    // GET    /api/v1/entities          - List entities
    // POST   /api/v1/entities          - Create entity
    // GET    /api/v1/entities/{id}     - Get entity
    // PUT    /api/v1/entities/{id}     - Update entity
    // DELETE /api/v1/entities/{id}     - Delete entity
    // GET    /api/v1/entities/{id}/relationships - Get relationships
}

// Request/Response models
type CreateEntityRequest struct {
    Title       string            `json:"title" validate:"required,max=255"`
    Content     string            `json:"content" validate:"max=10000"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type EntityResponse struct {
    ID          int64             `json:"id"`
    Title       string            `json:"title"`
    Content     string            `json:"content"`
    Metadata    map[string]interface{} `json:"metadata"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

#### Caching Strategy Template
```go
// Multi-level caching strategy
type CacheManager struct {
    l1Cache    *sync.Map          // In-memory cache
    l2Cache    *redis.Client      // Redis cache
    l3Cache    *memcached.Client  // Memcached
}

func (cm *CacheManager) Get(key string) (interface{}, bool) {
    // L1: Check in-memory cache
    if value, found := cm.l1Cache.Load(key); found {
        return value, true
    }
    
    // L2: Check Redis
    if value, err := cm.l2Cache.Get(key).Result(); err == nil {
        // Populate L1 cache
        cm.l1Cache.Store(key, value)
        return value, true
    }
    
    // L3: Check Memcached
    if value, err := cm.l3Cache.Get(key); err == nil {
        // Populate L2 and L1 caches
        cm.l2Cache.Set(key, value, time.Hour)
        cm.l1Cache.Store(key, value)
        return value, true
    }
    
    return nil, false
}

// Cache patterns
const (
    CachePatternWriteThrough  = "write_through"   // Write to cache and DB
    CachePatternWriteBack     = "write_back"      // Write to cache, async to DB
    CachePatternWriteAround   = "write_around"    // Write to DB, invalidate cache
    CachePatternCacheAside    = "cache_aside"     // Application manages cache
)
```

### ⚡ Phase 5: Scale & Optimize (8-12 minutes)

#### Scaling Strategies Checklist
```
🔄 Horizontal Scaling:
□ Database sharding strategies
□ Application server scaling
□ Load balancer scaling
□ CDN and edge computing

📈 Performance Optimization:
□ Database query optimization
□ Caching improvements
□ Connection pooling
□ Compression and minification

🛡️ Reliability & Availability:
□ Redundancy and failover
□ Circuit breakers
□ Health checks and monitoring
□ Disaster recovery

🔒 Security Scaling:
□ DDoS protection
□ Rate limiting
□ Authentication scaling
□ Data encryption
```

#### Bottleneck Identification Framework
```go
// Common bottlenecks and solutions
type BottleneckAnalysis struct {
    Component   string
    Symptoms    []string
    Solutions   []string
    Priority    int
}

var CommonBottlenecks = []BottleneckAnalysis{
    {
        Component: "Database",
        Symptoms:  []string{"High query latency", "Connection pool exhaustion", "Lock contention"},
        Solutions: []string{"Read replicas", "Sharding", "Query optimization", "Connection pooling"},
        Priority:  1,
    },
    {
        Component: "Application Server",
        Symptoms:  []string{"High CPU usage", "Memory exhaustion", "Thread pool saturation"},
        Solutions: []string{"Horizontal scaling", "Load balancing", "Async processing", "Caching"},
        Priority:  2,
    },
    {
        Component: "Network",
        Symptoms:  []string{"High bandwidth usage", "Latency spikes", "Packet loss"},
        Solutions: []string{"CDN", "Compression", "Connection pooling", "Edge computing"},
        Priority:  3,
    },
}

func IdentifyBottlenecks(metrics SystemMetrics) []BottleneckAnalysis {
    var bottlenecks []BottleneckAnalysis

    if metrics.DatabaseLatency > 100 { // ms
        bottlenecks = append(bottlenecks, CommonBottlenecks[0])
    }

    if metrics.CPUUtilization > 80 { // %
        bottlenecks = append(bottlenecks, CommonBottlenecks[1])
    }

    if metrics.NetworkLatency > 50 { // ms
        bottlenecks = append(bottlenecks, CommonBottlenecks[2])
    }

    return bottlenecks
}
```

### 🎯 Phase 6: Wrap-up & Questions (3-5 minutes)

#### Summary Template
```
📋 Design Summary:
"To summarize, I've designed a [system type] that:
- Handles [X] users with [Y] QPS through [scaling strategy]
- Uses [database choice] for data persistence with [consistency model]
- Implements [caching strategy] for performance
- Ensures [availability target] through [redundancy approach]
- Scales horizontally via [scaling mechanisms]

Key trade-offs made:
- [Trade-off 1]: Chose [option A] over [option B] because [reasoning]
- [Trade-off 2]: Prioritized [aspect X] over [aspect Y] due to [constraint]

Alternative approaches considered:
- [Alternative 1]: Would provide [benefit] but [drawback]
- [Alternative 2]: Better for [scenario] but not suitable because [reason]"
```

#### Questions to Ask Interviewer
```
🤔 Technical Questions:
□ "What's the team's experience with [proposed technologies]?"
□ "Are there existing systems this would need to integrate with?"
□ "What are the main technical challenges the team currently faces?"
□ "How does this fit into the broader system architecture?"

🏢 Business Questions:
□ "What's the expected timeline for building this system?"
□ "Are there specific compliance or regulatory requirements?"
□ "What's the team structure for a project like this?"
□ "How do you typically handle system migrations?"

🔍 Follow-up Questions:
□ "Would you like me to dive deeper into any specific component?"
□ "Are there any edge cases or scenarios I should consider?"
□ "How would you approach testing this system?"
□ "What monitoring and alerting would be most important?"
```

## 📋 Specialized Templates

### 🔄 Real-time System Template (Chat, Gaming, Live Streaming)

#### Additional Requirements to Clarify
```
⚡ Real-time Specific:
□ What's the acceptable latency for real-time features?
□ How many concurrent connections needed?
□ Do we need message ordering guarantees?
□ What happens when users go offline?
□ Are there different types of real-time events?

🔌 Connection Management:
□ WebSocket vs Server-Sent Events vs Long Polling?
□ How to handle connection drops and reconnections?
□ Do we need presence indicators (online/offline)?
□ How to scale WebSocket connections?
```

#### Real-time Architecture Components
```go
// WebSocket connection manager
type ConnectionManager struct {
    connections map[string]*websocket.Conn
    rooms       map[string][]string
    userSessions map[string][]string
    messageQueue MessageQueue
    presence     PresenceService
}

// Message routing for real-time systems
type MessageRouter struct {
    connectionMgr *ConnectionManager
    messageStore  MessageStore
    deliveryMgr   DeliveryManager
}

func (mr *MessageRouter) RouteMessage(message *Message) error {
    // Store message for durability
    if err := mr.messageStore.Store(message); err != nil {
        return err
    }

    // Route to online recipients
    for _, recipientID := range message.Recipients {
        if conn, online := mr.connectionMgr.GetConnection(recipientID); online {
            mr.sendToConnection(conn, message)
        } else {
            // Queue for offline delivery
            mr.deliveryMgr.QueueForOfflineDelivery(recipientID, message)
        }
    }

    return nil
}
```

### 📊 Data-Heavy System Template (Analytics, Search, ML)

#### Data-Specific Requirements
```
📈 Data Characteristics:
□ What's the data volume (GB, TB, PB)?
□ What's the data velocity (batch vs streaming)?
□ What's the data variety (structured, unstructured)?
□ How fresh does the data need to be?
□ What are the query patterns?

🔍 Analytics Requirements:
□ Real-time vs batch analytics?
□ What types of aggregations needed?
□ Do we need OLAP capabilities?
□ Are there machine learning requirements?
□ What's the query complexity?
```

#### Data Architecture Components
```go
// Data pipeline architecture
type DataPipeline struct {
    ingestion    DataIngestionService
    processing   DataProcessingService
    storage      DataStorageService
    serving      DataServingService
    monitoring   DataMonitoringService
}

// Lambda architecture for real-time + batch processing
type LambdaArchitecture struct {
    batchLayer    BatchProcessingLayer    // Hadoop, Spark
    speedLayer    StreamProcessingLayer   // Kafka, Storm
    servingLayer  ServingLayer           // HBase, Cassandra
}

// Data storage strategy
type DataStorageStrategy struct {
    hotStorage    Storage // Recent, frequently accessed data
    warmStorage   Storage // Older, occasionally accessed data
    coldStorage   Storage // Archive, rarely accessed data
    searchIndex   SearchEngine // Elasticsearch, Solr
}
```

### 🛒 E-commerce/Marketplace Template

#### E-commerce Specific Requirements
```
💰 Business Logic:
□ What types of products/services?
□ Do we need inventory management?
□ What payment methods to support?
□ Are there multiple sellers/vendors?
□ Do we need order tracking?

🔒 Transaction Requirements:
□ What's the consistency requirement for payments?
□ How to handle payment failures?
□ Do we need refunds and cancellations?
□ Are there fraud detection requirements?
□ What about tax calculations?
```

#### E-commerce Architecture Components
```go
// E-commerce service architecture
type EcommerceSystem struct {
    userService     UserService
    productService  ProductService
    inventoryService InventoryService
    orderService    OrderService
    paymentService  PaymentService
    shippingService ShippingService
    notificationService NotificationService
}

// Order processing workflow
type OrderWorkflow struct {
    steps []OrderStep
}

type OrderStep interface {
    Execute(order *Order) error
    Compensate(order *Order) error // For saga pattern
}

// Typical order steps:
// 1. Validate order
// 2. Check inventory
// 3. Reserve inventory
// 4. Process payment
// 5. Create shipment
// 6. Send confirmation
```

## 🎯 Template Usage Guidelines

### ✅ How to Use These Templates

1. **Pre-Interview Preparation**
   - Review relevant template based on expected problem type
   - Practice with template structure
   - Memorize key checklists and frameworks

2. **During Interview**
   - Use template as mental checklist
   - Adapt template to specific problem requirements
   - Don't rigidly follow if problem requires different approach

3. **Time Management**
   - Use time allocations as guidelines
   - Adjust based on interviewer preferences
   - Always leave time for questions and wrap-up

### 🎯 Template Customization

#### Company-Specific Adaptations
```
🔥 Meta/Facebook Style:
- Emphasize social graph and viral growth
- Focus on real-time features and engagement
- Discuss A/B testing and experimentation
- Consider global scale (billions of users)

🔥 Google Style:
- Emphasize search and information retrieval
- Focus on distributed systems and algorithms
- Discuss data processing at massive scale
- Consider reliability and fault tolerance

🔥 Amazon Style:
- Emphasize microservices and loose coupling
- Focus on cost optimization and efficiency
- Discuss backwards compatibility
- Consider operational excellence
```

## 🎯 Summary: Interview Template Mastery

### Complete Template Portfolio
- **Universal Template** - Works for any system design problem
- **Real-time Template** - Chat, gaming, live streaming systems
- **Data-Heavy Template** - Analytics, search, ML systems
- **E-commerce Template** - Shopping, marketplace, booking systems
- **Specialized Checklists** - Company-specific adaptations

### Key Template Benefits
- **Structured Approach** - Ensures comprehensive coverage
- **Time Management** - Proper allocation across all phases
- **Consistency** - Repeatable methodology for any problem
- **Completeness** - Covers all essential system design aspects
- **Flexibility** - Adaptable to different problem types and companies

**You now have professional-grade templates for excelling in any system design interview!** 🚀
