# 🎬 Netflix/Video Streaming System Design

## 🎯 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
```
✅ Core Features:
- Users can browse and search video content
- Users can stream videos in multiple qualities (480p, 720p, 1080p, 4K)
- Users can pause, resume, and seek within videos
- Users can create profiles and maintain watch history
- Users can rate and review content
- Content recommendation system
- Download videos for offline viewing
- Multiple device support (TV, mobile, web)

❌ Out of Scope:
- Live streaming
- User-generated content
- Social features (sharing, comments)
- Payment processing
- Content creation tools
```

### Non-Functional Requirements
```
Scale:
- 200M subscribers globally
- 1B hours watched per day
- 15K different titles
- 100K concurrent streams during peak
- 50 PB of video content storage

Performance:
- Video start time: < 2 seconds
- Buffering: < 0.1% of play time
- 99.9% availability
- Support for 4K streaming

Global:
- Serve users worldwide
- Content localization
- Regional content licensing
```

## 📊 Step 2: Capacity Estimation

### Back-of-Envelope Calculations
```go
// Global scale metrics
const (
    TotalSubscribers     = 200_000_000
    DailyActiveUsers     = 100_000_000  // 50% of subscribers watch daily
    HoursWatchedPerDay   = 1_000_000_000 // 1B hours
    PeakConcurrentUsers  = 100_000       // During prime time
    
    SecondsPerDay = 24 * 60 * 60
    SecondsPerHour = 60 * 60
)

// Video streaming calculations
var (
    AvgWatchTimePerUser = HoursWatchedPerDay / DailyActiveUsers // 10 hours per user
    PeakStreamingQPS    = PeakConcurrentUsers                   // 100K concurrent streams
    
    // Video quality bitrates (Mbps)
    BitrateSD  = 1.0  // 480p
    BitrateHD  = 3.0  // 720p
    BitrateFHD = 6.0  // 1080p
    Bitrate4K  = 25.0 // 4K
)

// Storage calculations
type StorageCalculation struct {
    VideoSizes      map[string]int64 // Quality -> Size per hour
    TotalContent    int64            // Total storage needed
    CDNStorage      int64            // CDN cache storage
    BackupStorage   int64            // Backup and redundancy
}

func CalculateStorage() StorageCalculation {
    calc := StorageCalculation{
        VideoSizes: map[string]int64{
            "480p":  450 * 1024 * 1024,  // 450MB per hour
            "720p":  1350 * 1024 * 1024, // 1.35GB per hour
            "1080p": 2700 * 1024 * 1024, // 2.7GB per hour
            "4K":    11250 * 1024 * 1024, // 11.25GB per hour
        },
    }
    
    // Average movie length: 2 hours, TV episode: 45 minutes
    avgContentLength := 1.5 // hours
    totalTitles := 15000
    
    // Calculate total storage for all qualities
    var totalPerTitle int64
    for _, size := range calc.VideoSizes {
        totalPerTitle += int64(float64(size) * avgContentLength)
    }
    
    calc.TotalContent = totalPerTitle * int64(totalTitles) // ~50 PB
    calc.CDNStorage = calc.TotalContent / 10               // 10% cached at edge
    calc.BackupStorage = calc.TotalContent * 2             // 2x for redundancy
    
    return calc
}

// Bandwidth calculations
func CalculateBandwidth() map[string]float64 {
    // Assume quality distribution during peak
    qualityDistribution := map[string]float64{
        "480p":  0.20, // 20%
        "720p":  0.40, // 40%
        "1080p": 0.30, // 30%
        "4K":    0.10, // 10%
    }
    
    bitrates := map[string]float64{
        "480p":  1.0,
        "720p":  3.0,
        "1080p": 6.0,
        "4K":    25.0,
    }
    
    totalBandwidth := 0.0
    for quality, percentage := range qualityDistribution {
        streams := float64(PeakConcurrentUsers) * percentage
        bandwidth := streams * bitrates[quality]
        totalBandwidth += bandwidth
    }
    
    return map[string]float64{
        "peak_bandwidth_mbps": totalBandwidth,
        "peak_bandwidth_gbps": totalBandwidth / 1000,
        "peak_bandwidth_tbps": totalBandwidth / 1000000,
    }
}
```

### Infrastructure Estimation
```go
type InfrastructureNeeds struct {
    OriginServers    int // Master video storage
    EdgeServers      int // CDN edge locations
    EncodingServers  int // Video processing
    DatabaseServers  int // Metadata and user data
    RecommendationServers int // ML recommendation system
}

func EstimateInfrastructure() InfrastructureNeeds {
    storage := CalculateStorage()
    bandwidth := CalculateBandwidth()
    
    // Origin servers (each handles 10TB)
    originServers := int(storage.TotalContent / (10 * 1024 * 1024 * 1024 * 1024))
    
    // Edge servers (global CDN)
    edgeServers := 1000 // Major cities worldwide
    
    // Encoding servers (each processes 10 concurrent videos)
    encodingServers := 50 // For new content processing
    
    // Database servers (sharded by user_id and content_id)
    databaseServers := 100
    
    // Recommendation servers (ML workload)
    recommendationServers := 20
    
    return InfrastructureNeeds{
        OriginServers:         originServers,
        EdgeServers:          edgeServers,
        EncodingServers:      encodingServers,
        DatabaseServers:      databaseServers,
        RecommendationServers: recommendationServers,
    }
}
```

## 🏗️ Step 3: High-Level Design

### System Architecture
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ Mobile Apps │  │  Smart TVs  │  │ Web Clients │  │   Gaming    │
│ (iOS/Android│  │ (Roku/Apple │  │ (React/Vue) │  │ Consoles    │
│  React Nat.)│  │  TV/Samsung)│  │             │  │ (PS/Xbox)   │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │                │
       └────────────────┼────────────────┼────────────────┘
                        │                │
              ┌─────────▼────────────────▼─────────┐
              │      Global Load Balancer         │
              │    (Geographic Distribution)      │
              └─────────┬─────────────────────────┘
                        │
              ┌─────────▼─────────┐
              │    API Gateway    │
              │  - Authentication │
              │  - Rate Limiting  │
              │  - Request Routing│
              └─────────┬─────────┘
                        │
    ┌───────────────────┼───────────────────┐
    │                   │                   │
┌───▼────┐  ┌──────▼──────┐  ┌──────▼──────┐
│ User   │  │  Content    │  │ Streaming   │
│Service │  │  Service    │  │  Service    │
│- Auth  │  │- Metadata   │  │- Video      │
│- Profile│  │- Catalog    │  │  Delivery   │
│- Billing│  │- Search     │  │- Quality    │
└───┬────┘  └──────┬──────┘  └──────┬──────┘
    │              │                │
    │      ┌───────▼───────┐        │
    │      │Recommendation │        │
    │      │   Service     │        │
    │      │- ML Models    │        │
    │      │- Personalize  │        │
    │      └───────┬───────┘        │
    │              │                │
    │      ┌───────▼───────┐        │
    │      │  Analytics    │        │
    │      │   Service     │        │
    │      │- View Tracking│        │
    │      │- A/B Testing  │        │
    │      └───────┬───────┘        │
    │              │                │
┌───▼──────────────▼────────────────▼───┐
│           Data & Storage Layer        │
├─────────┬─────────┬─────────┬─────────┤
│User DB  │Content  │Analytics│ Video   │
│(MySQL)  │Metadata │Data     │Storage  │
│         │(Cassand)│(Hadoop) │(S3/CDN) │
└─────────┴─────────┴─────────┴─────────┘
```

### Video Streaming Data Flow
```
┌─────────────────────────────────────────────────────────────┐
│                Video Streaming Flow                         │
└─────────────────────────────────────────────────────────────┘

User ──1──► API Gateway ──2──► Streaming Service ──3──► CDN
                                      │
                                      4
                                      ▼
                               Content Service ──5──► Metadata DB
                                      │
                                      6
                                      ▼
                               Analytics ──7──► View Tracking
                                Service

1. User requests video stream
2. Authentication & authorization
3. Get optimal CDN endpoint
4. Fetch video metadata
5. Get content information
6. Log viewing analytics
7. Track user behavior

┌─────────────────────────────────────────────────────────────┐
│              Content Upload & Processing Flow               │
└─────────────────────────────────────────────────────────────┘

Content ──1──► Upload Service ──2──► Video Processing ──3──► Transcoding
Creator                                   Pipeline              Farm
                                             │
                                             4
                                             ▼
                                      Quality Control ──5──► CDN Distribution
                                             │
                                             6
                                             ▼
                                      Metadata ──7──► Content Catalog
                                      Extraction

1. Content creator uploads video
2. Initiate processing pipeline
3. Parallel transcoding (multiple qualities)
4. Quality assurance checks
5. Distribute to global CDN
6. Extract metadata & thumbnails
7. Update content catalog
```

### Recommendation System Flow
```
┌─────────────────────────────────────────────────────────────┐
│                Recommendation Generation Flow               │
└─────────────────────────────────────────────────────────────┘

User ──1──► Recommendation ──2──► User Profile ──3──► ML Models
           Service                Service
              │                      │
              4                      │
              ▼                      │
       Content Filtering ──5─────────┘
              │
              6
              ▼
       Ranking Engine ──7──► Personalized ──8──► User Interface
                             Recommendations

1. User requests recommendations
2. Get user viewing history & preferences
3. Apply collaborative filtering & content-based models
4. Filter available content
5. Apply business rules & content policies
6. Rank recommendations by predicted engagement
7. Generate personalized list
8. Display to user with A/B testing
```

### Core Services
```go
// User Service - manages user accounts, profiles, subscriptions
type UserService struct {
    userRepo        UserRepository
    profileRepo     ProfileRepository
    subscriptionRepo SubscriptionRepository
    authService     AuthService
    cache           Cache
}

// Content Service - manages video metadata, catalogs
type ContentService struct {
    contentRepo     ContentRepository
    metadataRepo    MetadataRepository
    licenseService  LicenseService
    cache           Cache
}

// Streaming Service - handles video delivery and playback
type StreamingService struct {
    videoRepo       VideoRepository
    cdnService      CDNService
    qualitySelector QualitySelector
    analytics       AnalyticsService
}

// Recommendation Service - ML-based content recommendations
type RecommendationService struct {
    mlModels        MLModelService
    userBehavior    UserBehaviorService
    contentFeatures ContentFeatureService
    cache           Cache
}

// Video Processing Service - encodes and processes video content
type VideoProcessingService struct {
    encoder         VideoEncoder
    thumbnailGen    ThumbnailGenerator
    qualityGen      QualityGenerator
    storageService  StorageService
}
```

## 🗄️ Step 4: Database Design

### Data Models
```go
// User and Profile models
type User struct {
    ID              int64     `json:"id" db:"id"`
    Email           string    `json:"email" db:"email"`
    PasswordHash    string    `json:"-" db:"password_hash"`
    SubscriptionPlan string   `json:"subscription_plan" db:"subscription_plan"`
    Country         string    `json:"country" db:"country"`
    Language        string    `json:"language" db:"language"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Profile struct {
    ID          int64  `json:"id" db:"id"`
    UserID      int64  `json:"user_id" db:"user_id"`
    Name        string `json:"name" db:"name"`
    AvatarURL   string `json:"avatar_url" db:"avatar_url"`
    IsKids      bool   `json:"is_kids" db:"is_kids"`
    Language    string `json:"language" db:"language"`
    Preferences map[string]interface{} `json:"preferences" db:"preferences"`
}

// Content models
type Content struct {
    ID              int64     `json:"id" db:"id"`
    Title           string    `json:"title" db:"title"`
    Description     string    `json:"description" db:"description"`
    Type            string    `json:"type" db:"type"` // movie, series, documentary
    Genre           []string  `json:"genre" db:"genre"`
    ReleaseYear     int       `json:"release_year" db:"release_year"`
    Rating          string    `json:"rating" db:"rating"` // PG, PG-13, R, etc.
    Duration        int       `json:"duration" db:"duration"` // minutes
    ThumbnailURL    string    `json:"thumbnail_url" db:"thumbnail_url"`
    TrailerURL      string    `json:"trailer_url" db:"trailer_url"`
    AvailableRegions []string `json:"available_regions" db:"available_regions"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Episode struct {
    ID          int64  `json:"id" db:"id"`
    SeriesID    int64  `json:"series_id" db:"series_id"`
    SeasonNum   int    `json:"season_num" db:"season_num"`
    EpisodeNum  int    `json:"episode_num" db:"episode_num"`
    Title       string `json:"title" db:"title"`
    Description string `json:"description" db:"description"`
    Duration    int    `json:"duration" db:"duration"`
    ThumbnailURL string `json:"thumbnail_url" db:"thumbnail_url"`
}

// Video storage models
type VideoFile struct {
    ID          int64  `json:"id" db:"id"`
    ContentID   int64  `json:"content_id" db:"content_id"`
    EpisodeID   *int64 `json:"episode_id,omitempty" db:"episode_id"`
    Quality     string `json:"quality" db:"quality"` // 480p, 720p, 1080p, 4K
    Format      string `json:"format" db:"format"`   // mp4, webm
    FileSize    int64  `json:"file_size" db:"file_size"`
    Duration    int    `json:"duration" db:"duration"`
    Bitrate     int    `json:"bitrate" db:"bitrate"`
    StoragePath string `json:"storage_path" db:"storage_path"`
    CDNPath     string `json:"cdn_path" db:"cdn_path"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// User behavior models
type WatchHistory struct {
    ID          int64     `json:"id" db:"id"`
    UserID      int64     `json:"user_id" db:"user_id"`
    ProfileID   int64     `json:"profile_id" db:"profile_id"`
    ContentID   int64     `json:"content_id" db:"content_id"`
    EpisodeID   *int64    `json:"episode_id,omitempty" db:"episode_id"`
    WatchedAt   time.Time `json:"watched_at" db:"watched_at"`
    Duration    int       `json:"duration" db:"duration"` // seconds watched
    Progress    float64   `json:"progress" db:"progress"` // percentage completed
    Quality     string    `json:"quality" db:"quality"`
    DeviceType  string    `json:"device_type" db:"device_type"`
}

type Rating struct {
    UserID      int64     `json:"user_id" db:"user_id"`
    ProfileID   int64     `json:"profile_id" db:"profile_id"`
    ContentID   int64     `json:"content_id" db:"content_id"`
    Rating      int       `json:"rating" db:"rating"` // 1-5 stars
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### Database Schema
```sql
-- Users and Profiles (PostgreSQL)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    subscription_plan VARCHAR(50) NOT NULL,
    country VARCHAR(2) NOT NULL,
    language VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    is_kids BOOLEAN DEFAULT FALSE,
    language VARCHAR(10),
    preferences JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Content Catalog (PostgreSQL)
CREATE TABLE content (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL, -- movie, series, documentary
    genre TEXT[], -- Array of genres
    release_year INTEGER,
    rating VARCHAR(10),
    duration INTEGER, -- minutes
    thumbnail_url TEXT,
    trailer_url TEXT,
    available_regions TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_content_type_genre ON content(type, genre);
CREATE INDEX idx_content_release_year ON content(release_year);

-- Episodes for series
CREATE TABLE episodes (
    id BIGSERIAL PRIMARY KEY,
    series_id BIGINT REFERENCES content(id),
    season_num INTEGER NOT NULL,
    episode_num INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER,
    thumbnail_url TEXT,
    UNIQUE(series_id, season_num, episode_num)
);

-- Video Files (Cassandra for high read volume)
CREATE TABLE video_files (
    id BIGINT PRIMARY KEY,
    content_id BIGINT,
    episode_id BIGINT,
    quality TEXT, -- 480p, 720p, 1080p, 4K
    format TEXT,  -- mp4, webm
    file_size BIGINT,
    duration INTEGER,
    bitrate INTEGER,
    storage_path TEXT,
    cdn_path TEXT,
    created_at TIMESTAMP
);

-- User Behavior (Cassandra for high write volume)
CREATE TABLE watch_history (
    user_id BIGINT,
    profile_id BIGINT,
    watched_at TIMESTAMP,
    content_id BIGINT,
    episode_id BIGINT,
    duration INTEGER,
    progress DOUBLE,
    quality TEXT,
    device_type TEXT,
    PRIMARY KEY (user_id, watched_at)
) WITH CLUSTERING ORDER BY (watched_at DESC);

CREATE TABLE ratings (
    user_id BIGINT,
    profile_id BIGINT,
    content_id BIGINT,
    rating INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, content_id)
);
```

### Database Sharding Strategy
```go
type NetflixDatabaseSharding struct {
    userShards    []Database // Shard by user_id
    contentShards []Database // Shard by content_id
    behaviorShards []Database // Shard by user_id for watch history
}

// User data sharding
func (nds *NetflixDatabaseSharding) GetUserShard(userID int64) Database {
    shardIndex := userID % int64(len(nds.userShards))
    return nds.userShards[shardIndex]
}

// Content sharding - by content_id for even distribution
func (nds *NetflixDatabaseSharding) GetContentShard(contentID int64) Database {
    shardIndex := contentID % int64(len(nds.contentShards))
    return nds.contentShards[shardIndex]
}

// Behavior data sharding - by user_id for user-specific queries
func (nds *NetflixDatabaseSharding) GetBehaviorShard(userID int64) Database {
    shardIndex := userID % int64(len(nds.behaviorShards))
    return nds.behaviorShards[shardIndex]
}
```

## 🎥 Step 5: Video Processing Pipeline

### Video Encoding Service
```go
type VideoEncodingService struct {
    encoder       VideoEncoder
    storage       StorageService
    cdnService    CDNService
    eventBus      EventBus
    qualityConfig QualityConfiguration
}

type QualityConfiguration struct {
    Profiles []EncodingProfile
}

type EncodingProfile struct {
    Quality    string // 480p, 720p, 1080p, 4K
    Resolution string // 854x480, 1280x720, 1920x1080, 3840x2160
    Bitrate    int    // kbps
    Format     string // mp4, webm
    Codec      string // H.264, H.265, VP9, AV1
}

var DefaultQualityConfig = QualityConfiguration{
    Profiles: []EncodingProfile{
        {Quality: "480p", Resolution: "854x480", Bitrate: 1000, Format: "mp4", Codec: "H.264"},
        {Quality: "720p", Resolution: "1280x720", Bitrate: 3000, Format: "mp4", Codec: "H.264"},
        {Quality: "1080p", Resolution: "1920x1080", Bitrate: 6000, Format: "mp4", Codec: "H.264"},
        {Quality: "4K", Resolution: "3840x2160", Bitrate: 25000, Format: "mp4", Codec: "H.265"},
    },
}

func (ves *VideoEncodingService) ProcessVideo(ctx context.Context, videoID int64, sourceFile string) error {
    // Create encoding jobs for all quality profiles
    var encodingJobs []EncodingJob
    
    for _, profile := range ves.qualityConfig.Profiles {
        job := EncodingJob{
            VideoID:    videoID,
            SourceFile: sourceFile,
            Profile:    profile,
            Status:     "pending",
            CreatedAt:  time.Now(),
        }
        encodingJobs = append(encodingJobs, job)
    }
    
    // Process jobs in parallel
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 4) // Limit concurrent encoding
    
    for _, job := range encodingJobs {
        wg.Add(1)
        go func(j EncodingJob) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release
            
            if err := ves.encodeVideo(ctx, j); err != nil {
                log.Printf("Encoding failed for job %+v: %v", j, err)
            }
        }(job)
    }
    
    wg.Wait()
    
    // Generate thumbnails
    if err := ves.generateThumbnails(ctx, videoID, sourceFile); err != nil {
        log.Printf("Thumbnail generation failed: %v", err)
    }
    
    // Publish completion event
    ves.eventBus.Publish(ctx, VideoProcessingCompletedEvent{
        VideoID:   videoID,
        Qualities: ves.getCompletedQualities(videoID),
    })
    
    return nil
}

func (ves *VideoEncodingService) encodeVideo(ctx context.Context, job EncodingJob) error {
    outputFile := fmt.Sprintf("%s_%s.%s", job.VideoID, job.Profile.Quality, job.Profile.Format)
    
    // FFmpeg encoding command
    cmd := exec.CommandContext(ctx, "ffmpeg",
        "-i", job.SourceFile,
        "-c:v", ves.getVideoCodec(job.Profile.Codec),
        "-b:v", fmt.Sprintf("%dk", job.Profile.Bitrate),
        "-s", job.Profile.Resolution,
        "-c:a", "aac",
        "-b:a", "128k",
        "-f", job.Profile.Format,
        outputFile,
    )
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("encoding failed: %w", err)
    }
    
    // Upload to storage
    storageURL, err := ves.storage.Upload(ctx, outputFile)
    if err != nil {
        return fmt.Errorf("upload failed: %w", err)
    }
    
    // Update CDN
    cdnURL, err := ves.cdnService.Deploy(ctx, storageURL, job.Profile.Quality)
    if err != nil {
        return fmt.Errorf("CDN deployment failed: %w", err)
    }
    
    // Save video file metadata
    videoFile := VideoFile{
        ContentID:   job.VideoID,
        Quality:     job.Profile.Quality,
        Format:      job.Profile.Format,
        Bitrate:     job.Profile.Bitrate,
        StoragePath: storageURL,
        CDNPath:     cdnURL,
        CreatedAt:   time.Now(),
    }
    
    return ves.saveVideoFile(ctx, videoFile)
}
```

## 🌐 Step 6: Global CDN Architecture

### CDN Design for Video Streaming
```go
type NetflixCDN struct {
    originServers []OriginServer
    edgeServers   map[string][]EdgeServer // region -> servers
    popServers    map[string][]POPServer  // ISP -> servers
    loadBalancer  GlobalLoadBalancer
}

type OriginServer struct {
    ID       string
    Region   string
    Capacity int64 // TB
    Storage  StorageService
}

type EdgeServer struct {
    ID           string
    Location     string
    ISP          string
    Capacity     int64
    Cache        VideoCache
    HealthStatus string
}

type POPServer struct {
    ID       string
    ISP      string
    City     string
    Capacity int64
    Cache    VideoCache
}

// Global load balancing based on user location and content popularity
func (ncdn *NetflixCDN) GetOptimalServer(userLocation, contentID string, quality string) (*EdgeServer, error) {
    // 1. Find nearest edge servers
    nearbyServers := ncdn.getNearbyServers(userLocation)

    // 2. Check which servers have the content cached
    serversWithContent := ncdn.filterServersWithContent(nearbyServers, contentID, quality)

    if len(serversWithContent) > 0 {
        // 3. Select server with lowest load
        return ncdn.selectLeastLoadedServer(serversWithContent), nil
    }

    // 4. If no edge server has content, select best server and trigger cache
    bestServer := ncdn.selectBestServer(nearbyServers)
    go ncdn.cacheContent(bestServer, contentID, quality)

    return bestServer, nil
}

func (ncdn *NetflixCDN) cacheContent(server *EdgeServer, contentID, quality string) error {
    // Get content from origin or parent cache
    originServer := ncdn.getOriginServer(contentID)

    // Stream content to edge server
    return server.Cache.StreamFromOrigin(originServer, contentID, quality)
}

// Predictive caching based on viewing patterns
type PredictiveCaching struct {
    mlModel     MLModel
    viewingData ViewingAnalytics
    cachePolicy CachePolicy
}

func (pc *PredictiveCaching) PredictPopularContent(region string, timeWindow time.Duration) []ContentPrediction {
    // Analyze viewing patterns
    patterns := pc.viewingData.GetRegionalPatterns(region, timeWindow)

    // Use ML model to predict popular content
    predictions := pc.mlModel.PredictPopularity(patterns)

    return predictions
}

func (pc *PredictiveCaching) PreCacheContent(server *EdgeServer, predictions []ContentPrediction) {
    for _, prediction := range predictions {
        if prediction.Confidence > 0.8 { // High confidence threshold
            go server.Cache.PreCache(prediction.ContentID, prediction.Quality)
        }
    }
}
```

### Adaptive Bitrate Streaming (ABR)
```go
type AdaptiveBitrateService struct {
    qualitySelector QualitySelector
    bandwidthMonitor BandwidthMonitor
    bufferMonitor   BufferMonitor
}

type StreamingSession struct {
    UserID        int64
    ContentID     int64
    CurrentQuality string
    BufferHealth  float64 // seconds of content buffered
    Bandwidth     float64 // Mbps
    DeviceType    string
    StartTime     time.Time
}

func (abs *AdaptiveBitrateService) SelectOptimalQuality(session *StreamingSession) string {
    // Factors for quality selection
    factors := QualityFactors{
        AvailableBandwidth: session.Bandwidth,
        BufferHealth:      session.BufferHealth,
        DeviceCapability:  abs.getDeviceCapability(session.DeviceType),
        UserPreference:    abs.getUserPreference(session.UserID),
        NetworkStability:  abs.getNetworkStability(session.UserID),
    }

    return abs.qualitySelector.SelectQuality(factors)
}

func (abs *AdaptiveBitrateService) HandleQualityChange(session *StreamingSession, newQuality string) error {
    // Smooth quality transitions
    if abs.shouldUseSeamlessTransition(session.CurrentQuality, newQuality) {
        return abs.performSeamlessTransition(session, newQuality)
    }

    // Immediate quality change for significant bandwidth changes
    return abs.performImmediateTransition(session, newQuality)
}

// Quality selection algorithm
type QualitySelector struct {
    qualityLadder []QualityLevel
}

type QualityLevel struct {
    Quality         string
    MinBandwidth    float64 // Mbps
    RecommendedBandwidth float64
    Resolution      string
    Bitrate         int
}

func (qs *QualitySelector) SelectQuality(factors QualityFactors) string {
    // Start with bandwidth-based selection
    availableQualities := qs.filterByBandwidth(factors.AvailableBandwidth)

    // Consider buffer health
    if factors.BufferHealth < 5.0 { // Less than 5 seconds buffered
        // Select lower quality to improve buffer
        availableQualities = qs.filterForBufferRecovery(availableQualities)
    }

    // Consider device capability
    availableQualities = qs.filterByDeviceCapability(availableQualities, factors.DeviceCapability)

    // Apply user preference
    if factors.UserPreference == "data_saver" {
        return qs.selectDataSaverQuality(availableQualities)
    }

    // Select highest available quality
    return qs.selectHighestQuality(availableQualities)
}
```

## 🤖 Step 7: Recommendation System

### ML-Based Recommendation Engine
```go
type RecommendationEngine struct {
    collaborativeFilter CollaborativeFilteringModel
    contentBasedFilter  ContentBasedFilteringModel
    deepLearningModel   DeepLearningModel
    realTimeFeatures    RealTimeFeatureService
    cache              RecommendationCache
}

type UserFeatures struct {
    UserID          int64
    Age             int
    Country         string
    Language        string
    WatchHistory    []WatchEvent
    Ratings         []Rating
    GenrePreferences map[string]float64
    ActorPreferences map[string]float64
    DirectorPreferences map[string]float64
    WatchTimePatterns map[string]float64 // time_of_day -> preference
}

type ContentFeatures struct {
    ContentID   int64
    Genre       []string
    Actors      []string
    Directors   []string
    ReleaseYear int
    Duration    int
    Rating      string
    Popularity  float64
    TrendingScore float64
}

func (re *RecommendationEngine) GenerateRecommendations(ctx context.Context, userID int64, profileID int64) ([]Recommendation, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("recommendations:%d:%d", userID, profileID)
    if cached, found := re.cache.Get(cacheKey); found {
        return cached.([]Recommendation), nil
    }

    // Get user features
    userFeatures, err := re.getUserFeatures(ctx, userID, profileID)
    if err != nil {
        return nil, err
    }

    // Generate recommendations from different models
    var allRecommendations []Recommendation

    // 1. Collaborative Filtering (users with similar taste)
    cfRecommendations, err := re.collaborativeFilter.Recommend(ctx, userFeatures)
    if err == nil {
        allRecommendations = append(allRecommendations, cfRecommendations...)
    }

    // 2. Content-Based Filtering (similar content to what user liked)
    cbRecommendations, err := re.contentBasedFilter.Recommend(ctx, userFeatures)
    if err == nil {
        allRecommendations = append(allRecommendations, cbRecommendations...)
    }

    // 3. Deep Learning Model (neural collaborative filtering)
    dlRecommendations, err := re.deepLearningModel.Recommend(ctx, userFeatures)
    if err == nil {
        allRecommendations = append(allRecommendations, dlRecommendations...)
    }

    // 4. Trending and Popular Content
    trendingRecommendations := re.getTrendingContent(ctx, userFeatures.Country)
    allRecommendations = append(allRecommendations, trendingRecommendations...)

    // Combine and rank recommendations
    finalRecommendations := re.combineAndRankRecommendations(allRecommendations, userFeatures)

    // Cache results
    re.cache.Set(cacheKey, finalRecommendations, 1*time.Hour)

    return finalRecommendations, nil
}

// Real-time feature updates
func (re *RecommendationEngine) UpdateUserFeatures(ctx context.Context, event WatchEvent) {
    // Update user preferences based on viewing behavior
    userFeatures, err := re.getUserFeatures(ctx, event.UserID, event.ProfileID)
    if err != nil {
        return
    }

    // Update genre preferences
    content, _ := re.getContentFeatures(event.ContentID)
    for _, genre := range content.Genre {
        currentPref := userFeatures.GenrePreferences[genre]
        // Increase preference based on watch completion
        newPref := currentPref + (event.WatchPercentage * 0.1)
        userFeatures.GenrePreferences[genre] = math.Min(newPref, 1.0)
    }

    // Invalidate recommendation cache
    cacheKey := fmt.Sprintf("recommendations:%d:%d", event.UserID, event.ProfileID)
    re.cache.Delete(cacheKey)

    // Trigger async recommendation refresh
    go re.refreshRecommendations(ctx, event.UserID, event.ProfileID)
}
```

### Personalized Home Page
```go
type HomePageService struct {
    recommendationEngine RecommendationEngine
    contentService      ContentService
    trendingService     TrendingService
    personalizedRows    PersonalizedRowService
}

type HomePageLayout struct {
    Rows []ContentRow `json:"rows"`
}

type ContentRow struct {
    Title       string    `json:"title"`
    Type        string    `json:"type"` // trending, recommended, continue_watching, etc.
    Content     []Content `json:"content"`
    Priority    int       `json:"priority"`
    Personalized bool     `json:"personalized"`
}

func (hps *HomePageService) GenerateHomePage(ctx context.Context, userID, profileID int64) (*HomePageLayout, error) {
    var rows []ContentRow

    // 1. Continue Watching (highest priority)
    continueWatching, err := hps.getContinueWatching(ctx, userID, profileID)
    if err == nil && len(continueWatching) > 0 {
        rows = append(rows, ContentRow{
            Title:       "Continue Watching",
            Type:        "continue_watching",
            Content:     continueWatching,
            Priority:    1,
            Personalized: true,
        })
    }

    // 2. Personalized Recommendations
    recommendations, err := hps.recommendationEngine.GenerateRecommendations(ctx, userID, profileID)
    if err == nil {
        rows = append(rows, ContentRow{
            Title:       "Recommended for You",
            Type:        "recommended",
            Content:     hps.convertRecommendationsToContent(recommendations[:20]),
            Priority:    2,
            Personalized: true,
        })
    }

    // 3. Trending Now
    trending, err := hps.trendingService.GetTrending(ctx, hps.getUserCountry(userID))
    if err == nil {
        rows = append(rows, ContentRow{
            Title:       "Trending Now",
            Type:        "trending",
            Content:     trending,
            Priority:    3,
            Personalized: false,
        })
    }

    // 4. Genre-based rows (personalized)
    topGenres := hps.getUserTopGenres(ctx, userID, profileID)
    for i, genre := range topGenres {
        if i >= 3 { // Limit to top 3 genres
            break
        }

        genreContent, err := hps.contentService.GetContentByGenre(ctx, genre, 20)
        if err == nil {
            rows = append(rows, ContentRow{
                Title:       fmt.Sprintf("Popular %s", genre),
                Type:        "genre",
                Content:     genreContent,
                Priority:    4 + i,
                Personalized: true,
            })
        }
    }

    // 5. New Releases
    newReleases, err := hps.contentService.GetNewReleases(ctx, hps.getUserCountry(userID))
    if err == nil {
        rows = append(rows, ContentRow{
            Title:       "New Releases",
            Type:        "new_releases",
            Content:     newReleases,
            Priority:    10,
            Personalized: false,
        })
    }

    // Sort rows by priority
    sort.Slice(rows, func(i, j int) bool {
        return rows[i].Priority < rows[j].Priority
    })

    return &HomePageLayout{Rows: rows}, nil
}
```

## 📊 Step 8: Analytics & Monitoring

### Real-time Analytics Service
```go
type AnalyticsService struct {
    eventCollector  EventCollector
    streamProcessor StreamProcessor
    metricsStore    MetricsStore
    alertManager    AlertManager
}

type StreamingEvent struct {
    EventType   string                 `json:"event_type"`
    UserID      int64                  `json:"user_id"`
    ProfileID   int64                  `json:"profile_id"`
    ContentID   int64                  `json:"content_id"`
    SessionID   string                 `json:"session_id"`
    Timestamp   time.Time              `json:"timestamp"`
    Properties  map[string]interface{} `json:"properties"`
}

// Key streaming events
const (
    EventPlayStart    = "play_start"
    EventPlayPause    = "play_pause"
    EventPlayResume   = "play_resume"
    EventPlayStop     = "play_stop"
    EventQualityChange = "quality_change"
    EventBuffering    = "buffering"
    EventError        = "error"
    EventSeek         = "seek"
)

func (as *AnalyticsService) TrackStreamingEvent(ctx context.Context, event StreamingEvent) error {
    // Validate event
    if err := as.validateEvent(event); err != nil {
        return err
    }

    // Send to real-time processing
    if err := as.eventCollector.Collect(event); err != nil {
        return err
    }

    // Update real-time metrics
    as.updateRealTimeMetrics(event)

    // Check for alerts
    as.checkAlerts(event)

    return nil
}

func (as *AnalyticsService) updateRealTimeMetrics(event StreamingEvent) {
    switch event.EventType {
    case EventPlayStart:
        as.metricsStore.Increment("concurrent_streams")
        as.metricsStore.Increment("plays_started")

    case EventPlayStop:
        as.metricsStore.Decrement("concurrent_streams")

        // Calculate watch time
        if startTime, ok := event.Properties["start_time"].(time.Time); ok {
            watchTime := event.Timestamp.Sub(startTime)
            as.metricsStore.RecordDuration("watch_time", watchTime)
        }

    case EventBuffering:
        as.metricsStore.Increment("buffering_events")

        if duration, ok := event.Properties["duration"].(float64); ok {
            as.metricsStore.RecordValue("buffering_duration", duration)
        }

    case EventError:
        as.metricsStore.Increment("streaming_errors")

        if errorType, ok := event.Properties["error_type"].(string); ok {
            as.metricsStore.IncrementWithTags("errors_by_type", map[string]string{
                "error_type": errorType,
            })
        }
    }
}

// Quality of Experience (QoE) monitoring
type QoEMonitor struct {
    metricsCollector MetricsCollector
    alertThresholds  QoEThresholds
}

type QoEThresholds struct {
    MaxBufferingRatio    float64 // Max % of time spent buffering
    MaxStartupTime       time.Duration // Max time to start playback
    MinBitrateStability  float64 // Min % of time in optimal quality
    MaxErrorRate         float64 // Max error rate per session
}

func (qoe *QoEMonitor) CalculateQoEScore(sessionEvents []StreamingEvent) float64 {
    var (
        totalPlayTime    time.Duration
        totalBufferTime  time.Duration
        startupTime      time.Duration
        qualityChanges   int
        errors          int
        optimalQualityTime time.Duration
    )

    // Analyze session events
    for i, event := range sessionEvents {
        switch event.EventType {
        case EventPlayStart:
            if i == 0 { // First play event
                if startTime, ok := event.Properties["startup_time"].(time.Duration); ok {
                    startupTime = startTime
                }
            }

        case EventBuffering:
            if duration, ok := event.Properties["duration"].(time.Duration); ok {
                totalBufferTime += duration
            }

        case EventQualityChange:
            qualityChanges++

        case EventError:
            errors++
        }
    }

    // Calculate QoE components
    bufferingRatio := float64(totalBufferTime) / float64(totalPlayTime)
    startupScore := qoe.calculateStartupScore(startupTime)
    stabilityScore := qoe.calculateStabilityScore(qualityChanges, totalPlayTime)
    errorScore := qoe.calculateErrorScore(errors, len(sessionEvents))

    // Weighted QoE score (0-100)
    qoeScore := (startupScore*0.3 + (1-bufferingRatio)*0.4 + stabilityScore*0.2 + errorScore*0.1) * 100

    return math.Max(0, math.Min(100, qoeScore))
}
```

## 🔧 Step 9: API Implementation

### Streaming API Design
```go
type NetflixAPI struct {
    userService         UserService
    contentService      ContentService
    streamingService    StreamingService
    recommendationService RecommendationService
    analyticsService    AnalyticsService
    authMiddleware      AuthMiddleware
    rateLimiter        RateLimiter
}

func (api *NetflixAPI) SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    // Apply middleware
    authHandler := api.authMiddleware.Authenticate
    rateLimitHandler := api.rateLimiter.Middleware

    // Content discovery
    mux.HandleFunc("GET /api/v1/content/browse", rateLimitHandler(api.browseContent))
    mux.HandleFunc("GET /api/v1/content/search", rateLimitHandler(api.searchContent))
    mux.HandleFunc("GET /api/v1/content/{id}", rateLimitHandler(api.getContent))
    mux.HandleFunc("GET /api/v1/content/{id}/episodes", rateLimitHandler(api.getEpisodes))

    // Personalization
    mux.HandleFunc("GET /api/v1/recommendations", authHandler(rateLimitHandler(api.getRecommendations)))
    mux.HandleFunc("GET /api/v1/homepage", authHandler(rateLimitHandler(api.getHomePage)))

    // Streaming
    mux.HandleFunc("POST /api/v1/streaming/start", authHandler(rateLimitHandler(api.startStreaming)))
    mux.HandleFunc("GET /api/v1/streaming/{sessionId}/manifest", authHandler(api.getStreamingManifest))
    mux.HandleFunc("POST /api/v1/streaming/{sessionId}/events", authHandler(api.trackStreamingEvent))

    // User profiles
    mux.HandleFunc("GET /api/v1/profiles", authHandler(rateLimitHandler(api.getProfiles)))
    mux.HandleFunc("POST /api/v1/profiles", authHandler(rateLimitHandler(api.createProfile)))
    mux.HandleFunc("GET /api/v1/profiles/{id}/history", authHandler(rateLimitHandler(api.getWatchHistory)))

    return mux
}

func (api *NetflixAPI) startStreaming(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())

    var req StartStreamingRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    // Validate content access
    hasAccess, err := api.contentService.CheckUserAccess(r.Context(), userID, req.ContentID)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to check access")
        return
    }

    if !hasAccess {
        writeErrorResponse(w, http.StatusForbidden, "Content not available in your region")
        return
    }

    // Start streaming session
    session, err := api.streamingService.StartSession(r.Context(), StartSessionParams{
        UserID:     userID,
        ProfileID:  req.ProfileID,
        ContentID:  req.ContentID,
        EpisodeID:  req.EpisodeID,
        DeviceType: req.DeviceType,
        Quality:    req.PreferredQuality,
    })
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to start streaming")
        return
    }

    // Track analytics
    api.analyticsService.TrackStreamingEvent(r.Context(), StreamingEvent{
        EventType: EventPlayStart,
        UserID:    userID,
        ProfileID: req.ProfileID,
        ContentID: req.ContentID,
        SessionID: session.ID,
        Timestamp: time.Now(),
        Properties: map[string]interface{}{
            "device_type": req.DeviceType,
            "quality":     req.PreferredQuality,
        },
    })

    response := StartStreamingResponse{
        SessionID:   session.ID,
        ManifestURL: session.ManifestURL,
        ExpiresAt:   session.ExpiresAt,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (api *NetflixAPI) getStreamingManifest(w http.ResponseWriter, r *http.Request) {
    sessionID := r.PathValue("sessionId")
    userID := getUserIDFromContext(r.Context())

    // Validate session
    session, err := api.streamingService.GetSession(r.Context(), sessionID)
    if err != nil {
        writeErrorResponse(w, http.StatusNotFound, "Session not found")
        return
    }

    if session.UserID != userID {
        writeErrorResponse(w, http.StatusForbidden, "Access denied")
        return
    }

    // Generate adaptive streaming manifest (HLS/DASH)
    manifest, err := api.streamingService.GenerateManifest(r.Context(), session)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to generate manifest")
        return
    }

    w.Header().Set("Content-Type", "application/vnd.apple.mpegurl") // HLS
    w.Write([]byte(manifest))
}
```

## 🎯 Summary: Netflix System Design

### Architecture Highlights
- **Microservices**: User, Content, Streaming, Recommendation, Analytics services
- **Global CDN**: Multi-tier caching with predictive content placement
- **Video Processing**: Parallel encoding pipeline for multiple qualities
- **Adaptive Streaming**: Real-time quality adjustment based on network conditions
- **ML Recommendations**: Collaborative filtering + content-based + deep learning
- **Real-time Analytics**: Stream processing for QoE monitoring

### Scalability Features
- **200M users**: Sharded databases and distributed services
- **1B hours/day**: Global CDN with edge caching
- **100K concurrent streams**: Adaptive bitrate and load balancing
- **50 PB storage**: Distributed storage with intelligent caching
- **Global reach**: Regional content delivery and localization

### Performance Optimizations
- **Video start time**: < 2 seconds with predictive caching
- **Buffering**: < 0.1% with adaptive bitrate streaming
- **Recommendations**: < 100ms with multi-level caching
- **Search**: < 200ms with Elasticsearch
- **99.9% availability**: Multi-region deployment with failover

This design handles Netflix-scale traffic with optimal user experience worldwide! 🚀
