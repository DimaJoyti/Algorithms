# 🗄️ Database Systems Deep Dive

## 🎯 Database Selection Framework

### Decision Tree
```
1. Data Structure
   ├── Structured/Relational → SQL
   ├── Semi-structured → Document DB
   ├── Key-Value pairs → Key-Value Store
   ├── Graph relationships → Graph DB
   └── Time-series → Time-series DB

2. Consistency Requirements
   ├── Strong consistency → SQL (ACID)
   ├── Eventual consistency → NoSQL
   └── Weak consistency → Cache/In-memory

3. Scale Requirements
   ├── < 1TB, < 10K QPS → Single SQL instance
   ├── < 100TB, < 100K QPS → SQL with replicas
   └── > 100TB, > 100K QPS → NoSQL/Distributed
```

## 🏛️ SQL Databases (RDBMS)

### ACID Properties

#### Atomicity
```go
func TransferMoney(db *sql.DB, fromAccount, toAccount int64, amount decimal.Decimal) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // Rollback if not committed
    
    // Both operations must succeed or fail together
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccount)
    if err != nil {
        return err
    }
    
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccount)
    if err != nil {
        return err
    }
    
    return tx.Commit() // All or nothing
}
```

#### Consistency
```sql
-- Database constraints ensure consistency
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    balance DECIMAL(15,2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Foreign key constraints
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES accounts(id),
    to_account_id BIGINT REFERENCES accounts(id),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0)
);
```

#### Isolation Levels
```go
type IsolationLevel int

const (
    ReadUncommitted IsolationLevel = iota // Dirty reads possible
    ReadCommitted                         // Default in most DBs
    RepeatableRead                        // Phantom reads possible
    Serializable                          // Strictest, performance impact
)

// Example: Setting isolation level
func QueryWithIsolation(db *sql.DB, level IsolationLevel) error {
    tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Your queries here
    rows, err := tx.Query("SELECT * FROM accounts WHERE balance > 1000")
    // ... process rows
    
    return tx.Commit()
}
```

#### Durability
```
Write-Ahead Logging (WAL):
1. Log changes before applying to data files
2. Survive system crashes
3. Recovery using transaction logs

Synchronous vs Asynchronous commits:
- Sync: Guaranteed durability, higher latency
- Async: Better performance, small data loss risk
```

### SQL Scaling Strategies

#### Vertical Scaling
```go
type DatabaseConfig struct {
    MaxConnections    int    // Scale connection pool
    SharedBuffers     string // "256MB" -> "8GB"
    EffectiveCacheSize string // "4GB" -> "64GB"
    WorkMem          string // "4MB" -> "256MB"
    MaintenanceWorkMem string // "64MB" -> "2GB"
}

// PostgreSQL optimization example
func OptimizePostgreSQL(config DatabaseConfig) {
    // Connection pooling
    db.SetMaxOpenConns(config.MaxConnections)
    db.SetMaxIdleConns(config.MaxConnections / 2)
    db.SetConnMaxLifetime(time.Hour)
}
```

#### Read Replicas
```go
type DatabaseCluster struct {
    Master   *sql.DB
    Replicas []*sql.DB
    current  int
    mutex    sync.RWMutex
}

func (dc *DatabaseCluster) Write(query string, args ...interface{}) error {
    // All writes go to master
    _, err := dc.Master.Exec(query, args...)
    return err
}

func (dc *DatabaseCluster) Read(query string, args ...interface{}) (*sql.Rows, error) {
    // Round-robin read from replicas
    dc.mutex.Lock()
    replica := dc.Replicas[dc.current]
    dc.current = (dc.current + 1) % len(dc.Replicas)
    dc.mutex.Unlock()
    
    return replica.Query(query, args...)
}
```

#### Sharding (Horizontal Partitioning)
```go
type ShardedDatabase struct {
    shards []DatabaseShard
}

type DatabaseShard struct {
    ID       int
    Database *sql.DB
    MinRange int64
    MaxRange int64
}

// Range-based sharding
func (sd *ShardedDatabase) GetShardByUserID(userID int64) *DatabaseShard {
    for _, shard := range sd.shards {
        if userID >= shard.MinRange && userID <= shard.MaxRange {
            return &shard
        }
    }
    return nil
}

// Hash-based sharding
func (sd *ShardedDatabase) GetShardByHash(key string) *DatabaseShard {
    hash := crc32.ChecksumIEEE([]byte(key))
    shardIndex := int(hash) % len(sd.shards)
    return &sd.shards[shardIndex]
}

// Consistent hashing (similar to your algorithms collection)
func (sd *ShardedDatabase) GetShardConsistent(key string) *DatabaseShard {
    // Implementation using consistent hashing algorithm
    // Minimizes redistribution when adding/removing shards
    return sd.consistentHash.GetNode(key)
}
```

### SQL Database Types

#### PostgreSQL
```go
// Strengths: ACID, complex queries, JSON support, extensions
type PostgreSQLUseCase struct {
    Strengths []string
    UseCases  []string
    Scaling   ScalingOptions
}

var PostgreSQL = PostgreSQLUseCase{
    Strengths: []string{
        "Strong ACID compliance",
        "Advanced SQL features",
        "JSON/JSONB support",
        "Full-text search",
        "Extensible (PostGIS, etc.)",
    },
    UseCases: []string{
        "Financial systems",
        "E-commerce platforms",
        "Content management",
        "Analytics workloads",
    },
    Scaling: ScalingOptions{
        ReadReplicas: true,
        Sharding:     "Manual/Application-level",
        MaxSize:      "~10TB per instance",
    },
}
```

#### MySQL
```go
var MySQL = DatabaseProfile{
    Strengths: []string{
        "High performance",
        "Mature ecosystem",
        "Easy replication",
        "Wide adoption",
    },
    UseCases: []string{
        "Web applications",
        "OLTP systems",
        "Read-heavy workloads",
    },
    Considerations: []string{
        "Limited JSON support",
        "Storage engine choice important",
    },
}
```

## 🌐 NoSQL Databases

### Document Databases

#### MongoDB
```go
// Flexible schema, rich queries
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Email    string            `bson:"email"`
    Profile  UserProfile       `bson:"profile"`
    Settings map[string]interface{} `bson:"settings"`
}

type UserProfile struct {
    Name        string   `bson:"name"`
    Age         int      `bson:"age"`
    Interests   []string `bson:"interests"`
    Address     Address  `bson:"address"`
}

func (ur *UserRepository) CreateUser(user *User) error {
    collection := ur.db.Collection("users")
    result, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        return err
    }
    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

// Complex queries with aggregation pipeline
func (ur *UserRepository) GetUsersByInterest(interest string) ([]*User, error) {
    pipeline := []bson.M{
        {"$match": bson.M{"profile.interests": interest}},
        {"$sort": bson.M{"profile.age": 1}},
        {"$limit": 100},
    }
    
    cursor, err := ur.collection.Aggregate(context.Background(), pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    var users []*User
    if err = cursor.All(context.Background(), &users); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

### Key-Value Stores

#### Redis
```go
// In-memory, high performance
type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
    return &RedisCache{
        client: redis.NewClient(&redis.Options{
            Addr:     addr,
            Password: "",
            DB:       0,
        }),
    }
}

// String operations
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
    return r.client.Set(context.Background(), key, value, expiration).Err()
}

// Hash operations (for objects)
func (r *RedisCache) SetUser(userID string, user *User) error {
    userMap := map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
        "name":  user.Name,
    }
    return r.client.HMSet(context.Background(), fmt.Sprintf("user:%s", userID), userMap).Err()
}

// List operations (for feeds)
func (r *RedisCache) AddToFeed(userID string, postID string) error {
    key := fmt.Sprintf("feed:%s", userID)
    // Add to beginning of list
    r.client.LPush(context.Background(), key, postID)
    // Keep only latest 1000 posts
    return r.client.LTrim(context.Background(), key, 0, 999).Err()
}

// Set operations (for followers)
func (r *RedisCache) AddFollower(userID, followerID string) error {
    key := fmt.Sprintf("followers:%s", userID)
    return r.client.SAdd(context.Background(), key, followerID).Err()
}
```

#### DynamoDB
```go
// AWS managed, serverless scaling
type DynamoDBRepository struct {
    client *dynamodb.Client
    table  string
}

// Single-table design pattern
type UserItem struct {
    PK        string `dynamodbav:"PK"`        // USER#123
    SK        string `dynamodbav:"SK"`        // PROFILE
    Type      string `dynamodbav:"Type"`      // USER
    Email     string `dynamodbav:"Email"`
    Name      string `dynamodbav:"Name"`
    CreatedAt string `dynamodbav:"CreatedAt"`
}

type PostItem struct {
    PK        string `dynamodbav:"PK"`        // USER#123
    SK        string `dynamodbav:"SK"`        // POST#456
    Type      string `dynamodbav:"Type"`      // POST
    Content   string `dynamodbav:"Content"`
    Timestamp string `dynamodbav:"Timestamp"`
}

func (dr *DynamoDBRepository) GetUser(userID string) (*UserItem, error) {
    result, err := dr.client.GetItem(context.Background(), &dynamodb.GetItemInput{
        TableName: aws.String(dr.table),
        Key: map[string]types.AttributeValue{
            "PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
            "SK": &types.AttributeValueMemberS{Value: "PROFILE"},
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    var user UserItem
    err = attributevalue.UnmarshalMap(result.Item, &user)
    return &user, err
}
```

### Column-Family Databases

#### Cassandra
```go
// Write-optimized, linear scalability
type CassandraSession struct {
    session *gocql.Session
}

// Time-series data model
func (cs *CassandraSession) CreateTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS user_timeline (
            user_id UUID,
            post_id TIMEUUID,
            content TEXT,
            author_id UUID,
            created_at TIMESTAMP,
            PRIMARY KEY (user_id, post_id)
        ) WITH CLUSTERING ORDER BY (post_id DESC)`,
        
        `CREATE TABLE IF NOT EXISTS post_by_author (
            author_id UUID,
            post_id TIMEUUID,
            content TEXT,
            created_at TIMESTAMP,
            PRIMARY KEY (author_id, post_id)
        ) WITH CLUSTERING ORDER BY (post_id DESC)`,
    }
    
    for _, query := range queries {
        if err := cs.session.Query(query).Exec(); err != nil {
            return err
        }
    }
    return nil
}

// Write to multiple tables (denormalization)
func (cs *CassandraSession) CreatePost(post *Post) error {
    batch := cs.session.NewBatch(gocql.LoggedBatch)
    
    // Insert into author's posts
    batch.Query(`INSERT INTO post_by_author (author_id, post_id, content, created_at) 
                 VALUES (?, ?, ?, ?)`,
        post.AuthorID, post.ID, post.Content, post.CreatedAt)
    
    // Insert into followers' timelines (fan-out on write)
    followers, err := cs.GetFollowers(post.AuthorID)
    if err != nil {
        return err
    }
    
    for _, followerID := range followers {
        batch.Query(`INSERT INTO user_timeline (user_id, post_id, content, author_id, created_at) 
                     VALUES (?, ?, ?, ?, ?)`,
            followerID, post.ID, post.Content, post.AuthorID, post.CreatedAt)
    }
    
    return cs.session.ExecuteBatch(batch)
}
```

### Graph Databases

#### Neo4j
```go
// Relationship-focused queries
type Neo4jDriver struct {
    driver neo4j.Driver
}

func (nd *Neo4jDriver) CreateUser(user *User) error {
    session := nd.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()
    
    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        result, err := tx.Run(
            "CREATE (u:User {id: $id, name: $name, email: $email}) RETURN u",
            map[string]interface{}{
                "id":    user.ID,
                "name":  user.Name,
                "email": user.Email,
            })
        return result, err
    })
    
    return err
}

// Complex relationship queries
func (nd *Neo4jDriver) GetFriendsOfFriends(userID string) ([]*User, error) {
    session := nd.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()
    
    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        return tx.Run(`
            MATCH (u:User {id: $userID})-[:FOLLOWS]->(friend:User)-[:FOLLOWS]->(fof:User)
            WHERE fof.id <> $userID AND NOT (u)-[:FOLLOWS]->(fof)
            RETURN DISTINCT fof.id, fof.name, fof.email
            LIMIT 50`,
            map[string]interface{}{"userID": userID})
    })
    
    if err != nil {
        return nil, err
    }
    
    // Process results...
    return users, nil
}
```

## 🔄 Database Patterns

### CQRS (Command Query Responsibility Segregation)
```go
// Separate read and write models
type UserCommandService struct {
    writeDB *sql.DB // PostgreSQL for consistency
}

type UserQueryService struct {
    readDB *mongo.Database // MongoDB for flexible queries
    cache  *redis.Client   // Redis for fast reads
}

func (ucs *UserCommandService) CreateUser(cmd CreateUserCommand) error {
    // Write to SQL database
    _, err := ucs.writeDB.Exec(
        "INSERT INTO users (id, email, name) VALUES ($1, $2, $3)",
        cmd.ID, cmd.Email, cmd.Name)
    
    if err != nil {
        return err
    }
    
    // Publish event for read model update
    eventBus.Publish(UserCreatedEvent{
        ID:    cmd.ID,
        Email: cmd.Email,
        Name:  cmd.Name,
    })
    
    return nil
}

func (uqs *UserQueryService) HandleUserCreated(event UserCreatedEvent) {
    // Update read model (MongoDB)
    userDoc := bson.M{
        "_id":   event.ID,
        "email": event.Email,
        "name":  event.Name,
        "profile": bson.M{
            "created_at": time.Now(),
        },
    }
    
    uqs.readDB.Collection("users").InsertOne(context.Background(), userDoc)
    
    // Update cache
    uqs.cache.HMSet(context.Background(), fmt.Sprintf("user:%s", event.ID), map[string]interface{}{
        "email": event.Email,
        "name":  event.Name,
    })
}
```

### Event Sourcing
```go
// Store events instead of current state
type EventStore struct {
    db *sql.DB
}

type Event struct {
    ID          string    `json:"id"`
    AggregateID string    `json:"aggregate_id"`
    Type        string    `json:"type"`
    Data        []byte    `json:"data"`
    Version     int       `json:"version"`
    Timestamp   time.Time `json:"timestamp"`
}

func (es *EventStore) AppendEvent(aggregateID string, event Event) error {
    _, err := es.db.Exec(`
        INSERT INTO events (id, aggregate_id, type, data, version, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6)`,
        event.ID, aggregateID, event.Type, event.Data, event.Version, event.Timestamp)
    
    return err
}

func (es *EventStore) GetEvents(aggregateID string) ([]Event, error) {
    rows, err := es.db.Query(`
        SELECT id, aggregate_id, type, data, version, timestamp
        FROM events
        WHERE aggregate_id = $1
        ORDER BY version ASC`, aggregateID)
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var events []Event
    for rows.Next() {
        var event Event
        err := rows.Scan(&event.ID, &event.AggregateID, &event.Type,
            &event.Data, &event.Version, &event.Timestamp)
        if err != nil {
            return nil, err
        }
        events = append(events, event)
    }
    
    return events, nil
}

// Rebuild state from events
func (es *EventStore) ReplayEvents(aggregateID string) (*UserAggregate, error) {
    events, err := es.GetEvents(aggregateID)
    if err != nil {
        return nil, err
    }
    
    user := &UserAggregate{ID: aggregateID}
    for _, event := range events {
        user.Apply(event)
    }
    
    return user, nil
}
```

## 🎯 Database Selection Examples

### E-commerce Platform
```go
type EcommerceDatabases struct {
    UserProfiles    *sql.DB        // PostgreSQL - ACID for user data
    ProductCatalog  *mongo.Database // MongoDB - flexible product schemas
    Inventory       *sql.DB        // PostgreSQL - consistency critical
    Orders          *sql.DB        // PostgreSQL - ACID transactions
    Sessions        *redis.Client   // Redis - fast session storage
    Analytics       *clickhouse.DB  // ClickHouse - analytical queries
    Search          *elasticsearch.Client // Elasticsearch - full-text search
}
```

### Social Media Platform
```go
type SocialMediaDatabases struct {
    UserProfiles    *sql.DB           // PostgreSQL - structured user data
    Posts           *cassandra.Session // Cassandra - high write volume
    Timelines       *redis.Client     // Redis - fast timeline access
    Relationships   *neo4j.Driver     // Neo4j - complex social graph
    Messages        *cassandra.Session // Cassandra - chat history
    Notifications   *redis.Client     // Redis - real-time notifications
}
```

### IoT/Time-Series Platform
```go
type IoTDatabases struct {
    DeviceMetadata  *sql.DB           // PostgreSQL - device information
    TimeSeries      *influxdb.Client  // InfluxDB - sensor data
    RealTimeData    *redis.Client     // Redis - current device states
    Analytics       *clickhouse.DB    // ClickHouse - historical analysis
    Alerts          *elasticsearch.Client // Elasticsearch - log analysis
}
```

## 📊 Performance Characteristics

### Throughput Comparison (Approximate)
```
Database Type    | Read QPS  | Write QPS | Use Case
-----------------|-----------|-----------|------------------
PostgreSQL       | 50K       | 20K       | OLTP, Complex queries
MySQL            | 100K      | 50K       | Web applications
MongoDB          | 80K       | 40K       | Document storage
Cassandra        | 200K      | 300K      | Time-series, logs
Redis            | 1M        | 500K      | Caching, sessions
DynamoDB         | 40K       | 40K       | Serverless apps
Neo4j            | 30K       | 15K       | Graph relationships
```

### Latency Characteristics
```
Operation Type   | PostgreSQL | MongoDB | Cassandra | Redis
-----------------|------------|---------|-----------|-------
Point Read       | 1-5ms      | 1-3ms   | 1-2ms     | <1ms
Range Query      | 5-50ms     | 5-30ms  | 2-10ms    | 1-5ms
Complex Join     | 10-500ms   | N/A     | N/A       | N/A
Write            | 1-10ms     | 1-5ms   | <1ms      | <1ms
```

## 🎯 Interview Decision Framework

### Questions to Ask Yourself
```
1. Data Structure
   - Is the data relational?
   - Do I need complex queries?
   - Is the schema fixed or flexible?

2. Consistency Requirements
   - Can I tolerate eventual consistency?
   - Do I need ACID transactions?
   - What's the cost of inconsistency?

3. Scale Requirements
   - What's the expected data size?
   - What's the read/write ratio?
   - Do I need global distribution?

4. Performance Requirements
   - What's the acceptable latency?
   - What's the required throughput?
   - Are there seasonal spikes?

5. Operational Requirements
   - What's the team's expertise?
   - What's the operational overhead?
   - What's the budget constraint?
```

### Decision Matrix Example
```
Requirement      | Weight | PostgreSQL | MongoDB | Cassandra | Redis
-----------------|--------|------------|---------|-----------|-------
ACID Compliance  | High   | 10         | 3       | 2         | 1
Complex Queries  | High   | 10         | 7       | 3         | 2
Write Scalability| Medium | 4          | 6       | 10        | 9
Read Performance | High   | 7          | 8       | 9         | 10
Operational Cost | Medium | 8          | 6       | 4         | 7
-----------------|--------|------------|---------|-----------|-------
Weighted Score   |        | 8.2        | 6.4     | 6.0       | 6.8
```

## 🔗 Next Steps
- Study caching strategies and implementation patterns
- Learn about load balancing algorithms and CDN architecture
- Explore message queues and event streaming systems
- Practice database selection for different system design problems
