# 💬 WhatsApp/Chat System Design

## 🎯 Step 1: Requirements Clarification (5 minutes)

### Functional Requirements
```
✅ Core Features:
- Send and receive text messages in real-time
- One-on-one and group chat (up to 256 members)
- Message delivery status (sent, delivered, read)
- Online/offline status and last seen
- Media sharing (images, videos, documents)
- Message history and search
- End-to-end encryption
- Push notifications

❌ Out of Scope:
- Voice/video calls
- Stories/Status updates
- Payment features
- Business accounts
- Advanced bot integration
```

### Non-Functional Requirements
```
Scale:
- 2B active users globally
- 100B messages per day
- 50M concurrent users
- 1M new users per day
- Support for 180+ countries

Performance:
- Message delivery: < 100ms
- Message history load: < 200ms
- 99.99% availability
- End-to-end encryption
- Offline message delivery

Reliability:
- Message delivery guarantee
- No message loss
- Consistent message ordering
- Cross-device synchronization
```

## 📊 Step 2: Capacity Estimation

### Back-of-Envelope Calculations
```go
// Global scale metrics
const (
    TotalUsers        = 2_000_000_000  // 2B users
    DailyActiveUsers  = 1_500_000_000  // 1.5B DAU
    MessagesPerDay    = 100_000_000_000 // 100B messages/day
    ConcurrentUsers   = 50_000_000     // 50M concurrent
    GroupChats        = 500_000_000    // 500M group chats
    
    SecondsPerDay = 24 * 60 * 60
    AvgMessageSize = 100 // bytes (text message)
    MediaMessageSize = 1024 * 1024 // 1MB average
)

// Message patterns
var (
    MessagesPerSecond = MessagesPerDay / SecondsPerDay        // ~1.16M/sec
    PeakMessagesPerSecond = MessagesPerSecond * 3             // ~3.5M/sec
    MediaMessagesPerDay = MessagesPerDay / 10                 // 10% are media
    TextMessagesPerDay = MessagesPerDay - MediaMessagesPerDay // 90% are text
)

// Storage calculations
type StorageCalculation struct {
    TextMessages  int64 // Text message storage
    MediaMessages int64 // Media file storage
    UserData      int64 // User profiles and metadata
    MessageIndex  int64 // Search and indexing
    Metadata      int64 // Delivery status, timestamps
}

func CalculateStorage() StorageCalculation {
    calc := StorageCalculation{}
    
    // Text messages: 90B messages/day * 100 bytes * 365 days * 2 years
    calc.TextMessages = int64(TextMessagesPerDay) * AvgMessageSize * 365 * 2 // ~6.6 PB
    
    // Media messages: 10B messages/day * 1MB * 365 days * 2 years
    calc.MediaMessages = int64(MediaMessagesPerDay) * MediaMessageSize * 365 * 2 // ~7.3 EB
    
    // User data: 2B users * 1KB each
    calc.UserData = TotalUsers * 1024 // ~2 TB
    
    // Message metadata: 100B messages/day * 200 bytes * 365 days * 2 years
    calc.Metadata = MessagesPerDay * 200 * 365 * 2 // ~14.6 PB
    
    // Search index: ~10% of text message size
    calc.MessageIndex = calc.TextMessages / 10 // ~660 TB
    
    return calc
}

// Real-time connection requirements
func CalculateConnections() map[string]int64 {
    return map[string]int64{
        "websocket_connections":    ConcurrentUsers,     // 50M concurrent connections
        "message_queue_throughput": PeakMessagesPerSecond, // 3.5M messages/sec
        "database_writes":          PeakMessagesPerSecond, // 3.5M writes/sec
        "cache_operations":         PeakMessagesPerSecond * 5, // 17.5M cache ops/sec
        "push_notifications":       PeakMessagesPerSecond / 2, // 1.75M notifications/sec
    }
}
```

### Infrastructure Estimation
```go
type InfrastructureNeeds struct {
    MessageServers      int // Handle message routing
    WebSocketServers    int // Maintain persistent connections
    DatabaseServers     int // Message and user data storage
    CacheServers        int // Redis for real-time data
    MediaServers        int // Handle media uploads/downloads
    NotificationServers int // Push notification service
    EncryptionServers   int // End-to-end encryption
}

func EstimateInfrastructure() InfrastructureNeeds {
    // WebSocket servers (each handles 10K connections)
    websocketServers := int(ConcurrentUsers / 10000) // ~5,000 servers
    
    // Message servers (each handles 5K messages/sec)
    messageServers := int(PeakMessagesPerSecond / 5000) // ~700 servers
    
    // Database servers (sharded by user_id and chat_id)
    databaseServers := 1000 // Heavily sharded for write performance
    
    // Cache servers (Redis cluster, each 64GB)
    cacheServers := 500 // For online users, recent messages, delivery status
    
    // Media servers (each handles 1K uploads/sec)
    mediaServers := int(PeakMessagesPerSecond / 10 / 1000) // ~350 servers
    
    // Notification servers
    notificationServers := 100 // For push notifications
    
    // Encryption servers (for key management)
    encryptionServers := 50
    
    return InfrastructureNeeds{
        MessageServers:      messageServers,
        WebSocketServers:    websocketServers,
        DatabaseServers:     databaseServers,
        CacheServers:        cacheServers,
        MediaServers:        mediaServers,
        NotificationServers: notificationServers,
        EncryptionServers:   encryptionServers,
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
                                           ├── [Chat Service]
                                           ├── [Message Service]
                                           ├── [Media Service]
                                           ├── [Notification Service]
                                           ├── [Encryption Service]
                                           └── [Presence Service]
                                                    │
                        [Message Queue] ── [Cache Layer] ── [Database Layer]
                                │
                        [WebSocket Gateway] ── [Connection Pool]
```

### Core Services
```go
// User Service - manages user accounts and contacts
type UserService struct {
    userRepo        UserRepository
    contactRepo     ContactRepository
    authService     AuthService
    encryptionService EncryptionService
    cache           Cache
}

// Chat Service - manages chat rooms and group metadata
type ChatService struct {
    chatRepo        ChatRepository
    memberRepo      MemberRepository
    permissionService PermissionService
    cache           Cache
    eventBus        EventBus
}

// Message Service - handles message delivery and storage
type MessageService struct {
    messageRepo     MessageRepository
    deliveryService DeliveryService
    encryptionService EncryptionService
    messageQueue    MessageQueue
    cache           Cache
}

// Media Service - handles file uploads and downloads
type MediaService struct {
    mediaRepo       MediaRepository
    storageService  StorageService
    compressionService CompressionService
    cdnService      CDNService
}

// Notification Service - push notifications
type NotificationService struct {
    notificationRepo NotificationRepository
    pushService     PushService
    templateService TemplateService
    eventBus        EventBus
}

// Presence Service - online status and last seen
type PresenceService struct {
    presenceRepo    PresenceRepository
    connectionManager ConnectionManager
    cache           Cache
    eventBus        EventBus
}

// WebSocket Gateway - manages real-time connections
type WebSocketGateway struct {
    connectionPool  ConnectionPool
    messageRouter   MessageRouter
    authService     AuthService
    presenceService PresenceService
}
```

## 🗄️ Step 4: Database Design

### Data Models
```go
// User models
type User struct {
    ID              int64     `json:"id" db:"id"`
    PhoneNumber     string    `json:"phone_number" db:"phone_number"`
    Username        string    `json:"username" db:"username"`
    DisplayName     string    `json:"display_name" db:"display_name"`
    ProfileImageURL string    `json:"profile_image_url" db:"profile_image_url"`
    Status          string    `json:"status" db:"status"` // online, offline, away
    LastSeen        time.Time `json:"last_seen" db:"last_seen"`
    PublicKey       string    `json:"public_key" db:"public_key"` // For E2E encryption
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Contact struct {
    UserID      int64     `json:"user_id" db:"user_id"`
    ContactID   int64     `json:"contact_id" db:"contact_id"`
    DisplayName string    `json:"display_name" db:"display_name"` // Custom name
    IsBlocked   bool      `json:"is_blocked" db:"is_blocked"`
    AddedAt     time.Time `json:"added_at" db:"added_at"`
}

// Chat models
type Chat struct {
    ID          int64     `json:"id" db:"id"`
    Type        string    `json:"type" db:"type"` // direct, group
    Name        string    `json:"name,omitempty" db:"name"` // For group chats
    Description string    `json:"description,omitempty" db:"description"`
    ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
    CreatedBy   int64     `json:"created_by" db:"created_by"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ChatMember struct {
    ChatID    int64     `json:"chat_id" db:"chat_id"`
    UserID    int64     `json:"user_id" db:"user_id"`
    Role      string    `json:"role" db:"role"` // admin, member
    JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
    LeftAt    *time.Time `json:"left_at,omitempty" db:"left_at"`
}

// Message models
type Message struct {
    ID          int64     `json:"id" db:"id"`
    ChatID      int64     `json:"chat_id" db:"chat_id"`
    SenderID    int64     `json:"sender_id" db:"sender_id"`
    Type        string    `json:"type" db:"type"` // text, image, video, document
    Content     string    `json:"content" db:"content"` // Encrypted content
    MediaURL    string    `json:"media_url,omitempty" db:"media_url"`
    ReplyToID   *int64    `json:"reply_to_id,omitempty" db:"reply_to_id"`
    IsEdited    bool      `json:"is_edited" db:"is_edited"`
    IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Message delivery tracking
type MessageDelivery struct {
    MessageID   int64     `json:"message_id" db:"message_id"`
    UserID      int64     `json:"user_id" db:"user_id"`
    Status      string    `json:"status" db:"status"` // sent, delivered, read
    DeliveredAt *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
    ReadAt      *time.Time `json:"read_at,omitempty" db:"read_at"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Media models
type MediaFile struct {
    ID          int64     `json:"id" db:"id"`
    MessageID   int64     `json:"message_id" db:"message_id"`
    FileName    string    `json:"file_name" db:"file_name"`
    FileSize    int64     `json:"file_size" db:"file_size"`
    MimeType    string    `json:"mime_type" db:"mime_type"`
    StoragePath string    `json:"storage_path" db:"storage_path"`
    CDNPath     string    `json:"cdn_path" db:"cdn_path"`
    Thumbnail   string    `json:"thumbnail,omitempty" db:"thumbnail"`
    Duration    *int      `json:"duration,omitempty" db:"duration"` // For videos/audio
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Presence models
type UserPresence struct {
    UserID       int64     `json:"user_id" db:"user_id"`
    Status       string    `json:"status" db:"status"` // online, offline, away
    LastSeen     time.Time `json:"last_seen" db:"last_seen"`
    DeviceType   string    `json:"device_type" db:"device_type"`
    ConnectionID string    `json:"connection_id" db:"connection_id"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
```

### Database Schema
```sql
-- Users table (PostgreSQL)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    profile_image_url TEXT,
    status VARCHAR(20) DEFAULT 'offline',
    last_seen TIMESTAMP DEFAULT NOW(),
    public_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_username ON users(username);

-- Contacts table
CREATE TABLE contacts (
    user_id BIGINT REFERENCES users(id),
    contact_id BIGINT REFERENCES users(id),
    display_name VARCHAR(100),
    is_blocked BOOLEAN DEFAULT FALSE,
    added_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, contact_id)
);

-- Chats table
CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL, -- direct, group
    name VARCHAR(100),
    description TEXT,
    image_url TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chats_type ON chats(type);
CREATE INDEX idx_chats_created_by ON chats(created_by);

-- Chat members table
CREATE TABLE chat_members (
    chat_id BIGINT REFERENCES chats(id),
    user_id BIGINT REFERENCES users(id),
    role VARCHAR(10) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP,
    PRIMARY KEY (chat_id, user_id)
);

CREATE INDEX idx_chat_members_user_id ON chat_members(user_id);

-- Messages table (Cassandra for high write volume)
CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    chat_id BIGINT,
    sender_id BIGINT,
    type VARCHAR(20) NOT NULL,
    content TEXT,
    media_url TEXT,
    reply_to_id BIGINT,
    is_edited BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Partition by chat_id for better performance
CREATE INDEX idx_messages_chat_id_created_at ON messages(chat_id, created_at DESC);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);

-- Message delivery table (Cassandra)
CREATE TABLE message_delivery (
    message_id BIGINT,
    user_id BIGINT,
    status VARCHAR(20) NOT NULL,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);

-- Media files table
CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT REFERENCES messages(id),
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    storage_path TEXT NOT NULL,
    cdn_path TEXT,
    thumbnail TEXT,
    duration INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_media_files_message_id ON media_files(message_id);

-- User presence (Redis for real-time data)
-- Stored as key-value pairs in Redis:
-- user:presence:{user_id} -> {status, last_seen, device_type, connection_id}
```

### Database Sharding Strategy
```go
type WhatsAppDatabaseSharding struct {
    userShards    []Database // Shard by user_id
    chatShards    []Database // Shard by chat_id
    messageShards []Database // Shard by chat_id for locality
}

// User data sharding - by user_id for user-specific queries
func (wds *WhatsAppDatabaseSharding) GetUserShard(userID int64) Database {
    shardIndex := userID % int64(len(wds.userShards))
    return wds.userShards[shardIndex]
}

// Chat data sharding - by chat_id for chat-specific queries
func (wds *WhatsAppDatabaseSharding) GetChatShard(chatID int64) Database {
    shardIndex := chatID % int64(len(wds.chatShards))
    return wds.chatShards[shardIndex]
}

// Message sharding - by chat_id to keep chat messages together
func (wds *WhatsAppDatabaseSharding) GetMessageShard(chatID int64) Database {
    shardIndex := chatID % int64(len(wds.messageShards))
    return wds.messageShards[shardIndex]
}
```

## 📱 Step 5: Real-time Message Delivery

### WebSocket Connection Management
```go
type WebSocketGateway struct {
    connectionPool  *ConnectionPool
    messageRouter   *MessageRouter
    authService     AuthService
    presenceService PresenceService
    rateLimiter     RateLimiter
}

type ConnectionPool struct {
    connections map[int64]*UserConnection // userID -> connection
    mutex       sync.RWMutex
    metrics     ConnectionMetrics
}

type UserConnection struct {
    UserID       int64
    ConnectionID string
    WebSocket    *websocket.Conn
    DeviceType   string
    LastActivity time.Time
    SendChannel  chan []byte
    IsActive     bool
    mutex        sync.Mutex
}

func (wsg *WebSocketGateway) HandleConnection(w http.ResponseWriter, r *http.Request) {
    // Authenticate user
    userID, err := wsg.authService.ValidateToken(r.Header.Get("Authorization"))
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Upgrade to WebSocket
    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()

    // Create user connection
    connectionID := generateConnectionID()
    userConn := &UserConnection{
        UserID:       userID,
        ConnectionID: connectionID,
        WebSocket:    conn,
        DeviceType:   r.Header.Get("Device-Type"),
        LastActivity: time.Now(),
        SendChannel:  make(chan []byte, 1000),
        IsActive:     true,
    }

    // Register connection
    wsg.connectionPool.AddConnection(userID, userConn)
    defer wsg.connectionPool.RemoveConnection(userID, connectionID)

    // Update user presence
    wsg.presenceService.SetUserOnline(userID, connectionID, userConn.DeviceType)
    defer wsg.presenceService.SetUserOffline(userID, connectionID)

    // Start message sender goroutine
    go wsg.messageSender(userConn)

    // Handle incoming messages
    wsg.messageReceiver(userConn)
}

func (wsg *WebSocketGateway) messageSender(userConn *UserConnection) {
    ticker := time.NewTicker(30 * time.Second) // Ping interval
    defer ticker.Stop()

    for {
        select {
        case message := <-userConn.SendChannel:
            userConn.mutex.Lock()
            if userConn.IsActive {
                err := userConn.WebSocket.WriteMessage(websocket.TextMessage, message)
                if err != nil {
                    userConn.IsActive = false
                    userConn.mutex.Unlock()
                    return
                }
                userConn.LastActivity = time.Now()
            }
            userConn.mutex.Unlock()

        case <-ticker.C:
            // Send ping to keep connection alive
            userConn.mutex.Lock()
            if userConn.IsActive {
                err := userConn.WebSocket.WriteMessage(websocket.PingMessage, nil)
                if err != nil {
                    userConn.IsActive = false
                    userConn.mutex.Unlock()
                    return
                }
            }
            userConn.mutex.Unlock()
        }
    }
}

func (wsg *WebSocketGateway) messageReceiver(userConn *UserConnection) {
    for {
        messageType, data, err := userConn.WebSocket.ReadMessage()
        if err != nil {
            break
        }

        userConn.LastActivity = time.Now()

        switch messageType {
        case websocket.TextMessage:
            wsg.handleIncomingMessage(userConn, data)
        case websocket.PongMessage:
            // Connection is alive
        }
    }
}

func (wsg *WebSocketGateway) handleIncomingMessage(userConn *UserConnection, data []byte) {
    var incomingMsg IncomingMessage
    if err := json.Unmarshal(data, &incomingMsg); err != nil {
        return
    }

    // Rate limiting
    if !wsg.rateLimiter.Allow(userConn.UserID) {
        return
    }

    // Route message based on type
    switch incomingMsg.Type {
    case "send_message":
        wsg.messageRouter.RouteMessage(userConn.UserID, incomingMsg)
    case "message_read":
        wsg.messageRouter.MarkMessageRead(userConn.UserID, incomingMsg)
    case "typing":
        wsg.messageRouter.HandleTypingIndicator(userConn.UserID, incomingMsg)
    }
}
```

### Message Routing and Delivery
```go
type MessageRouter struct {
    messageService  MessageService
    connectionPool  *ConnectionPool
    messageQueue    MessageQueue
    deliveryService DeliveryService
}

func (mr *MessageRouter) RouteMessage(senderID int64, incomingMsg IncomingMessage) {
    // Create message record
    message := &Message{
        ChatID:    incomingMsg.ChatID,
        SenderID:  senderID,
        Type:      incomingMsg.MessageType,
        Content:   incomingMsg.Content,
        CreatedAt: time.Now(),
    }

    // Save message to database
    if err := mr.messageService.SaveMessage(context.Background(), message); err != nil {
        log.Printf("Failed to save message: %v", err)
        return
    }

    // Get chat members
    members, err := mr.messageService.GetChatMembers(context.Background(), incomingMsg.ChatID)
    if err != nil {
        log.Printf("Failed to get chat members: %v", err)
        return
    }

    // Route to online users immediately
    onlineMembers := make([]int64, 0)
    offlineMembers := make([]int64, 0)

    for _, memberID := range members {
        if memberID == senderID {
            continue // Don't send to sender
        }

        if mr.connectionPool.IsUserOnline(memberID) {
            onlineMembers = append(onlineMembers, memberID)
            mr.deliverMessageToUser(memberID, message)
        } else {
            offlineMembers = append(offlineMembers, memberID)
        }
    }

    // Queue messages for offline users
    if len(offlineMembers) > 0 {
        mr.messageQueue.QueueOfflineDelivery(message, offlineMembers)
    }

    // Update delivery status
    mr.deliveryService.UpdateDeliveryStatus(message.ID, onlineMembers, "delivered")
}

func (mr *MessageRouter) deliverMessageToUser(userID int64, message *Message) {
    connections := mr.connectionPool.GetUserConnections(userID)

    outgoingMsg := OutgoingMessage{
        Type:      "new_message",
        MessageID: message.ID,
        ChatID:    message.ChatID,
        SenderID:  message.SenderID,
        Content:   message.Content,
        Timestamp: message.CreatedAt,
    }

    messageData, _ := json.Marshal(outgoingMsg)

    // Send to all user's active connections
    for _, conn := range connections {
        select {
        case conn.SendChannel <- messageData:
            // Message queued for delivery
        default:
            // Channel full, connection might be slow
            log.Printf("Send channel full for user %d", userID)
        }
    }
}

// Message delivery guarantees
type DeliveryService struct {
    deliveryRepo    DeliveryRepository
    retryQueue      RetryQueue
    notificationService NotificationService
}

func (ds *DeliveryService) EnsureDelivery(message *Message, recipientIDs []int64) {
    for _, recipientID := range recipientIDs {
        delivery := &MessageDelivery{
            MessageID: message.ID,
            UserID:    recipientID,
            Status:    "sent",
            CreatedAt: time.Now(),
        }

        ds.deliveryRepo.Create(context.Background(), delivery)

        // Set retry timer for undelivered messages
        go ds.scheduleRetry(delivery, 30*time.Second)
    }
}

func (ds *DeliveryService) scheduleRetry(delivery *MessageDelivery, delay time.Duration) {
    time.Sleep(delay)

    // Check if message was delivered
    currentDelivery, err := ds.deliveryRepo.GetByMessageAndUser(
        context.Background(),
        delivery.MessageID,
        delivery.UserID,
    )
    if err != nil {
        return
    }

    if currentDelivery.Status == "sent" {
        // Still not delivered, send push notification
        ds.notificationService.SendPushNotification(context.Background(), PushNotificationRequest{
            UserID:    delivery.UserID,
            MessageID: delivery.MessageID,
            Type:      "new_message",
        })

        // Schedule another retry
        if delay < 5*time.Minute {
            go ds.scheduleRetry(delivery, delay*2) // Exponential backoff
        }
    }
}
```

## 🔐 Step 6: End-to-End Encryption

### Encryption Service Implementation
```go
type EncryptionService struct {
    keyRepo         KeyRepository
    cryptoProvider  CryptoProvider
    keyExchange     KeyExchangeService
}

// Signal Protocol implementation for E2E encryption
type SignalProtocol struct {
    identityKeyPair *IdentityKeyPair
    preKeys         map[int]*PreKey
    signedPreKey    *SignedPreKey
    sessionStore    SessionStore
}

func (es *EncryptionService) InitializeUserKeys(userID int64) error {
    // Generate identity key pair
    identityKeyPair, err := es.cryptoProvider.GenerateIdentityKeyPair()
    if err != nil {
        return err
    }

    // Generate pre-keys (one-time keys)
    preKeys := make([]*PreKey, 100)
    for i := 0; i < 100; i++ {
        preKey, err := es.cryptoProvider.GeneratePreKey(i)
        if err != nil {
            return err
        }
        preKeys[i] = preKey
    }

    // Generate signed pre-key
    signedPreKey, err := es.cryptoProvider.GenerateSignedPreKey(identityKeyPair.PrivateKey)
    if err != nil {
        return err
    }

    // Store keys
    userKeys := &UserKeys{
        UserID:          userID,
        IdentityKeyPair: identityKeyPair,
        PreKeys:         preKeys,
        SignedPreKey:    signedPreKey,
        CreatedAt:       time.Now(),
    }

    return es.keyRepo.SaveUserKeys(context.Background(), userKeys)
}

func (es *EncryptionService) EncryptMessage(senderID, recipientID int64, plaintext string) (*EncryptedMessage, error) {
    // Get or create session
    session, err := es.getOrCreateSession(senderID, recipientID)
    if err != nil {
        return nil, err
    }

    // Encrypt message
    ciphertext, err := session.Encrypt([]byte(plaintext))
    if err != nil {
        return nil, err
    }

    return &EncryptedMessage{
        SenderID:     senderID,
        RecipientID:  recipientID,
        Ciphertext:   ciphertext,
        SessionID:    session.ID,
        EncryptedAt:  time.Now(),
    }, nil
}

func (es *EncryptionService) DecryptMessage(encryptedMsg *EncryptedMessage) (string, error) {
    // Get session
    session, err := es.getSession(encryptedMsg.SenderID, encryptedMsg.RecipientID)
    if err != nil {
        return "", err
    }

    // Decrypt message
    plaintext, err := session.Decrypt(encryptedMsg.Ciphertext)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

// Group chat encryption (using sender keys)
func (es *EncryptionService) EncryptGroupMessage(senderID int64, chatID int64, plaintext string) (*GroupEncryptedMessage, error) {
    // Get chat members
    members, err := es.getChatMembers(chatID)
    if err != nil {
        return nil, err
    }

    // Encrypt for each member
    encryptedForMembers := make(map[int64][]byte)

    for _, memberID := range members {
        if memberID == senderID {
            continue
        }

        encryptedMsg, err := es.EncryptMessage(senderID, memberID, plaintext)
        if err != nil {
            continue // Skip failed encryptions
        }

        encryptedForMembers[memberID] = encryptedMsg.Ciphertext
    }

    return &GroupEncryptedMessage{
        SenderID:            senderID,
        ChatID:              chatID,
        EncryptedForMembers: encryptedForMembers,
        EncryptedAt:         time.Now(),
    }, nil
}
```

## 📲 Step 7: Push Notification System

### Notification Service Implementation
```go
type NotificationService struct {
    notificationRepo NotificationRepository
    pushProviders    map[string]PushProvider // iOS, Android, Web
    templateService  TemplateService
    userPreferences  UserPreferencesService
    rateLimiter     NotificationRateLimiter
}

type PushProvider interface {
    SendNotification(ctx context.Context, request PushRequest) error
    SendBulkNotifications(ctx context.Context, requests []PushRequest) error
}

func (ns *NotificationService) SendMessageNotification(ctx context.Context, message *Message, recipientID int64) error {
    // Check if user is online
    if ns.isUserOnline(recipientID) {
        return nil // Don't send push if user is online
    }

    // Check user notification preferences
    preferences, err := ns.userPreferences.GetPreferences(ctx, recipientID)
    if err != nil || !preferences.MessageNotifications {
        return nil
    }

    // Rate limiting to prevent spam
    if !ns.rateLimiter.Allow(recipientID) {
        return nil
    }

    // Get user devices
    devices, err := ns.getUserDevices(ctx, recipientID)
    if err != nil {
        return err
    }

    // Create notification content
    notification := ns.createMessageNotification(ctx, message, recipientID)

    // Send to all user devices
    for _, device := range devices {
        pushRequest := PushRequest{
            DeviceToken: device.Token,
            Platform:    device.Platform,
            Title:       notification.Title,
            Body:        notification.Body,
            Data:        notification.Data,
            Badge:       ns.getUnreadCount(ctx, recipientID),
        }

        provider := ns.pushProviders[device.Platform]
        go provider.SendNotification(ctx, pushRequest)
    }

    return nil
}

func (ns *NotificationService) createMessageNotification(ctx context.Context, message *Message, recipientID int64) *NotificationContent {
    // Get sender info
    sender, _ := ns.getSenderInfo(ctx, message.SenderID)

    // Get chat info
    chat, _ := ns.getChatInfo(ctx, message.ChatID)

    var title, body string

    if chat.Type == "direct" {
        title = sender.DisplayName
        body = ns.getMessagePreview(message)
    } else {
        title = chat.Name
        body = fmt.Sprintf("%s: %s", sender.DisplayName, ns.getMessagePreview(message))
    }

    return &NotificationContent{
        Title: title,
        Body:  body,
        Data: map[string]interface{}{
            "chat_id":    message.ChatID,
            "message_id": message.ID,
            "sender_id":  message.SenderID,
            "type":       "new_message",
        },
    }
}

// iOS Push Notification Provider (APNs)
type APNsProvider struct {
    client   *apns2.Client
    bundleID string
}

func (ap *APNsProvider) SendNotification(ctx context.Context, request PushRequest) error {
    notification := &apns2.Notification{
        DeviceToken: request.DeviceToken,
        Topic:       ap.bundleID,
        Payload: &payload.Payload{
            Alert: &payload.Alert{
                Title: request.Title,
                Body:  request.Body,
            },
            Badge:            &request.Badge,
            Sound:            "default",
            MutableContent:   1,
            ContentAvailable: 1,
            Custom:           request.Data,
        },
    }

    response, err := ap.client.Push(notification)
    if err != nil {
        return err
    }

    if response.StatusCode != 200 {
        return fmt.Errorf("APNs error: %s", response.Reason)
    }

    return nil
}

// Android Push Notification Provider (FCM)
type FCMProvider struct {
    client *messaging.Client
}

func (fp *FCMProvider) SendNotification(ctx context.Context, request PushRequest) error {
    message := &messaging.Message{
        Token: request.DeviceToken,
        Notification: &messaging.Notification{
            Title: request.Title,
            Body:  request.Body,
        },
        Data: convertToStringMap(request.Data),
        Android: &messaging.AndroidConfig{
            Priority: "high",
            Notification: &messaging.AndroidNotification{
                ChannelID: "messages",
                Priority:  messaging.PriorityHigh,
            },
        },
    }

    _, err := fp.client.Send(ctx, message)
    return err
}
```

### Smart Notification Batching
```go
type NotificationBatcher struct {
    batches     map[int64]*NotificationBatch // userID -> batch
    mutex       sync.RWMutex
    flushInterval time.Duration
}

type NotificationBatch struct {
    UserID       int64
    Messages     []*Message
    LastUpdated  time.Time
    FlushTimer   *time.Timer
}

func (nb *NotificationBatcher) AddMessage(userID int64, message *Message) {
    nb.mutex.Lock()
    defer nb.mutex.Unlock()

    batch, exists := nb.batches[userID]
    if !exists {
        batch = &NotificationBatch{
            UserID:      userID,
            Messages:    make([]*Message, 0),
            LastUpdated: time.Now(),
        }
        nb.batches[userID] = batch
    }

    batch.Messages = append(batch.Messages, message)
    batch.LastUpdated = time.Now()

    // Reset flush timer
    if batch.FlushTimer != nil {
        batch.FlushTimer.Stop()
    }

    batch.FlushTimer = time.AfterFunc(nb.flushInterval, func() {
        nb.flushBatch(userID)
    })
}

func (nb *NotificationBatcher) flushBatch(userID int64) {
    nb.mutex.Lock()
    batch, exists := nb.batches[userID]
    if !exists {
        nb.mutex.Unlock()
        return
    }

    messages := make([]*Message, len(batch.Messages))
    copy(messages, batch.Messages)
    delete(nb.batches, userID)
    nb.mutex.Unlock()

    // Create batched notification
    if len(messages) == 1 {
        // Single message notification
        nb.sendSingleMessageNotification(userID, messages[0])
    } else {
        // Multiple messages notification
        nb.sendBatchedNotification(userID, messages)
    }
}

func (nb *NotificationBatcher) sendBatchedNotification(userID int64, messages []*Message) {
    // Group by chat
    chatMessages := make(map[int64][]*Message)
    for _, msg := range messages {
        chatMessages[msg.ChatID] = append(chatMessages[msg.ChatID], msg)
    }

    var title, body string
    if len(chatMessages) == 1 {
        // Messages from single chat
        for chatID, msgs := range chatMessages {
            chat, _ := nb.getChatInfo(chatID)
            title = chat.Name
            body = fmt.Sprintf("%d new messages", len(msgs))
        }
    } else {
        // Messages from multiple chats
        title = "WhatsApp"
        body = fmt.Sprintf("%d new messages from %d chats", len(messages), len(chatMessages))
    }

    // Send notification
    nb.sendNotification(userID, title, body, map[string]interface{}{
        "type":         "batched_messages",
        "message_count": len(messages),
        "chat_count":   len(chatMessages),
    })
}
```

## 📁 Step 8: Media Handling

### Media Service Implementation
```go
type MediaService struct {
    mediaRepo       MediaRepository
    storageService  StorageService
    compressionService CompressionService
    cdnService      CDNService
    virusScanner    VirusScanner
}

func (ms *MediaService) UploadMedia(ctx context.Context, request MediaUploadRequest) (*MediaUploadResult, error) {
    // Validate file
    if err := ms.validateMediaFile(request); err != nil {
        return nil, err
    }

    // Virus scan
    if err := ms.virusScanner.ScanFile(request.FileData); err != nil {
        return nil, ErrVirusDetected
    }

    // Generate unique file ID
    fileID := generateFileID()

    // Process based on media type
    var processedFiles []ProcessedFile

    switch request.MediaType {
    case "image":
        processedFiles = ms.processImage(request.FileData, fileID)
    case "video":
        processedFiles = ms.processVideo(request.FileData, fileID)
    case "audio":
        processedFiles = ms.processAudio(request.FileData, fileID)
    case "document":
        processedFiles = ms.processDocument(request.FileData, fileID)
    }

    // Upload to storage
    var uploadResults []UploadResult
    for _, file := range processedFiles {
        result, err := ms.storageService.Upload(ctx, file)
        if err != nil {
            continue // Skip failed uploads
        }
        uploadResults = append(uploadResults, result)
    }

    // Save media metadata
    mediaFile := &MediaFile{
        ID:          fileID,
        MessageID:   request.MessageID,
        FileName:    request.FileName,
        FileSize:    request.FileSize,
        MimeType:    request.MimeType,
        StoragePath: uploadResults[0].StoragePath,
        CDNPath:     uploadResults[0].CDNPath,
        CreatedAt:   time.Now(),
    }

    if err := ms.mediaRepo.Save(ctx, mediaFile); err != nil {
        return nil, err
    }

    return &MediaUploadResult{
        FileID:   fileID,
        URL:      uploadResults[0].CDNPath,
        FileSize: mediaFile.FileSize,
    }, nil
}

func (ms *MediaService) processImage(imageData []byte, fileID string) []ProcessedFile {
    var files []ProcessedFile

    // Original image
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "original",
        Data:     imageData,
        FileName: fmt.Sprintf("%s_original.jpg", fileID),
    })

    // Compressed version for mobile
    compressed := ms.compressionService.CompressImage(imageData, CompressionOptions{
        Quality:   80,
        MaxWidth:  1200,
        MaxHeight: 1200,
    })
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "compressed",
        Data:     compressed,
        FileName: fmt.Sprintf("%s_compressed.jpg", fileID),
    })

    // Thumbnail
    thumbnail := ms.compressionService.GenerateThumbnail(imageData, ThumbnailOptions{
        Width:  200,
        Height: 200,
    })
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "thumbnail",
        Data:     thumbnail,
        FileName: fmt.Sprintf("%s_thumb.jpg", fileID),
    })

    return files
}

func (ms *MediaService) processVideo(videoData []byte, fileID string) []ProcessedFile {
    var files []ProcessedFile

    // Original video
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "original",
        Data:     videoData,
        FileName: fmt.Sprintf("%s_original.mp4", fileID),
    })

    // Compressed version
    compressed := ms.compressionService.CompressVideo(videoData, VideoCompressionOptions{
        Bitrate:    "1000k",
        Resolution: "720p",
        Format:     "mp4",
    })
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "compressed",
        Data:     compressed,
        FileName: fmt.Sprintf("%s_compressed.mp4", fileID),
    })

    // Generate video thumbnail
    thumbnail := ms.compressionService.GenerateVideoThumbnail(videoData, 0) // First frame
    files = append(files, ProcessedFile{
        ID:       fileID,
        Variant:  "thumbnail",
        Data:     thumbnail,
        FileName: fmt.Sprintf("%s_thumb.jpg", fileID),
    })

    return files
}

// Progressive media loading
func (ms *MediaService) GetMediaWithProgressive(ctx context.Context, fileID string, quality string) (*MediaResponse, error) {
    mediaFile, err := ms.mediaRepo.GetByID(ctx, fileID)
    if err != nil {
        return nil, err
    }

    var cdnPath string
    switch quality {
    case "thumbnail":
        cdnPath = mediaFile.ThumbnailPath
    case "compressed":
        cdnPath = mediaFile.CompressedPath
    case "original":
        cdnPath = mediaFile.CDNPath
    default:
        cdnPath = mediaFile.CompressedPath // Default to compressed
    }

    return &MediaResponse{
        FileID:   fileID,
        URL:      cdnPath,
        FileSize: mediaFile.FileSize,
        MimeType: mediaFile.MimeType,
    }, nil
}
```

## 🔍 Step 9: Message Search and Indexing

### Search Service Implementation
```go
type SearchService struct {
    searchIndex     SearchIndex // Elasticsearch
    messageService  MessageService
    encryptionService EncryptionService
}

func (ss *SearchService) IndexMessage(ctx context.Context, message *Message) error {
    // Decrypt message for indexing (server-side search)
    decryptedContent, err := ss.encryptionService.DecryptForIndexing(message.Content)
    if err != nil {
        return err // Skip indexing if decryption fails
    }

    // Create search document
    doc := SearchDocument{
        MessageID:   message.ID,
        ChatID:      message.ChatID,
        SenderID:    message.SenderID,
        Content:     decryptedContent,
        MessageType: message.Type,
        Timestamp:   message.CreatedAt,
    }

    return ss.searchIndex.Index(ctx, doc)
}

func (ss *SearchService) SearchMessages(ctx context.Context, userID int64, query SearchQuery) (*SearchResults, error) {
    // Get user's accessible chats
    accessibleChats, err := ss.getUserAccessibleChats(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Build Elasticsearch query
    esQuery := map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "multi_match": map[string]interface{}{
                        "query":  query.Text,
                        "fields": []string{"content^2", "sender_name"},
                    },
                },
                {
                    "terms": map[string]interface{}{
                        "chat_id": accessibleChats,
                    },
                },
            },
            "filter": ss.buildFilters(query),
        },
    }

    // Add sorting
    sort := []map[string]interface{}{
        {"timestamp": map[string]string{"order": "desc"}},
    }

    // Execute search
    results, err := ss.searchIndex.Search(ctx, SearchRequest{
        Query: esQuery,
        Sort:  sort,
        From:  query.Offset,
        Size:  query.Limit,
    })
    if err != nil {
        return nil, err
    }

    // Convert to message objects
    var messages []*Message
    for _, hit := range results.Hits {
        message, err := ss.messageService.GetMessage(ctx, hit.MessageID)
        if err != nil {
            continue
        }
        messages = append(messages, message)
    }

    return &SearchResults{
        Messages:   messages,
        Total:      results.Total,
        HasMore:    results.Total > query.Offset+query.Limit,
        SearchTime: results.SearchTime,
    }, nil
}

func (ss *SearchService) buildFilters(query SearchQuery) []map[string]interface{} {
    var filters []map[string]interface{}

    // Date range filter
    if !query.StartDate.IsZero() || !query.EndDate.IsZero() {
        dateRange := map[string]interface{}{}
        if !query.StartDate.IsZero() {
            dateRange["gte"] = query.StartDate
        }
        if !query.EndDate.IsZero() {
            dateRange["lte"] = query.EndDate
        }

        filters = append(filters, map[string]interface{}{
            "range": map[string]interface{}{
                "timestamp": dateRange,
            },
        })
    }

    // Message type filter
    if query.MessageType != "" {
        filters = append(filters, map[string]interface{}{
            "term": map[string]interface{}{
                "message_type": query.MessageType,
            },
        })
    }

    // Sender filter
    if query.SenderID != 0 {
        filters = append(filters, map[string]interface{}{
            "term": map[string]interface{}{
                "sender_id": query.SenderID,
            },
        })
    }

    return filters
}
```

## 👥 Step 10: Presence and Typing Indicators

### Presence Service Implementation
```go
type PresenceService struct {
    presenceRepo    PresenceRepository
    cache           Cache
    connectionManager ConnectionManager
    eventBus        EventBus
}

func (ps *PresenceService) SetUserOnline(userID int64, connectionID, deviceType string) {
    presence := &UserPresence{
        UserID:       userID,
        Status:       "online",
        LastSeen:     time.Now(),
        DeviceType:   deviceType,
        ConnectionID: connectionID,
        UpdatedAt:    time.Now(),
    }

    // Update cache
    cacheKey := fmt.Sprintf("presence:%d", userID)
    ps.cache.Set(cacheKey, presence, 5*time.Minute)

    // Update database (async)
    go ps.presenceRepo.Upsert(context.Background(), presence)

    // Notify contacts
    ps.notifyContactsOfPresenceChange(userID, "online")
}

func (ps *PresenceService) SetUserOffline(userID int64, connectionID string) {
    // Check if user has other active connections
    activeConnections := ps.connectionManager.GetUserConnections(userID)

    if len(activeConnections) <= 1 { // Only this connection
        presence := &UserPresence{
            UserID:    userID,
            Status:    "offline",
            LastSeen:  time.Now(),
            UpdatedAt: time.Now(),
        }

        // Update cache
        cacheKey := fmt.Sprintf("presence:%d", userID)
        ps.cache.Set(cacheKey, presence, 24*time.Hour)

        // Update database (async)
        go ps.presenceRepo.Upsert(context.Background(), presence)

        // Notify contacts
        ps.notifyContactsOfPresenceChange(userID, "offline")
    }
}

func (ps *PresenceService) GetUserPresence(userID int64) (*UserPresence, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("presence:%d", userID)
    if cached, found := ps.cache.Get(cacheKey); found {
        return cached.(*UserPresence), nil
    }

    // Fallback to database
    return ps.presenceRepo.GetByUserID(context.Background(), userID)
}

// Typing indicators
type TypingIndicatorService struct {
    cache           Cache
    connectionManager ConnectionManager
    rateLimiter     RateLimiter
}

func (tis *TypingIndicatorService) StartTyping(userID, chatID int64) {
    // Rate limiting to prevent spam
    if !tis.rateLimiter.Allow(fmt.Sprintf("typing:%d", userID)) {
        return
    }

    // Set typing status
    typingKey := fmt.Sprintf("typing:%d:%d", chatID, userID)
    tis.cache.Set(typingKey, true, 5*time.Second)

    // Notify other chat members
    tis.notifyTypingStatus(chatID, userID, true)

    // Auto-stop typing after 5 seconds
    time.AfterFunc(5*time.Second, func() {
        tis.StopTyping(userID, chatID)
    })
}

func (tis *TypingIndicatorService) StopTyping(userID, chatID int64) {
    typingKey := fmt.Sprintf("typing:%d:%d", chatID, userID)
    tis.cache.Delete(typingKey)

    // Notify other chat members
    tis.notifyTypingStatus(chatID, userID, false)
}

func (tis *TypingIndicatorService) GetTypingUsers(chatID int64) []int64 {
    pattern := fmt.Sprintf("typing:%d:*", chatID)
    keys := tis.cache.GetKeysByPattern(pattern)

    var typingUsers []int64
    for _, key := range keys {
        parts := strings.Split(key, ":")
        if len(parts) == 3 {
            if userID, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
                typingUsers = append(typingUsers, userID)
            }
        }
    }

    return typingUsers
}

func (tis *TypingIndicatorService) notifyTypingStatus(chatID, userID int64, isTyping bool) {
    // Get chat members
    members := tis.getChatMembers(chatID)

    notification := TypingNotification{
        ChatID:    chatID,
        UserID:    userID,
        IsTyping:  isTyping,
        Timestamp: time.Now(),
    }

    notificationData, _ := json.Marshal(notification)

    // Send to all members except the typer
    for _, memberID := range members {
        if memberID != userID {
            connections := tis.connectionManager.GetUserConnections(memberID)
            for _, conn := range connections {
                select {
                case conn.SendChannel <- notificationData:
                default:
                    // Channel full, skip
                }
            }
        }
    }
}
```

## 🎯 Summary: WhatsApp System Design

### Architecture Highlights
- **Real-time Messaging**: WebSocket connections with message delivery guarantees
- **End-to-End Encryption**: Signal Protocol implementation for secure messaging
- **Scalable Storage**: Sharded databases with Cassandra for high write volume
- **Smart Notifications**: Batched push notifications with user preferences
- **Media Processing**: Multi-variant media with progressive loading
- **Advanced Search**: Elasticsearch integration with encrypted content indexing
- **Presence System**: Real-time online status and typing indicators

### Scalability Features
- **2B users**: Heavily sharded databases and distributed architecture
- **100B messages/day**: Event-driven architecture with message queues
- **50M concurrent connections**: Connection pooling and WebSocket optimization
- **Global deployment**: Multi-region with data locality
- **99.99% availability**: Redundancy and failover mechanisms

### Performance Optimizations
- **Message delivery**: < 100ms with direct WebSocket routing
- **Message history**: < 200ms with optimized database queries
- **Media loading**: Progressive loading with CDN
- **Search**: < 500ms with Elasticsearch indexing
- **Presence updates**: Real-time with Redis caching

This design handles WhatsApp-scale messaging with optimal real-time performance and security! 🚀
