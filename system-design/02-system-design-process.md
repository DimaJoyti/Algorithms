# 🎯 System Design Interview Process

## 🎪 The 45-Minute System Design Interview

### Timeline Breakdown
- **5 minutes**: Requirements clarification
- **10 minutes**: High-level design
- **15 minutes**: Detailed design
- **10 minutes**: Scale and optimize
- **5 minutes**: Wrap up and questions

## 📋 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
**Ask the right questions to understand WHAT the system should do**

#### Example: Design Twitter
```
❓ Questions to Ask:
- Can users post tweets? (Core feature)
- Can users follow other users? (Social graph)
- Can users see a timeline of tweets? (Feed generation)
- Can users search tweets? (Search functionality)
- Can users like/retweet? (Engagement features)
- Do we need direct messaging? (Scope boundary)
- Do we need notifications? (Real-time features)

✅ Clarified Requirements:
- Users can post tweets (280 characters)
- Users can follow/unfollow others
- Users see home timeline (following + own tweets)
- Users can like and retweet
- Basic search functionality
- No DMs or notifications (out of scope)
```

### Non-Functional Requirements
**Ask about system qualities and constraints**

#### Performance Requirements
```
- How many users? (Scale)
- How many tweets per day? (Write load)
- How many timeline views per day? (Read load)
- What's the read/write ratio? (Caching strategy)
- What's acceptable latency? (Performance target)
- What's required availability? (Reliability target)
```

#### Example Numbers for Twitter
```
- 300M active users
- 200M tweets per day
- 100B timeline views per day
- Read/Write ratio: 500:1
- Timeline latency: < 200ms
- Availability: 99.9%
```

### Constraints and Assumptions
```
- Global service or specific region?
- Mobile app, web, or both?
- Any compliance requirements?
- Budget constraints?
- Technology preferences?
```

## 🏗️ Step 2: High-Level Design (10 minutes)

### Core Components Identification

#### For Twitter Example:
```
1. User Service - manage user profiles
2. Tweet Service - handle tweet creation/storage
3. Timeline Service - generate user timelines
4. Follow Service - manage social graph
5. Search Service - tweet search functionality
```

### Basic Architecture Diagram
```
[Mobile App] ──┐
               ├── [Load Balancer] ── [API Gateway]
[Web App] ─────┘                           │
                                          ├── [User Service]
                                          ├── [Tweet Service]
                                          ├── [Timeline Service]
                                          ├── [Follow Service]
                                          └── [Search Service]
                                                    │
                                          [Database Layer]
```

### API Design
**Define key endpoints early**

```go
// User Service APIs
POST /api/v1/users/register
GET  /api/v1/users/{userId}
POST /api/v1/users/{userId}/follow
DELETE /api/v1/users/{userId}/follow

// Tweet Service APIs
POST /api/v1/tweets
GET  /api/v1/tweets/{tweetId}
POST /api/v1/tweets/{tweetId}/like
POST /api/v1/tweets/{tweetId}/retweet

// Timeline Service APIs
GET  /api/v1/users/{userId}/timeline
GET  /api/v1/users/{userId}/timeline/home
```

### Data Models
**Define core entities**

```go
type User struct {
    ID          int64     `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    CreatedAt   time.Time `json:"created_at"`
    FollowersCount int    `json:"followers_count"`
    FollowingCount int    `json:"following_count"`
}

type Tweet struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    LikesCount    int   `json:"likes_count"`
    RetweetsCount int   `json:"retweets_count"`
}

type Follow struct {
    FollowerID  int64     `json:"follower_id"`
    FollowingID int64     `json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## 🔧 Step 3: Detailed Design (15 minutes)

### Database Design

#### SQL vs NoSQL Decision
```
Twitter Example Analysis:
- Social graph: Complex relationships → SQL (PostgreSQL)
- Tweets: High write volume, simple queries → NoSQL (Cassandra)
- User profiles: ACID properties needed → SQL (PostgreSQL)
- Timeline cache: Fast reads → Redis
```

#### Database Schema
```sql
-- Users table (PostgreSQL)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    followers_count INT DEFAULT 0,
    following_count INT DEFAULT 0
);

-- Follows table (PostgreSQL)
CREATE TABLE follows (
    follower_id BIGINT REFERENCES users(id),
    following_id BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

-- Tweets table (Cassandra-like structure)
CREATE TABLE tweets (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,
    content TEXT,
    created_at TIMESTAMP,
    likes_count INT,
    retweets_count INT
);
```

### Service Implementation Details

#### Tweet Service
```go
type TweetService struct {
    tweetRepo   TweetRepository
    userRepo    UserRepository
    timelineService TimelineService
    cache       Cache
}

func (ts *TweetService) CreateTweet(ctx context.Context, req CreateTweetRequest) (*Tweet, error) {
    // 1. Validate input
    if len(req.Content) > 280 {
        return nil, ErrTweetTooLong
    }
    
    // 2. Create tweet
    tweet := &Tweet{
        ID:        generateTweetID(), // Snowflake ID
        UserID:    req.UserID,
        Content:   req.Content,
        CreatedAt: time.Now(),
    }
    
    // 3. Store in database
    if err := ts.tweetRepo.Create(ctx, tweet); err != nil {
        return nil, err
    }
    
    // 4. Async: Update followers' timelines
    go ts.timelineService.FanOutTweet(ctx, tweet)
    
    return tweet, nil
}
```

#### Timeline Generation Strategies

##### Push Model (Fan-out on Write)
```go
func (ts *TimelineService) FanOutTweet(ctx context.Context, tweet *Tweet) error {
    // Get all followers
    followers, err := ts.followRepo.GetFollowers(ctx, tweet.UserID)
    if err != nil {
        return err
    }
    
    // Add tweet to each follower's timeline
    for _, followerID := range followers {
        timelineKey := fmt.Sprintf("timeline:%d", followerID)
        ts.cache.LPush(timelineKey, tweet.ID)
        ts.cache.LTrim(timelineKey, 0, 999) // Keep latest 1000 tweets
    }
    
    return nil
}
```

##### Pull Model (Fan-out on Read)
```go
func (ts *TimelineService) GetTimeline(ctx context.Context, userID int64) ([]*Tweet, error) {
    // Get users that this user follows
    following, err := ts.followRepo.GetFollowing(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Get recent tweets from all followed users
    var allTweets []*Tweet
    for _, followedUserID := range following {
        tweets, err := ts.tweetRepo.GetUserTweets(ctx, followedUserID, 100)
        if err != nil {
            continue // Handle gracefully
        }
        allTweets = append(allTweets, tweets...)
    }
    
    // Sort by timestamp (merge k sorted lists algorithm)
    sort.Slice(allTweets, func(i, j int) bool {
        return allTweets[i].CreatedAt.After(allTweets[j].CreatedAt)
    })
    
    return allTweets[:min(len(allTweets), 50)], nil
}
```

### Caching Strategy
```go
type CacheLayer struct {
    userCache     Cache // User profiles
    tweetCache    Cache // Individual tweets
    timelineCache Cache // Pre-computed timelines
    searchCache   Cache // Search results
}

// Cache patterns
func (c *CacheLayer) GetTweet(tweetID int64) (*Tweet, error) {
    // 1. Check cache first
    if tweet, found := c.tweetCache.Get(fmt.Sprintf("tweet:%d", tweetID)); found {
        return tweet.(*Tweet), nil
    }
    
    // 2. Cache miss - get from database
    tweet, err := c.tweetRepo.GetByID(tweetID)
    if err != nil {
        return nil, err
    }
    
    // 3. Store in cache
    c.tweetCache.Set(fmt.Sprintf("tweet:%d", tweetID), tweet, 1*time.Hour)
    
    return tweet, nil
}
```

## ⚡ Step 4: Scale and Optimize (10 minutes)

### Identify Bottlenecks
```
1. Database writes (200M tweets/day)
2. Timeline generation (100B views/day)
3. Search queries (complex text search)
4. Hot users (celebrities with millions of followers)
```

### Scaling Solutions

#### Database Scaling
```
1. Read Replicas: Route reads to replicas
2. Sharding: Partition by user_id or tweet_id
3. Caching: Redis for hot data
4. Denormalization: Store computed values
```

#### Timeline Optimization
```
1. Hybrid approach: Push for normal users, pull for celebrities
2. Pre-compute timelines for active users
3. Use message queues for async processing
4. Implement timeline pagination
```

#### Search Optimization
```
1. Elasticsearch for full-text search
2. Search result caching
3. Auto-complete with Trie data structure
4. Trending topics with real-time analytics
```

### Architecture Evolution
```
Initial: Monolith → Microservices → Event-Driven → Serverless (where appropriate)
```

## 🎯 Step 5: Wrap Up (5 minutes)

### Summary Checklist
- ✅ Addressed all functional requirements
- ✅ Met non-functional requirements
- ✅ Identified and solved bottlenecks
- ✅ Discussed trade-offs
- ✅ Considered failure scenarios

### Common Follow-up Questions
```
1. How would you handle service failures?
2. How would you monitor the system?
3. How would you deploy updates?
4. How would you handle data consistency?
5. How would you secure the system?
```

## 🎪 Interview Tips

### Do's
- ✅ Ask clarifying questions
- ✅ Start simple, then add complexity
- ✅ Explain your reasoning
- ✅ Draw diagrams
- ✅ Discuss trade-offs
- ✅ Consider edge cases

### Don'ts
- ❌ Jump into details immediately
- ❌ Design for maximum scale from start
- ❌ Ignore the interviewer's hints
- ❌ Get stuck on one component
- ❌ Forget about failure scenarios

## 🔗 Next Steps
- Practice with real system design problems
- Learn back-of-envelope calculations
- Study trade-offs and constraints
