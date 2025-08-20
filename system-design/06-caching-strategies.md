# 🚀 Caching Strategies & Implementation

## 🎯 Why Caching Matters

### Performance Impact
```
Without Cache: Database query = 100ms
With Cache: Memory lookup = 1ms
Performance improvement: 100x faster

Cost Impact:
- Database server: $500/month
- Cache server: $50/month
- Cost reduction: 90% for read-heavy workloads
```

### Cache Hit Rate Economics
```go
// Cache effectiveness calculation
func CalculateCacheEffectiveness(hitRate float64, cacheLatency, dbLatency time.Duration) time.Duration {
    effectiveLatency := hitRate*float64(cacheLatency) + (1-hitRate)*float64(dbLatency)
    return time.Duration(effectiveLatency)
}

// Example: 80% hit rate
// Effective latency = 0.8 * 1ms + 0.2 * 100ms = 20.8ms
// Still 5x improvement over no cache
```

## 🏗️ Cache Levels & Hierarchy

### Multi-Level Caching Architecture
```go
type CacheHierarchy struct {
    L1 *LRUCache      // Application memory (fastest, smallest)
    L2 *RedisCluster  // Distributed cache (fast, medium)
    L3 *CDN           // Edge cache (global, largest)
    DB *Database      // Source of truth (slowest)
}

func (ch *CacheHierarchy) Get(key string) (interface{}, error) {
    // L1 Cache (in-memory)
    if value, found := ch.L1.Get(key); found {
        return value, nil
    }
    
    // L2 Cache (Redis)
    if value, err := ch.L2.Get(key); err == nil {
        ch.L1.Set(key, value, 5*time.Minute) // Populate L1
        return value, nil
    }
    
    // L3 Cache (CDN) - for static content
    if value, err := ch.L3.Get(key); err == nil {
        ch.L2.Set(key, value, 1*time.Hour)   // Populate L2
        ch.L1.Set(key, value, 5*time.Minute) // Populate L1
        return value, nil
    }
    
    // Database (source of truth)
    value, err := ch.DB.Get(key)
    if err != nil {
        return nil, err
    }
    
    // Populate all cache levels
    go ch.L3.Set(key, value, 24*time.Hour)
    go ch.L2.Set(key, value, 1*time.Hour)
    ch.L1.Set(key, value, 5*time.Minute)
    
    return value, nil
}
```

### Application-Level Caching
```go
// LRU Cache implementation (similar to your algorithms collection)
type LRUCache struct {
    capacity int
    cache    map[string]*Node
    head     *Node
    tail     *Node
    mutex    sync.RWMutex
}

type Node struct {
    key   string
    value interface{}
    prev  *Node
    next  *Node
}

func NewLRUCache(capacity int) *LRUCache {
    lru := &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*Node),
    }
    
    // Create dummy head and tail
    lru.head = &Node{}
    lru.tail = &Node{}
    lru.head.next = lru.tail
    lru.tail.prev = lru.head
    
    return lru
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
    lru.mutex.Lock()
    defer lru.mutex.Unlock()
    
    if node, found := lru.cache[key]; found {
        // Move to head (most recently used)
        lru.moveToHead(node)
        return node.value, true
    }
    
    return nil, false
}

func (lru *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
    lru.mutex.Lock()
    defer lru.mutex.Unlock()
    
    if node, found := lru.cache[key]; found {
        // Update existing node
        node.value = value
        lru.moveToHead(node)
    } else {
        // Create new node
        newNode := &Node{key: key, value: value}
        
        if len(lru.cache) >= lru.capacity {
            // Remove least recently used
            tail := lru.removeTail()
            delete(lru.cache, tail.key)
        }
        
        lru.cache[key] = newNode
        lru.addToHead(newNode)
    }
    
    // Handle TTL with goroutine
    if ttl > 0 {
        go func() {
            time.Sleep(ttl)
            lru.Delete(key)
        }()
    }
}
```

## 🔄 Cache Patterns

### Cache Pattern Visual Overview
```
┌─────────────────────────────────────────────────────────────┐
│                    Cache Pattern Types                     │
└─────────────────────────────────────────────────────────────┘

1. Cache-Aside (Lazy Loading)
   App ──1──► Cache ──2(miss)──► Database ──3──► App ──4──► Cache

2. Write-Through
   App ──1──► Cache ──2──► Database

3. Write-Behind (Write-Back)
   App ──1──► Cache ──2(async)──► Database

4. Refresh-Ahead
   Cache ──1(predict)──► Database ──2──► Cache

┌─────────────────┬─────────────────┬─────────────────┐
│  Cache-Aside    │  Write-Through  │  Write-Behind   │
├─────────────────┼─────────────────┼─────────────────┤
│ • App manages   │ • Cache manages │ • Cache manages │
│   cache         │   writes        │   async writes  │
│ • Lazy loading  │ • Sync writes   │ • Better perf   │
│ • Cache misses  │ • Consistency   │ • Risk of loss  │
│ • Simple logic  │ • Write penalty │ • Complex logic │
└─────────────────┴─────────────────┴─────────────────┘
```

### 1. Cache-Aside (Lazy Loading)
```
┌─────────────────────────────────────────────────────────────┐
│                Cache-Aside Pattern Flow                    │
└─────────────────────────────────────────────────────────────┘

Read Flow:
┌─────────┐    ┌─────────┐    ┌─────────┐
│   App   │──1─►│  Cache  │    │Database │
│         │◄─2─│         │    │         │
│         │    │         │    │         │
│         │──3─┼─────────┼──4─►│         │
│         │◄─5─┼─────────┼────│         │
│         │──6─►│         │    │         │
└─────────┘    └─────────┘    └─────────┘

1. App checks cache
2. Cache miss returned
3. App queries database
4. Database returns data
5. App receives data
6. App stores in cache

Write Flow:
┌─────────┐    ┌─────────┐    ┌─────────┐
│   App   │──1─┼─────────┼──2─►│Database │
│         │    │  Cache  │    │         │
│         │──3─►│         │    │         │
└─────────┘    └─────────┘    └─────────┘

1. App writes to database
2. Database confirms write
3. App invalidates cache
```
```go
type CacheAsideService struct {
    cache Cache
    db    Database
}

func (cas *CacheAsideService) GetUser(userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // 1. Check cache first
    if cached, found := cas.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }
    
    // 2. Cache miss - load from database
    user, err := cas.db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Store in cache for future requests
    cas.cache.Set(cacheKey, user, 1*time.Hour)
    
    return user, nil
}

func (cas *CacheAsideService) UpdateUser(user *User) error {
    // 1. Update database
    if err := cas.db.UpdateUser(user); err != nil {
        return err
    }
    
    // 2. Invalidate cache (let next read populate it)
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    cas.cache.Delete(cacheKey)
    
    return nil
}

// Pros: Simple, cache failures don't affect system
// Cons: Cache miss penalty, potential stale data
```

### 2. Write-Through Cache
```go
type WriteThroughService struct {
    cache Cache
    db    Database
}

func (wts *WriteThroughService) UpdateUser(user *User) error {
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    
    // 1. Update database first
    if err := wts.db.UpdateUser(user); err != nil {
        return err
    }
    
    // 2. Update cache immediately
    wts.cache.Set(cacheKey, user, 1*time.Hour)
    
    return nil
}

func (wts *WriteThroughService) GetUser(userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Cache should always have fresh data
    if cached, found := wts.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }
    
    // Fallback to database (shouldn't happen often)
    return wts.db.GetUser(userID)
}

// Pros: Cache always consistent, no cache miss penalty for writes
// Cons: Write latency increased, unused data cached
```

### 3. Write-Behind (Write-Back)
```go
type WriteBehindService struct {
    cache     Cache
    db        Database
    writeQueue chan WriteOperation
    batchSize  int
}

type WriteOperation struct {
    Key   string
    Value interface{}
    Op    string // "UPDATE", "DELETE"
}

func NewWriteBehindService(cache Cache, db Database) *WriteBehindService {
    wbs := &WriteBehindService{
        cache:     cache,
        db:        db,
        writeQueue: make(chan WriteOperation, 10000),
        batchSize: 100,
    }
    
    // Start background writer
    go wbs.backgroundWriter()
    
    return wbs
}

func (wbs *WriteBehindService) UpdateUser(user *User) error {
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    
    // 1. Update cache immediately
    wbs.cache.Set(cacheKey, user, 1*time.Hour)
    
    // 2. Queue database write
    select {
    case wbs.writeQueue <- WriteOperation{
        Key:   cacheKey,
        Value: user,
        Op:    "UPDATE",
    }:
        return nil
    default:
        // Queue full - fallback to synchronous write
        return wbs.db.UpdateUser(user)
    }
}

func (wbs *WriteBehindService) backgroundWriter() {
    batch := make([]WriteOperation, 0, wbs.batchSize)
    ticker := time.NewTicker(5 * time.Second)
    
    for {
        select {
        case op := <-wbs.writeQueue:
            batch = append(batch, op)
            
            if len(batch) >= wbs.batchSize {
                wbs.flushBatch(batch)
                batch = batch[:0]
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                wbs.flushBatch(batch)
                batch = batch[:0]
            }
        }
    }
}

func (wbs *WriteBehindService) flushBatch(batch []WriteOperation) {
    // Batch write to database
    for _, op := range batch {
        switch op.Op {
        case "UPDATE":
            wbs.db.UpdateUser(op.Value.(*User))
        case "DELETE":
            wbs.db.DeleteUser(op.Key)
        }
    }
}

// Pros: Low write latency, high write throughput
// Cons: Data loss risk, complex consistency
```

### 4. Refresh-Ahead
```go
type RefreshAheadService struct {
    cache       Cache
    db          Database
    refresher   *time.Ticker
    hotKeys     map[string]time.Time
    mutex       sync.RWMutex
}

func (ras *RefreshAheadService) GetUser(userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // Track hot keys
    ras.mutex.Lock()
    ras.hotKeys[cacheKey] = time.Now()
    ras.mutex.Unlock()
    
    if cached, found := ras.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }
    
    // Cache miss - load from database
    user, err := ras.db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    ras.cache.Set(cacheKey, user, 1*time.Hour)
    return user, nil
}

func (ras *RefreshAheadService) refreshHotKeys() {
    ras.mutex.RLock()
    hotKeys := make([]string, 0, len(ras.hotKeys))
    for key, lastAccess := range ras.hotKeys {
        if time.Since(lastAccess) < 10*time.Minute {
            hotKeys = append(hotKeys, key)
        }
    }
    ras.mutex.RUnlock()
    
    // Refresh hot keys before they expire
    for _, key := range hotKeys {
        go func(k string) {
            userID := strings.TrimPrefix(k, "user:")
            if user, err := ras.db.GetUser(userID); err == nil {
                ras.cache.Set(k, user, 1*time.Hour)
            }
        }(key)
    }
}

// Pros: No cache miss penalty for hot data
// Cons: Additional complexity, may refresh unused data
```

## 🌐 Distributed Caching

### Redis Cluster Implementation
```go
type RedisCluster struct {
    clients []*redis.Client
    slots   []int // Consistent hashing slots
}

func NewRedisCluster(addresses []string) *RedisCluster {
    clients := make([]*redis.Client, len(addresses))
    for i, addr := range addresses {
        clients[i] = redis.NewClient(&redis.Options{
            Addr: addr,
        })
    }
    
    return &RedisCluster{
        clients: clients,
        slots:   make([]int, 16384), // Redis cluster has 16384 slots
    }
}

func (rc *RedisCluster) getClient(key string) *redis.Client {
    // Calculate hash slot (similar to your consistent hashing algorithm)
    slot := crc16(key) % 16384
    clientIndex := rc.slots[slot]
    return rc.clients[clientIndex]
}

func (rc *RedisCluster) Set(key string, value interface{}, expiration time.Duration) error {
    client := rc.getClient(key)
    return client.Set(context.Background(), key, value, expiration).Err()
}

func (rc *RedisCluster) Get(key string) (string, error) {
    client := rc.getClient(key)
    return client.Get(context.Background(), key).Result()
}

// Handle node failures
func (rc *RedisCluster) SetWithFailover(key string, value interface{}, expiration time.Duration) error {
    client := rc.getClient(key)
    
    err := client.Set(context.Background(), key, value, expiration).Err()
    if err != nil {
        // Try replica or next node
        for i := 0; i < len(rc.clients); i++ {
            if err = rc.clients[i].Set(context.Background(), key, value, expiration).Err(); err == nil {
                break
            }
        }
    }
    
    return err
}
```

### Cache Partitioning Strategies
```go
// 1. Hash-based partitioning
func (c *DistributedCache) GetPartition(key string) int {
    hash := fnv.New32a()
    hash.Write([]byte(key))
    return int(hash.Sum32()) % len(c.partitions)
}

// 2. Range-based partitioning
func (c *DistributedCache) GetPartitionByRange(key string) int {
    if key < "m" {
        return 0 // a-l
    } else {
        return 1 // m-z
    }
}

// 3. Consistent hashing (from your algorithms collection)
type ConsistentHash struct {
    ring     map[uint32]string
    sortedKeys []uint32
    replicas int
}

func (ch *ConsistentHash) GetNode(key string) string {
    if len(ch.ring) == 0 {
        return ""
    }
    
    hash := ch.hashKey(key)
    idx := sort.Search(len(ch.sortedKeys), func(i int) bool {
        return ch.sortedKeys[i] >= hash
    })
    
    if idx == len(ch.sortedKeys) {
        idx = 0
    }
    
    return ch.ring[ch.sortedKeys[idx]]
}
```

## 📊 Cache Optimization Techniques

### Cache Warming
```go
type CacheWarmer struct {
    cache Cache
    db    Database
}

func (cw *CacheWarmer) WarmUserCache() error {
    // Pre-populate cache with frequently accessed users
    popularUsers, err := cw.db.GetPopularUsers(1000)
    if err != nil {
        return err
    }
    
    // Warm cache in batches to avoid overwhelming the system
    batchSize := 50
    for i := 0; i < len(popularUsers); i += batchSize {
        end := min(i+batchSize, len(popularUsers))
        batch := popularUsers[i:end]
        
        for _, user := range batch {
            cacheKey := fmt.Sprintf("user:%s", user.ID)
            cw.cache.Set(cacheKey, user, 2*time.Hour)
        }
        
        // Rate limiting
        time.Sleep(100 * time.Millisecond)
    }
    
    return nil
}

// Scheduled cache warming
func (cw *CacheWarmer) StartScheduledWarming() {
    ticker := time.NewTicker(6 * time.Hour)
    go func() {
        for range ticker.C {
            cw.WarmUserCache()
        }
    }()
}
```

### Cache Compression
```go
type CompressedCache struct {
    cache Cache
}

func (cc *CompressedCache) Set(key string, value interface{}, expiration time.Duration) error {
    // Serialize and compress
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    var compressed bytes.Buffer
    writer := gzip.NewWriter(&compressed)
    writer.Write(data)
    writer.Close()
    
    return cc.cache.Set(key, compressed.Bytes(), expiration)
}

func (cc *CompressedCache) Get(key string) (interface{}, error) {
    compressed, found := cc.cache.Get(key)
    if !found {
        return nil, ErrCacheMiss
    }
    
    // Decompress and deserialize
    reader, err := gzip.NewReader(bytes.NewReader(compressed.([]byte)))
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    var data bytes.Buffer
    data.ReadFrom(reader)
    
    var result interface{}
    err = json.Unmarshal(data.Bytes(), &result)
    return result, err
}

// Compression effectiveness
// Text data: 60-80% size reduction
// JSON data: 40-60% size reduction
// Already compressed data: 0-10% reduction
```

## 🎯 Cache Invalidation Strategies

### Time-Based Expiration (TTL)
```go
type TTLCache struct {
    data map[string]CacheItem
    mutex sync.RWMutex
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

func (tc *TTLCache) Set(key string, value interface{}, ttl time.Duration) {
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    tc.data[key] = CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (tc *TTLCache) Get(key string) (interface{}, bool) {
    tc.mutex.RLock()
    defer tc.mutex.RUnlock()
    
    item, found := tc.data[key]
    if !found || time.Now().After(item.ExpiresAt) {
        return nil, false
    }
    
    return item.Value, true
}

// Background cleanup
func (tc *TTLCache) startCleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for range ticker.C {
            tc.cleanup()
        }
    }()
}

func (tc *TTLCache) cleanup() {
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    now := time.Now()
    for key, item := range tc.data {
        if now.After(item.ExpiresAt) {
            delete(tc.data, key)
        }
    }
}
```

### Event-Based Invalidation
```go
type EventBasedCache struct {
    cache     Cache
    eventBus  EventBus
}

func (ebc *EventBasedCache) Initialize() {
    // Subscribe to relevant events
    ebc.eventBus.Subscribe("user.updated", ebc.handleUserUpdated)
    ebc.eventBus.Subscribe("user.deleted", ebc.handleUserDeleted)
    ebc.eventBus.Subscribe("post.created", ebc.handlePostCreated)
}

func (ebc *EventBasedCache) handleUserUpdated(event UserUpdatedEvent) {
    // Invalidate user cache
    userKey := fmt.Sprintf("user:%s", event.UserID)
    ebc.cache.Delete(userKey)
    
    // Invalidate related caches
    timelineKey := fmt.Sprintf("timeline:%s", event.UserID)
    ebc.cache.Delete(timelineKey)
}

func (ebc *EventBasedCache) handlePostCreated(event PostCreatedEvent) {
    // Invalidate author's timeline
    authorTimelineKey := fmt.Sprintf("timeline:%s", event.AuthorID)
    ebc.cache.Delete(authorTimelineKey)
    
    // Invalidate followers' timelines
    followers, _ := ebc.getFollowers(event.AuthorID)
    for _, followerID := range followers {
        followerTimelineKey := fmt.Sprintf("timeline:%s", followerID)
        ebc.cache.Delete(followerTimelineKey)
    }
}
```

## 🔗 Next Steps
- Study load balancing algorithms and CDN architecture
- Learn about message queues and event streaming systems
- Explore API design and communication patterns
