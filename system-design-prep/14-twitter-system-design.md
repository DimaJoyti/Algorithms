# 🐦 Twitter/X System Design

## 🎯 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
```
✅ Core Features:
- Users can post tweets (280 characters)
- Users can follow/unfollow other users
- Users can view home timeline (tweets from followed users)
- Users can view user timeline (tweets from specific user)
- Users can like and retweet posts
- Users can search tweets
- Basic user profiles

❌ Out of Scope:
- Direct messaging
- Push notifications
- Advanced media (videos > 2MB)
- Trending topics
- Ads system
```

### Non-Functional Requirements
```
Scale:
- 300M monthly active users
- 150M daily active users
- 400M tweets per day
- 100B timeline views per day
- Read/Write ratio: 100:1

Performance:
- Timeline generation: < 200ms
- Tweet posting: < 100ms
- Search: < 500ms

Availability:
- 99.9% uptime
- Global service
```

## 📊 Step 2: Capacity Estimation

### Back-of-Envelope Calculations
```go
// Daily metrics
const (
    DailyActiveUsers = 150_000_000
    TweetsPerDay     = 400_000_000
    TimelineViews    = 100_000_000_000
    
    SecondsPerDay = 24 * 60 * 60 // ~86,400
)

// QPS calculations
var (
    TweetWriteQPS = TweetsPerDay / SecondsPerDay        // ~4,630 QPS
    TimelineReadQPS = TimelineViews / SecondsPerDay     // ~1,157,407 QPS
    PeakWriteQPS = TweetWriteQPS * 3                    // ~13,890 QPS
    PeakReadQPS = TimelineReadQPS * 3                   // ~3,472,222 QPS
)

// Storage estimation
type StorageCalculation struct {
    TweetSize      int // 280 chars + metadata = ~300 bytes
    UserSize       int // Profile data = ~1KB
    MediaSize      int // Average 200KB per media tweet
    
    DailyStorage   int64 // Tweets + media
    YearlyStorage  int64 // With replication
}

func CalculateStorage() StorageCalculation {
    calc := StorageCalculation{
        TweetSize: 300,
        UserSize:  1024,
        MediaSize: 200 * 1024,
    }
    
    // Daily storage
    textTweets := int64(TweetsPerDay) * 0.8 * int64(calc.TweetSize)      // 80% text
    mediaTweets := int64(TweetsPerDay) * 0.2 * int64(calc.MediaSize)     // 20% media
    calc.DailyStorage = textTweets + mediaTweets                         // ~16GB/day
    
    // 5-year storage with 3x replication
    calc.YearlyStorage = calc.DailyStorage * 365 * 5 * 3                 // ~87TB
    
    return calc
}

// Memory for caching
func CalculateCacheMemory() int64 {
    // Cache 20% of daily timeline reads
    cacheHitRate := 0.2
    avgTimelineSize := 50 * 300 // 50 tweets * 300 bytes each
    
    dailyCacheReads := int64(float64(TimelineViews) * cacheHitRate)
    return dailyCacheReads * int64(avgTimelineSize) // ~300TB cache needed
}
```

### Infrastructure Estimation
```go
// Server requirements
type InfrastructureNeeds struct {
    WebServers    int // Handle API requests
    DBServers     int // Database instances
    CacheServers  int // Redis/Memcached
    CDNNodes      int // Global content delivery
}

func EstimateInfrastructure() InfrastructureNeeds {
    // Assuming each server handles 10K QPS
    webServers := int(PeakReadQPS / 10000) // ~347 servers
    
    // Database sharding - each shard handles 5K QPS
    dbServers := int(PeakReadQPS / 5000) // ~694 DB instances
    
    // Cache servers - 64GB RAM each
    totalCacheMemory := CalculateCacheMemory()
    cacheServers := int(totalCacheMemory / (64 * 1024 * 1024 * 1024)) // ~4,687 cache servers
    
    return InfrastructureNeeds{
        WebServers:   webServers,
        DBServers:    dbServers,
        CacheServers: cacheServers,
        CDNNodes:     50, // Global distribution
    }
}
```

## 🏗️ Step 3: High-Level Design

### System Architecture
```
[Mobile Apps] ──┐
                ├── [Load Balancer] ── [API Gateway]
[Web Apps] ─────┘                           │
                                           ├── [User Service]
                                           ├── [Tweet Service]
                                           ├── [Timeline Service]
                                           ├── [Follow Service]
                                           ├── [Search Service]
                                           └── [Media Service]
                                                    │
                                    [Message Queue] ── [Cache Layer] ── [Database Layer]
```

### Core Services
```go
// User Service - manages user profiles and authentication
type UserService struct {
    userRepo     UserRepository
    authService  AuthService
    cache        Cache
    eventBus     EventBus
}

// Tweet Service - handles tweet creation and retrieval
type TweetService struct {
    tweetRepo    TweetRepository
    mediaService MediaService
    cache        Cache
    eventBus     EventBus
}

// Timeline Service - generates user timelines
type TimelineService struct {
    timelineRepo TimelineRepository
    tweetService TweetService
    followService FollowService
    cache        Cache
}

// Follow Service - manages user relationships
type FollowService struct {
    followRepo   FollowRepository
    cache        Cache
    eventBus     EventBus
}

// Search Service - handles tweet search
type SearchService struct {
    searchIndex  SearchIndex // Elasticsearch
    tweetService TweetService
}
```

## 🗄️ Step 4: Database Design

### Data Models
```go
// User model
type User struct {
    ID              int64     `json:"id" db:"id"`
    Username        string    `json:"username" db:"username"`
    Email           string    `json:"email" db:"email"`
    DisplayName     string    `json:"display_name" db:"display_name"`
    Bio             string    `json:"bio" db:"bio"`
    ProfileImageURL string    `json:"profile_image_url" db:"profile_image_url"`
    FollowersCount  int       `json:"followers_count" db:"followers_count"`
    FollowingCount  int       `json:"following_count" db:"following_count"`
    TweetsCount     int       `json:"tweets_count" db:"tweets_count"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Tweet model
type Tweet struct {
    ID              int64     `json:"id" db:"id"`
    UserID          int64     `json:"user_id" db:"user_id"`
    Content         string    `json:"content" db:"content"`
    MediaURLs       []string  `json:"media_urls" db:"media_urls"`
    RetweetID       *int64    `json:"retweet_id,omitempty" db:"retweet_id"`
    ReplyToID       *int64    `json:"reply_to_id,omitempty" db:"reply_to_id"`
    LikesCount      int       `json:"likes_count" db:"likes_count"`
    RetweetsCount   int       `json:"retweets_count" db:"retweets_count"`
    RepliesCount    int       `json:"replies_count" db:"replies_count"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Follow relationship
type Follow struct {
    FollowerID  int64     `json:"follower_id" db:"follower_id"`
    FollowingID int64     `json:"following_id" db:"following_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Like relationship
type Like struct {
    UserID    int64     `json:"user_id" db:"user_id"`
    TweetID   int64     `json:"tweet_id" db:"tweet_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### Database Schema
```sql
-- Users table (PostgreSQL)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    profile_image_url TEXT,
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    tweets_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Tweets table (Cassandra-like structure for high write volume)
CREATE TABLE tweets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    media_urls TEXT[], -- JSON array
    retweet_id BIGINT,
    reply_to_id BIGINT,
    likes_count INTEGER DEFAULT 0,
    retweets_count INTEGER DEFAULT 0,
    replies_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tweets_user_id ON tweets(user_id, created_at DESC);
CREATE INDEX idx_tweets_created_at ON tweets(created_at DESC);

-- Follows table (PostgreSQL)
CREATE TABLE follows (
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);

-- Likes table (Cassandra for high write volume)
CREATE TABLE likes (
    user_id BIGINT NOT NULL,
    tweet_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, tweet_id)
);

CREATE INDEX idx_likes_tweet_id ON likes(tweet_id);
```

### Database Sharding Strategy
```go
type DatabaseSharding struct {
    userShards  []Database // Shard by user_id
    tweetShards []Database // Shard by tweet_id or user_id
}

// User sharding - by user ID
func (ds *DatabaseSharding) GetUserShard(userID int64) Database {
    shardIndex := userID % int64(len(ds.userShards))
    return ds.userShards[shardIndex]
}

// Tweet sharding - by user ID for better locality
func (ds *DatabaseSharding) GetTweetShard(userID int64) Database {
    shardIndex := userID % int64(len(ds.tweetShards))
    return ds.tweetShards[shardIndex]
}

// Alternative: Tweet sharding by tweet ID for even distribution
func (ds *DatabaseSharding) GetTweetShardByID(tweetID int64) Database {
    shardIndex := tweetID % int64(len(ds.tweetShards))
    return ds.tweetShards[shardIndex]
}
```

## 🔄 Step 5: Timeline Generation

### Timeline Generation Strategies

#### 1. Pull Model (Fan-out on Read)
```go
type PullTimelineService struct {
    tweetService  TweetService
    followService FollowService
    cache         Cache
}

func (pts *PullTimelineService) GenerateTimeline(ctx context.Context, userID int64, limit int) ([]*Tweet, error) {
    cacheKey := fmt.Sprintf("timeline:pull:%d", userID)
    
    // Check cache first
    if cached, found := pts.cache.Get(cacheKey); found {
        return cached.([]*Tweet), nil
    }
    
    // Get users that this user follows
    following, err := pts.followService.GetFollowing(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Add user's own tweets
    following = append(following, userID)
    
    // Get recent tweets from all followed users
    var allTweets []*Tweet
    for _, followedUserID := range following {
        tweets, err := pts.tweetService.GetUserTweets(ctx, followedUserID, 100)
        if err != nil {
            continue // Handle gracefully
        }
        allTweets = append(allTweets, tweets...)
    }
    
    // Sort by timestamp (merge k sorted lists algorithm)
    sort.Slice(allTweets, func(i, j int) bool {
        return allTweets[i].CreatedAt.After(allTweets[j].CreatedAt)
    })
    
    // Take top N tweets
    if len(allTweets) > limit {
        allTweets = allTweets[:limit]
    }
    
    // Cache result
    pts.cache.Set(cacheKey, allTweets, 5*time.Minute)
    
    return allTweets, nil
}

// Pros: Storage efficient, good for users with many follows
// Cons: High latency, expensive for active users
```

#### 2. Push Model (Fan-out on Write)
```go
type PushTimelineService struct {
    followService FollowService
    cache         Cache
    eventBus      EventBus
}

func (pts *PushTimelineService) HandleTweetCreated(ctx context.Context, event TweetCreatedEvent) error {
    // Get all followers of the tweet author
    followers, err := pts.followService.GetFollowers(ctx, event.UserID)
    if err != nil {
        return err
    }
    
    // Fan-out to all followers' timelines
    for _, followerID := range followers {
        timelineKey := fmt.Sprintf("timeline:push:%d", followerID)
        
        // Add tweet to beginning of timeline
        pts.cache.LPush(timelineKey, event.TweetID)
        
        // Keep only latest 1000 tweets
        pts.cache.LTrim(timelineKey, 0, 999)
        
        // Set expiration
        pts.cache.Expire(timelineKey, 24*time.Hour)
    }
    
    return nil
}

func (pts *PushTimelineService) GetTimeline(ctx context.Context, userID int64, limit int) ([]*Tweet, error) {
    timelineKey := fmt.Sprintf("timeline:push:%d", userID)
    
    // Get tweet IDs from cache
    tweetIDs, err := pts.cache.LRange(timelineKey, 0, limit-1)
    if err != nil {
        return nil, err
    }
    
    // Fetch tweet details
    var tweets []*Tweet
    for _, tweetID := range tweetIDs {
        tweet, err := pts.tweetService.GetTweet(ctx, tweetID)
        if err != nil {
            continue // Handle gracefully
        }
        tweets = append(tweets, tweet)
    }
    
    return tweets, nil
}

// Pros: Fast reads, pre-computed timelines
// Cons: High storage cost, expensive writes for celebrities
```

#### 3. Hybrid Model (Best of Both)
```go
type HybridTimelineService struct {
    pushService  *PushTimelineService
    pullService  *PullTimelineService
    userService  UserService
}

func (hts *HybridTimelineService) GenerateTimeline(ctx context.Context, userID int64, limit int) ([]*Tweet, error) {
    user, err := hts.userService.GetUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Use push model for regular users
    if user.FollowingCount <= 1000 {
        return hts.pushService.GetTimeline(ctx, userID, limit)
    }
    
    // Use pull model for users following many people
    return hts.pullService.GenerateTimeline(ctx, userID, limit)
}

func (hts *HybridTimelineService) HandleTweetCreated(ctx context.Context, event TweetCreatedEvent) error {
    user, err := hts.userService.GetUser(ctx, event.UserID)
    if err != nil {
        return err
    }
    
    // Don't fan-out for celebrities (>1M followers)
    if user.FollowersCount > 1_000_000 {
        return nil
    }
    
    // Use push model for regular users
    return hts.pushService.HandleTweetCreated(ctx, event)
}

// Pros: Optimized for different user types
// Cons: More complex implementation
```

## 🔍 Step 6: Search Implementation

### Elasticsearch Integration
```go
type SearchService struct {
    esClient *elasticsearch.Client
    tweetService TweetService
}

func (ss *SearchService) IndexTweet(ctx context.Context, tweet *Tweet) error {
    doc := map[string]interface{}{
        "id":         tweet.ID,
        "user_id":    tweet.UserID,
        "content":    tweet.Content,
        "created_at": tweet.CreatedAt,
        "likes_count": tweet.LikesCount,
        "retweets_count": tweet.RetweetsCount,
    }
    
    body, err := json.Marshal(doc)
    if err != nil {
        return err
    }
    
    req := esapi.IndexRequest{
        Index:      "tweets",
        DocumentID: strconv.FormatInt(tweet.ID, 10),
        Body:       bytes.NewReader(body),
        Refresh:    "true",
    }
    
    res, err := req.Do(ctx, ss.esClient)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    
    return nil
}

func (ss *SearchService) SearchTweets(ctx context.Context, query string, limit int) ([]*Tweet, error) {
    searchQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":  query,
                "fields": []string{"content^2", "user_id"}, // Boost content matches
            },
        },
        "sort": []map[string]interface{}{
            {"created_at": map[string]string{"order": "desc"}},
        },
        "size": limit,
    }
    
    body, err := json.Marshal(searchQuery)
    if err != nil {
        return nil, err
    }
    
    res, err := ss.esClient.Search(
        ss.esClient.Search.WithContext(ctx),
        ss.esClient.Search.WithIndex("tweets"),
        ss.esClient.Search.WithBody(bytes.NewReader(body)),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    
    // Parse results and fetch full tweet details
    var searchResults SearchResults
    if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
        return nil, err
    }
    
    var tweets []*Tweet
    for _, hit := range searchResults.Hits.Hits {
        tweetID, _ := strconv.ParseInt(hit.Source.ID, 10, 64)
        tweet, err := ss.tweetService.GetTweet(ctx, tweetID)
        if err != nil {
            continue
        }
        tweets = append(tweets, tweet)
    }
    
    return tweets, nil
}
```

## 🚀 Step 7: Caching Strategy

### Multi-Level Caching
```go
type TwitterCacheStrategy struct {
    L1Cache *LRUCache      // Application memory (user sessions, hot tweets)
    L2Cache *RedisCluster  // Distributed cache (timelines, user data)
    CDN     *CDNService    // Edge cache (media, static content)
}

// Cache hot tweets (trending, viral content)
func (tcs *TwitterCacheStrategy) CacheHotTweet(tweetID int64, tweet *Tweet) {
    key := fmt.Sprintf("hot_tweet:%d", tweetID)

    // Cache in all levels for maximum performance
    tcs.L1Cache.Set(key, tweet, 10*time.Minute)
    tcs.L2Cache.Set(key, tweet, 1*time.Hour)
}

// Cache user timelines with different TTLs
func (tcs *TwitterCacheStrategy) CacheTimeline(userID int64, tweets []*Tweet, timelineType string) {
    key := fmt.Sprintf("timeline:%s:%d", timelineType, userID)

    var ttl time.Duration
    switch timelineType {
    case "home":
        ttl = 5 * time.Minute  // Home timeline changes frequently
    case "user":
        ttl = 15 * time.Minute // User timeline more stable
    default:
        ttl = 10 * time.Minute
    }

    tcs.L2Cache.Set(key, tweets, ttl)
}

// Cache user data with write-through pattern
func (tcs *TwitterCacheStrategy) CacheUser(user *User) {
    key := fmt.Sprintf("user:%d", user.ID)

    // Cache user profile for 1 hour
    tcs.L2Cache.Set(key, user, 1*time.Hour)

    // Cache frequently accessed users in L1
    if user.FollowersCount > 100000 { // Popular users
        tcs.L1Cache.Set(key, user, 30*time.Minute)
    }
}
```

### Cache Invalidation Strategy
```go
type CacheInvalidationService struct {
    cache    Cache
    eventBus EventBus
}

func (cis *CacheInvalidationService) HandleTweetCreated(ctx context.Context, event TweetCreatedEvent) {
    // Invalidate author's user timeline
    authorTimelineKey := fmt.Sprintf("timeline:user:%d", event.UserID)
    cis.cache.Delete(authorTimelineKey)

    // For push model: timelines are updated, not invalidated
    // For pull model: invalidate followers' home timelines
    if event.UserFollowersCount < 1000000 { // Not a celebrity
        followers, _ := cis.getFollowers(event.UserID)
        for _, followerID := range followers {
            homeTimelineKey := fmt.Sprintf("timeline:home:%d", followerID)
            cis.cache.Delete(homeTimelineKey)
        }
    }
}

func (cis *CacheInvalidationService) HandleUserUpdated(ctx context.Context, event UserUpdatedEvent) {
    // Invalidate user cache
    userKey := fmt.Sprintf("user:%d", event.UserID)
    cis.cache.Delete(userKey)

    // Invalidate any cached tweets that include user info
    // This is expensive, so we might use eventual consistency here
}
```

## 📱 Step 8: Real-time Features

### WebSocket Implementation for Real-time Updates
```go
type RealtimeService struct {
    connections map[int64]*websocket.Conn // userID -> connection
    mutex       sync.RWMutex
    eventBus    EventBus
}

func (rs *RealtimeService) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromToken(r)

    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()

    // Register connection
    rs.mutex.Lock()
    rs.connections[userID] = conn
    rs.mutex.Unlock()

    // Clean up on disconnect
    defer func() {
        rs.mutex.Lock()
        delete(rs.connections, userID)
        rs.mutex.Unlock()
    }()

    // Keep connection alive
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func (rs *RealtimeService) HandleTweetCreated(ctx context.Context, event TweetCreatedEvent) {
    // Get followers of the tweet author
    followers, err := rs.getFollowers(event.UserID)
    if err != nil {
        return
    }

    // Send real-time update to online followers
    notification := RealtimeNotification{
        Type: "new_tweet",
        Data: map[string]interface{}{
            "tweet_id": event.TweetID,
            "user_id":  event.UserID,
            "content":  event.Content,
        },
    }

    rs.mutex.RLock()
    for _, followerID := range followers {
        if conn, exists := rs.connections[followerID]; exists {
            go func(c *websocket.Conn) {
                c.WriteJSON(notification)
            }(conn)
        }
    }
    rs.mutex.RUnlock()
}
```

## 📊 Step 9: Media Storage & CDN

### Media Service Implementation
```go
type MediaService struct {
    s3Client    *s3.Client
    cdnService  CDNService
    imageProcessor ImageProcessor
}

func (ms *MediaService) UploadMedia(ctx context.Context, userID int64, file io.Reader, contentType string) (*MediaUploadResult, error) {
    // Generate unique filename
    mediaID := generateMediaID()
    filename := fmt.Sprintf("media/%d/%s", userID, mediaID)

    // Process image (resize, compress)
    processedImage, err := ms.imageProcessor.Process(file, ImageProcessingOptions{
        MaxWidth:  1200,
        MaxHeight: 1200,
        Quality:   85,
        Format:    "jpeg",
    })
    if err != nil {
        return nil, err
    }

    // Upload to S3
    _, err = ms.s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String("twitter-media"),
        Key:         aws.String(filename),
        Body:        processedImage,
        ContentType: aws.String(contentType),
        ACL:         types.ObjectCannedACLPublicRead,
    })
    if err != nil {
        return nil, err
    }

    // Generate CDN URL
    cdnURL := fmt.Sprintf("https://cdn.twitter.com/%s", filename)

    return &MediaUploadResult{
        MediaID: mediaID,
        URL:     cdnURL,
        Size:    processedImage.Len(),
    }, nil
}

// CDN configuration for global media delivery
type CDNConfiguration struct {
    Origins []CDNOrigin
    Rules   []CDNRule
}

var TwitterCDNConfig = CDNConfiguration{
    Origins: []CDNOrigin{
        {Name: "media-origin", URL: "https://s3.amazonaws.com/twitter-media"},
        {Name: "static-origin", URL: "https://s3.amazonaws.com/twitter-static"},
    },
    Rules: []CDNRule{
        {
            Path:        "/media/*",
            Origin:      "media-origin",
            CacheTTL:    24 * time.Hour,
            Compression: true,
        },
        {
            Path:        "/static/*",
            Origin:      "static-origin",
            CacheTTL:    7 * 24 * time.Hour,
            Compression: true,
        },
    },
}
```

## 🔧 Step 10: API Implementation

### RESTful API Design
```go
type TwitterAPI struct {
    userService     UserService
    tweetService    TweetService
    timelineService TimelineService
    followService   FollowService
    searchService   SearchService
    authMiddleware  AuthMiddleware
    rateLimiter     RateLimiter
}

func (api *TwitterAPI) SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    // Apply middleware
    authHandler := api.authMiddleware.Authenticate
    rateLimitHandler := api.rateLimiter.Middleware

    // User endpoints
    mux.HandleFunc("GET /api/v1/users/{id}", authHandler(rateLimitHandler(api.getUser)))
    mux.HandleFunc("PUT /api/v1/users/{id}", authHandler(rateLimitHandler(api.updateUser)))

    // Tweet endpoints
    mux.HandleFunc("POST /api/v1/tweets", authHandler(rateLimitHandler(api.createTweet)))
    mux.HandleFunc("GET /api/v1/tweets/{id}", rateLimitHandler(api.getTweet))
    mux.HandleFunc("DELETE /api/v1/tweets/{id}", authHandler(rateLimitHandler(api.deleteTweet)))

    // Timeline endpoints
    mux.HandleFunc("GET /api/v1/users/{id}/timeline", rateLimitHandler(api.getUserTimeline))
    mux.HandleFunc("GET /api/v1/timeline/home", authHandler(rateLimitHandler(api.getHomeTimeline)))

    // Social endpoints
    mux.HandleFunc("POST /api/v1/users/{id}/follow", authHandler(rateLimitHandler(api.followUser)))
    mux.HandleFunc("DELETE /api/v1/users/{id}/follow", authHandler(rateLimitHandler(api.unfollowUser)))
    mux.HandleFunc("POST /api/v1/tweets/{id}/like", authHandler(rateLimitHandler(api.likeTweet)))

    // Search endpoints
    mux.HandleFunc("GET /api/v1/search/tweets", rateLimitHandler(api.searchTweets))

    return mux
}

func (api *TwitterAPI) createTweet(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())

    var req CreateTweetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validate tweet content
    if len(req.Content) == 0 || len(req.Content) > 280 {
        writeErrorResponse(w, http.StatusBadRequest, "Tweet content must be 1-280 characters")
        return
    }

    // Create tweet
    tweet, err := api.tweetService.CreateTweet(r.Context(), CreateTweetParams{
        UserID:    userID,
        Content:   req.Content,
        MediaURLs: req.MediaURLs,
        ReplyToID: req.ReplyToID,
    })
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to create tweet")
        return
    }

    // Return created tweet
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tweet)
}

func (api *TwitterAPI) getHomeTimeline(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())

    // Parse pagination parameters
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 || limit > 100 {
        limit = 20
    }

    cursor := r.URL.Query().Get("cursor")

    // Get timeline
    timeline, nextCursor, err := api.timelineService.GetHomeTimeline(r.Context(), userID, limit, cursor)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to get timeline")
        return
    }

    response := TimelineResponse{
        Tweets:     timeline,
        NextCursor: nextCursor,
        HasMore:    nextCursor != "",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## 📈 Step 11: Monitoring & Observability

### Key Metrics to Track
```go
type TwitterMetrics struct {
    // Business metrics
    TweetsPerSecond     prometheus.Counter
    TimelineGeneration  prometheus.Histogram
    SearchLatency       prometheus.Histogram
    UserEngagement      prometheus.Counter

    // Technical metrics
    DatabaseConnections prometheus.Gauge
    CacheHitRate       prometheus.Gauge
    QueueDepth         prometheus.Gauge
    ErrorRate          prometheus.Counter
}

func (tm *TwitterMetrics) RecordTweetCreated(userID int64, duration time.Duration) {
    tm.TweetsPerSecond.Inc()

    // Track by user type
    labels := prometheus.Labels{"user_type": "regular"}
    if tm.isVerifiedUser(userID) {
        labels["user_type"] = "verified"
    }

    tm.TweetsPerSecond.With(labels).Inc()
}

func (tm *TwitterMetrics) RecordTimelineGeneration(timelineType string, duration time.Duration, cacheHit bool) {
    labels := prometheus.Labels{
        "timeline_type": timelineType,
        "cache_hit":     strconv.FormatBool(cacheHit),
    }

    tm.TimelineGeneration.With(labels).Observe(duration.Seconds())
}
```

### Health Checks
```go
type TwitterHealthChecker struct {
    dbHealth     DatabaseHealthCheck
    cacheHealth  CacheHealthCheck
    searchHealth SearchHealthCheck
}

func (thc *TwitterHealthChecker) CheckHealth(ctx context.Context) HealthStatus {
    checks := map[string]HealthCheck{
        "database": &thc.dbHealth,
        "cache":    &thc.cacheHealth,
        "search":   &thc.searchHealth,
    }

    results := make(map[string]HealthStatus)
    overallHealthy := true

    for name, check := range checks {
        status := check.Check(ctx)
        results[name] = status

        if status.Status != "healthy" {
            overallHealthy = false
        }
    }

    overallStatus := "healthy"
    if !overallHealthy {
        overallStatus = "unhealthy"
    }

    return HealthStatus{
        Status:  overallStatus,
        Details: results,
    }
}
```

## 🎯 Summary: Twitter System Design

### Architecture Highlights
- **Microservices**: User, Tweet, Timeline, Follow, Search, Media services
- **Database**: PostgreSQL for users/follows, Cassandra for tweets/likes
- **Caching**: Multi-level (L1/L2/CDN) with smart invalidation
- **Timeline**: Hybrid push/pull model optimized for different user types
- **Search**: Elasticsearch for full-text search with real-time indexing
- **Media**: S3 + CDN for global media delivery
- **Real-time**: WebSockets for live updates

### Scalability Features
- **Horizontal scaling**: Microservices + database sharding
- **Caching strategy**: Hot data cached at multiple levels
- **CDN**: Global media and static content delivery
- **Load balancing**: Geographic and service-based routing
- **Async processing**: Event-driven architecture for heavy operations

### Performance Optimizations
- **Timeline generation**: < 200ms with hybrid approach
- **Tweet creation**: < 100ms with async fan-out
- **Search**: < 500ms with Elasticsearch
- **Media delivery**: < 50ms with CDN
- **Cache hit rates**: > 80% for hot data

This design can handle 300M users, 400M tweets/day, and 100B timeline views/day with 99.9% availability! 🚀
