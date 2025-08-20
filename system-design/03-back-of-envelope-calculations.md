# 🧮 Back-of-Envelope Calculations

## 🎯 Why This Matters
Back-of-envelope calculations are crucial for:
- Estimating system capacity requirements
- Making informed architectural decisions
- Demonstrating quantitative thinking in interviews
- Validating if your design can handle the expected load

## 📊 Essential Numbers to Memorize

### Latency Numbers (2024 Updated)
```
L1 cache reference                           0.5 ns
Branch mispredict                            5   ns
L2 cache reference                           7   ns
Mutex lock/unlock                           25   ns
Main memory reference                      100   ns
Compress 1K bytes with Zippy             3,000   ns  =   3 μs
Send 1K bytes over 1 Gbps network       10,000   ns  =  10 μs
Read 4K randomly from SSD              150,000   ns  = 150 μs
Read 1 MB sequentially from memory      250,000   ns  = 250 μs
Round trip within same datacenter       500,000   ns  = 500 μs
Read 1 MB sequentially from SSD      1,000,000   ns  =   1 ms
Disk seek                           10,000,000   ns  =  10 ms
Read 1 MB sequentially from disk    20,000,000   ns  =  20 ms
Send packet CA→Netherlands→CA      150,000,000   ns  = 150 ms
```

### Storage and Bandwidth
```
1 byte = 8 bits
1 KB = 1,000 bytes (10^3)
1 MB = 1,000 KB = 1,000,000 bytes (10^6)
1 GB = 1,000 MB = 1,000,000,000 bytes (10^9)
1 TB = 1,000 GB = 1,000,000,000,000 bytes (10^12)
1 PB = 1,000 TB = 1,000,000,000,000,000 bytes (10^15)

ASCII character = 1 byte
Unicode character = 2-4 bytes (UTF-8)
Integer (32-bit) = 4 bytes
Integer (64-bit) = 8 bytes
UUID = 16 bytes
```

### Time Conversions
```
1 second = 1,000 milliseconds (ms)
1 second = 1,000,000 microseconds (μs)
1 second = 1,000,000,000 nanoseconds (ns)

1 day = 86,400 seconds ≈ 10^5 seconds
1 month = 30 days ≈ 2.6 * 10^6 seconds
1 year = 365 days ≈ 3.2 * 10^7 seconds
```

## 🔢 Calculation Framework

### Step 1: Define the Problem
```
Example: Design a URL shortener like TinyURL
- 100 million URLs shortened per day
- 10:1 read/write ratio
- 5 years of data retention
- 100 characters average URL length
```

### Step 2: Calculate QPS (Queries Per Second)

#### Write QPS
```
URLs per day = 100 million = 10^8
Seconds per day = 24 * 60 * 60 = 86,400 ≈ 10^5

Write QPS = 10^8 / 10^5 = 10^3 = 1,000 QPS
Peak Write QPS = 2 * 1,000 = 2,000 QPS (2x safety factor)
```

#### Read QPS
```
Read/Write ratio = 10:1
Read QPS = 10 * 1,000 = 10,000 QPS
Peak Read QPS = 2 * 10,000 = 20,000 QPS
```

### Step 3: Storage Estimation

#### Data Size per Entry
```go
type URLEntry struct {
    ShortURL    string    // 7 characters = 7 bytes
    LongURL     string    // 100 characters average = 100 bytes
    UserID      int64     // 8 bytes
    CreatedAt   time.Time // 8 bytes
    ExpiresAt   time.Time // 8 bytes
    ClickCount  int32     // 4 bytes
}
// Total per entry ≈ 135 bytes ≈ 150 bytes (with overhead)
```

#### Total Storage Calculation
```
Daily storage = 10^8 URLs * 150 bytes = 1.5 * 10^10 bytes = 15 GB/day
Monthly storage = 15 GB * 30 = 450 GB/month
Yearly storage = 15 GB * 365 = 5.5 TB/year
5-year storage = 5.5 TB * 5 = 27.5 TB

With replication (3x) = 27.5 TB * 3 = 82.5 TB
With indexing overhead (2x) = 82.5 TB * 2 = 165 TB
```

### Step 4: Bandwidth Estimation

#### Write Bandwidth
```
Write QPS = 1,000
Data per write = 150 bytes
Write bandwidth = 1,000 * 150 bytes = 150 KB/s
```

#### Read Bandwidth
```
Read QPS = 10,000
Data per read = 150 bytes (assuming full URL data)
Read bandwidth = 10,000 * 150 bytes = 1.5 MB/s
```

### Step 5: Memory/Cache Estimation

#### Cache Size (80/20 Rule)
```
Daily reads = 10,000 QPS * 86,400 seconds = 8.64 * 10^8 reads
Cache 20% of daily reads = 8.64 * 10^8 * 0.2 = 1.73 * 10^8 entries
Cache memory = 1.73 * 10^8 * 150 bytes = 26 GB

With cache overhead and metadata = 26 GB * 1.5 = 39 GB ≈ 40 GB
```

## 🏗️ Server Estimation

### Application Servers
```go
// Assumptions for Go servers (based on your microservices architecture)
type ServerCapacity struct {
    CPU_Cores    int     // 8 cores
    RAM_GB       int     // 32 GB
    QPS_Per_Core int     // 1000 QPS per core for simple operations
    Max_QPS      int     // 8000 QPS per server
}

// Calculate servers needed
Total_QPS := 20000 // Peak read QPS
QPS_Per_Server := 8000
Servers_Needed := ceil(20000 / 8000) = 3 servers

// Add redundancy and load distribution
Total_Servers := 3 * 2 = 6 servers (2x for redundancy)
```

### Database Servers

#### Read Replicas
```
Read QPS = 20,000
QPS per DB server = 5,000 (conservative for complex queries)
Read replicas needed = ceil(20,000 / 5,000) = 4 replicas
```

#### Write Master
```
Write QPS = 2,000
Single master can typically handle 5,000-10,000 writes/sec
1 master is sufficient
```

## 📈 Real-World Examples

### Example 1: Instagram-like Photo Sharing

#### Requirements
```
- 500M users, 1M daily active users
- Each user uploads 2 photos per day on average
- Each user views 20 photos per day
- Average photo size: 200KB
- Metadata per photo: 1KB
```

#### Calculations
```go
// Daily uploads
daily_uploads := 1_000_000 * 2 = 2_000_000 photos/day
upload_qps := 2_000_000 / 86_400 = 23 QPS

// Daily views  
daily_views := 1_000_000 * 20 = 20_000_000 views/day
view_qps := 20_000_000 / 86_400 = 231 QPS

// Storage
photo_storage_daily := 2_000_000 * 200_000 bytes = 400 GB/day
metadata_storage_daily := 2_000_000 * 1_000 bytes = 2 GB/day
total_daily_storage := 402 GB/day

// 5-year storage
total_storage := 402 GB * 365 * 5 = 733 TB
with_replication := 733 TB * 3 = 2.2 PB

// Bandwidth
upload_bandwidth := 23 QPS * 200 KB = 4.6 MB/s
view_bandwidth := 231 QPS * 200 KB = 46.2 MB/s
```

### Example 2: Chat System like WhatsApp

#### Requirements
```
- 2B users, 500M daily active users
- Each user sends 40 messages per day
- Each message: 100 bytes average
- Message retention: 1 year
```

#### Calculations
```go
// Message volume
daily_messages := 500_000_000 * 40 = 20_000_000_000 messages/day
message_qps := 20_000_000_000 / 86_400 = 231_481 QPS

// Storage
message_size := 100 // bytes
daily_storage := 20_000_000_000 * 100 = 2 TB/day
yearly_storage := 2 TB * 365 = 730 TB
with_replication := 730 TB * 3 = 2.2 PB

// Bandwidth
message_bandwidth := 231_481 * 100 bytes = 23.1 MB/s

// Servers needed (assuming 10,000 QPS per server)
servers_needed := ceil(231_481 / 10_000) = 24 servers
with_redundancy := 24 * 2 = 48 servers
```

## 🎯 Optimization Techniques

### Compression
```
Text compression ratio: 3:1 to 10:1
Image compression: 5:1 to 20:1
Video compression: 100:1 to 1000:1

Example: 100 byte message compressed to 30 bytes
Storage savings = 70%
Bandwidth savings = 70%
```

### Caching Hit Rates
```
L1 Cache: 95% hit rate
L2 Cache: 85% hit rate  
Application Cache: 80% hit rate
CDN Cache: 90% hit rate

Effective latency with 80% cache hit:
= 0.8 * cache_latency + 0.2 * database_latency
= 0.8 * 1ms + 0.2 * 100ms = 20.8ms
```

### Database Optimization
```
Indexing overhead: 10-30% additional storage
Sharding: Linear scalability (ideally)
Replication: 3x storage, improved read performance
```

## 🧪 Practice Problems

### Problem 1: Design YouTube
```
Requirements:
- 2B users, 100M daily active users
- Users upload 500 hours of video per minute
- Users watch 1B hours of video per day
- Average video: 50MB, 10 minutes long

Calculate:
1. Upload QPS and bandwidth
2. View QPS and bandwidth  
3. Storage requirements for 5 years
4. CDN bandwidth requirements
5. Number of servers needed
```

### Problem 2: Design Uber
```
Requirements:
- 100M users, 10M daily active users
- 10M rides per day
- Each ride generates 100 location updates
- Location update: 50 bytes
- Trip data: 1KB per ride

Calculate:
1. Location update QPS
2. Storage for location data
3. Real-time processing requirements
4. Database and cache sizing
```

## 🎯 Key Tips for Interviews

### Rounding Rules
```
✅ Use powers of 10 for simplicity
✅ 86,400 seconds/day ≈ 10^5 seconds
✅ 1 million ≈ 10^6
✅ 1 billion ≈ 10^9

❌ Don't get lost in precise calculations
❌ Don't spend too much time on exact numbers
```

### Show Your Work
```go
// Good approach - show reasoning
daily_users := 1_000_000
messages_per_user := 50
daily_messages := daily_users * messages_per_user // 50M messages/day
seconds_per_day := 24 * 60 * 60 // ~100K seconds
qps := daily_messages / seconds_per_day // ~500 QPS

// Add safety margin
peak_qps := qps * 2 // 1000 QPS
```

### Validate Your Numbers
```
Sanity checks:
- Does 1000 QPS seem reasonable for this system?
- Is 100TB storage realistic for this use case?
- Can a single server handle this load?
- Are we over-engineering or under-estimating?
```

## 🔗 Next Steps
- Practice with different system types
- Learn about trade-offs and constraints
- Study real system architectures and their numbers
