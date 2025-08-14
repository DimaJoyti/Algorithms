# ⚡ Performance & Optimization

## 🎯 Performance Optimization Overview

Performance optimization is the practice of making systems faster, more efficient, and more scalable. It involves identifying bottlenecks, optimizing algorithms, improving resource utilization, and implementing caching strategies.

### Key Performance Metrics
```
✅ Latency Metrics:
- Response time (P50, P95, P99)
- Time to first byte (TTFB)
- Round-trip time (RTT)
- Processing time

✅ Throughput Metrics:
- Requests per second (RPS/QPS)
- Transactions per second (TPS)
- Bandwidth utilization
- Data transfer rates

✅ Resource Utilization:
- CPU usage and efficiency
- Memory consumption
- Disk I/O operations
- Network bandwidth

✅ Scalability Metrics:
- Horizontal scaling efficiency
- Load distribution
- Resource contention
- Queue lengths
```

## 🧠 CPU Optimization

### Algorithm Optimization
```go
// Inefficient O(n²) algorithm
func findDuplicatesNaive(nums []int) []int {
    var duplicates []int
    seen := make(map[int]bool)
    
    for i := 0; i < len(nums); i++ {
        for j := i + 1; j < len(nums); j++ {
            if nums[i] == nums[j] && !seen[nums[i]] {
                duplicates = append(duplicates, nums[i])
                seen[nums[i]] = true
            }
        }
    }
    return duplicates
}

// Optimized O(n) algorithm
func findDuplicatesOptimized(nums []int) []int {
    seen := make(map[int]bool)
    duplicates := make(map[int]bool)
    var result []int
    
    for _, num := range nums {
        if seen[num] && !duplicates[num] {
            result = append(result, num)
            duplicates[num] = true
        }
        seen[num] = true
    }
    return result
}

// CPU-intensive task optimization with goroutines
type WorkerPool struct {
    workerCount int
    jobQueue    chan Job
    resultQueue chan Result
    workers     []*Worker
}

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Data  interface{}
    Error error
}

func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobQueue:    make(chan Job, queueSize),
        resultQueue: make(chan Result, queueSize),
        workers:     make([]*Worker, workerCount),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workerCount; i++ {
        worker := &Worker{
            id:          i,
            jobQueue:    wp.jobQueue,
            resultQueue: wp.resultQueue,
            quit:        make(chan bool),
        }
        wp.workers[i] = worker
        go worker.Start()
    }
}

type Worker struct {
    id          int
    jobQueue    chan Job
    resultQueue chan Result
    quit        chan bool
}

func (w *Worker) Start() {
    for {
        select {
        case job := <-w.jobQueue:
            result := w.processJob(job)
            w.resultQueue <- result
        case <-w.quit:
            return
        }
    }
}

func (w *Worker) processJob(job Job) Result {
    // Simulate CPU-intensive work
    start := time.Now()
    
    // Example: Complex mathematical computation
    result := w.performComplexCalculation(job.Data)
    
    return Result{
        JobID: job.ID,
        Data:  result,
        Error: nil,
    }
}

// SIMD optimization example (using assembly or specialized libraries)
func vectorAddOptimized(a, b []float64) []float64 {
    if len(a) != len(b) {
        panic("vectors must have same length")
    }
    
    result := make([]float64, len(a))
    
    // Process in chunks of 4 for SIMD optimization
    chunkSize := 4
    for i := 0; i < len(a)-chunkSize+1; i += chunkSize {
        // This would use SIMD instructions in optimized implementation
        for j := 0; j < chunkSize; j++ {
            result[i+j] = a[i+j] + b[i+j]
        }
    }
    
    // Handle remaining elements
    for i := len(a) - (len(a) % chunkSize); i < len(a); i++ {
        result[i] = a[i] + b[i]
    }
    
    return result
}
```

### CPU Cache Optimization
```go
// Cache-friendly data structures
type CacheFriendlyMatrix struct {
    data   []float64
    rows   int
    cols   int
}

func NewCacheFriendlyMatrix(rows, cols int) *CacheFriendlyMatrix {
    return &CacheFriendlyMatrix{
        data: make([]float64, rows*cols),
        rows: rows,
        cols: cols,
    }
}

func (m *CacheFriendlyMatrix) Get(row, col int) float64 {
    return m.data[row*m.cols+col]
}

func (m *CacheFriendlyMatrix) Set(row, col int, value float64) {
    m.data[row*m.cols+col] = value
}

// Cache-friendly matrix multiplication
func (m *CacheFriendlyMatrix) Multiply(other *CacheFriendlyMatrix) *CacheFriendlyMatrix {
    if m.cols != other.rows {
        panic("incompatible matrix dimensions")
    }
    
    result := NewCacheFriendlyMatrix(m.rows, other.cols)
    
    // Block-wise multiplication for better cache locality
    blockSize := 64 // Optimize based on cache line size
    
    for i := 0; i < m.rows; i += blockSize {
        for j := 0; j < other.cols; j += blockSize {
            for k := 0; k < m.cols; k += blockSize {
                // Process block
                iMax := min(i+blockSize, m.rows)
                jMax := min(j+blockSize, other.cols)
                kMax := min(k+blockSize, m.cols)
                
                for ii := i; ii < iMax; ii++ {
                    for jj := j; jj < jMax; jj++ {
                        sum := result.Get(ii, jj)
                        for kk := k; kk < kMax; kk++ {
                            sum += m.Get(ii, kk) * other.Get(kk, jj)
                        }
                        result.Set(ii, jj, sum)
                    }
                }
            }
        }
    }
    
    return result
}

// CPU profiling and optimization
type CPUProfiler struct {
    samples    []CPUSample
    startTime  time.Time
    sampleRate time.Duration
}

type CPUSample struct {
    Timestamp time.Time
    CPUUsage  float64
    Function  string
    Duration  time.Duration
}

func (cp *CPUProfiler) StartProfiling(sampleRate time.Duration) {
    cp.startTime = time.Now()
    cp.sampleRate = sampleRate
    cp.samples = make([]CPUSample, 0)
    
    go cp.collectSamples()
}

func (cp *CPUProfiler) collectSamples() {
    ticker := time.NewTicker(cp.sampleRate)
    defer ticker.Stop()
    
    for range ticker.C {
        sample := CPUSample{
            Timestamp: time.Now(),
            CPUUsage:  cp.getCurrentCPUUsage(),
        }
        cp.samples = append(cp.samples, sample)
    }
}

func (cp *CPUProfiler) ProfileFunction(name string, fn func()) {
    start := time.Now()
    fn()
    duration := time.Since(start)
    
    sample := CPUSample{
        Timestamp: start,
        Function:  name,
        Duration:  duration,
    }
    cp.samples = append(cp.samples, sample)
}
```

## 💾 Memory Optimization

### Memory Pool Pattern
```go
// Object pool for reducing GC pressure
type ObjectPool struct {
    pool sync.Pool
    new  func() interface{}
}

func NewObjectPool(newFunc func() interface{}) *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{New: newFunc},
        new:  newFunc,
    }
}

func (op *ObjectPool) Get() interface{} {
    return op.pool.Get()
}

func (op *ObjectPool) Put(obj interface{}) {
    // Reset object state before returning to pool
    if resetter, ok := obj.(interface{ Reset() }); ok {
        resetter.Reset()
    }
    op.pool.Put(obj)
}

// Buffer pool for byte slices
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 4096) // 4KB initial capacity
    },
}

func GetBuffer() []byte {
    return bufferPool.Get().([]byte)
}

func PutBuffer(buf []byte) {
    if cap(buf) > 64*1024 { // Don't pool very large buffers
        return
    }
    buf = buf[:0] // Reset length but keep capacity
    bufferPool.Put(buf)
}

// Memory-efficient string builder
type MemoryEfficientBuilder struct {
    chunks [][]byte
    size   int
}

func (meb *MemoryEfficientBuilder) WriteString(s string) {
    chunk := []byte(s)
    meb.chunks = append(meb.chunks, chunk)
    meb.size += len(chunk)
}

func (meb *MemoryEfficientBuilder) WriteByte(b byte) {
    meb.chunks = append(meb.chunks, []byte{b})
    meb.size++
}

func (meb *MemoryEfficientBuilder) String() string {
    if len(meb.chunks) == 0 {
        return ""
    }
    
    // Allocate final buffer once
    result := make([]byte, 0, meb.size)
    for _, chunk := range meb.chunks {
        result = append(result, chunk...)
    }
    
    return string(result)
}

// Memory monitoring and optimization
type MemoryMonitor struct {
    samples    []MemorySample
    threshold  uint64
    callbacks  []MemoryCallback
}

type MemorySample struct {
    Timestamp    time.Time
    HeapAlloc    uint64
    HeapSys      uint64
    HeapIdle     uint64
    HeapInuse    uint64
    GCCycles     uint32
}

func (mm *MemoryMonitor) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            sample := MemorySample{
                Timestamp: time.Now(),
                HeapAlloc: m.HeapAlloc,
                HeapSys:   m.HeapSys,
                HeapIdle:  m.HeapIdle,
                HeapInuse: m.HeapInuse,
                GCCycles:  m.NumGC,
            }
            
            mm.samples = append(mm.samples, sample)
            
            // Check threshold
            if sample.HeapAlloc > mm.threshold {
                mm.triggerCallbacks(sample)
            }
        }
    }()
}

func (mm *MemoryMonitor) triggerCallbacks(sample MemorySample) {
    for _, callback := range mm.callbacks {
        go callback.OnMemoryThreshold(sample)
    }
}

// Zero-copy string operations
func ZeroCopySubstring(s string, start, end int) string {
    if start < 0 || end > len(s) || start > end {
        return ""
    }
    
    // Use unsafe to avoid copying
    return s[start:end]
}

// Memory-mapped file for large data processing
type MemoryMappedFile struct {
    file   *os.File
    data   []byte
    size   int64
}

func NewMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
    file, err := os.OpenFile(filename, os.O_RDWR, 0644)
    if err != nil {
        return nil, err
    }
    
    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }
    
    // Memory map the file (platform-specific implementation)
    data, err := syscall.Mmap(int(file.Fd()), 0, int(stat.Size()), 
        syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
    if err != nil {
        file.Close()
        return nil, err
    }
    
    return &MemoryMappedFile{
        file: file,
        data: data,
        size: stat.Size(),
    }, nil
}

func (mmf *MemoryMappedFile) Read(offset int64, length int) ([]byte, error) {
    if offset+int64(length) > mmf.size {
        return nil, fmt.Errorf("read beyond file size")
    }
    
    return mmf.data[offset:offset+int64(length)], nil
}

func (mmf *MemoryMappedFile) Close() error {
    if err := syscall.Munmap(mmf.data); err != nil {
        return err
    }
    return mmf.file.Close()
}
```

## 🌐 Network Optimization

### Connection Pooling and Reuse
```go
// HTTP connection pool with keep-alive
type HTTPConnectionPool struct {
    client      *http.Client
    transport   *http.Transport
    maxIdleConns int
    maxConnsPerHost int
    idleTimeout time.Duration
}

func NewHTTPConnectionPool(config PoolConfig) *HTTPConnectionPool {
    transport := &http.Transport{
        MaxIdleConns:        config.MaxIdleConns,
        MaxIdleConnsPerHost: config.MaxConnsPerHost,
        IdleConnTimeout:     config.IdleTimeout,
        DisableKeepAlives:   false,
        
        // TCP optimization
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        
        // TLS optimization
        TLSHandshakeTimeout: 10 * time.Second,
        
        // HTTP/2 support
        ForceAttemptHTTP2: true,
    }
    
    client := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    
    return &HTTPConnectionPool{
        client:    client,
        transport: transport,
        maxIdleConns: config.MaxIdleConns,
        maxConnsPerHost: config.MaxConnsPerHost,
        idleTimeout: config.IdleTimeout,
    }
}

// Request batching for reducing network overhead
type RequestBatcher struct {
    batchSize     int
    flushInterval time.Duration
    requests      []Request
    mutex         sync.Mutex
    flushTimer    *time.Timer
    processor     BatchProcessor
}

func NewRequestBatcher(batchSize int, flushInterval time.Duration, processor BatchProcessor) *RequestBatcher {
    rb := &RequestBatcher{
        batchSize:     batchSize,
        flushInterval: flushInterval,
        requests:      make([]Request, 0, batchSize),
        processor:     processor,
    }
    
    rb.resetFlushTimer()
    return rb
}

func (rb *RequestBatcher) AddRequest(req Request) {
    rb.mutex.Lock()
    defer rb.mutex.Unlock()
    
    rb.requests = append(rb.requests, req)
    
    if len(rb.requests) >= rb.batchSize {
        rb.flush()
    }
}

func (rb *RequestBatcher) flush() {
    if len(rb.requests) == 0 {
        return
    }
    
    batch := make([]Request, len(rb.requests))
    copy(batch, rb.requests)
    rb.requests = rb.requests[:0]
    
    rb.resetFlushTimer()
    
    // Process batch asynchronously
    go rb.processor.ProcessBatch(batch)
}

func (rb *RequestBatcher) resetFlushTimer() {
    if rb.flushTimer != nil {
        rb.flushTimer.Stop()
    }
    
    rb.flushTimer = time.AfterFunc(rb.flushInterval, func() {
        rb.mutex.Lock()
        rb.flush()
        rb.mutex.Unlock()
    })
}

// Network compression
type CompressionHandler struct {
    compressionLevel int
    threshold        int // Minimum size to compress
}

func (ch *CompressionHandler) CompressResponse(w http.ResponseWriter, r *http.Request, data []byte) error {
    // Check if client accepts compression
    if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Write(data)
        return nil
    }
    
    // Only compress if data is large enough
    if len(data) < ch.threshold {
        w.Write(data)
        return nil
    }
    
    // Set compression headers
    w.Header().Set("Content-Encoding", "gzip")
    w.Header().Set("Vary", "Accept-Encoding")
    
    // Compress data
    gz, err := gzip.NewWriterLevel(w, ch.compressionLevel)
    if err != nil {
        return err
    }
    defer gz.Close()
    
    _, err = gz.Write(data)
    return err
}
```

## 🗄️ Database Query Optimization

### Query Performance Tuning
```go
// Query optimizer with execution plan analysis
type QueryOptimizer struct {
    db          *sql.DB
    cache       QueryCache
    analyzer    QueryAnalyzer
    indexHints  IndexHintManager
}

func (qo *QueryOptimizer) OptimizeQuery(query string, params []interface{}) (*OptimizedQuery, error) {
    // Parse and analyze query
    parsed := qo.analyzer.ParseQuery(query)

    // Check for common anti-patterns
    if issues := qo.analyzer.FindPerformanceIssues(parsed); len(issues) > 0 {
        return nil, fmt.Errorf("query performance issues: %v", issues)
    }

    // Generate execution plan
    plan, err := qo.getExecutionPlan(query, params)
    if err != nil {
        return nil, err
    }

    // Suggest optimizations
    optimizations := qo.suggestOptimizations(plan)

    // Apply automatic optimizations
    optimizedQuery := qo.applyOptimizations(query, optimizations)

    return &OptimizedQuery{
        OriginalQuery:  query,
        OptimizedQuery: optimizedQuery,
        ExecutionPlan:  plan,
        Optimizations:  optimizations,
    }, nil
}

func (qo *QueryOptimizer) getExecutionPlan(query string, params []interface{}) (*ExecutionPlan, error) {
    explainQuery := "EXPLAIN ANALYZE " + query
    rows, err := qo.db.Query(explainQuery, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var plan ExecutionPlan
    for rows.Next() {
        var step ExecutionStep
        err := rows.Scan(&step.Operation, &step.Cost, &step.Rows, &step.Time)
        if err != nil {
            return nil, err
        }
        plan.Steps = append(plan.Steps, step)
    }

    return &plan, nil
}

// Index optimization
type IndexOptimizer struct {
    db              *sql.DB
    queryLog        QueryLog
    indexAnalyzer   IndexAnalyzer
    usageStats      IndexUsageStats
}

func (io *IndexOptimizer) AnalyzeIndexUsage() (*IndexAnalysisReport, error) {
    // Get current indexes
    indexes, err := io.getCurrentIndexes()
    if err != nil {
        return nil, err
    }

    // Analyze query patterns
    queryPatterns := io.queryLog.GetQueryPatterns()

    var report IndexAnalysisReport

    // Find unused indexes
    for _, index := range indexes {
        if !io.usageStats.IsIndexUsed(index.Name) {
            report.UnusedIndexes = append(report.UnusedIndexes, index)
        }
    }

    // Suggest missing indexes
    for _, pattern := range queryPatterns {
        if suggestedIndex := io.suggestIndexForPattern(pattern); suggestedIndex != nil {
            report.SuggestedIndexes = append(report.SuggestedIndexes, *suggestedIndex)
        }
    }

    // Find duplicate indexes
    duplicates := io.findDuplicateIndexes(indexes)
    report.DuplicateIndexes = duplicates

    return &report, nil
}

func (io *IndexOptimizer) suggestIndexForPattern(pattern QueryPattern) *SuggestedIndex {
    // Analyze WHERE clauses
    whereColumns := pattern.WhereColumns
    orderByColumns := pattern.OrderByColumns

    if len(whereColumns) == 0 {
        return nil
    }

    // Create composite index suggestion
    var indexColumns []string
    indexColumns = append(indexColumns, whereColumns...)

    // Add ORDER BY columns to avoid sorting
    for _, col := range orderByColumns {
        if !contains(indexColumns, col) {
            indexColumns = append(indexColumns, col)
        }
    }

    return &SuggestedIndex{
        TableName: pattern.TableName,
        Columns:   indexColumns,
        Type:      "BTREE",
        Reason:    fmt.Sprintf("Optimize queries with WHERE %v", whereColumns),
    }
}

// Connection pooling optimization
type DatabaseConnectionPool struct {
    db              *sql.DB
    maxOpenConns    int
    maxIdleConns    int
    connMaxLifetime time.Duration
    connMaxIdleTime time.Duration
    metrics         PoolMetrics
}

func NewDatabaseConnectionPool(config DBPoolConfig) *DatabaseConnectionPool {
    db, err := sql.Open(config.Driver, config.DSN)
    if err != nil {
        panic(err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

    pool := &DatabaseConnectionPool{
        db:              db,
        maxOpenConns:    config.MaxOpenConns,
        maxIdleConns:    config.MaxIdleConns,
        connMaxLifetime: config.ConnMaxLifetime,
        connMaxIdleTime: config.ConnMaxIdleTime,
        metrics:         NewPoolMetrics(),
    }

    // Start monitoring
    go pool.monitorPool()

    return pool
}

func (dcp *DatabaseConnectionPool) monitorPool() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        stats := dcp.db.Stats()

        dcp.metrics.OpenConnections.Set(float64(stats.OpenConnections))
        dcp.metrics.InUseConnections.Set(float64(stats.InUse))
        dcp.metrics.IdleConnections.Set(float64(stats.Idle))
        dcp.metrics.WaitCount.Add(float64(stats.WaitCount))
        dcp.metrics.WaitDuration.Add(float64(stats.WaitDuration.Milliseconds()))

        // Alert if pool is under pressure
        if float64(stats.InUse)/float64(dcp.maxOpenConns) > 0.8 {
            log.Printf("Connection pool utilization high: %d/%d", stats.InUse, dcp.maxOpenConns)
        }
    }
}

// Read replica routing
type ReadReplicaRouter struct {
    master   *sql.DB
    replicas []*sql.DB
    selector ReplicaSelector
    health   HealthChecker
}

func (rrr *ReadReplicaRouter) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    // Route read queries to replicas
    if rrr.isReadQuery(query) {
        replica := rrr.selector.SelectReplica(rrr.replicas)
        if replica != nil && rrr.health.IsHealthy(replica) {
            return replica.QueryContext(ctx, query, args...)
        }
    }

    // Fallback to master
    return rrr.master.QueryContext(ctx, query, args...)
}

func (rrr *ReadReplicaRouter) isReadQuery(query string) bool {
    query = strings.TrimSpace(strings.ToUpper(query))
    return strings.HasPrefix(query, "SELECT")
}
```

### Advanced Caching Strategies
```go
// Multi-level cache hierarchy
type MultiLevelCache struct {
    l1Cache *LRUCache      // In-memory cache
    l2Cache *RedisCache    // Distributed cache
    l3Cache *DatabaseCache // Database cache
    metrics CacheMetrics
}

func (mlc *MultiLevelCache) Get(key string) (interface{}, bool) {
    // L1 Cache (fastest)
    if value, found := mlc.l1Cache.Get(key); found {
        mlc.metrics.L1Hits.Inc()
        return value, true
    }

    // L2 Cache (distributed)
    if value, found := mlc.l2Cache.Get(key); found {
        mlc.metrics.L2Hits.Inc()
        // Promote to L1
        mlc.l1Cache.Set(key, value, 5*time.Minute)
        return value, true
    }

    // L3 Cache (database)
    if value, found := mlc.l3Cache.Get(key); found {
        mlc.metrics.L3Hits.Inc()
        // Promote to L2 and L1
        mlc.l2Cache.Set(key, value, 30*time.Minute)
        mlc.l1Cache.Set(key, value, 5*time.Minute)
        return value, true
    }

    mlc.metrics.Misses.Inc()
    return nil, false
}

func (mlc *MultiLevelCache) Set(key string, value interface{}, ttl time.Duration) {
    // Set in all levels with appropriate TTLs
    mlc.l1Cache.Set(key, value, min(ttl, 5*time.Minute))
    mlc.l2Cache.Set(key, value, min(ttl, 30*time.Minute))
    mlc.l3Cache.Set(key, value, ttl)
}

// Cache warming strategy
type CacheWarmer struct {
    cache       Cache
    dataSource  DataSource
    scheduler   Scheduler
    patterns    []WarmingPattern
}

type WarmingPattern struct {
    KeyPattern  string
    Frequency   time.Duration
    Priority    int
    Condition   func() bool
}

func (cw *CacheWarmer) StartWarming() {
    for _, pattern := range cw.patterns {
        go cw.warmPattern(pattern)
    }
}

func (cw *CacheWarmer) warmPattern(pattern WarmingPattern) {
    ticker := time.NewTicker(pattern.Frequency)
    defer ticker.Stop()

    for range ticker.C {
        if pattern.Condition != nil && !pattern.Condition() {
            continue
        }

        keys := cw.generateKeysFromPattern(pattern.KeyPattern)
        for _, key := range keys {
            go cw.warmKey(key, pattern.Priority)
        }
    }
}

func (cw *CacheWarmer) warmKey(key string, priority int) {
    // Check if key is already cached
    if _, found := cw.cache.Get(key); found {
        return
    }

    // Fetch from data source
    value, err := cw.dataSource.Get(key)
    if err != nil {
        return
    }

    // Cache with priority-based TTL
    ttl := cw.calculateTTL(priority)
    cw.cache.Set(key, value, ttl)
}

// Intelligent cache invalidation
type SmartCacheInvalidator struct {
    cache           Cache
    dependencyGraph DependencyGraph
    eventBus        EventBus
    invalidationLog InvalidationLog
}

func (sci *SmartCacheInvalidator) InvalidateByDependency(changedKey string) {
    // Find all keys that depend on the changed key
    dependentKeys := sci.dependencyGraph.GetDependents(changedKey)

    var invalidatedKeys []string
    for _, key := range dependentKeys {
        if sci.cache.Delete(key) {
            invalidatedKeys = append(invalidatedKeys, key)
        }
    }

    // Log invalidation for analysis
    sci.invalidationLog.Record(InvalidationEvent{
        TriggerKey:      changedKey,
        InvalidatedKeys: invalidatedKeys,
        Timestamp:       time.Now(),
    })

    // Publish invalidation event
    sci.eventBus.Publish(CacheInvalidationEvent{
        Keys:      invalidatedKeys,
        Reason:    "dependency_change",
        Timestamp: time.Now(),
    })
}

// Cache-aside pattern with write-through
type CacheAsideManager struct {
    cache      Cache
    dataStore  DataStore
    serializer Serializer
    metrics    CacheMetrics
}

func (cam *CacheAsideManager) Get(key string) (interface{}, error) {
    // Try cache first
    if cached, found := cam.cache.Get(key); found {
        cam.metrics.Hits.Inc()
        return cam.serializer.Deserialize(cached.([]byte))
    }

    // Cache miss - fetch from data store
    cam.metrics.Misses.Inc()
    value, err := cam.dataStore.Get(key)
    if err != nil {
        return nil, err
    }

    // Serialize and cache
    serialized, err := cam.serializer.Serialize(value)
    if err == nil {
        cam.cache.Set(key, serialized, 1*time.Hour)
    }

    return value, nil
}

func (cam *CacheAsideManager) Set(key string, value interface{}) error {
    // Write to data store first
    if err := cam.dataStore.Set(key, value); err != nil {
        return err
    }

    // Update cache
    serialized, err := cam.serializer.Serialize(value)
    if err == nil {
        cam.cache.Set(key, serialized, 1*time.Hour)
    }

    return nil
}
```

### Load Balancing Optimization
```go
// Weighted round-robin load balancer
type WeightedRoundRobinBalancer struct {
    servers         []Server
    currentWeights  []int
    effectiveWeights []int
    totalWeight     int
    mutex           sync.RWMutex
}

func NewWeightedRoundRobinBalancer(servers []Server) *WeightedRoundRobinBalancer {
    balancer := &WeightedRoundRobinBalancer{
        servers:         servers,
        currentWeights:  make([]int, len(servers)),
        effectiveWeights: make([]int, len(servers)),
    }

    // Initialize weights
    for i, server := range servers {
        balancer.effectiveWeights[i] = server.Weight
        balancer.totalWeight += server.Weight
    }

    return balancer
}

func (wrr *WeightedRoundRobinBalancer) NextServer() *Server {
    wrr.mutex.Lock()
    defer wrr.mutex.Unlock()

    if len(wrr.servers) == 0 {
        return nil
    }

    // Find server with highest current weight
    bestIndex := -1
    for i := range wrr.servers {
        wrr.currentWeights[i] += wrr.effectiveWeights[i]

        if bestIndex == -1 || wrr.currentWeights[i] > wrr.currentWeights[bestIndex] {
            bestIndex = i
        }
    }

    // Reduce current weight by total weight
    wrr.currentWeights[bestIndex] -= wrr.totalWeight

    return &wrr.servers[bestIndex]
}

// Least connections load balancer
type LeastConnectionsBalancer struct {
    servers []ServerWithConnections
    mutex   sync.RWMutex
}

type ServerWithConnections struct {
    Server
    ActiveConnections int64
}

func (lcb *LeastConnectionsBalancer) NextServer() *Server {
    lcb.mutex.RLock()
    defer lcb.mutex.RUnlock()

    if len(lcb.servers) == 0 {
        return nil
    }

    // Find server with least connections
    bestIndex := 0
    for i := 1; i < len(lcb.servers); i++ {
        if lcb.servers[i].ActiveConnections < lcb.servers[bestIndex].ActiveConnections {
            bestIndex = i
        }
    }

    return &lcb.servers[bestIndex].Server
}

func (lcb *LeastConnectionsBalancer) OnConnectionStart(serverID string) {
    lcb.mutex.Lock()
    defer lcb.mutex.Unlock()

    for i := range lcb.servers {
        if lcb.servers[i].ID == serverID {
            atomic.AddInt64(&lcb.servers[i].ActiveConnections, 1)
            break
        }
    }
}

func (lcb *LeastConnectionsBalancer) OnConnectionEnd(serverID string) {
    lcb.mutex.Lock()
    defer lcb.mutex.Unlock()

    for i := range lcb.servers {
        if lcb.servers[i].ID == serverID {
            atomic.AddInt64(&lcb.servers[i].ActiveConnections, -1)
            break
        }
    }
}

// Consistent hashing for sticky sessions
type ConsistentHashBalancer struct {
    ring     map[uint64]string
    servers  map[string]Server
    replicas int
    hashFunc hash.Hash64
}

func NewConsistentHashBalancer(servers []Server, replicas int) *ConsistentHashBalancer {
    chb := &ConsistentHashBalancer{
        ring:     make(map[uint64]string),
        servers:  make(map[string]Server),
        replicas: replicas,
        hashFunc: fnv.New64a(),
    }

    for _, server := range servers {
        chb.AddServer(server)
    }

    return chb
}

func (chb *ConsistentHashBalancer) AddServer(server Server) {
    chb.servers[server.ID] = server

    // Add virtual nodes
    for i := 0; i < chb.replicas; i++ {
        virtualKey := fmt.Sprintf("%s:%d", server.ID, i)
        hash := chb.hash(virtualKey)
        chb.ring[hash] = server.ID
    }
}

func (chb *ConsistentHashBalancer) GetServer(key string) *Server {
    if len(chb.ring) == 0 {
        return nil
    }

    hash := chb.hash(key)

    // Find first server clockwise from hash
    var keys []uint64
    for k := range chb.ring {
        keys = append(keys, k)
    }
    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

    for _, k := range keys {
        if k >= hash {
            serverID := chb.ring[k]
            server := chb.servers[serverID]
            return &server
        }
    }

    // Wrap around to first server
    serverID := chb.ring[keys[0]]
    server := chb.servers[serverID]
    return &server
}

func (chb *ConsistentHashBalancer) hash(key string) uint64 {
    chb.hashFunc.Reset()
    chb.hashFunc.Write([]byte(key))
    return chb.hashFunc.Sum64()
}
```

## 📊 Auto-scaling and Performance Monitoring

### Auto-scaling Implementation
```go
// Horizontal Pod Autoscaler (HPA) implementation
type HorizontalAutoscaler struct {
    targetCPU       float64
    targetMemory    float64
    minReplicas     int
    maxReplicas     int
    scaleUpCooldown time.Duration
    scaleDownCooldown time.Duration
    lastScaleTime   time.Time
    metrics         MetricsCollector
    scaler          PodScaler
}

func (ha *HorizontalAutoscaler) AutoScale(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            ha.evaluateScaling()
        }
    }
}

func (ha *HorizontalAutoscaler) evaluateScaling() {
    currentMetrics := ha.metrics.GetCurrentMetrics()
    currentReplicas := ha.scaler.GetCurrentReplicas()

    // Calculate desired replicas based on CPU utilization
    cpuDesiredReplicas := ha.calculateDesiredReplicas(
        currentReplicas,
        currentMetrics.CPUUtilization,
        ha.targetCPU,
    )

    // Calculate desired replicas based on memory utilization
    memoryDesiredReplicas := ha.calculateDesiredReplicas(
        currentReplicas,
        currentMetrics.MemoryUtilization,
        ha.targetMemory,
    )

    // Take the maximum to ensure both constraints are met
    desiredReplicas := max(cpuDesiredReplicas, memoryDesiredReplicas)

    // Apply min/max constraints
    desiredReplicas = max(ha.minReplicas, min(ha.maxReplicas, desiredReplicas))

    // Check cooldown period
    if ha.shouldScale(currentReplicas, desiredReplicas) {
        ha.scaler.ScaleTo(desiredReplicas)
        ha.lastScaleTime = time.Now()
    }
}

func (ha *HorizontalAutoscaler) calculateDesiredReplicas(current int, currentUtilization, targetUtilization float64) int {
    if targetUtilization == 0 {
        return current
    }

    ratio := currentUtilization / targetUtilization
    return int(math.Ceil(float64(current) * ratio))
}

func (ha *HorizontalAutoscaler) shouldScale(current, desired int) bool {
    if current == desired {
        return false
    }

    timeSinceLastScale := time.Since(ha.lastScaleTime)

    if desired > current {
        // Scale up
        return timeSinceLastScale >= ha.scaleUpCooldown
    } else {
        // Scale down
        return timeSinceLastScale >= ha.scaleDownCooldown
    }
}

// Vertical Pod Autoscaler (VPA) for resource optimization
type VerticalAutoscaler struct {
    targetUtilization float64
    safetyMargin      float64
    metrics           MetricsCollector
    resourceUpdater   ResourceUpdater
    recommendations   []ResourceRecommendation
}

func (va *VerticalAutoscaler) GenerateRecommendations() []ResourceRecommendation {
    historicalMetrics := va.metrics.GetHistoricalMetrics(7 * 24 * time.Hour) // 7 days

    var recommendations []ResourceRecommendation

    // Analyze CPU usage patterns
    cpuP95 := calculatePercentile(historicalMetrics.CPUUsage, 95)
    cpuRecommendation := cpuP95 * (1 + va.safetyMargin)

    recommendations = append(recommendations, ResourceRecommendation{
        ResourceType: "cpu",
        Current:      historicalMetrics.CurrentCPULimit,
        Recommended:  cpuRecommendation,
        Confidence:   va.calculateConfidence(historicalMetrics.CPUUsage),
    })

    // Analyze memory usage patterns
    memoryP95 := calculatePercentile(historicalMetrics.MemoryUsage, 95)
    memoryRecommendation := memoryP95 * (1 + va.safetyMargin)

    recommendations = append(recommendations, ResourceRecommendation{
        ResourceType: "memory",
        Current:      historicalMetrics.CurrentMemoryLimit,
        Recommended:  memoryRecommendation,
        Confidence:   va.calculateConfidence(historicalMetrics.MemoryUsage),
    })

    va.recommendations = recommendations
    return recommendations
}

// Predictive scaling based on historical patterns
type PredictiveScaler struct {
    predictor       TimeSeriesPredictor
    scaler          PodScaler
    lookAheadWindow time.Duration
    confidence      float64
}

func (ps *PredictiveScaler) PredictAndScale(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            ps.performPredictiveScaling()
        }
    }
}

func (ps *PredictiveScaler) performPredictiveScaling() {
    // Predict load for the next period
    prediction := ps.predictor.Predict(ps.lookAheadWindow)

    if prediction.Confidence < ps.confidence {
        return // Not confident enough in prediction
    }

    // Calculate required replicas for predicted load
    requiredReplicas := ps.calculateReplicasForLoad(prediction.Value)
    currentReplicas := ps.scaler.GetCurrentReplicas()

    // Only scale up proactively (scale down reactively)
    if requiredReplicas > currentReplicas {
        ps.scaler.ScaleTo(requiredReplicas)
    }
}
```

### Performance Monitoring and Alerting
```go
// Comprehensive performance monitoring system
type PerformanceMonitor struct {
    metrics     MetricsCollector
    alerting    AlertingSystem
    dashboard   DashboardUpdater
    thresholds  PerformanceThresholds
    collectors  []MetricCollector
}

type PerformanceThresholds struct {
    CPUWarning    float64
    CPUCritical   float64
    MemoryWarning float64
    MemoryCritical float64
    LatencyP95    time.Duration
    LatencyP99    time.Duration
    ErrorRate     float64
}

func (pm *PerformanceMonitor) StartMonitoring(ctx context.Context) {
    // Start metric collection
    for _, collector := range pm.collectors {
        go collector.Start(ctx)
    }

    // Start monitoring loop
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            pm.checkThresholds()
            pm.updateDashboard()
        }
    }
}

func (pm *PerformanceMonitor) checkThresholds() {
    currentMetrics := pm.metrics.GetCurrentMetrics()

    // Check CPU thresholds
    if currentMetrics.CPUUtilization > pm.thresholds.CPUCritical {
        pm.alerting.SendAlert(Alert{
            Level:   Critical,
            Message: fmt.Sprintf("CPU utilization critical: %.2f%%", currentMetrics.CPUUtilization),
            Metric:  "cpu_utilization",
            Value:   currentMetrics.CPUUtilization,
        })
    } else if currentMetrics.CPUUtilization > pm.thresholds.CPUWarning {
        pm.alerting.SendAlert(Alert{
            Level:   Warning,
            Message: fmt.Sprintf("CPU utilization high: %.2f%%", currentMetrics.CPUUtilization),
            Metric:  "cpu_utilization",
            Value:   currentMetrics.CPUUtilization,
        })
    }

    // Check memory thresholds
    if currentMetrics.MemoryUtilization > pm.thresholds.MemoryCritical {
        pm.alerting.SendAlert(Alert{
            Level:   Critical,
            Message: fmt.Sprintf("Memory utilization critical: %.2f%%", currentMetrics.MemoryUtilization),
            Metric:  "memory_utilization",
            Value:   currentMetrics.MemoryUtilization,
        })
    }

    // Check latency thresholds
    if currentMetrics.LatencyP99 > pm.thresholds.LatencyP99 {
        pm.alerting.SendAlert(Alert{
            Level:   Critical,
            Message: fmt.Sprintf("P99 latency too high: %v", currentMetrics.LatencyP99),
            Metric:  "latency_p99",
            Value:   float64(currentMetrics.LatencyP99.Milliseconds()),
        })
    }

    // Check error rate
    if currentMetrics.ErrorRate > pm.thresholds.ErrorRate {
        pm.alerting.SendAlert(Alert{
            Level:   Critical,
            Message: fmt.Sprintf("Error rate too high: %.2f%%", currentMetrics.ErrorRate*100),
            Metric:  "error_rate",
            Value:   currentMetrics.ErrorRate,
        })
    }
}

// Application Performance Monitoring (APM)
type APMTracer struct {
    tracer       opentracing.Tracer
    sampler      Sampler
    reporter     Reporter
    baggage      map[string]string
}

func (apm *APMTracer) TraceRequest(ctx context.Context, operationName string, fn func(context.Context) error) error {
    // Start span
    span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
    defer span.Finish()

    // Add tags
    span.SetTag("component", "api")
    span.SetTag("span.kind", "server")

    // Add baggage
    for key, value := range apm.baggage {
        span.SetBaggageItem(key, value)
    }

    start := time.Now()
    err := fn(ctx)
    duration := time.Since(start)

    // Record metrics
    span.SetTag("duration_ms", duration.Milliseconds())

    if err != nil {
        span.SetTag("error", true)
        span.LogFields(
            log.String("event", "error"),
            log.String("message", err.Error()),
        )
    }

    return err
}

// Real-time performance dashboard
type PerformanceDashboard struct {
    metrics     MetricsStore
    websockets  WebSocketManager
    charts      []ChartConfig
    alerts      AlertManager
}

func (pd *PerformanceDashboard) UpdateDashboard() {
    currentMetrics := pd.metrics.GetLatestMetrics()

    dashboardData := DashboardData{
        Timestamp: time.Now(),
        Metrics:   currentMetrics,
        Charts:    pd.generateChartData(),
        Alerts:    pd.alerts.GetActiveAlerts(),
    }

    // Broadcast to all connected clients
    pd.websockets.Broadcast(dashboardData)
}

func (pd *PerformanceDashboard) generateChartData() []ChartData {
    var charts []ChartData

    for _, config := range pd.charts {
        data := pd.metrics.GetTimeSeriesData(config.Metric, config.TimeRange)

        chartData := ChartData{
            Name:   config.Name,
            Type:   config.Type,
            Data:   data,
            Labels: pd.generateLabels(config.TimeRange),
        }

        charts = append(charts, chartData)
    }

    return charts
}
```

### Benchmarking and Load Testing
```go
// Load testing framework
type LoadTester struct {
    target      string
    concurrency int
    duration    time.Duration
    rampUp      time.Duration
    scenarios   []TestScenario
    results     *TestResults
}

type TestScenario struct {
    Name        string
    Weight      int
    RequestFunc func() (*http.Request, error)
    Validator   func(*http.Response) error
}

func (lt *LoadTester) RunLoadTest() (*TestResults, error) {
    results := &TestResults{
        StartTime:    time.Now(),
        Requests:     make([]RequestResult, 0),
        Errors:       make([]error, 0),
    }

    // Create worker pool
    workers := make(chan struct{}, lt.concurrency)
    var wg sync.WaitGroup

    // Ramp up gradually
    rampUpInterval := lt.rampUp / time.Duration(lt.concurrency)

    ctx, cancel := context.WithTimeout(context.Background(), lt.duration)
    defer cancel()

    // Start workers
    for i := 0; i < lt.concurrency; i++ {
        time.Sleep(rampUpInterval)

        wg.Add(1)
        go func() {
            defer wg.Done()
            lt.worker(ctx, workers, results)
        }()
    }

    wg.Wait()

    results.EndTime = time.Now()
    results.Duration = results.EndTime.Sub(results.StartTime)

    return results, nil
}

func (lt *LoadTester) worker(ctx context.Context, workers chan struct{}, results *TestResults) {
    client := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }

    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Select scenario based on weight
            scenario := lt.selectScenario()

            // Execute request
            start := time.Now()
            req, err := scenario.RequestFunc()
            if err != nil {
                results.AddError(err)
                continue
            }

            resp, err := client.Do(req)
            duration := time.Since(start)

            result := RequestResult{
                Scenario:   scenario.Name,
                Duration:   duration,
                StatusCode: 0,
                Error:      err,
            }

            if resp != nil {
                result.StatusCode = resp.StatusCode
                resp.Body.Close()

                // Validate response
                if scenario.Validator != nil {
                    if validationErr := scenario.Validator(resp); validationErr != nil {
                        result.Error = validationErr
                    }
                }
            }

            results.AddResult(result)
        }
    }
}

// Benchmark comparison framework
type BenchmarkSuite struct {
    benchmarks []Benchmark
    iterations int
    warmup     int
}

type Benchmark struct {
    Name string
    Func func() error
}

func (bs *BenchmarkSuite) RunBenchmarks() *BenchmarkResults {
    results := &BenchmarkResults{
        Results: make(map[string]BenchmarkResult),
    }

    for _, benchmark := range bs.benchmarks {
        result := bs.runSingleBenchmark(benchmark)
        results.Results[benchmark.Name] = result
    }

    return results
}

func (bs *BenchmarkSuite) runSingleBenchmark(benchmark Benchmark) BenchmarkResult {
    // Warmup
    for i := 0; i < bs.warmup; i++ {
        benchmark.Func()
    }

    // Actual benchmark
    durations := make([]time.Duration, bs.iterations)

    for i := 0; i < bs.iterations; i++ {
        start := time.Now()
        err := benchmark.Func()
        duration := time.Since(start)

        if err != nil {
            return BenchmarkResult{
                Name:  benchmark.Name,
                Error: err,
            }
        }

        durations[i] = duration
    }

    // Calculate statistics
    return BenchmarkResult{
        Name:       benchmark.Name,
        Iterations: bs.iterations,
        Mean:       calculateMean(durations),
        Median:     calculateMedian(durations),
        P95:        calculatePercentile(durations, 95),
        P99:        calculatePercentile(durations, 99),
        Min:        calculateMin(durations),
        Max:        calculateMax(durations),
    }
}
```

## 🎯 Summary: Performance & Optimization

### Comprehensive Optimization Coverage
- **CPU Optimization** - Algorithm efficiency, worker pools, SIMD, cache-friendly structures
- **Memory Optimization** - Object pools, memory mapping, GC optimization, monitoring
- **Network Optimization** - Connection pooling, batching, compression, HTTP/2
- **Database Optimization** - Query tuning, indexing, connection pooling, read replicas
- **Caching Strategies** - Multi-level hierarchies, warming, intelligent invalidation
- **Load Balancing** - Multiple algorithms for different scenarios
- **Auto-scaling** - Horizontal, vertical, and predictive scaling
- **Monitoring & Alerting** - Real-time performance tracking and alerting
- **Load Testing** - Comprehensive benchmarking and performance validation

### Advanced Performance Patterns
- **Object Pooling** - Reduce garbage collection overhead
- **Memory Mapping** - Efficient large file processing
- **Request Batching** - Minimize network round trips
- **Cache Hierarchies** - Optimize for different access patterns
- **Smart Invalidation** - Dependency-based cache management
- **Predictive Scaling** - Proactive resource allocation
- **APM Tracing** - Distributed request tracing
- **Load Testing** - Realistic performance validation

### Real-world Applications
These optimization techniques directly apply to all the systems we've designed:
- **Twitter**: Feed generation optimization, caching strategies
- **Netflix**: Video streaming optimization, CDN performance
- **Uber**: Real-time matching optimization, location indexing
- **WhatsApp**: Message delivery optimization, connection pooling
- **TinyURL**: High-throughput read optimization, caching
- **Search Engine**: Query optimization, indexing performance

You now have **expert-level performance optimization knowledge** covering all aspects of system performance! 🚀
