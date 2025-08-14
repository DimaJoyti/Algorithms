# ⚖️ Trade-offs & Constraints in System Design

## 🎯 Understanding Trade-offs
Every system design decision involves trade-offs. There's no perfect solution - only solutions that are optimal for specific requirements and constraints.

## 🔺 The CAP Theorem

### Definition
In a distributed system, you can only guarantee **2 out of 3**:
- **Consistency**: All nodes see the same data simultaneously
- **Availability**: System remains operational 100% of the time
- **Partition Tolerance**: System continues despite network failures

### Real-World Applications

#### CP Systems (Consistency + Partition Tolerance)
```
Examples: MongoDB, Redis Cluster, HBase
Trade-off: Sacrifice availability during network partitions
Use case: Financial systems, inventory management

When network partition occurs:
- System becomes unavailable
- Ensures data consistency
- No conflicting writes
```

#### AP Systems (Availability + Partition Tolerance)
```
Examples: Cassandra, DynamoDB, CouchDB
Trade-off: Sacrifice consistency for availability
Use case: Social media feeds, content delivery

When network partition occurs:
- System remains available
- May serve stale data
- Eventual consistency model
```

#### CA Systems (Consistency + Availability)
```
Examples: Traditional RDBMS (PostgreSQL, MySQL)
Trade-off: Cannot handle network partitions
Use case: Single-datacenter applications

Note: In distributed systems, network partitions are inevitable,
so true CA systems don't exist in practice.
```

## 🔄 Consistency Models

### Strong Consistency
```go
// Example: Bank account balance
func TransferMoney(from, to AccountID, amount Money) error {
    tx := db.BeginTransaction()
    defer tx.Rollback()
    
    // Both operations must succeed or fail together
    if err := tx.DebitAccount(from, amount); err != nil {
        return err
    }
    if err := tx.CreditAccount(to, amount); err != nil {
        return err
    }
    
    return tx.Commit() // Atomic operation
}

Guarantees:
✅ All reads receive the most recent write
✅ No stale data
❌ Higher latency
❌ Lower availability during failures
```

### Eventual Consistency
```go
// Example: Social media post likes
func LikePost(userID, postID string) error {
    // Update multiple replicas asynchronously
    go updateReplica1(postID, incrementLikes)
    go updateReplica2(postID, incrementLikes)
    go updateReplica3(postID, incrementLikes)
    
    return nil // Return immediately
}

Guarantees:
✅ High availability
✅ Low latency
❌ Temporary inconsistency
❌ Complex conflict resolution
```

### Weak Consistency
```go
// Example: Live video streaming viewer count
func UpdateViewerCount(streamID string, count int) {
    // Best effort update, some updates may be lost
    cache.Set(fmt.Sprintf("viewers:%s", streamID), count, 30*time.Second)
}

Use cases:
- Real-time analytics
- Live dashboards
- Gaming leaderboards
```

## ⚡ Latency vs Throughput

### Latency
**Definition**: Time to process a single request

```go
// Low latency optimization
func GetUserProfile(userID string) (*User, error) {
    // Check L1 cache (1ms)
    if user, found := l1Cache.Get(userID); found {
        return user.(*User), nil
    }
    
    // Check L2 cache (5ms)
    if user, found := l2Cache.Get(userID); found {
        l1Cache.Set(userID, user, 1*time.Minute)
        return user.(*User), nil
    }
    
    // Database query (100ms)
    user, err := db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache for future requests
    l1Cache.Set(userID, user, 1*time.Minute)
    l2Cache.Set(userID, user, 1*time.Hour)
    
    return user, nil
}

Optimization techniques:
- Multi-level caching
- Connection pooling
- Async processing
- Geographic distribution
```

### Throughput
**Definition**: Number of requests processed per unit time

```go
// High throughput optimization
func ProcessBatch(requests []Request) error {
    // Batch processing for higher throughput
    const batchSize = 1000
    
    for i := 0; i < len(requests); i += batchSize {
        end := min(i+batchSize, len(requests))
        batch := requests[i:end]
        
        // Process batch in parallel
        go func(batch []Request) {
            db.BatchInsert(batch)
        }(batch)
    }
    
    return nil
}

Optimization techniques:
- Batch processing
- Connection pooling
- Horizontal scaling
- Async processing
```

### The Trade-off
```
High Latency + High Throughput: Batch processing systems
Low Latency + Low Throughput: Real-time systems
Low Latency + High Throughput: Expensive, requires optimization
High Latency + Low Throughput: Poor design (avoid)
```

## 💾 Storage Trade-offs

### SQL vs NoSQL

#### SQL Databases (RDBMS)
```sql
-- Strong schema, ACID properties
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL
);

Advantages:
✅ ACID compliance
✅ Complex queries (JOINs)
✅ Strong consistency
✅ Mature ecosystem

Disadvantages:
❌ Vertical scaling limitations
❌ Schema rigidity
❌ Complex sharding
```

#### NoSQL Databases

##### Document Stores (MongoDB, CouchDB)
```javascript
// Flexible schema
{
  "_id": "user123",
  "email": "user@example.com",
  "profile": {
    "name": "John Doe",
    "preferences": ["tech", "sports"],
    "metadata": {
      "lastLogin": "2024-01-15T10:30:00Z"
    }
  }
}

Use cases:
- Content management
- User profiles
- Product catalogs
```

##### Key-Value Stores (Redis, DynamoDB)
```go
// Simple key-value operations
cache.Set("user:123:profile", userProfile, 1*time.Hour)
cache.Set("session:abc123", sessionData, 30*time.Minute)

Use cases:
- Caching
- Session storage
- Real-time recommendations
```

##### Column-Family (Cassandra, HBase)
```cql
-- Optimized for write-heavy workloads
CREATE TABLE user_activities (
    user_id UUID,
    timestamp TIMESTAMP,
    activity_type TEXT,
    data TEXT,
    PRIMARY KEY (user_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

Use cases:
- Time-series data
- IoT data
- Analytics
```

##### Graph Databases (Neo4j, Amazon Neptune)
```cypher
// Complex relationship queries
MATCH (user:User)-[:FOLLOWS]->(friend:User)-[:LIKES]->(post:Post)
WHERE user.id = "123"
RETURN post

Use cases:
- Social networks
- Recommendation engines
- Fraud detection
```

## 🔄 Caching Trade-offs

### Cache Patterns

#### Cache-Aside (Lazy Loading)
```go
func GetUser(userID string) (*User, error) {
    // Check cache first
    if user, found := cache.Get(userID); found {
        return user.(*User), nil
    }
    
    // Cache miss - load from database
    user, err := db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    cache.Set(userID, user, 1*time.Hour)
    return user, nil
}

Advantages:
✅ Only cache requested data
✅ Cache failures don't affect system
❌ Cache miss penalty
❌ Stale data possible
```

#### Write-Through
```go
func UpdateUser(user *User) error {
    // Update database first
    if err := db.UpdateUser(user); err != nil {
        return err
    }
    
    // Update cache
    cache.Set(user.ID, user, 1*time.Hour)
    return nil
}

Advantages:
✅ Cache always consistent
✅ No cache miss penalty for writes
❌ Write latency increased
❌ Unused data cached
```

#### Write-Behind (Write-Back)
```go
func UpdateUser(user *User) error {
    // Update cache immediately
    cache.Set(user.ID, user, 1*time.Hour)
    
    // Async database update
    go func() {
        db.UpdateUser(user)
    }()
    
    return nil
}

Advantages:
✅ Low write latency
✅ High write throughput
❌ Data loss risk
❌ Complex consistency
```

## 🌐 Network Trade-offs

### Synchronous vs Asynchronous Communication

#### Synchronous (HTTP, gRPC)
```go
func ProcessOrder(orderID string) error {
    // Synchronous calls - blocking
    user, err := userService.GetUser(order.UserID)
    if err != nil {
        return err
    }
    
    inventory, err := inventoryService.CheckStock(order.ProductID)
    if err != nil {
        return err
    }
    
    payment, err := paymentService.ProcessPayment(order.Amount)
    if err != nil {
        return err
    }
    
    return nil
}

Advantages:
✅ Simple error handling
✅ Immediate consistency
✅ Easy debugging
❌ Higher latency
❌ Cascading failures
❌ Resource blocking
```

#### Asynchronous (Message Queues, Events)
```go
func ProcessOrder(orderID string) error {
    // Publish events asynchronously
    eventBus.Publish("order.created", OrderCreatedEvent{
        OrderID: orderID,
        UserID:  order.UserID,
        Amount:  order.Amount,
    })
    
    return nil // Return immediately
}

// Event handlers process asynchronously
func HandleOrderCreated(event OrderCreatedEvent) {
    // Process in background
    go userService.UpdateUserStats(event.UserID)
    go inventoryService.ReserveStock(event.ProductID)
    go paymentService.ProcessPayment(event.Amount)
}

Advantages:
✅ Low latency
✅ High throughput
✅ Fault tolerance
❌ Complex error handling
❌ Eventual consistency
❌ Debugging complexity
```

## 🔐 Security vs Performance

### Authentication Trade-offs

#### JWT vs Session-based
```go
// JWT - Stateless but larger
type JWTClaims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Roles    []string `json:"roles"`
    jwt.StandardClaims
}

// Session - Stateful but smaller
type Session struct {
    ID     string
    UserID string
    Expiry time.Time
}

JWT Advantages:
✅ Stateless
✅ Scalable
❌ Larger tokens
❌ Hard to revoke

Session Advantages:
✅ Smaller overhead
✅ Easy revocation
❌ Server state required
❌ Scaling complexity
```

### Encryption Trade-offs
```go
// Performance impact of encryption
func BenchmarkEncryption(b *testing.B) {
    data := make([]byte, 1024)
    
    // No encryption: ~1000 ns/op
    b.Run("NoEncryption", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = data
        }
    })
    
    // AES encryption: ~5000 ns/op
    b.Run("AES", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            encrypted := aes.Encrypt(data)
            _ = encrypted
        }
    })
}

Trade-offs:
Security ↑ Performance ↓
Compliance requirements vs Speed
```

## 📊 Monitoring vs Performance

### Observability Overhead
```go
// Tracing adds latency
func ProcessRequest(ctx context.Context, req Request) error {
    // Start span (adds ~1-5μs)
    span, ctx := tracer.Start(ctx, "process_request")
    defer span.End()
    
    // Add attributes (adds ~1μs per attribute)
    span.SetAttributes(
        attribute.String("user_id", req.UserID),
        attribute.String("request_type", req.Type),
    )
    
    // Business logic
    return processBusinessLogic(ctx, req)
}

Trade-offs:
Observability ↑ Performance ↓
Debugging capability vs Latency
```

## 🎯 Making Trade-off Decisions

### Decision Framework
```
1. Identify Requirements
   - Functional requirements
   - Non-functional requirements (SLAs)
   - Business constraints

2. List Options
   - Different architectural approaches
   - Technology choices
   - Implementation strategies

3. Evaluate Trade-offs
   - Performance implications
   - Complexity costs
   - Operational overhead
   - Future scalability

4. Make Decision
   - Document reasoning
   - Plan for monitoring
   - Prepare for evolution
```

### Example: Choosing Database for Social Media
```
Requirements:
- 100M users
- 1B posts per day
- Real-time feeds
- Global distribution

Options Analysis:
1. PostgreSQL
   ✅ Strong consistency
   ✅ Complex queries
   ❌ Scaling limitations
   ❌ Global latency

2. Cassandra
   ✅ High write throughput
   ✅ Global distribution
   ❌ Limited query flexibility
   ❌ Eventual consistency

3. Hybrid Approach
   ✅ PostgreSQL for user data
   ✅ Cassandra for posts/feeds
   ❌ Increased complexity
   ❌ Data synchronization

Decision: Hybrid approach
Reasoning: Optimize for each use case
```

## 🔗 Next Steps
- Study specific system components (databases, caches, load balancers)
- Practice making trade-off decisions in system design problems
- Learn about microservices and event-driven architectures
