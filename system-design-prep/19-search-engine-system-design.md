# 🔍 Search Engine System Design

## 🎯 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
```
✅ Core Features:
- Web crawling and indexing of billions of web pages
- Fast text search with relevance ranking
- Support for different query types (phrase, boolean, wildcard)
- Auto-complete and query suggestions
- Search result snippets and highlighting
- Image, video, and news search
- Personalized search results
- Safe search filtering

❌ Out of Scope:
- Voice search
- Real-time search (live tweets, etc.)
- Shopping/product search
- Advanced AI features (ChatGPT-style responses)
- Advertising platform
```

### Non-Functional Requirements
```
Scale:
- Index 50B web pages
- Handle 100K search queries per second
- Crawl 10M new pages per day
- Support 1B users globally
- Store 100PB of web content

Performance:
- Search latency: < 200ms
- Index freshness: < 24 hours for popular sites
- Crawl rate: 1000 pages/second per crawler
- 99.99% availability
- Global search with regional relevance

Quality:
- High relevance ranking
- Spam detection and filtering
- Duplicate content detection
- Multi-language support
```

## 📊 Step 2: Capacity Estimation

### Back-of-Envelope Calculations
```go
// Global scale metrics
const (
    TotalWebPages       = 50_000_000_000  // 50B pages indexed
    DailyQueries        = 8_640_000_000   // 100K QPS * 86400 seconds
    NewPagesPerDay      = 10_000_000      // 10M new pages daily
    CrawlersNeeded      = 1000            // To crawl 1000 pages/sec each
    GlobalUsers         = 1_000_000_000   // 1B users
    
    SecondsPerDay = 24 * 60 * 60
    AvgPageSize   = 50 * 1024 // 50KB average page size
    AvgQueryLength = 3        // 3 words per query
)

// QPS calculations
var (
    SearchQPS     = DailyQueries / SecondsPerDay        // ~100K QPS
    PeakSearchQPS = SearchQPS * 3                       // ~300K QPS
    CrawlQPS      = NewPagesPerDay / SecondsPerDay      // ~115 QPS
    IndexingQPS   = CrawlQPS * 2                        // ~230 QPS (including updates)
)

// Storage calculations
type StorageCalculation struct {
    RawWebContent   int64 // Original web pages
    ProcessedContent int64 // Cleaned and processed content
    InvertedIndex   int64 // Search index
    LinkGraph       int64 // Web graph for PageRank
    UserData        int64 // Search history and preferences
    CacheData       int64 // Query results cache
}

func CalculateStorage() StorageCalculation {
    calc := StorageCalculation{}
    
    // Raw web content: 50B pages * 50KB each
    calc.RawWebContent = TotalWebPages * AvgPageSize // ~2.5 PB
    
    // Processed content: ~30% of raw (after cleaning, deduplication)
    calc.ProcessedContent = int64(float64(calc.RawWebContent) * 0.3) // ~750 TB
    
    // Inverted index: ~20% of processed content
    calc.InvertedIndex = int64(float64(calc.ProcessedContent) * 0.2) // ~150 TB
    
    // Link graph: 50B pages * 100 bytes per page (URL + links metadata)
    calc.LinkGraph = TotalWebPages * 100 // ~5 TB
    
    // User data: 1B users * 10KB each (search history, preferences)
    calc.UserData = GlobalUsers * 10 * 1024 // ~10 TB
    
    // Cache data: Popular queries and results
    calc.CacheData = 50 * 1024 * 1024 * 1024 * 1024 // ~50 TB
    
    return calc
}

// Processing requirements
func CalculateProcessing() map[string]int64 {
    return map[string]int64{
        "crawl_operations_per_sec":    CrawlQPS,           // ~115/sec
        "index_operations_per_sec":    IndexingQPS,        // ~230/sec
        "search_operations_per_sec":   PeakSearchQPS,      // ~300K/sec
        "pagerank_calculations":       TotalWebPages / (24 * 3600), // Daily PageRank update
        "ml_ranking_predictions":      PeakSearchQPS * 10, // ~3M/sec (top 10 results)
    }
}
```

### Infrastructure Estimation
```go
type InfrastructureNeeds struct {
    CrawlerServers    int // Web crawling
    IndexingServers   int // Content processing and indexing
    SearchServers     int // Query processing
    RankingServers    int // ML-based ranking
    StorageServers    int // Distributed storage
    CacheServers      int // Query result caching
}

func EstimateInfrastructure() InfrastructureNeeds {
    storage := CalculateStorage()
    
    // Crawler servers (each runs 10 crawlers at 100 pages/sec each)
    crawlerServers := CrawlersNeeded / 10 // ~100 servers
    
    // Indexing servers (each processes 50 pages/sec)
    indexingServers := int(IndexingQPS / 50) + 10 // ~15 servers
    
    // Search servers (each handles 1K QPS)
    searchServers := int(PeakSearchQPS / 1000) // ~300 servers
    
    // Ranking servers (ML inference, each handles 10K predictions/sec)
    rankingServers := int(PeakSearchQPS * 10 / 10000) // ~300 servers
    
    // Storage servers (each stores 10TB)
    totalStorage := storage.RawWebContent + storage.ProcessedContent + storage.InvertedIndex
    storageServers := int(totalStorage / (10 * 1024 * 1024 * 1024 * 1024)) // ~340 servers
    
    // Cache servers (Redis cluster, each 64GB)
    cacheServers := int(storage.CacheData / (64 * 1024 * 1024 * 1024)) // ~800 servers
    
    return InfrastructureNeeds{
        CrawlerServers:  crawlerServers,
        IndexingServers: indexingServers,
        SearchServers:   searchServers,
        RankingServers:  rankingServers,
        StorageServers:  storageServers,
        CacheServers:    cacheServers,
    }
}
```

## 🏗️ Step 3: High-Level Design

### System Architecture
```
[Web Browsers] ──┐
                 ├── [Load Balancer] ── [API Gateway]
[Mobile Apps] ───┘                           │
                                             ├── [Search Service]
                                             ├── [Autocomplete Service]
                                             ├── [Ranking Service]
                                             └── [Analytics Service]
                                                      │
                              [Cache Layer] ── [Search Index] ── [Storage Layer]
                                    │
[Web Crawlers] ── [Content Processing] ── [Indexing Pipeline]
       │
[Link Graph] ── [PageRank Calculator] ── [Quality Signals]
```

### Core Services
```go
// Web Crawler Service - discovers and fetches web pages
type WebCrawlerService struct {
    crawlerPool     CrawlerPool
    urlQueue        URLQueue
    robotsParser    RobotsParser
    contentFilter   ContentFilter
    duplicateDetector DuplicateDetector
}

// Indexing Service - processes content and builds search index
type IndexingService struct {
    contentProcessor ContentProcessor
    textExtractor    TextExtractor
    languageDetector LanguageDetector
    indexBuilder     IndexBuilder
    invertedIndex    InvertedIndex
}

// Search Service - handles search queries and returns results
type SearchService struct {
    queryProcessor   QueryProcessor
    indexSearcher    IndexSearcher
    rankingService   RankingService
    resultFormatter  ResultFormatter
    cache           Cache
}

// Ranking Service - ML-based relevance ranking
type RankingService struct {
    pageRankCalculator PageRankCalculator
    mlRankingModel     MLRankingModel
    qualitySignals     QualitySignalsService
    personalizer       PersonalizationService
}

// Autocomplete Service - query suggestions
type AutocompleteService struct {
    queryTrie        QueryTrie
    popularQueries   PopularQueriesService
    personalizedSuggestions PersonalizedSuggestionsService
    cache           Cache
}
```

## 🗄️ Step 4: Database Design

### Data Models
```go
// Web page models
type WebPage struct {
    ID              int64     `json:"id" db:"id"`
    URL             string    `json:"url" db:"url"`
    Title           string    `json:"title" db:"title"`
    Content         string    `json:"content" db:"content"`
    CleanedContent  string    `json:"cleaned_content" db:"cleaned_content"`
    Language        string    `json:"language" db:"language"`
    ContentType     string    `json:"content_type" db:"content_type"`
    LastCrawled     time.Time `json:"last_crawled" db:"last_crawled"`
    LastModified    time.Time `json:"last_modified" db:"last_modified"`
    PageRank        float64   `json:"page_rank" db:"page_rank"`
    QualityScore    float64   `json:"quality_score" db:"quality_score"`
    InboundLinks    int       `json:"inbound_links" db:"inbound_links"`
    OutboundLinks   int       `json:"outbound_links" db:"outbound_links"`
    IsSpam          bool      `json:"is_spam" db:"is_spam"`
    IsActive        bool      `json:"is_active" db:"is_active"`
}

// Link graph models
type Link struct {
    ID         int64  `json:"id" db:"id"`
    FromPageID int64  `json:"from_page_id" db:"from_page_id"`
    ToPageID   int64  `json:"to_page_id" db:"to_page_id"`
    AnchorText string `json:"anchor_text" db:"anchor_text"`
    LinkType   string `json:"link_type" db:"link_type"` // internal, external
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Search index models
type IndexTerm struct {
    Term       string  `json:"term" db:"term"`
    PageID     int64   `json:"page_id" db:"page_id"`
    Frequency  int     `json:"frequency" db:"frequency"`
    Position   []int   `json:"position" db:"position"` // Word positions in document
    TfIdf      float64 `json:"tf_idf" db:"tf_idf"`
    FieldType  string  `json:"field_type" db:"field_type"` // title, content, meta
}

// Query models
type SearchQuery struct {
    ID            int64     `json:"id" db:"id"`
    QueryText     string    `json:"query_text" db:"query_text"`
    NormalizedQuery string  `json:"normalized_query" db:"normalized_query"`
    UserID        string    `json:"user_id,omitempty" db:"user_id"`
    IPAddress     string    `json:"ip_address" db:"ip_address"`
    UserAgent     string    `json:"user_agent" db:"user_agent"`
    Language      string    `json:"language" db:"language"`
    Country       string    `json:"country" db:"country"`
    ResultsCount  int       `json:"results_count" db:"results_count"`
    ClickedResults []int64  `json:"clicked_results" db:"clicked_results"`
    SearchTime    time.Duration `json:"search_time" db:"search_time"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// Autocomplete models
type QuerySuggestion struct {
    Prefix      string  `json:"prefix" db:"prefix"`
    Suggestion  string  `json:"suggestion" db:"suggestion"`
    Frequency   int64   `json:"frequency" db:"frequency"`
    Score       float64 `json:"score" db:"score"`
    Language    string  `json:"language" db:"language"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// User behavior models
type UserInteraction struct {
    ID          int64     `json:"id" db:"id"`
    UserID      string    `json:"user_id" db:"user_id"`
    QueryID     int64     `json:"query_id" db:"query_id"`
    PageID      int64     `json:"page_id" db:"page_id"`
    Action      string    `json:"action" db:"action"` // click, dwell, bounce
    Position    int       `json:"position" db:"position"` // Result position
    DwellTime   int       `json:"dwell_time" db:"dwell_time"` // Seconds on page
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
```

### Database Schema
```sql
-- Web pages table (PostgreSQL with partitioning)
CREATE TABLE web_pages (
    id BIGSERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    content TEXT,
    cleaned_content TEXT,
    language VARCHAR(10),
    content_type VARCHAR(50),
    last_crawled TIMESTAMP,
    last_modified TIMESTAMP,
    page_rank DOUBLE PRECISION DEFAULT 0.0,
    quality_score DOUBLE PRECISION DEFAULT 0.0,
    inbound_links INTEGER DEFAULT 0,
    outbound_links INTEGER DEFAULT 0,
    is_spam BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
) PARTITION BY HASH (id);

-- Create partitions for better performance
CREATE TABLE web_pages_0 PARTITION OF web_pages FOR VALUES WITH (MODULUS 10, REMAINDER 0);
CREATE TABLE web_pages_1 PARTITION OF web_pages FOR VALUES WITH (MODULUS 10, REMAINDER 1);
-- ... create 8 more partitions

CREATE INDEX idx_web_pages_url ON web_pages(url);
CREATE INDEX idx_web_pages_last_crawled ON web_pages(last_crawled);
CREATE INDEX idx_web_pages_page_rank ON web_pages(page_rank DESC);
CREATE INDEX idx_web_pages_language ON web_pages(language);

-- Links table (for PageRank calculation)
CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    from_page_id BIGINT REFERENCES web_pages(id),
    to_page_id BIGINT REFERENCES web_pages(id),
    anchor_text TEXT,
    link_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_links_from_page ON links(from_page_id);
CREATE INDEX idx_links_to_page ON links(to_page_id);

-- Inverted index (Elasticsearch/Solr for better performance)
-- Stored in specialized search engines, but represented here for completeness
CREATE TABLE index_terms (
    term VARCHAR(100),
    page_id BIGINT,
    frequency INTEGER,
    positions INTEGER[],
    tf_idf DOUBLE PRECISION,
    field_type VARCHAR(20),
    PRIMARY KEY (term, page_id)
);

-- Partition by term prefix for better distribution
CREATE INDEX idx_index_terms_page_id ON index_terms(page_id);
CREATE INDEX idx_index_terms_tf_idf ON index_terms(tf_idf DESC);

-- Search queries (Cassandra for high write volume)
CREATE TABLE search_queries (
    id BIGINT PRIMARY KEY,
    query_text TEXT,
    normalized_query TEXT,
    user_id TEXT,
    ip_address INET,
    user_agent TEXT,
    language VARCHAR(10),
    country VARCHAR(2),
    results_count INTEGER,
    clicked_results BIGINT[],
    search_time_ms INTEGER,
    created_at TIMESTAMP
);

CREATE INDEX idx_search_queries_normalized ON search_queries(normalized_query);
CREATE INDEX idx_search_queries_created_at ON search_queries(created_at DESC);

-- Query suggestions (Redis for fast autocomplete)
-- Stored as sorted sets in Redis:
-- ZADD suggestions:{prefix} {score} {suggestion}

-- User interactions (for ML training)
CREATE TABLE user_interactions (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT,
    query_id BIGINT,
    page_id BIGINT,
    action VARCHAR(20),
    position INTEGER,
    dwell_time INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_interactions_user_id ON user_interactions(user_id);
CREATE INDEX idx_user_interactions_query_id ON user_interactions(query_id);
CREATE INDEX idx_user_interactions_created_at ON user_interactions(created_at DESC);
```

## 🕷️ Step 5: Web Crawling System

### Distributed Web Crawler Implementation
```go
type WebCrawlerService struct {
    crawlerPool     *CrawlerPool
    urlQueue        URLQueue
    robotsCache     RobotsCache
    contentFilter   ContentFilter
    duplicateDetector DuplicateDetector
    rateLimiter     RateLimiter
}

type CrawlerPool struct {
    crawlers    []*WebCrawler
    coordinator *CrawlerCoordinator
    metrics     CrawlerMetrics
}

type WebCrawler struct {
    ID              string
    httpClient      *http.Client
    urlQueue        URLQueue
    contentProcessor ContentProcessor
    robotsParser    RobotsParser
    maxConcurrency  int
    crawlDelay      time.Duration
}

func (wc *WebCrawler) Start(ctx context.Context) {
    semaphore := make(chan struct{}, wc.maxConcurrency)

    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Get next URL to crawl
            urlToCrawl, err := wc.urlQueue.Dequeue()
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }

            // Acquire semaphore
            semaphore <- struct{}{}

            go func(url string) {
                defer func() { <-semaphore }()
                wc.crawlURL(ctx, url)
            }(urlToCrawl.URL)
        }
    }
}

func (wc *WebCrawler) crawlURL(ctx context.Context, url string) error {
    // Check robots.txt
    if !wc.robotsParser.IsAllowed(url, "Googlebot") {
        return ErrRobotsDisallowed
    }

    // Rate limiting per domain
    domain := extractDomain(url)
    if !wc.rateLimiter.Allow(domain) {
        // Re-queue for later
        wc.urlQueue.Enqueue(URLToCrawl{URL: url, Priority: 1}, 5*time.Minute)
        return ErrRateLimited
    }

    // Fetch page
    resp, err := wc.httpClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Check content type
    contentType := resp.Header.Get("Content-Type")
    if !wc.contentFilter.IsAllowed(contentType) {
        return ErrContentTypeNotAllowed
    }

    // Read content with size limit
    content, err := wc.readWithLimit(resp.Body, 10*1024*1024) // 10MB limit
    if err != nil {
        return err
    }

    // Check for duplicates
    contentHash := calculateHash(content)
    if wc.duplicateDetector.IsDuplicate(contentHash) {
        return ErrDuplicateContent
    }

    // Process content
    page := &WebPage{
        URL:          url,
        Content:      string(content),
        ContentType:  contentType,
        LastCrawled:  time.Now(),
        LastModified: parseLastModified(resp.Header.Get("Last-Modified")),
    }

    // Extract links for future crawling
    links := wc.extractLinks(content, url)
    for _, link := range links {
        wc.urlQueue.Enqueue(URLToCrawl{
            URL:      link.URL,
            Priority: calculatePriority(link),
            Depth:    link.Depth + 1,
        }, 0)
    }

    // Send for indexing
    return wc.contentProcessor.Process(ctx, page)
}

// URL Queue with priority
type URLQueue interface {
    Enqueue(url URLToCrawl, delay time.Duration) error
    Dequeue() (*URLToCrawl, error)
    Size() int64
}

type URLToCrawl struct {
    URL         string    `json:"url"`
    Priority    int       `json:"priority"` // 1-10, higher is more important
    Depth       int       `json:"depth"`
    Referrer    string    `json:"referrer"`
    ScheduledAt time.Time `json:"scheduled_at"`
    Attempts    int       `json:"attempts"`
}

// Redis-based URL Queue implementation
type RedisURLQueue struct {
    client      *redis.Client
    queueKey    string
    delayedKey  string
}

func (ruq *RedisURLQueue) Enqueue(url URLToCrawl, delay time.Duration) error {
    urlData, err := json.Marshal(url)
    if err != nil {
        return err
    }

    if delay > 0 {
        // Add to delayed queue
        score := float64(time.Now().Add(delay).Unix())
        return ruq.client.ZAdd(ruq.delayedKey, &redis.Z{
            Score:  score,
            Member: string(urlData),
        }).Err()
    }

    // Add to priority queue
    return ruq.client.ZAdd(ruq.queueKey, &redis.Z{
        Score:  float64(url.Priority),
        Member: string(urlData),
    }).Err()
}

func (ruq *RedisURLQueue) Dequeue() (*URLToCrawl, error) {
    // First, move ready URLs from delayed queue
    ruq.moveReadyURLs()

    // Get highest priority URL
    result, err := ruq.client.ZPopMax(ruq.queueKey).Result()
    if err != nil {
        return nil, err
    }

    if len(result) == 0 {
        return nil, ErrQueueEmpty
    }

    var url URLToCrawl
    err = json.Unmarshal([]byte(result[0].Member.(string)), &url)
    return &url, err
}

func (ruq *RedisURLQueue) moveReadyURLs() {
    now := float64(time.Now().Unix())

    // Get URLs ready to be crawled
    readyURLs, _ := ruq.client.ZRangeByScore(ruq.delayedKey, &redis.ZRangeBy{
        Min: "0",
        Max: fmt.Sprintf("%f", now),
    }).Result()

    if len(readyURLs) == 0 {
        return
    }

    // Move to main queue
    pipe := ruq.client.Pipeline()
    for _, urlData := range readyURLs {
        var url URLToCrawl
        if json.Unmarshal([]byte(urlData), &url) == nil {
            pipe.ZAdd(ruq.queueKey, &redis.Z{
                Score:  float64(url.Priority),
                Member: urlData,
            })
        }
    }

    // Remove from delayed queue
    pipe.ZRemRangeByScore(ruq.delayedKey, "0", fmt.Sprintf("%f", now))
    pipe.Exec()
}
```

### Robots.txt Parser and Politeness
```go
type RobotsParser struct {
    cache       Cache
    httpClient  *http.Client
    defaultDelay time.Duration
}

type RobotsRules struct {
    UserAgent    string
    Disallowed   []string
    Allowed      []string
    CrawlDelay   time.Duration
    Sitemap      []string
    CachedAt     time.Time
}

func (rp *RobotsParser) IsAllowed(url, userAgent string) bool {
    domain := extractDomain(url)
    rules := rp.getRobotsRules(domain)

    path := extractPath(url)

    // Check disallowed patterns
    for _, pattern := range rules.Disallowed {
        if rp.matchesPattern(path, pattern) {
            return false
        }
    }

    // Check explicitly allowed patterns
    for _, pattern := range rules.Allowed {
        if rp.matchesPattern(path, pattern) {
            return true
        }
    }

    return true
}

func (rp *RobotsParser) GetCrawlDelay(domain string) time.Duration {
    rules := rp.getRobotsRules(domain)
    if rules.CrawlDelay > 0 {
        return rules.CrawlDelay
    }
    return rp.defaultDelay
}

func (rp *RobotsParser) getRobotsRules(domain string) *RobotsRules {
    cacheKey := fmt.Sprintf("robots:%s", domain)

    // Check cache first
    if cached, found := rp.cache.Get(cacheKey); found {
        return cached.(*RobotsRules)
    }

    // Fetch robots.txt
    robotsURL := fmt.Sprintf("https://%s/robots.txt", domain)
    resp, err := rp.httpClient.Get(robotsURL)
    if err != nil {
        // Default permissive rules
        return &RobotsRules{
            UserAgent:  "*",
            CrawlDelay: rp.defaultDelay,
        }
    }
    defer resp.Body.Close()

    content, err := io.ReadAll(resp.Body)
    if err != nil {
        return &RobotsRules{CrawlDelay: rp.defaultDelay}
    }

    rules := rp.parseRobotsTxt(string(content))

    // Cache for 24 hours
    rp.cache.Set(cacheKey, rules, 24*time.Hour)

    return rules
}

func (rp *RobotsParser) parseRobotsTxt(content string) *RobotsRules {
    rules := &RobotsRules{
        CrawlDelay: rp.defaultDelay,
        CachedAt:   time.Now(),
    }

    lines := strings.Split(content, "\n")
    currentUserAgent := ""

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, ":", 2)
        if len(parts) != 2 {
            continue
        }

        directive := strings.ToLower(strings.TrimSpace(parts[0]))
        value := strings.TrimSpace(parts[1])

        switch directive {
        case "user-agent":
            currentUserAgent = value
        case "disallow":
            if currentUserAgent == "*" || currentUserAgent == "Googlebot" {
                rules.Disallowed = append(rules.Disallowed, value)
            }
        case "allow":
            if currentUserAgent == "*" || currentUserAgent == "Googlebot" {
                rules.Allowed = append(rules.Allowed, value)
            }
        case "crawl-delay":
            if delay, err := strconv.Atoi(value); err == nil {
                rules.CrawlDelay = time.Duration(delay) * time.Second
            }
        case "sitemap":
            rules.Sitemap = append(rules.Sitemap, value)
        }
    }

    return rules
}
```

## 📚 Step 6: Indexing System

### Content Processing Pipeline
```go
type ContentProcessor struct {
    textExtractor    TextExtractor
    languageDetector LanguageDetector
    contentCleaner   ContentCleaner
    indexBuilder     IndexBuilder
    linkExtractor    LinkExtractor
}

func (cp *ContentProcessor) Process(ctx context.Context, page *WebPage) error {
    // Extract text content
    extractedText, err := cp.textExtractor.Extract(page.Content, page.ContentType)
    if err != nil {
        return err
    }

    // Clean and normalize content
    cleanedContent := cp.contentCleaner.Clean(extractedText)
    page.CleanedContent = cleanedContent

    // Detect language
    language := cp.languageDetector.Detect(cleanedContent)
    page.Language = language

    // Extract title and meta information
    page.Title = cp.extractTitle(page.Content)

    // Extract links
    links := cp.linkExtractor.Extract(page.Content, page.URL)

    // Build search index
    if err := cp.indexBuilder.IndexPage(ctx, page); err != nil {
        return err
    }

    // Store links for PageRank calculation
    if err := cp.storeLinks(ctx, page.ID, links); err != nil {
        return err
    }

    return nil
}

// Inverted Index Builder
type IndexBuilder struct {
    invertedIndex InvertedIndex
    tokenizer     Tokenizer
    stemmer       Stemmer
    stopWords     StopWordFilter
}

func (ib *IndexBuilder) IndexPage(ctx context.Context, page *WebPage) error {
    // Tokenize content
    tokens := ib.tokenizer.Tokenize(page.CleanedContent)

    // Remove stop words
    filteredTokens := ib.stopWords.Filter(tokens)

    // Stem words
    stemmedTokens := ib.stemmer.StemTokens(filteredTokens)

    // Calculate term frequencies
    termFreqs := ib.calculateTermFrequencies(stemmedTokens)

    // Index title with higher weight
    titleTokens := ib.tokenizer.Tokenize(page.Title)
    titleTermFreqs := ib.calculateTermFrequencies(titleTokens)

    // Build index entries
    var indexEntries []IndexEntry

    // Content terms
    for term, freq := range termFreqs {
        entry := IndexEntry{
            Term:      term,
            PageID:    page.ID,
            Frequency: freq,
            FieldType: "content",
            Positions: ib.findPositions(stemmedTokens, term),
        }
        indexEntries = append(indexEntries, entry)
    }

    // Title terms (with boost)
    for term, freq := range titleTermFreqs {
        entry := IndexEntry{
            Term:      term,
            PageID:    page.ID,
            Frequency: freq * 3, // Title boost
            FieldType: "title",
            Positions: ib.findPositions(titleTokens, term),
        }
        indexEntries = append(indexEntries, entry)
    }

    // Store in inverted index
    return ib.invertedIndex.AddEntries(ctx, indexEntries)
}

// Distributed Inverted Index
type InvertedIndex interface {
    AddEntries(ctx context.Context, entries []IndexEntry) error
    Search(ctx context.Context, terms []string) ([]SearchResult, error)
    GetTermFrequency(ctx context.Context, term string) (int64, error)
}

type ElasticsearchIndex struct {
    client *elasticsearch.Client
    index  string
}

func (ei *ElasticsearchIndex) AddEntries(ctx context.Context, entries []IndexEntry) error {
    var buf bytes.Buffer

    for _, entry := range entries {
        // Create document
        doc := map[string]interface{}{
            "term":       entry.Term,
            "page_id":    entry.PageID,
            "frequency":  entry.Frequency,
            "field_type": entry.FieldType,
            "positions":  entry.Positions,
            "tf_idf":     entry.TfIdf,
        }

        // Add to bulk request
        meta := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": ei.index,
                "_id":    fmt.Sprintf("%s_%d", entry.Term, entry.PageID),
            },
        }

        metaJSON, _ := json.Marshal(meta)
        docJSON, _ := json.Marshal(doc)

        buf.Write(metaJSON)
        buf.WriteByte('\n')
        buf.Write(docJSON)
        buf.WriteByte('\n')
    }

    // Execute bulk request
    res, err := ei.client.Bulk(
        bytes.NewReader(buf.Bytes()),
        ei.client.Bulk.WithIndex(ei.index),
        ei.client.Bulk.WithContext(ctx),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    return nil
}
```

## 📊 Step 7: PageRank Algorithm

### PageRank Calculation Implementation
```go
type PageRankCalculator struct {
    linkGraph     LinkGraph
    dampingFactor float64
    maxIterations int
    convergenceThreshold float64
    distributedCompute   DistributedCompute
}

func (prc *PageRankCalculator) CalculatePageRank(ctx context.Context) error {
    // Get all pages
    pages, err := prc.linkGraph.GetAllPages()
    if err != nil {
        return err
    }

    // Initialize PageRank values
    pageRanks := make(map[int64]float64)
    numPages := float64(len(pages))
    initialValue := 1.0 / numPages

    for _, pageID := range pages {
        pageRanks[pageID] = initialValue
    }

    // Iterative calculation
    for iteration := 0; iteration < prc.maxIterations; iteration++ {
        newPageRanks := make(map[int64]float64)

        // Distribute computation across workers
        results := prc.distributedCompute.ProcessInParallel(pages, func(pageID int64) float64 {
            return prc.calculatePageRankForPage(pageID, pageRanks)
        })

        // Collect results
        maxDiff := 0.0
        for pageID, newRank := range results {
            oldRank := pageRanks[pageID]
            newPageRanks[pageID] = newRank

            diff := math.Abs(newRank - oldRank)
            if diff > maxDiff {
                maxDiff = diff
            }
        }

        pageRanks = newPageRanks

        // Check convergence
        if maxDiff < prc.convergenceThreshold {
            log.Printf("PageRank converged after %d iterations", iteration+1)
            break
        }
    }

    // Store results
    return prc.storePageRankResults(ctx, pageRanks)
}

func (prc *PageRankCalculator) calculatePageRankForPage(pageID int64, currentRanks map[int64]float64) float64 {
    // Get incoming links
    incomingLinks, err := prc.linkGraph.GetIncomingLinks(pageID)
    if err != nil {
        return 0.0
    }

    rankSum := 0.0
    for _, linkingPageID := range incomingLinks {
        linkingPageRank := currentRanks[linkingPageID]
        outgoingLinksCount := prc.linkGraph.GetOutgoingLinksCount(linkingPageID)

        if outgoingLinksCount > 0 {
            rankSum += linkingPageRank / float64(outgoingLinksCount)
        }
    }

    // PageRank formula: PR(A) = (1-d)/N + d * sum(PR(T)/C(T))
    numPages := float64(len(currentRanks))
    return (1.0-prc.dampingFactor)/numPages + prc.dampingFactor*rankSum
}

// Distributed PageRank for large graphs
type DistributedPageRank struct {
    workers       []PageRankWorker
    coordinator   PageRankCoordinator
    graphPartitioner GraphPartitioner
}

func (dpr *DistributedPageRank) CalculateDistributed(ctx context.Context) error {
    // Partition graph across workers
    partitions := dpr.graphPartitioner.PartitionGraph()

    // Initialize workers
    for i, worker := range dpr.workers {
        worker.Initialize(partitions[i])
    }

    // Iterative calculation with message passing
    for iteration := 0; iteration < maxIterations; iteration++ {
        // Each worker calculates local PageRank
        var wg sync.WaitGroup
        for _, worker := range dpr.workers {
            wg.Add(1)
            go func(w PageRankWorker) {
                defer wg.Done()
                w.CalculateLocalPageRank()
            }(worker)
        }
        wg.Wait()

        // Exchange boundary information
        dpr.coordinator.ExchangeBoundaryData()

        // Check global convergence
        if dpr.coordinator.HasConverged() {
            break
        }
    }

    // Collect and merge results
    return dpr.coordinator.CollectResults()
}
```

### Link Graph Storage
```go
type LinkGraph interface {
    GetAllPages() ([]int64, error)
    GetIncomingLinks(pageID int64) ([]int64, error)
    GetOutgoingLinksCount(pageID int64) int
    AddLink(fromPageID, toPageID int64, anchorText string) error
}

type DistributedLinkGraph struct {
    shards []LinkGraphShard
    hasher ConsistentHasher
}

func (dlg *DistributedLinkGraph) GetIncomingLinks(pageID int64) ([]int64, error) {
    // Query all shards since incoming links can come from any page
    var allIncomingLinks []int64

    for _, shard := range dlg.shards {
        links, err := shard.GetIncomingLinks(pageID)
        if err != nil {
            continue // Handle gracefully
        }
        allIncomingLinks = append(allIncomingLinks, links...)
    }

    return allIncomingLinks, nil
}

func (dlg *DistributedLinkGraph) AddLink(fromPageID, toPageID int64, anchorText string) error {
    // Store link in shard based on source page
    shard := dlg.getShard(fromPageID)
    return shard.AddLink(fromPageID, toPageID, anchorText)
}

func (dlg *DistributedLinkGraph) getShard(pageID int64) LinkGraphShard {
    shardIndex := dlg.hasher.Hash(pageID) % len(dlg.shards)
    return dlg.shards[shardIndex]
}
```

## 🔍 Step 8: Search Query Processing

### Query Processing Pipeline
```go
type SearchService struct {
    queryProcessor   QueryProcessor
    indexSearcher    IndexSearcher
    rankingService   RankingService
    resultFormatter  ResultFormatter
    cache           Cache
    analytics       SearchAnalytics
}

func (ss *SearchService) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
    startTime := time.Now()

    // Check cache first
    cacheKey := ss.generateCacheKey(request)
    if cached, found := ss.cache.Get(cacheKey); found {
        response := cached.(*SearchResponse)
        response.SearchTime = time.Since(startTime)
        response.CacheHit = true
        return response, nil
    }

    // Process query
    processedQuery, err := ss.queryProcessor.Process(request.Query, request.Language)
    if err != nil {
        return nil, err
    }

    // Search index
    searchResults, err := ss.indexSearcher.Search(ctx, SearchQuery{
        Terms:     processedQuery.Terms,
        Phrases:   processedQuery.Phrases,
        Filters:   request.Filters,
        Language:  request.Language,
        Country:   request.Country,
        Limit:     request.Limit,
        Offset:    request.Offset,
    })
    if err != nil {
        return nil, err
    }

    // Rank results
    rankedResults, err := ss.rankingService.RankResults(ctx, RankingRequest{
        Results:   searchResults,
        Query:     processedQuery,
        UserID:    request.UserID,
        Context:   request.Context,
    })
    if err != nil {
        return nil, err
    }

    // Format results
    formattedResults := ss.resultFormatter.Format(rankedResults, processedQuery)

    response := &SearchResponse{
        Results:     formattedResults,
        TotalCount:  len(searchResults),
        SearchTime:  time.Since(startTime),
        Query:       request.Query,
        Suggestions: ss.generateSuggestions(processedQuery),
        CacheHit:    false,
    }

    // Cache results
    ss.cache.Set(cacheKey, response, 10*time.Minute)

    // Track analytics
    go ss.analytics.TrackSearch(ctx, request, response)

    return response, nil
}

// Query Processor
type QueryProcessor struct {
    tokenizer     Tokenizer
    stemmer       Stemmer
    stopWords     StopWordFilter
    spellChecker  SpellChecker
    synonyms      SynonymExpander
}

type ProcessedQuery struct {
    OriginalQuery string
    Terms         []string
    Phrases       []string
    Operators     []QueryOperator
    Language      string
    Intent        QueryIntent
}

func (qp *QueryProcessor) Process(query, language string) (*ProcessedQuery, error) {
    processed := &ProcessedQuery{
        OriginalQuery: query,
        Language:      language,
    }

    // Normalize query
    normalizedQuery := qp.normalizeQuery(query)

    // Detect query intent
    processed.Intent = qp.detectIntent(normalizedQuery)

    // Extract phrases (quoted text)
    phrases, remainingQuery := qp.extractPhrases(normalizedQuery)
    processed.Phrases = phrases

    // Tokenize remaining query
    tokens := qp.tokenizer.Tokenize(remainingQuery)

    // Remove stop words (but keep for phrase queries)
    filteredTokens := qp.stopWords.Filter(tokens)

    // Spell correction
    correctedTokens := qp.spellChecker.Correct(filteredTokens, language)

    // Stem terms
    stemmedTokens := qp.stemmer.StemTokens(correctedTokens)
    processed.Terms = stemmedTokens

    // Expand with synonyms
    expandedTerms := qp.synonyms.Expand(stemmedTokens, language)
    processed.Terms = append(processed.Terms, expandedTerms...)

    return processed, nil
}

func (qp *QueryProcessor) detectIntent(query string) QueryIntent {
    query = strings.ToLower(query)

    // Navigational queries
    if strings.Contains(query, "site:") || strings.Contains(query, "www.") {
        return IntentNavigational
    }

    // Transactional queries
    transactionalKeywords := []string{"buy", "purchase", "order", "price", "shop"}
    for _, keyword := range transactionalKeywords {
        if strings.Contains(query, keyword) {
            return IntentTransactional
        }
    }

    // Question queries
    questionWords := []string{"what", "how", "why", "when", "where", "who"}
    for _, word := range questionWords {
        if strings.HasPrefix(query, word) {
            return IntentInformational
        }
    }

    return IntentInformational // Default
}
```

### Index Searcher Implementation
```go
type IndexSearcher struct {
    invertedIndex InvertedIndex
    pageRepo      PageRepository
    tfIdfCalculator TfIdfCalculator
}

func (is *IndexSearcher) Search(ctx context.Context, query SearchQuery) ([]SearchResult, error) {
    var allResults []SearchResult

    // Search for each term
    for _, term := range query.Terms {
        termResults, err := is.searchTerm(ctx, term, query)
        if err != nil {
            continue // Handle gracefully
        }
        allResults = append(allResults, termResults...)
    }

    // Search for phrases
    for _, phrase := range query.Phrases {
        phraseResults, err := is.searchPhrase(ctx, phrase, query)
        if err != nil {
            continue
        }
        allResults = append(allResults, phraseResults...)
    }

    // Merge and deduplicate results
    mergedResults := is.mergeResults(allResults)

    // Apply filters
    filteredResults := is.applyFilters(mergedResults, query.Filters)

    // Calculate relevance scores
    scoredResults := is.calculateRelevanceScores(filteredResults, query)

    return scoredResults, nil
}

func (is *IndexSearcher) searchTerm(ctx context.Context, term string, query SearchQuery) ([]SearchResult, error) {
    // Get posting list for term
    postings, err := is.invertedIndex.GetPostings(ctx, term)
    if err != nil {
        return nil, err
    }

    var results []SearchResult
    for _, posting := range postings {
        // Get page information
        page, err := is.pageRepo.GetByID(ctx, posting.PageID)
        if err != nil {
            continue
        }

        // Skip inactive or spam pages
        if !page.IsActive || page.IsSpam {
            continue
        }

        // Calculate TF-IDF score
        tfIdf := is.tfIdfCalculator.Calculate(posting.Frequency, posting.PageID, term)

        result := SearchResult{
            PageID:      posting.PageID,
            URL:         page.URL,
            Title:       page.Title,
            Snippet:     is.generateSnippet(page.CleanedContent, term),
            Score:       tfIdf,
            PageRank:    page.PageRank,
            MatchedTerms: []string{term},
        }

        results = append(results, result)
    }

    return results, nil
}

func (is *IndexSearcher) searchPhrase(ctx context.Context, phrase string, query SearchQuery) ([]SearchResult, error) {
    phraseTerms := strings.Fields(phrase)
    if len(phraseTerms) < 2 {
        return is.searchTerm(ctx, phrase, query)
    }

    // Get posting lists for all terms in phrase
    var termPostings [][]Posting
    for _, term := range phraseTerms {
        postings, err := is.invertedIndex.GetPostings(ctx, term)
        if err != nil {
            return nil, err
        }
        termPostings = append(termPostings, postings)
    }

    // Find documents containing all terms
    commonDocs := is.findCommonDocuments(termPostings)

    var results []SearchResult
    for _, docID := range commonDocs {
        // Verify phrase exists in document
        if is.verifyPhraseInDocument(ctx, docID, phraseTerms) {
            page, err := is.pageRepo.GetByID(ctx, docID)
            if err != nil {
                continue
            }

            result := SearchResult{
                PageID:      docID,
                URL:         page.URL,
                Title:       page.Title,
                Snippet:     is.generateSnippet(page.CleanedContent, phrase),
                Score:       is.calculatePhraseScore(docID, phraseTerms),
                PageRank:    page.PageRank,
                MatchedTerms: phraseTerms,
            }

            results = append(results, result)
        }
    }

    return results, nil
}

// TF-IDF Calculator
type TfIdfCalculator struct {
    invertedIndex InvertedIndex
    totalDocs     int64
}

func (calc *TfIdfCalculator) Calculate(termFreq int, docID int64, term string) float64 {
    // Term Frequency (TF)
    tf := float64(termFreq)

    // Document Frequency (DF)
    df, _ := calc.invertedIndex.GetDocumentFrequency(term)

    // Inverse Document Frequency (IDF)
    idf := math.Log(float64(calc.totalDocs) / float64(df))

    // TF-IDF Score
    return tf * idf
}
```

## 🤖 Step 9: ML-Based Ranking System

### Machine Learning Ranking Service
```go
type RankingService struct {
    mlModel         MLRankingModel
    featureExtractor FeatureExtractor
    pageRankService PageRankService
    qualitySignals  QualitySignalsService
    personalizer    PersonalizationService
    cache          Cache
}

func (rs *RankingService) RankResults(ctx context.Context, request RankingRequest) ([]RankedResult, error) {
    var rankedResults []RankedResult

    for _, result := range request.Results {
        // Extract features for ML model
        features := rs.featureExtractor.Extract(ctx, FeatureRequest{
            Query:    request.Query,
            Document: result,
            UserID:   request.UserID,
            Context:  request.Context,
        })

        // Get ML relevance score
        mlScore := rs.mlModel.PredictRelevance(features)

        // Combine with other signals
        finalScore := rs.combineScores(CombineScoresRequest{
            MLScore:      mlScore,
            TfIdfScore:   result.Score,
            PageRank:     result.PageRank,
            QualityScore: rs.qualitySignals.GetQualityScore(result.PageID),
            PersonalizedScore: rs.personalizer.GetPersonalizedScore(request.UserID, result.PageID),
        })

        rankedResult := RankedResult{
            SearchResult: result,
            FinalScore:   finalScore,
            MLScore:      mlScore,
            Explanation:  rs.generateExplanation(features, mlScore),
        }

        rankedResults = append(rankedResults, rankedResult)
    }

    // Sort by final score
    sort.Slice(rankedResults, func(i, j int) bool {
        return rankedResults[i].FinalScore > rankedResults[j].FinalScore
    })

    return rankedResults, nil
}

// Feature Extractor for ML Model
type FeatureExtractor struct {
    textAnalyzer    TextAnalyzer
    urlAnalyzer     URLAnalyzer
    contentAnalyzer ContentAnalyzer
}

type MLFeatures struct {
    // Query-Document Features
    QueryTermMatches    float64 `json:"query_term_matches"`
    TitleMatches       float64 `json:"title_matches"`
    ExactPhraseMatches float64 `json:"exact_phrase_matches"`
    TfIdfScore         float64 `json:"tf_idf_score"`

    // Document Quality Features
    PageRank           float64 `json:"page_rank"`
    InboundLinks       float64 `json:"inbound_links"`
    ContentLength      float64 `json:"content_length"`
    ContentQuality     float64 `json:"content_quality"`
    URLQuality         float64 `json:"url_quality"`

    // Freshness Features
    LastModified       float64 `json:"last_modified"` // Days since last modified
    CrawlFreshness     float64 `json:"crawl_freshness"`

    // User Context Features
    Language           string  `json:"language"`
    Country            string  `json:"country"`
    DeviceType         string  `json:"device_type"`
    TimeOfDay          int     `json:"time_of_day"`

    // Behavioral Features
    ClickThroughRate   float64 `json:"click_through_rate"`
    DwellTime          float64 `json:"dwell_time"`
    BounceRate         float64 `json:"bounce_rate"`
}

func (fe *FeatureExtractor) Extract(ctx context.Context, request FeatureRequest) *MLFeatures {
    features := &MLFeatures{}

    // Query-Document matching features
    features.QueryTermMatches = fe.calculateTermMatches(request.Query.Terms, request.Document)
    features.TitleMatches = fe.calculateTitleMatches(request.Query.Terms, request.Document.Title)
    features.ExactPhraseMatches = fe.calculatePhraseMatches(request.Query.Phrases, request.Document)
    features.TfIdfScore = request.Document.Score

    // Document quality features
    features.PageRank = request.Document.PageRank
    features.InboundLinks = float64(request.Document.InboundLinks)
    features.ContentLength = float64(len(request.Document.Content))
    features.ContentQuality = fe.contentAnalyzer.AnalyzeQuality(request.Document.Content)
    features.URLQuality = fe.urlAnalyzer.AnalyzeQuality(request.Document.URL)

    // Freshness features
    features.LastModified = float64(time.Since(request.Document.LastModified).Hours() / 24)
    features.CrawlFreshness = float64(time.Since(request.Document.LastCrawled).Hours() / 24)

    // Context features
    features.Language = request.Context.Language
    features.Country = request.Context.Country
    features.DeviceType = request.Context.DeviceType
    features.TimeOfDay = time.Now().Hour()

    // Behavioral features (from historical data)
    behaviorData := fe.getBehaviorData(request.Document.PageID)
    features.ClickThroughRate = behaviorData.CTR
    features.DwellTime = behaviorData.AvgDwellTime
    features.BounceRate = behaviorData.BounceRate

    return features
}

// ML Ranking Model
type MLRankingModel struct {
    model       *tensorflow.SavedModel
    scaler      FeatureScaler
    modelPath   string
    version     string
}

func (mlm *MLRankingModel) PredictRelevance(features *MLFeatures) float64 {
    // Convert features to tensor
    featureVector := mlm.featuresToVector(features)

    // Scale features
    scaledFeatures := mlm.scaler.Scale(featureVector)

    // Create input tensor
    inputTensor, err := tensorflow.NewTensor([][]float32{scaledFeatures})
    if err != nil {
        return 0.5 // Default score
    }

    // Run inference
    results, err := mlm.model.Session.Run(
        map[tensorflow.Output]*tensorflow.Tensor{
            mlm.model.Graph.Operation("input").Output(0): inputTensor,
        },
        []tensorflow.Output{
            mlm.model.Graph.Operation("output").Output(0),
        },
        nil,
    )

    if err != nil {
        return 0.5 // Default score
    }

    // Extract prediction
    predictions := results[0].Value().([][]float32)
    return float64(predictions[0][0])
}

func (mlm *MLRankingModel) featuresToVector(features *MLFeatures) []float32 {
    return []float32{
        float32(features.QueryTermMatches),
        float32(features.TitleMatches),
        float32(features.ExactPhraseMatches),
        float32(features.TfIdfScore),
        float32(features.PageRank),
        float32(features.InboundLinks),
        float32(features.ContentLength),
        float32(features.ContentQuality),
        float32(features.URLQuality),
        float32(features.LastModified),
        float32(features.CrawlFreshness),
        float32(features.TimeOfDay),
        float32(features.ClickThroughRate),
        float32(features.DwellTime),
        float32(features.BounceRate),
    }
}
```

## 🎯 Summary: Search Engine System Design

### Architecture Highlights
- **Distributed Web Crawling** - 1000 crawlers processing 1000 pages/sec each with politeness
- **Inverted Index** - Elasticsearch-based distributed indexing with TF-IDF scoring
- **PageRank Algorithm** - Distributed graph computation with iterative convergence
- **ML-Based Ranking** - TensorFlow model with 15+ features for relevance scoring
- **Query Processing** - Intent detection, spell correction, synonym expansion
- **Real-time Search** - Sub-200ms query processing with multi-level caching

### Scalability Features
- **50B web pages** indexed with distributed storage and processing
- **100K search QPS** with horizontal scaling and caching
- **Global deployment** with regional relevance and language support
- **99.99% availability** with redundancy and failover mechanisms

### Performance Optimizations
- **Search latency**: < 200ms with optimized indexing and caching
- **Crawl efficiency**: 1000 pages/sec per crawler with politeness policies
- **Index freshness**: < 24 hours for popular sites with priority queuing
- **Cache hit rate**: > 80% for popular queries with intelligent caching

This design handles Google-scale search with optimal relevance and performance! 🚀
