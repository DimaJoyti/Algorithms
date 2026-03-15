# System Design Interview Guide

> Architecture patterns and frameworks for senior interviews

---

## Overview

### What's Evaluated

```
┌─────────────────────────────────────────────────────────────┐
│                    SYSTEM DESIGN INTERVIEW                   │
├─────────────────────────────────────────────────────────────┤
│  Requirements (20%)    │  Constraints & Assumptions (20%)    │
│  - Clarification       │  - Scale estimation                 │
│  - Functional/Non-func │  - Trade-offs                       │
├─────────────────────────────────────────────────────────────┤
│  Architecture (40%)    │  Deep Dive (20%)                    │
│  - Components          │  - Specific component detail        │
│  - Data flow           │  - Optimization                     │
│  - Scalability         │  - Edge cases                       │
└─────────────────────────────────────────────────────────────┘
```

### Interview Structure (45 min)

```
0-5 min:   Requirements gathering + clarification
5-10 min:  Define scope + capacity estimation
10-25 min: High-level architecture
25-35 min: Deep dive into components
35-40 min: Bottlenecks + optimizations
40-45 min: Wrap up + questions
```

---

## RESHADE Framework

### Step-by-Step Approach

```
R - Requirements (Functional + Non-functional)
E - Estimation (Capacity, storage, bandwidth)
S - Storage (Database choice, schema)
H - High-level Architecture (Components)
A - API Design (Endpoints, data models)
D - Deep Dive (Specific components)
E - Edge Cases + Optimizations
```

---

## 1. Requirements Gathering

### Clarifying Questions Template

```
Functional Requirements:
□ What are the core features?
□ Who are the users?
□ What devices/platforms?
□ What's the read/write ratio?
□ Real-time or eventual consistency?

Non-Functional Requirements:
□ What scale? (DAU, requests/sec)
□ Latency requirements?
□ Availability target? (99.9%?)
□ Consistency vs availability?
□ Cost constraints?
```

### Example: Design URL Shortener

```
Functional:
- Shorten long URL → short URL
- Redirect short URL → long URL
- Custom aliases?
- Analytics?

Non-Functional:
- Highly available (redirects must work)
- Low latency
- Unique short URLs
- 100M URLs/month

Out of Scope:
- User accounts
- Password protection
```

---

## 2. Capacity Estimation

### Quick Math Reference

```
Power of 2:
2^10 = 1 KB (Thousand)
2^20 = 1 MB (Million)
2^30 = 1 GB (Billion)
2^40 = 1 TB (Trillion)

Time:
1 day = 24 × 60 × 60 = 86,400 seconds
1 month ≈ 30 days ≈ 2.5M seconds
```

### Estimation Template

```
Assumptions:
- DAU: 1M
- Requests per user per day: 100
- Read:Write ratio: 100:1

QPS (Queries Per Second):
- Daily requests = 1M × 100 = 100M
- QPS = 100M / 86400 ≈ 1,200
- Peak QPS = QPS × 2-3 ≈ 3,600

Storage (per year):
- Write QPS = 1,200 / 100 = 12
- Daily writes = 12 × 86400 ≈ 1M
- Data per write: 1KB
- Daily storage = 1M × 1KB = 1GB
- Yearly storage = 365GB ≈ 0.4TB

Bandwidth:
- Incoming = QPS × request size
- Outgoing = QPS × response size
```

---

## 3. Storage Design

### Database Selection

| Type | Use Case | Examples |
|------|----------|----------|
| Relational | Structured, ACID, joins | PostgreSQL, MySQL |
| Document | Flexible schema, nested | MongoDB, CouchDB |
| Key-Value | Simple lookups, cache | Redis, DynamoDB |
| Column | Time series, analytics | Cassandra, HBase |
| Graph | Relationships | Neo4j, Neptune |
| Search | Full-text search | Elasticsearch, Solr |

### When to Use What

```
Relational (SQL):
✅ Strong consistency required
✅ Complex queries/joins
✅ Structured data
✅ Financial transactions

NoSQL - Document:
✅ Flexible schema
✅ Nested/hierarchical data
✅ Rapid iteration
❌ Complex joins

NoSQL - Key-Value:
✅ Simple lookups
✅ High throughput
✅ Caching layer
❌ Range queries

NoSQL - Column:
✅ Time series data
✅ High write volume
✅ Analytics
❌ Strong consistency

NoSQL - Graph:
✅ Social networks
✅ Recommendations
✅ Fraud detection
```

### Common Schemas

**URL Shortener:**
```
Table: urls
- id (BIGINT, auto-increment)
- short_code (VARCHAR, indexed)
- original_url (TEXT)
- created_at (TIMESTAMP)
- expiry_date (TIMESTAMP)
- user_id (BIGINT, nullable)
```

**Twitter Timeline:**
```
Table: users
- id, username, email, created_at

Table: tweets
- id, user_id, content, created_at

Table: follows
- follower_id, followee_id, created_at

Table: likes
- user_id, tweet_id, created_at
```

---

## 4. High-Level Architecture

### Basic Components

```
                    ┌─────────────┐
                    │   Client    │
                    │  (Mobile/   │
                    │   Web)      │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │    Load     │
                    │  Balancer   │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐     │     ┌──────▼──────┐
       │   API       │     │     │   API       │
       │  Server 1   │     │     │  Server N   │
       └──────┬──────┘     │     └──────┬──────┘
              │            │            │
              └────────────┼────────────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐     │     ┌──────▼──────┐
       │   Primary   │     │     │   Replica   │
       │   Database  ├─────┴─────┤   Database  │
       └─────────────┘           └─────────────┘
```

### Common Patterns

**1. Microservices**
```
Client → API Gateway → [Service A, Service B, Service C]
                        ↓              ↓              ↓
                      DB A           DB B           DB C
```

**2. Event-Driven**
```
Producer → Message Queue → Consumer
              ↓
           [Consumer 1, Consumer 2, Consumer N]
```

**3. CQRS (Command Query Responsibility Segregation)**
```
Write API → Write DB
              ↓ (Event)
            Event Bus
              ↓
Read API ← Read DB (Materialized View)
```

---

## 5. API Design

### RESTful Endpoints

```
GET    /api/v1/urls           # List URLs
POST   /api/v1/urls           # Create short URL
GET    /api/v1/urls/{code}    # Get original URL
DELETE /api/v1/urls/{code}    # Delete URL
GET    /api/v1/urls/{code}/stats  # Get analytics
```

### Response Formats

```json
// Create short URL
POST /api/v1/urls
Request:
{
  "original_url": "https://example.com/very/long/url",
  "custom_alias": "my-link",  // optional
  "expiry_days": 30           // optional
}

Response:
{
  "short_url": "https://tiny.url/my-link",
  "code": "my-link",
  "created_at": "2024-01-15T10:30:00Z",
  "expires_at": "2024-02-14T10:30:00Z"
}

// Error response
{
  "error": {
    "code": "INVALID_URL",
    "message": "The provided URL is not valid"
  }
}
```

---

## 6. Scalability Patterns

### Horizontal vs Vertical Scaling

```
Vertical (Scale Up):
+ Simpler
+ No code changes
- Hardware limits
- Single point of failure

Horizontal (Scale Out):
+ No theoretical limit
+ Fault tolerance
- Complexity
- Load balancing needed
```

### Caching Strategies

```
1. Cache-Aside (Lazy Loading)
   App → Check Cache → Miss → DB → Update Cache → Return

2. Write-Through
   App → Update Cache → Update DB → Return

3. Write-Behind (Write-Back)
   App → Update Cache → Return (Async: Update DB)

4. Refresh-Ahead
   Background refresh before expiry
```

### Cache Locations

```
Client Cache (Browser)
    ↓
CDN (Edge caching)
    ↓
Load Balancer Cache
    ↓
Application Cache (In-memory)
    ↓
Distributed Cache (Redis/Memcached)
    ↓
Database
```

### Database Scaling

```
1. Read Replicas
   - Write to primary
   - Read from replicas
   - Eventual consistency

2. Sharding (Horizontal Partitioning)
   - Split data across servers
   - Shard key selection critical
   - No cross-shard joins

3. Vertical Partitioning
   - Split tables by function
   - Frequently accessed together

4. Connection Pooling
   - Reuse connections
   - Limit max connections
```

---

## 7. Common System Designs

### URL Shortener (TinyURL)

```
Requirements:
- Generate short URL
- Redirect to original
- 100M URLs/month
- Custom aliases

Estimation:
- Write QPS: 40
- Read QPS: 4000
- Storage: 300GB/year

Architecture:
┌─────────┐     ┌────────────┐     ┌──────────┐
│ Client  │────▶│   Load     │────▶│  API     │
└─────────┘     │  Balancer  │     │ Servers  │
                └────────────┘     └────┬─────┘
                                        │
                    ┌───────────────────┼───────────────────┐
                    │                   │                   │
              ┌─────▼─────┐       ┌─────▼─────┐       ┌─────▼─────┐
              │   Redis   │       │  MySQL    │       │   CDN     │
              │  (Cache)  │       │ (Counter) │       │ (Static)  │
              └───────────┘       └───────────┘       └───────────┘

Key Components:
1. ID Generator (Distributed counter / ZooKeeper)
2. Base62 encoding for short code
3. Redis for hot URLs cache
4. Database for persistence
```

### Rate Limiter

```
Algorithms:

1. Token Bucket
   - Fixed rate tokens added
   - Request consumes token
   - Allows burst

2. Leaky Bucket
   - Fixed rate processing
   - Queue requests
   - Smooth traffic

3. Sliding Window Log
   - Track timestamps
   - Count in window
   - Memory intensive

4. Sliding Window Counter
   - Hybrid approach
   - Current + previous window
   - Memory efficient

Architecture:
┌─────────┐     ┌────────────┐     ┌──────────┐
│ Request │────▶│    Rate    │────▶│  Redis   │
└─────────┘     │  Limiter   │     │ (Counter)│
                │ Middleware │     └──────────┘
                └────────────┘

Implementation:
- Store counters in Redis
- Use sliding window
- Distributed with Redis
```

### Chat System (WhatsApp)

```
Requirements:
- 1-on-1 chat
- Group chat (100 members)
- Online status
- Message delivery status

Architecture:
┌─────────┐     ┌────────────┐     ┌──────────┐
│ Mobile  │◀───▶│  WebSocket │◀───▶│  Chat    │
│ Client  │     │  Server    │     │ Service  │
└─────────┘     └────────────┘     └────┬─────┘
                                        │
              ┌─────────────────────────┼─────────────────────┐
              │                         │                     │
        ┌─────▼─────┐           ┌───────▼───────┐     ┌───────▼───────┐
        │  Message  │           │   Presence    │     │ Notification  │
        │   Queue   │           │   Service     │     │   Service     │
        │ (Kafka)   │           │  (Redis)      │     │   (FCM/APN)   │
        └───────────┘           └───────────────┘     └───────────────┘

Components:
1. Connection Manager (WebSocket)
2. Message Queue (Kafka)
3. Message Storage (Cassandra)
4. Presence Service (Redis)
5. Push Notifications
```

### News Feed (Facebook/Twitter)

```
Requirements:
- Feed generation
- Follow/unfollow
- Like/comment
- Timeline scroll

Two Approaches:

1. Pull (Fan-out on Read)
   - Get followed user IDs
   - Fetch their recent posts
   - Merge and sort
   - Cache result

2. Push (Fan-out on Write)
   - On post: push to all followers' feeds
   - Pre-computed timelines
   - Fast reads, slower writes

Hybrid:
- Push for regular users
- Pull for celebrities (many followers)

Architecture:
┌─────────┐     ┌────────────┐     ┌──────────┐
│ Client  │────▶│    Feed    │────▶│  Feed    │
└─────────┘     │  Service   │     │  Cache   │
                └────────────┘     │ (Redis)  │
                        │          └──────────┘
                        │
              ┌─────────┼─────────┐
              │         │         │
        ┌─────▼─────┐ ┌─▼───┐ ┌───▼───┐
        │   Post    │ │User │ │Graph  │
        │  Service  │ │ Svc │ │  Svc  │
        └───────────┘ └─────┘ └───────┘
```

---

## 8. Trade-offs & Decisions

### CAP Theorem

```
C - Consistency (All nodes see same data)
A - Availability (System always responds)
P - Partition Tolerance (Network failures tolerated)

You can only pick 2:
- CP: Strong consistency, may be unavailable
- AP: Always available, may return stale data
- CA: Not possible in distributed systems
```

### Common Trade-offs

```
Consistency vs Latency
- Strong: Slower reads
- Eventual: Faster but stale data possible

SQL vs NoSQL
- SQL: ACID, joins, structured
- NoSQL: Scale, flexible schema, simple queries

Cache vs Freshness
- Cache: Fast but stale
- No cache: Fresh but slow

Batch vs Real-time
- Batch: Efficient, delayed
- Real-time: Immediate, resource intensive

Monolith vs Microservices
- Monolith: Simple, harder to scale
- Microservices: Complex, independent scaling
```

---

## 9. Deep Dive Topics

### Load Balancing

```
Algorithms:
- Round Robin
- Weighted Round Robin
- Least Connections
- IP Hash (Sticky sessions)
- Least Response Time

Health Checks:
- Active: Periodic requests
- Passive: Monitor responses

Layers:
- Layer 4 (Transport): TCP/UDP
- Layer 7 (Application): HTTP headers, paths
```

### Database Indexing

```
B-Tree Index:
- Default for most DBs
- Range queries
- O(log n) lookup

Hash Index:
- Exact match only
- O(1) lookup
- No range queries

Composite Index:
- Multiple columns
- Order matters
- Leftmost prefix rule

When to Index:
✅ Frequent WHERE clauses
✅ JOIN columns
✅ ORDER BY columns
❌ Frequent updates
❌ Low cardinality
```

### Message Queues

```
Use Cases:
- Async processing
- Load leveling
- Decoupling services
- Retry logic

Options:
- Kafka: High throughput, persistent
- RabbitMQ: Feature-rich, flexible routing
- SQS: Managed, simple
- Redis: Simple, in-memory

Patterns:
- Point-to-Point: One consumer
- Pub/Sub: Multiple consumers
- Dead Letter Queue: Failed messages
```

---

## 10. Quick Reference

### Availability Numbers

```
99%     = 3.65 days downtime/year
99.9%   = 8.76 hours downtime/year
99.99%  = 52.56 minutes downtime/year
99.999% = 5.26 minutes downtime/year
```

### Latency Numbers

```
L1 cache reference         0.5 ns
L2 cache reference         7 ns
Main memory reference    100 ns
SSD random read        16,000 ns
Read 1MB from SSD     500,000 ns
Network same datacenter  0.5 ms
Read 1MB from network     5 ms
Disk seek                10 ms
```

### Scalability Checklist

```
□ Stateless services
□ Horizontal scaling
□ Load balancing
□ Caching (multiple layers)
□ Database sharding/replication
□ Async processing
□ CDN for static content
□ Rate limiting
□ Monitoring & alerting
□ Graceful degradation
```

---

## 11. Practice Problems

### Beginner

1. Design a URL shortener
2. Design a pastebin
3. Design a rate limiter
4. Design a key-value store
5. Design a unique ID generator

### Intermediate

6. Design Twitter/X
7. Design Instagram
8. Design Netflix
9. Design Uber
10. Design WhatsApp

### Advanced

11. Design Google Search
12. Design YouTube
13. Design Google Docs (collaborative editing)
14. Design a distributed cache
15. Design a CDN

---

## Summary Framework

### 5-Minute Pitch Structure

```
1. Problem Statement (30 sec)
   "I'm designing X which does Y for Z users"

2. Requirements (1 min)
   "Key features are A, B, C"
   "Scale is X users, Y requests/sec"

3. High-Level Design (1.5 min)
   "Core components: A, B, C"
   "Data flows like this..."

4. Deep Dive (1.5 min)
   "Let me explain X in more detail"
   "For Y, I chose Z because..."

5. Trade-offs (30 sec)
   "Main trade-off is A vs B"
   "I prioritized X because..."
```

---

*Remember: Communication > Perfect Solution. Think out loud, ask questions, show reasoning.*
