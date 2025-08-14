# ⚖️ Load Balancing & CDN Architecture

## 🎯 Load Balancing Fundamentals

### Why Load Balancing?
```
Single Server Limitations:
- CPU: 100% utilization = bottleneck
- Memory: Limited RAM capacity
- Network: Bandwidth constraints
- Availability: Single point of failure

Load Balancer Benefits:
- Distribute traffic across multiple servers
- Improve availability and fault tolerance
- Enable horizontal scaling
- Optimize resource utilization
```

### Load Balancer Types

#### Layer 4 (Transport Layer) Load Balancing
```go
// TCP/UDP level load balancing
type L4LoadBalancer struct {
    servers []Server
    current int
    mutex   sync.Mutex
}

type Server struct {
    Address string
    Port    int
    Weight  int
    Health  HealthStatus
}

func (lb *L4LoadBalancer) RouteConnection(conn net.Conn) error {
    server := lb.selectServer()
    
    // Forward TCP connection to selected server
    serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Address, server.Port))
    if err != nil {
        return err
    }
    
    // Bidirectional proxy
    go io.Copy(serverConn, conn)
    go io.Copy(conn, serverConn)
    
    return nil
}

// Pros: Fast, protocol agnostic, low latency
// Cons: No application-level intelligence
```

#### Layer 7 (Application Layer) Load Balancing
```go
// HTTP/HTTPS level load balancing
type L7LoadBalancer struct {
    servers []Server
    rules   []RoutingRule
}

type RoutingRule struct {
    Path      string
    Method    string
    Headers   map[string]string
    ServerIDs []string
}

func (lb *L7LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Route based on application-level information
    server := lb.selectServerByRule(r)
    
    // Create proxy request
    proxyURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", server.Address, server.Port))
    proxy := httputil.NewSingleHostReverseProxy(proxyURL)
    
    // Add custom headers
    r.Header.Set("X-Forwarded-For", r.RemoteAddr)
    r.Header.Set("X-Forwarded-Proto", "http")
    
    proxy.ServeHTTP(w, r)
}

func (lb *L7LoadBalancer) selectServerByRule(r *http.Request) *Server {
    // Route based on path
    if strings.HasPrefix(r.URL.Path, "/api/") {
        return lb.getAPIServer()
    }
    
    // Route based on user type (from JWT)
    if userType := r.Header.Get("X-User-Type"); userType == "premium" {
        return lb.getPremiumServer()
    }
    
    // Route based on geographic location
    if region := r.Header.Get("X-Region"); region != "" {
        return lb.getRegionalServer(region)
    }
    
    // Default routing
    return lb.getDefaultServer()
}

// Pros: Application-aware routing, SSL termination, content-based routing
// Cons: Higher latency, more CPU intensive
```

## 🔄 Load Balancing Algorithms

### 1. Round Robin
```go
type RoundRobinBalancer struct {
    servers []Server
    current int
    mutex   sync.Mutex
}

func (rb *RoundRobinBalancer) NextServer() *Server {
    rb.mutex.Lock()
    defer rb.mutex.Unlock()
    
    server := &rb.servers[rb.current]
    rb.current = (rb.current + 1) % len(rb.servers)
    
    return server
}

// Time Complexity: O(1)
// Use Case: Servers have similar capacity
// Pros: Simple, fair distribution
// Cons: Doesn't consider server load
```

### 2. Weighted Round Robin
```go
type WeightedRoundRobinBalancer struct {
    servers []WeightedServer
    mutex   sync.Mutex
}

type WeightedServer struct {
    Server
    Weight        int
    CurrentWeight int
}

func (wrb *WeightedRoundRobinBalancer) NextServer() *Server {
    wrb.mutex.Lock()
    defer wrb.mutex.Unlock()
    
    totalWeight := 0
    var selected *WeightedServer
    
    for i := range wrb.servers {
        server := &wrb.servers[i]
        server.CurrentWeight += server.Weight
        totalWeight += server.Weight
        
        if selected == nil || server.CurrentWeight > selected.CurrentWeight {
            selected = server
        }
    }
    
    if selected != nil {
        selected.CurrentWeight -= totalWeight
    }
    
    return &selected.Server
}

// Example: Servers with weights [5, 1, 1]
// Distribution: Server1=5/7, Server2=1/7, Server3=1/7
```

### 3. Least Connections
```go
type LeastConnectionsBalancer struct {
    servers []ServerWithConnections
    mutex   sync.RWMutex
}

type ServerWithConnections struct {
    Server
    ActiveConnections int32
}

func (lcb *LeastConnectionsBalancer) NextServer() *Server {
    lcb.mutex.RLock()
    defer lcb.mutex.RUnlock()
    
    var selected *ServerWithConnections
    minConnections := int32(math.MaxInt32)
    
    for i := range lcb.servers {
        server := &lcb.servers[i]
        connections := atomic.LoadInt32(&server.ActiveConnections)
        
        if connections < minConnections {
            minConnections = connections
            selected = server
        }
    }
    
    if selected != nil {
        atomic.AddInt32(&selected.ActiveConnections, 1)
    }
    
    return &selected.Server
}

func (lcb *LeastConnectionsBalancer) ReleaseConnection(server *Server) {
    lcb.mutex.RLock()
    defer lcb.mutex.RUnlock()
    
    for i := range lcb.servers {
        if lcb.servers[i].Address == server.Address {
            atomic.AddInt32(&lcb.servers[i].ActiveConnections, -1)
            break
        }
    }
}

// Time Complexity: O(n)
// Optimization: Use min-heap for O(log n)
// Use Case: Long-lived connections, varying request processing times
```

### 4. Weighted Least Connections
```go
func (wlcb *WeightedLeastConnectionsBalancer) NextServer() *Server {
    var selected *WeightedServerWithConnections
    minRatio := float64(math.MaxFloat64)
    
    for i := range wlcb.servers {
        server := &wlcb.servers[i]
        connections := atomic.LoadInt32(&server.ActiveConnections)
        ratio := float64(connections) / float64(server.Weight)
        
        if ratio < minRatio {
            minRatio = ratio
            selected = server
        }
    }
    
    return &selected.Server
}

// Considers both connection count and server capacity
```

### 5. IP Hash
```go
type IPHashBalancer struct {
    servers []Server
}

func (ihb *IPHashBalancer) NextServer(clientIP string) *Server {
    hash := fnv.New32a()
    hash.Write([]byte(clientIP))
    index := int(hash.Sum32()) % len(ihb.servers)
    
    return &ihb.servers[index]
}

// Pros: Session affinity, consistent routing
// Cons: Uneven distribution if clients are not evenly distributed
```

### 6. Consistent Hashing (Advanced)
```go
// Implementation similar to your algorithms collection
type ConsistentHashBalancer struct {
    ring     map[uint32]string
    sortedKeys []uint32
    servers  map[string]*Server
    replicas int
}

func (chb *ConsistentHashBalancer) AddServer(server *Server) {
    chb.servers[server.Address] = server
    
    for i := 0; i < chb.replicas; i++ {
        key := chb.hashKey(fmt.Sprintf("%s:%d", server.Address, i))
        chb.ring[key] = server.Address
        chb.sortedKeys = append(chb.sortedKeys, key)
    }
    
    sort.Slice(chb.sortedKeys, func(i, j int) bool {
        return chb.sortedKeys[i] < chb.sortedKeys[j]
    })
}

func (chb *ConsistentHashBalancer) NextServer(key string) *Server {
    if len(chb.ring) == 0 {
        return nil
    }
    
    hash := chb.hashKey(key)
    idx := sort.Search(len(chb.sortedKeys), func(i int) bool {
        return chb.sortedKeys[i] >= hash
    })
    
    if idx == len(chb.sortedKeys) {
        idx = 0
    }
    
    serverAddr := chb.ring[chb.sortedKeys[idx]]
    return chb.servers[serverAddr]
}

// Pros: Minimal redistribution when servers added/removed
// Use Case: Distributed caching, database sharding
```

## 🏥 Health Checking

### Active Health Checks
```go
type HealthChecker struct {
    servers   []*Server
    interval  time.Duration
    timeout   time.Duration
    unhealthyThreshold int
    healthyThreshold   int
}

func (hc *HealthChecker) StartHealthChecks() {
    ticker := time.NewTicker(hc.interval)
    
    go func() {
        for range ticker.C {
            for _, server := range hc.servers {
                go hc.checkServerHealth(server)
            }
        }
    }()
}

func (hc *HealthChecker) checkServerHealth(server *Server) {
    client := &http.Client{Timeout: hc.timeout}
    
    resp, err := client.Get(fmt.Sprintf("http://%s:%d/health", server.Address, server.Port))
    if err != nil || resp.StatusCode != 200 {
        server.FailureCount++
        if server.FailureCount >= hc.unhealthyThreshold {
            server.Health = Unhealthy
        }
    } else {
        server.SuccessCount++
        if server.SuccessCount >= hc.healthyThreshold {
            server.Health = Healthy
            server.FailureCount = 0
        }
    }
    
    if resp != nil {
        resp.Body.Close()
    }
}

// Health check endpoints should be lightweight
func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Check database connectivity
    if err := db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    // Check external dependencies
    if err := checkExternalServices(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

### Passive Health Checks
```go
type PassiveHealthChecker struct {
    servers map[string]*ServerHealth
    mutex   sync.RWMutex
}

type ServerHealth struct {
    Server
    ErrorRate    float64
    ResponseTime time.Duration
    LastError    time.Time
}

func (phc *PassiveHealthChecker) RecordRequest(serverAddr string, duration time.Duration, err error) {
    phc.mutex.Lock()
    defer phc.mutex.Unlock()
    
    health, exists := phc.servers[serverAddr]
    if !exists {
        return
    }
    
    // Update response time (exponential moving average)
    alpha := 0.1
    health.ResponseTime = time.Duration(float64(health.ResponseTime)*(1-alpha) + float64(duration)*alpha)
    
    // Update error rate
    if err != nil {
        health.ErrorRate = health.ErrorRate*0.9 + 0.1 // Increase error rate
        health.LastError = time.Now()
        
        // Mark as unhealthy if error rate too high
        if health.ErrorRate > 0.1 { // 10% error rate threshold
            health.Health = Unhealthy
        }
    } else {
        health.ErrorRate = health.ErrorRate * 0.95 // Decrease error rate
        
        // Mark as healthy if error rate low and no recent errors
        if health.ErrorRate < 0.01 && time.Since(health.LastError) > 1*time.Minute {
            health.Health = Healthy
        }
    }
}
```

## 🌐 Content Delivery Network (CDN)

### CDN Architecture
```go
type CDN struct {
    originServer string
    edgeServers  map[string]*EdgeServer
    cache        Cache
}

type EdgeServer struct {
    Location    string
    Capacity    int64
    Cache       Cache
    Bandwidth   int64
}

func (cdn *CDN) ServeContent(w http.ResponseWriter, r *http.Request) {
    // 1. Determine closest edge server
    clientIP := r.Header.Get("X-Real-IP")
    edgeServer := cdn.getClosestEdgeServer(clientIP)
    
    // 2. Check edge cache
    cacheKey := cdn.generateCacheKey(r.URL.Path, r.Header)
    if content, found := edgeServer.Cache.Get(cacheKey); found {
        cdn.serveFromCache(w, content)
        return
    }
    
    // 3. Check regional cache
    if content, found := cdn.cache.Get(cacheKey); found {
        // Populate edge cache
        edgeServer.Cache.Set(cacheKey, content, 1*time.Hour)
        cdn.serveFromCache(w, content)
        return
    }
    
    // 4. Fetch from origin server
    content, err := cdn.fetchFromOrigin(r.URL.Path)
    if err != nil {
        http.Error(w, "Internal Server Error", 500)
        return
    }
    
    // 5. Cache at multiple levels
    cdn.cache.Set(cacheKey, content, 24*time.Hour)
    edgeServer.Cache.Set(cacheKey, content, 1*time.Hour)
    
    cdn.serveContent(w, content)
}

func (cdn *CDN) getClosestEdgeServer(clientIP string) *EdgeServer {
    // Simplified geolocation-based selection
    location := cdn.getLocationFromIP(clientIP)
    
    var closest *EdgeServer
    minDistance := math.MaxFloat64
    
    for _, server := range cdn.edgeServers {
        distance := cdn.calculateDistance(location, server.Location)
        if distance < minDistance {
            minDistance = distance
            closest = server
        }
    }
    
    return closest
}
```

### CDN Caching Strategies

#### Cache Control Headers
```go
func (cdn *CDN) setCacheHeaders(w http.ResponseWriter, contentType string) {
    switch {
    case strings.HasPrefix(contentType, "image/"):
        // Images: Cache for 1 year
        w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
        
    case strings.HasPrefix(contentType, "text/css"):
        // CSS: Cache for 1 month
        w.Header().Set("Cache-Control", "public, max-age=2592000")
        
    case strings.HasPrefix(contentType, "application/javascript"):
        // JavaScript: Cache for 1 month
        w.Header().Set("Cache-Control", "public, max-age=2592000")
        
    case strings.HasPrefix(contentType, "text/html"):
        // HTML: Cache for 5 minutes
        w.Header().Set("Cache-Control", "public, max-age=300")
        
    default:
        // Default: No cache
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    }
    
    // Add ETag for validation
    etag := cdn.generateETag(content)
    w.Header().Set("ETag", etag)
}
```

#### Cache Invalidation
```go
type CDNInvalidation struct {
    cdn     *CDN
    eventBus EventBus
}

func (ci *CDNInvalidation) Initialize() {
    ci.eventBus.Subscribe("content.updated", ci.handleContentUpdated)
    ci.eventBus.Subscribe("content.deleted", ci.handleContentDeleted)
}

func (ci *CDNInvalidation) handleContentUpdated(event ContentUpdatedEvent) {
    // Invalidate specific content
    pattern := fmt.Sprintf("/content/%s/*", event.ContentID)
    ci.invalidatePattern(pattern)
    
    // Invalidate related content
    if event.ContentType == "user-profile" {
        userPattern := fmt.Sprintf("/users/%s/*", event.UserID)
        ci.invalidatePattern(userPattern)
    }
}

func (ci *CDNInvalidation) invalidatePattern(pattern string) {
    // Invalidate on all edge servers
    for _, edgeServer := range ci.cdn.edgeServers {
        go func(server *EdgeServer) {
            server.Cache.InvalidatePattern(pattern)
        }(edgeServer)
    }
    
    // Invalidate regional cache
    ci.cdn.cache.InvalidatePattern(pattern)
}

// Batch invalidation for efficiency
func (ci *CDNInvalidation) batchInvalidate(patterns []string) {
    for _, edgeServer := range ci.cdn.edgeServers {
        go func(server *EdgeServer, p []string) {
            for _, pattern := range p {
                server.Cache.InvalidatePattern(pattern)
            }
        }(edgeServer, patterns)
    }
}
```

### CDN Performance Optimization

#### Compression
```go
func (cdn *CDN) compressContent(content []byte, contentType string) []byte {
    if !cdn.shouldCompress(contentType) {
        return content
    }
    
    var compressed bytes.Buffer
    writer := gzip.NewWriter(&compressed)
    writer.Write(content)
    writer.Close()
    
    // Only use compression if it reduces size significantly
    if compressed.Len() < len(content)*0.9 {
        return compressed.Bytes()
    }
    
    return content
}

func (cdn *CDN) shouldCompress(contentType string) bool {
    compressibleTypes := []string{
        "text/html",
        "text/css",
        "application/javascript",
        "application/json",
        "text/xml",
    }
    
    for _, t := range compressibleTypes {
        if strings.Contains(contentType, t) {
            return true
        }
    }
    
    return false
}
```

#### Image Optimization
```go
func (cdn *CDN) optimizeImage(imageData []byte, format string, quality int) []byte {
    // Resize based on device type
    if cdn.isMobileRequest() {
        return cdn.resizeImage(imageData, 800, 600)
    }
    
    // Convert to WebP for supported browsers
    if cdn.supportsWebP() {
        return cdn.convertToWebP(imageData, quality)
    }
    
    // Optimize JPEG quality
    if format == "jpeg" {
        return cdn.optimizeJPEG(imageData, quality)
    }
    
    return imageData
}
```

## 🔗 Next Steps
- Study message queues and event streaming systems
- Learn about API design and communication patterns
- Explore microservices architecture patterns
