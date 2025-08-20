# 🔗 URL Shortener (TinyURL) System Design

## 🎯 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
```
✅ Core Features:
- Shorten long URLs to short URLs (e.g., https://example.com/very/long/url -> https://tinyurl.com/abc123)
- Redirect short URLs to original URLs
- Custom aliases for short URLs (optional)
- URL expiration (optional)
- Analytics: click tracking, geographic data, referrer info
- Bulk URL shortening API
- URL preview before redirect

❌ Out of Scope:
- User authentication (anonymous service)
- URL content filtering/moderation
- Real-time collaboration features
- Advanced security features (beyond basic validation)
```

### Non-Functional Requirements
```
Scale:
- 100M URLs shortened per month
- 100:1 read/write ratio (10B redirects per month)
- 500 new URLs per second
- 50K redirects per second
- 100M total URLs in system

Performance:
- URL shortening: < 200ms
- URL redirect: < 100ms
- 99.9% availability
- Global CDN support

Storage:
- Store URLs for 5+ years
- Handle URLs up to 2048 characters
- Short URLs should be 6-8 characters
```

## 📊 Step 2: Capacity Estimation

### Back-of-Envelope Calculations
```go
// Scale metrics
const (
    URLsPerMonth        = 100_000_000   // 100M URLs/month
    ReadWriteRatio      = 100           // 100:1 read/write ratio
    RedirectsPerMonth   = URLsPerMonth * ReadWriteRatio // 10B redirects/month
    TotalURLsInSystem   = 100_000_000   // 100M URLs stored
    
    SecondsPerMonth = 30 * 24 * 60 * 60 // ~2.6M seconds
    URLLifetimeYears = 5
)

// QPS calculations
var (
    URLCreationQPS = URLsPerMonth / SecondsPerMonth        // ~38 QPS
    RedirectQPS    = RedirectsPerMonth / SecondsPerMonth   // ~3,800 QPS
    PeakCreationQPS = URLCreationQPS * 10                  // ~380 QPS (peak)
    PeakRedirectQPS = RedirectQPS * 10                     // ~38,000 QPS (peak)
)

// Storage calculations
type StorageCalculation struct {
    URLData      int64 // Original and short URL storage
    AnalyticsData int64 // Click tracking and analytics
    CacheData    int64 // Hot URLs in cache
    IndexData    int64 // Database indexes
}

func CalculateStorage() StorageCalculation {
    calc := StorageCalculation{}
    
    // URL data: 100M URLs * (2KB original + 100 bytes metadata)
    avgOriginalURLSize := 2048 // bytes
    metadataSize := 100        // bytes (short_url, created_at, expires_at, etc.)
    calc.URLData = int64(TotalURLsInSystem) * int64(avgOriginalURLSize + metadataSize) // ~215 GB
    
    // Analytics data: 10B clicks * 200 bytes each
    clickRecordSize := 200 // bytes (timestamp, IP, user_agent, referrer, etc.)
    totalClicks := int64(RedirectsPerMonth) * 12 * int64(URLLifetimeYears) // 5 years of data
    calc.AnalyticsData = totalClicks * int64(clickRecordSize) // ~12 TB
    
    // Cache data: 20% of URLs cached (hot URLs)
    cacheHitRate := 0.2
    calc.CacheData = int64(float64(calc.URLData) * cacheHitRate) // ~43 GB
    
    // Index data: ~10% of main data
    calc.IndexData = calc.URLData / 10 // ~21.5 GB
    
    return calc
}

// Bandwidth calculations
func CalculateBandwidth() map[string]float64 {
    // Average request/response sizes
    shortenRequestSize := 2048 + 100   // URL + headers
    shortenResponseSize := 100         // Short URL + headers
    redirectRequestSize := 100         // Short URL + headers  
    redirectResponseSize := 200        // Redirect headers
    
    // Bandwidth calculations (bytes/second)
    inboundBandwidth := float64(PeakCreationQPS * shortenRequestSize + PeakRedirectQPS * redirectRequestSize)
    outboundBandwidth := float64(PeakCreationQPS * shortenResponseSize + PeakRedirectQPS * redirectResponseSize)
    
    return map[string]float64{
        "peak_inbound_mbps":  inboundBandwidth / (1024 * 1024),  // ~8 Mbps
        "peak_outbound_mbps": outboundBandwidth / (1024 * 1024), // ~8 Mbps
        "daily_bandwidth_gb": (inboundBandwidth + outboundBandwidth) * 24 * 3600 / (1024 * 1024 * 1024), // ~14 GB/day
    }
}
```

### Infrastructure Estimation
```go
type InfrastructureNeeds struct {
    WebServers      int // Handle API requests
    DatabaseServers int // URL and analytics storage
    CacheServers    int // Redis for hot URLs
    CDNNodes        int // Global redirect performance
}

func EstimateInfrastructure() InfrastructureNeeds {
    // Web servers (each handles 1K QPS)
    webServers := int(PeakRedirectQPS / 1000) + 5 // ~43 servers
    
    // Database servers (read replicas for high read load)
    databaseServers := 10 // 1 master + 9 read replicas
    
    // Cache servers (Redis cluster)
    storage := CalculateStorage()
    cacheMemoryNeeded := storage.CacheData
    cacheServers := int(cacheMemoryNeeded / (64 * 1024 * 1024 * 1024)) + 1 // 64GB per server, ~1 server
    
    // CDN nodes for global performance
    cdnNodes := 50 // Major cities worldwide
    
    return InfrastructureNeeds{
        WebServers:      webServers,
        DatabaseServers: databaseServers,
        CacheServers:    cacheServers,
        CDNNodes:        cdnNodes,
    }
}
```

## 🏗️ Step 3: High-Level Design

### System Architecture
```
[Web Apps] ──┐
             ├── [Load Balancer] ── [API Gateway]
[Mobile Apps]──┘                         │
                                         ├── [URL Shortening Service]
                                         ├── [URL Redirect Service]
                                         ├── [Analytics Service]
                                         └── [Custom Alias Service]
                                                  │
                              [Cache Layer] ── [Database Layer]
                                    │
                              [CDN Network] ── [Global Edge Servers]
```

### Core Services
```go
// URL Shortening Service - creates short URLs
type URLShorteningService struct {
    urlRepo         URLRepository
    encodingService EncodingService
    cache           Cache
    analyticsService AnalyticsService
}

// URL Redirect Service - handles redirects
type URLRedirectService struct {
    urlRepo         URLRepository
    cache           Cache
    analyticsService AnalyticsService
    cdnService      CDNService
}

// Analytics Service - tracks clicks and generates reports
type AnalyticsService struct {
    clickRepo       ClickRepository
    aggregationService AggregationService
    cache           Cache
}

// Encoding Service - generates short URL codes
type EncodingService struct {
    counterService  CounterService
    base62Encoder   Base62Encoder
    customAliasRepo CustomAliasRepository
}

// Custom Alias Service - handles custom short URLs
type CustomAliasService struct {
    aliasRepo       AliasRepository
    urlRepo         URLRepository
    cache           Cache
}
```

## 🗄️ Step 4: Database Design

### Data Models
```go
// URL models
type URL struct {
    ID          int64     `json:"id" db:"id"`
    ShortCode   string    `json:"short_code" db:"short_code"`
    OriginalURL string    `json:"original_url" db:"original_url"`
    CustomAlias string    `json:"custom_alias,omitempty" db:"custom_alias"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
    ClickCount  int64     `json:"click_count" db:"click_count"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatorIP   string    `json:"creator_ip" db:"creator_ip"`
}

// Analytics models
type ClickEvent struct {
    ID          int64     `json:"id" db:"id"`
    ShortCode   string    `json:"short_code" db:"short_code"`
    IPAddress   string    `json:"ip_address" db:"ip_address"`
    UserAgent   string    `json:"user_agent" db:"user_agent"`
    Referrer    string    `json:"referrer" db:"referrer"`
    Country     string    `json:"country" db:"country"`
    City        string    `json:"city" db:"city"`
    DeviceType  string    `json:"device_type" db:"device_type"`
    Browser     string    `json:"browser" db:"browser"`
    OS          string    `json:"os" db:"os"`
    ClickedAt   time.Time `json:"clicked_at" db:"clicked_at"`
}

// Aggregated analytics for faster queries
type DailyStats struct {
    ShortCode   string    `json:"short_code" db:"short_code"`
    Date        time.Time `json:"date" db:"date"`
    ClickCount  int64     `json:"click_count" db:"click_count"`
    UniqueIPs   int64     `json:"unique_ips" db:"unique_ips"`
    TopCountry  string    `json:"top_country" db:"top_country"`
    TopReferrer string    `json:"top_referrer" db:"top_referrer"`
}

// Custom alias tracking
type CustomAlias struct {
    Alias       string    `json:"alias" db:"alias"`
    ShortCode   string    `json:"short_code" db:"short_code"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    IsReserved  bool      `json:"is_reserved" db:"is_reserved"`
}
```

### Database Schema
```sql
-- URLs table (PostgreSQL)
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    custom_alias VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    click_count BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    creator_ip INET
);

CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_custom_alias ON urls(custom_alias) WHERE custom_alias IS NOT NULL;
CREATE INDEX idx_urls_created_at ON urls(created_at);
CREATE INDEX idx_urls_expires_at ON urls(expires_at) WHERE expires_at IS NOT NULL;

-- Click events table (Cassandra for high write volume)
CREATE TABLE click_events (
    id BIGINT PRIMARY KEY,
    short_code VARCHAR(10),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(2),
    city VARCHAR(100),
    device_type VARCHAR(20),
    browser VARCHAR(50),
    os VARCHAR(50),
    clicked_at TIMESTAMP
);

-- Partition by short_code and time for better performance
CREATE INDEX idx_click_events_short_code_time ON click_events(short_code, clicked_at DESC);
CREATE INDEX idx_click_events_clicked_at ON click_events(clicked_at DESC);

-- Daily aggregated stats (for analytics dashboard)
CREATE TABLE daily_stats (
    short_code VARCHAR(10),
    date DATE,
    click_count BIGINT,
    unique_ips BIGINT,
    top_country VARCHAR(2),
    top_referrer TEXT,
    PRIMARY KEY (short_code, date)
);

-- Custom aliases table
CREATE TABLE custom_aliases (
    alias VARCHAR(50) PRIMARY KEY,
    short_code VARCHAR(10) REFERENCES urls(short_code),
    created_at TIMESTAMP DEFAULT NOW(),
    is_reserved BOOLEAN DEFAULT FALSE
);

-- Counter table for generating unique IDs
CREATE TABLE counters (
    name VARCHAR(50) PRIMARY KEY,
    value BIGINT DEFAULT 0
);

INSERT INTO counters (name, value) VALUES ('url_counter', 1000000); -- Start from 1M
```

## 🔢 Step 5: URL Encoding Algorithms

### Base62 Encoding Implementation
```go
type Base62Encoder struct {
    alphabet string
    base     int
}

func NewBase62Encoder() *Base62Encoder {
    return &Base62Encoder{
        alphabet: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
        base:     62,
    }
}

func (b62 *Base62Encoder) Encode(num int64) string {
    if num == 0 {
        return string(b62.alphabet[0])
    }
    
    var result []byte
    for num > 0 {
        remainder := num % int64(b62.base)
        result = append([]byte{b62.alphabet[remainder]}, result...)
        num = num / int64(b62.base)
    }
    
    return string(result)
}

func (b62 *Base62Encoder) Decode(encoded string) int64 {
    var result int64 = 0
    base := int64(b62.base)
    
    for _, char := range encoded {
        result = result * base
        
        // Find character position in alphabet
        for i, alphabetChar := range b62.alphabet {
            if char == alphabetChar {
                result += int64(i)
                break
            }
        }
    }
    
    return result
}

// Counter-based approach for guaranteed uniqueness
type CounterService struct {
    db    Database
    cache Cache
    mutex sync.Mutex
}

func (cs *CounterService) GetNextID() (int64, error) {
    cs.mutex.Lock()
    defer cs.mutex.Unlock()
    
    // Try to get from cache first
    if cached, found := cs.cache.Get("url_counter"); found {
        nextID := cached.(int64) + 1
        cs.cache.Set("url_counter", nextID, 1*time.Hour)
        
        // Update database every 1000 IDs
        if nextID%1000 == 0 {
            go cs.updateDatabaseCounter(nextID)
        }
        
        return nextID, nil
    }
    
    // Get from database
    var currentValue int64
    err := cs.db.QueryRow("SELECT value FROM counters WHERE name = 'url_counter'").Scan(&currentValue)
    if err != nil {
        return 0, err
    }
    
    nextID := currentValue + 1
    
    // Update database
    _, err = cs.db.Exec("UPDATE counters SET value = $1 WHERE name = 'url_counter'", nextID)
    if err != nil {
        return 0, err
    }
    
    // Cache the value
    cs.cache.Set("url_counter", nextID, 1*time.Hour)
    
    return nextID, nil
}

// Alternative: Hash-based approach (for distributed systems)
type HashBasedEncoder struct {
    hasher hash.Hash
    base62 *Base62Encoder
}

func (hbe *HashBasedEncoder) GenerateShortCode(originalURL string) string {
    hbe.hasher.Reset()
    hbe.hasher.Write([]byte(originalURL))
    hashBytes := hbe.hasher.Sum(nil)
    
    // Convert first 8 bytes to int64
    hashInt := int64(0)
    for i := 0; i < 8 && i < len(hashBytes); i++ {
        hashInt = (hashInt << 8) | int64(hashBytes[i])
    }
    
    // Make positive
    if hashInt < 0 {
        hashInt = -hashInt
    }
    
    return hbe.base62.Encode(hashInt)[:6] // Take first 6 characters
}
```

### Collision Handling
```go
type EncodingService struct {
    counterService  *CounterService
    base62Encoder   *Base62Encoder
    hashEncoder     *HashBasedEncoder
    urlRepo         URLRepository
    maxRetries      int
}

func (es *EncodingService) GenerateShortCode(originalURL string, customAlias string) (string, error) {
    // If custom alias provided, validate and use it
    if customAlias != "" {
        return es.handleCustomAlias(customAlias, originalURL)
    }
    
    // Try hash-based approach first (faster)
    shortCode := es.hashEncoder.GenerateShortCode(originalURL)
    
    // Check for collision
    exists, err := es.urlRepo.ShortCodeExists(shortCode)
    if err != nil {
        return "", err
    }
    
    if !exists {
        return shortCode, nil
    }
    
    // Collision detected, try with counter-based approach
    for i := 0; i < es.maxRetries; i++ {
        id, err := es.counterService.GetNextID()
        if err != nil {
            return "", err
        }
        
        shortCode = es.base62Encoder.Encode(id)
        
        exists, err := es.urlRepo.ShortCodeExists(shortCode)
        if err != nil {
            return "", err
        }
        
        if !exists {
            return shortCode, nil
        }
    }
    
    return "", ErrMaxRetriesExceeded
}

func (es *EncodingService) handleCustomAlias(alias, originalURL string) (string, error) {
    // Validate alias format
    if !es.isValidAlias(alias) {
        return "", ErrInvalidAlias
    }
    
    // Check if alias is reserved
    if es.isReservedAlias(alias) {
        return "", ErrReservedAlias
    }
    
    // Check if alias already exists
    exists, err := es.urlRepo.CustomAliasExists(alias)
    if err != nil {
        return "", err
    }
    
    if exists {
        return "", ErrAliasAlreadyExists
    }
    
    return alias, nil
}

func (es *EncodingService) isValidAlias(alias string) bool {
    // Check length
    if len(alias) < 3 || len(alias) > 50 {
        return false
    }
    
    // Check characters (alphanumeric and hyphens only)
    matched, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", alias)
    return matched
}

func (es *EncodingService) isReservedAlias(alias string) bool {
    reservedWords := []string{
        "api", "www", "admin", "help", "about", "contact", "terms", "privacy",
        "login", "register", "dashboard", "settings", "profile", "home",
    }
    
    lowerAlias := strings.ToLower(alias)
    for _, reserved := range reservedWords {
        if lowerAlias == reserved {
            return true
        }
    }
    
    return false
}
```

## ⚡ Step 6: Caching Strategy

### Multi-Level Caching Implementation
```go
type CachingService struct {
    l1Cache *LRUCache      // Application-level cache (hot URLs)
    l2Cache *RedisCluster  // Distributed cache
    cdnCache *CDNService   // Edge cache for redirects
}

func (cs *CachingService) GetURL(shortCode string) (*URL, error) {
    // L1 Cache (Application memory)
    if url, found := cs.l1Cache.Get(shortCode); found {
        return url.(*URL), nil
    }

    // L2 Cache (Redis)
    if urlData, found := cs.l2Cache.Get(fmt.Sprintf("url:%s", shortCode)); found {
        var url URL
        json.Unmarshal([]byte(urlData), &url)

        // Store in L1 cache
        cs.l1Cache.Set(shortCode, &url, 10*time.Minute)

        return &url, nil
    }

    return nil, ErrURLNotFound
}

func (cs *CachingService) CacheURL(url *URL) {
    // Cache in L1 (hot URLs)
    cs.l1Cache.Set(url.ShortCode, url, 10*time.Minute)

    // Cache in L2 (Redis) with longer TTL
    urlData, _ := json.Marshal(url)
    cs.l2Cache.Set(fmt.Sprintf("url:%s", url.ShortCode), string(urlData), 1*time.Hour)

    // Cache popular URLs at CDN edge
    if url.ClickCount > 1000 { // Popular URL threshold
        cs.cdnCache.CacheRedirect(url.ShortCode, url.OriginalURL, 24*time.Hour)
    }
}

// Cache warming for popular URLs
func (cs *CachingService) WarmCache() {
    // Get top 10K most clicked URLs
    popularURLs, err := cs.getPopularURLs(10000)
    if err != nil {
        return
    }

    for _, url := range popularURLs {
        cs.CacheURL(url)
    }
}

// Cache invalidation
func (cs *CachingService) InvalidateURL(shortCode string) {
    cs.l1Cache.Delete(shortCode)
    cs.l2Cache.Delete(fmt.Sprintf("url:%s", shortCode))
    cs.cdnCache.InvalidateRedirect(shortCode)
}
```

### Cache Partitioning Strategy
```go
type PartitionedCache struct {
    partitions []*RedisClient
    hasher     hash.Hash32
}

func (pc *PartitionedCache) getPartition(key string) *RedisClient {
    pc.hasher.Reset()
    pc.hasher.Write([]byte(key))
    hash := pc.hasher.Sum32()

    partitionIndex := hash % uint32(len(pc.partitions))
    return pc.partitions[partitionIndex]
}

func (pc *PartitionedCache) Get(key string) (string, bool) {
    partition := pc.getPartition(key)
    return partition.Get(key)
}

func (pc *PartitionedCache) Set(key, value string, ttl time.Duration) error {
    partition := pc.getPartition(key)
    return partition.Set(key, value, ttl)
}
```

## 📊 Step 7: Analytics System

### Real-time Analytics Implementation
```go
type AnalyticsService struct {
    clickRepo       ClickRepository
    statsRepo       StatsRepository
    eventStream     EventStream
    aggregator      RealTimeAggregator
    geoService      GeoLocationService
    userAgentParser UserAgentParser
}

func (as *AnalyticsService) TrackClick(shortCode, ipAddress, userAgent, referrer string) error {
    // Parse user agent
    deviceInfo := as.userAgentParser.Parse(userAgent)

    // Get geo location
    geoInfo, _ := as.geoService.GetLocation(ipAddress)

    // Create click event
    clickEvent := &ClickEvent{
        ShortCode:  shortCode,
        IPAddress:  ipAddress,
        UserAgent:  userAgent,
        Referrer:   referrer,
        Country:    geoInfo.Country,
        City:       geoInfo.City,
        DeviceType: deviceInfo.DeviceType,
        Browser:    deviceInfo.Browser,
        OS:         deviceInfo.OS,
        ClickedAt:  time.Now(),
    }

    // Store click event (async)
    go as.clickRepo.Save(clickEvent)

    // Send to real-time stream
    as.eventStream.Publish("click_events", clickEvent)

    // Update real-time counters
    as.aggregator.IncrementClick(shortCode)

    return nil
}

// Real-time aggregation
type RealTimeAggregator struct {
    cache       Cache
    batchSize   int
    flushInterval time.Duration
    buffer      map[string]*ClickBuffer
    mutex       sync.RWMutex
}

type ClickBuffer struct {
    ShortCode   string
    ClickCount  int64
    UniqueIPs   map[string]bool
    Countries   map[string]int64
    Referrers   map[string]int64
    LastUpdated time.Time
}

func (rta *RealTimeAggregator) IncrementClick(shortCode string) {
    rta.mutex.Lock()
    defer rta.mutex.Unlock()

    buffer, exists := rta.buffer[shortCode]
    if !exists {
        buffer = &ClickBuffer{
            ShortCode:   shortCode,
            UniqueIPs:   make(map[string]bool),
            Countries:   make(map[string]int64),
            Referrers:   make(map[string]int64),
            LastUpdated: time.Now(),
        }
        rta.buffer[shortCode] = buffer
    }

    buffer.ClickCount++
    buffer.LastUpdated = time.Now()

    // Flush if buffer is full
    if buffer.ClickCount >= int64(rta.batchSize) {
        go rta.flushBuffer(shortCode)
    }
}

func (rta *RealTimeAggregator) flushBuffer(shortCode string) {
    rta.mutex.Lock()
    buffer, exists := rta.buffer[shortCode]
    if !exists {
        rta.mutex.Unlock()
        return
    }

    // Copy buffer data
    bufferCopy := *buffer

    // Reset buffer
    rta.buffer[shortCode] = &ClickBuffer{
        ShortCode: shortCode,
        UniqueIPs: make(map[string]bool),
        Countries: make(map[string]int64),
        Referrers: make(map[string]int64),
    }
    rta.mutex.Unlock()

    // Update database
    rta.updateDailyStats(shortCode, &bufferCopy)

    // Update cache
    rta.updateRealTimeStats(shortCode, &bufferCopy)
}

// Analytics API
func (as *AnalyticsService) GetURLAnalytics(shortCode string, timeRange TimeRange) (*URLAnalytics, error) {
    analytics := &URLAnalytics{
        ShortCode: shortCode,
        TimeRange: timeRange,
    }

    // Get basic stats
    stats, err := as.statsRepo.GetStats(shortCode, timeRange.Start, timeRange.End)
    if err != nil {
        return nil, err
    }

    analytics.TotalClicks = stats.TotalClicks
    analytics.UniqueClicks = stats.UniqueClicks

    // Get geographic distribution
    geoStats, _ := as.statsRepo.GetGeoStats(shortCode, timeRange.Start, timeRange.End)
    analytics.CountryStats = geoStats

    // Get referrer stats
    referrerStats, _ := as.statsRepo.GetReferrerStats(shortCode, timeRange.Start, timeRange.End)
    analytics.ReferrerStats = referrerStats

    // Get time series data
    timeSeriesData, _ := as.statsRepo.GetTimeSeriesData(shortCode, timeRange.Start, timeRange.End, timeRange.Granularity)
    analytics.TimeSeries = timeSeriesData

    return analytics, nil
}
```

### Analytics Data Models
```go
type URLAnalytics struct {
    ShortCode      string                    `json:"short_code"`
    TimeRange      TimeRange                 `json:"time_range"`
    TotalClicks    int64                     `json:"total_clicks"`
    UniqueClicks   int64                     `json:"unique_clicks"`
    CountryStats   []CountryStats            `json:"country_stats"`
    ReferrerStats  []ReferrerStats           `json:"referrer_stats"`
    DeviceStats    []DeviceStats             `json:"device_stats"`
    TimeSeries     []TimeSeriesPoint         `json:"time_series"`
}

type CountryStats struct {
    Country    string  `json:"country"`
    ClickCount int64   `json:"click_count"`
    Percentage float64 `json:"percentage"`
}

type ReferrerStats struct {
    Referrer   string  `json:"referrer"`
    ClickCount int64   `json:"click_count"`
    Percentage float64 `json:"percentage"`
}

type TimeSeriesPoint struct {
    Timestamp  time.Time `json:"timestamp"`
    ClickCount int64     `json:"click_count"`
}
```

## 🛡️ Step 8: Security & Rate Limiting

### Rate Limiting Implementation
```go
type RateLimiter struct {
    cache       Cache
    rules       []RateLimitRule
    ipWhitelist map[string]bool
}

type RateLimitRule struct {
    Name        string
    Limit       int
    Window      time.Duration
    KeyPattern  string // "ip:{ip}", "user:{user_id}", "global"
}

func (rl *RateLimiter) IsAllowed(request *RateLimitRequest) (bool, error) {
    for _, rule := range rl.rules {
        key := rl.buildKey(rule.KeyPattern, request)

        // Check whitelist
        if rl.ipWhitelist[request.IPAddress] {
            continue
        }

        // Get current count
        currentCount, _ := rl.cache.Get(key)
        count := 0
        if currentCount != nil {
            count = currentCount.(int)
        }

        if count >= rule.Limit {
            return false, fmt.Errorf("rate limit exceeded for rule: %s", rule.Name)
        }

        // Increment counter
        rl.cache.Increment(key, 1, rule.Window)
    }

    return true, nil
}

// URL shortening rate limits
var ShorteningRateLimits = []RateLimitRule{
    {
        Name:       "ip_per_minute",
        Limit:      10,
        Window:     1 * time.Minute,
        KeyPattern: "shorten:ip:{ip}",
    },
    {
        Name:       "ip_per_hour",
        Limit:      100,
        Window:     1 * time.Hour,
        KeyPattern: "shorten:ip:{ip}",
    },
    {
        Name:       "global_per_second",
        Limit:      1000,
        Window:     1 * time.Second,
        KeyPattern: "shorten:global",
    },
}

// Redirect rate limits (more lenient)
var RedirectRateLimits = []RateLimitRule{
    {
        Name:       "ip_per_minute",
        Limit:      1000,
        Window:     1 * time.Minute,
        KeyPattern: "redirect:ip:{ip}",
    },
}
```

### Security Measures
```go
type SecurityService struct {
    urlValidator    URLValidator
    malwareScanner  MalwareScanner
    spamDetector    SpamDetector
    blacklistRepo   BlacklistRepository
}

func (ss *SecurityService) ValidateURL(originalURL string) error {
    // Basic URL validation
    if !ss.urlValidator.IsValid(originalURL) {
        return ErrInvalidURL
    }

    // Check URL length
    if len(originalURL) > 2048 {
        return ErrURLTooLong
    }

    // Check against blacklist
    if ss.blacklistRepo.IsBlacklisted(originalURL) {
        return ErrBlacklistedURL
    }

    // Malware scanning (async for performance)
    go ss.malwareScanner.ScanURL(originalURL)

    // Spam detection
    if ss.spamDetector.IsSpam(originalURL) {
        return ErrSpamURL
    }

    return nil
}

type URLValidator struct {
    allowedSchemes []string
    blockedDomains []string
}

func (uv *URLValidator) IsValid(urlStr string) bool {
    parsedURL, err := url.Parse(urlStr)
    if err != nil {
        return false
    }

    // Check scheme
    schemeAllowed := false
    for _, scheme := range uv.allowedSchemes {
        if parsedURL.Scheme == scheme {
            schemeAllowed = true
            break
        }
    }

    if !schemeAllowed {
        return false
    }

    // Check blocked domains
    for _, domain := range uv.blockedDomains {
        if strings.Contains(parsedURL.Host, domain) {
            return false
        }
    }

    return true
}
```

## 🌐 Step 9: Global CDN Integration

### CDN Service Implementation
```go
type CDNService struct {
    edgeServers    map[string]EdgeServer // region -> server
    originServers  []OriginServer
    cachePolicy    CachePolicy
    healthChecker  HealthChecker
}

func (cdn *CDNService) HandleRedirect(shortCode, clientIP string) (*RedirectResponse, error) {
    // Determine best edge server
    edgeServer := cdn.getBestEdgeServer(clientIP)

    // Try edge cache first
    if originalURL, found := edgeServer.GetCachedRedirect(shortCode); found {
        return &RedirectResponse{
            OriginalURL: originalURL,
            CacheHit:    true,
            ServerID:    edgeServer.ID,
        }, nil
    }

    // Fallback to origin
    originalURL, err := cdn.getFromOrigin(shortCode)
    if err != nil {
        return nil, err
    }

    // Cache at edge for future requests
    go edgeServer.CacheRedirect(shortCode, originalURL, cdn.cachePolicy.TTL)

    return &RedirectResponse{
        OriginalURL: originalURL,
        CacheHit:    false,
        ServerID:    edgeServer.ID,
    }, nil
}

func (cdn *CDNService) getBestEdgeServer(clientIP string) EdgeServer {
    // Get client location
    clientLocation := cdn.getClientLocation(clientIP)

    // Find nearest edge server
    var bestServer EdgeServer
    minDistance := math.MaxFloat64

    for region, server := range cdn.edgeServers {
        if !cdn.healthChecker.IsHealthy(server.ID) {
            continue
        }

        distance := cdn.calculateDistance(clientLocation, server.Location)
        if distance < minDistance {
            minDistance = distance
            bestServer = server
        }
    }

    return bestServer
}

// Edge server implementation
type EdgeServer struct {
    ID       string
    Location GeoLocation
    Cache    EdgeCache
    Metrics  ServerMetrics
}

func (es *EdgeServer) GetCachedRedirect(shortCode string) (string, bool) {
    return es.Cache.Get(fmt.Sprintf("redirect:%s", shortCode))
}

func (es *EdgeServer) CacheRedirect(shortCode, originalURL string, ttl time.Duration) {
    es.Cache.Set(fmt.Sprintf("redirect:%s", shortCode), originalURL, ttl)
    es.Metrics.IncrementCacheWrites()
}
```

### Cache Invalidation Strategy
```go
type CacheInvalidationService struct {
    cdnService    CDNService
    edgeServers   []EdgeServer
    eventBus      EventBus
}

func (cis *CacheInvalidationService) InvalidateURL(shortCode string) error {
    // Create invalidation event
    event := CacheInvalidationEvent{
        ShortCode:   shortCode,
        Timestamp:   time.Now(),
        RequestID:   generateRequestID(),
    }

    // Publish to all edge servers
    for _, server := range cis.edgeServers {
        go cis.invalidateOnServer(server, event)
    }

    // Publish event for monitoring
    cis.eventBus.Publish("cache.invalidation", event)

    return nil
}

func (cis *CacheInvalidationService) invalidateOnServer(server EdgeServer, event CacheInvalidationEvent) {
    cacheKey := fmt.Sprintf("redirect:%s", event.ShortCode)
    server.Cache.Delete(cacheKey)

    // Log invalidation
    log.Printf("Invalidated %s on server %s", cacheKey, server.ID)
}
```

## 🔧 Step 10: API Implementation

### URL Shortener API
```go
type URLShortenerAPI struct {
    shorteningService *URLShorteningService
    redirectService   *URLRedirectService
    analyticsService  *AnalyticsService
    securityService   *SecurityService
    rateLimiter       *RateLimiter
}

func (api *URLShortenerAPI) SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    // Core endpoints
    mux.HandleFunc("POST /api/v1/shorten", api.shortenURL)
    mux.HandleFunc("GET /{shortCode}", api.redirectURL)
    mux.HandleFunc("GET /api/v1/analytics/{shortCode}", api.getAnalytics)

    // Utility endpoints
    mux.HandleFunc("GET /api/v1/preview/{shortCode}", api.previewURL)
    mux.HandleFunc("POST /api/v1/bulk-shorten", api.bulkShortenURLs)

    return mux
}

func (api *URLShortenerAPI) shortenURL(w http.ResponseWriter, r *http.Request) {
    // Rate limiting
    rateLimitReq := &RateLimitRequest{
        IPAddress: getClientIP(r),
        UserAgent: r.UserAgent(),
    }

    if allowed, err := api.rateLimiter.IsAllowed(rateLimitReq); !allowed {
        writeErrorResponse(w, http.StatusTooManyRequests, err.Error())
        return
    }

    // Parse request
    var req ShortenURLRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validate URL
    if err := api.securityService.ValidateURL(req.OriginalURL); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    // Create short URL
    result, err := api.shorteningService.ShortenURL(r.Context(), ShortenURLParams{
        OriginalURL: req.OriginalURL,
        CustomAlias: req.CustomAlias,
        ExpiresAt:   req.ExpiresAt,
        CreatorIP:   getClientIP(r),
    })

    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to shorten URL")
        return
    }

    response := ShortenURLResponse{
        ShortCode:   result.ShortCode,
        ShortURL:    fmt.Sprintf("https://tinyurl.com/%s", result.ShortCode),
        OriginalURL: req.OriginalURL,
        CreatedAt:   result.CreatedAt,
        ExpiresAt:   result.ExpiresAt,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (api *URLShortenerAPI) redirectURL(w http.ResponseWriter, r *http.Request) {
    shortCode := r.PathValue("shortCode")

    // Get original URL
    originalURL, err := api.redirectService.GetOriginalURL(r.Context(), shortCode)
    if err != nil {
        if errors.Is(err, ErrURLNotFound) {
            http.NotFound(w, r)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Track click (async)
    go api.analyticsService.TrackClick(
        shortCode,
        getClientIP(r),
        r.UserAgent(),
        r.Referer(),
    )

    // Redirect
    http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (api *URLShortenerAPI) getAnalytics(w http.ResponseWriter, r *http.Request) {
    shortCode := r.PathValue("shortCode")

    // Parse time range
    timeRange := parseTimeRange(r.URL.Query())

    // Get analytics
    analytics, err := api.analyticsService.GetURLAnalytics(shortCode, timeRange)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to get analytics")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(analytics)
}
```

## 🎯 Summary: URL Shortener System Design

### Architecture Highlights
- **Base62 Encoding** - Efficient short code generation with collision handling
- **Multi-level Caching** - L1/L2/CDN caching for optimal redirect performance
- **Real-time Analytics** - Click tracking with geographic and device insights
- **Global CDN** - Edge servers for sub-100ms redirect times worldwide
- **Rate Limiting** - Multi-tier protection against abuse
- **Security Validation** - URL validation, malware scanning, spam detection

### Scalability Features
- **100M URLs** - Counter-based and hash-based encoding strategies
- **50K redirects/sec** - Distributed caching and CDN optimization
- **Global deployment** - Edge servers in major cities worldwide
- **High availability** - 99.9% uptime with redundancy and failover

### Performance Optimizations
- **URL shortening**: < 200ms with optimized encoding algorithms
- **URL redirect**: < 100ms with multi-level caching
- **Analytics queries**: < 500ms with pre-aggregated data
- **Cache hit rate**: > 90% for popular URLs
- **CDN coverage**: < 50ms latency globally

This design handles TinyURL-scale traffic with optimal performance and reliability! 🚀
