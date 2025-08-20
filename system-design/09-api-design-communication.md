# 🔌 API Design & Communication Patterns

## 🎯 API Design Principles

### RESTful API Design
```go
// Resource-based URLs
type UserAPI struct {
    userService UserService
}

// Good RESTful design
func (api *UserAPI) SetupRoutes(mux *http.ServeMux) {
    // Users collection
    mux.HandleFunc("GET /api/v1/users", api.listUsers)           // List users
    mux.HandleFunc("POST /api/v1/users", api.createUser)        // Create user
    
    // Individual user resource
    mux.HandleFunc("GET /api/v1/users/{id}", api.getUser)       // Get user
    mux.HandleFunc("PUT /api/v1/users/{id}", api.updateUser)    // Update user
    mux.HandleFunc("DELETE /api/v1/users/{id}", api.deleteUser) // Delete user
    
    // Nested resources
    mux.HandleFunc("GET /api/v1/users/{id}/posts", api.getUserPosts)
    mux.HandleFunc("POST /api/v1/users/{id}/posts", api.createUserPost)
    
    // Actions on resources (when CRUD isn't enough)
    mux.HandleFunc("POST /api/v1/users/{id}/follow", api.followUser)
    mux.HandleFunc("DELETE /api/v1/users/{id}/follow", api.unfollowUser)
}

func (api *UserAPI) createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Validate request
    if err := req.Validate(); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }
    
    // Create user
    user, err := api.userService.CreateUser(r.Context(), req)
    if err != nil {
        if errors.Is(err, ErrUserAlreadyExists) {
            writeErrorResponse(w, http.StatusConflict, "User already exists")
            return
        }
        writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
        return
    }
    
    // Return created resource
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Location", fmt.Sprintf("/api/v1/users/%s", user.ID))
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Consistent error response format
type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Code    string            `json:"code,omitempty"`
    Details map[string]string `json:"details,omitempty"`
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := ErrorResponse{
        Error:   http.StatusText(statusCode),
        Message: message,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

### API Versioning Strategies
```go
// 1. URL Path Versioning (Recommended)
mux.HandleFunc("GET /api/v1/users/{id}", api.getUserV1)
mux.HandleFunc("GET /api/v2/users/{id}", api.getUserV2)

// 2. Header Versioning
func (api *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {
    version := r.Header.Get("API-Version")
    
    switch version {
    case "v2":
        api.getUserV2(w, r)
    default:
        api.getUserV1(w, r)
    }
}

// 3. Query Parameter Versioning
// GET /api/users/123?version=v2

// Version compatibility handling
type UserResponseV1 struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponseV2 struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Profile  Profile   `json:"profile"`
    Settings Settings  `json:"settings"`
}

func (api *UserAPI) getUserV1(w http.ResponseWriter, r *http.Request) {
    user, err := api.userService.GetUser(r.Context(), getUserID(r))
    if err != nil {
        // Handle error
        return
    }
    
    // Convert to V1 format
    response := UserResponseV1{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

### Pagination & Filtering
```go
type PaginationParams struct {
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    Cursor   string `json:"cursor,omitempty"`
}

type FilterParams struct {
    Name     string `json:"name,omitempty"`
    Email    string `json:"email,omitempty"`
    Status   string `json:"status,omitempty"`
    SortBy   string `json:"sort_by,omitempty"`
    SortDir  string `json:"sort_dir,omitempty"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination struct {
        Page       int    `json:"page"`
        PageSize   int    `json:"page_size"`
        Total      int64  `json:"total"`
        TotalPages int    `json:"total_pages"`
        NextCursor string `json:"next_cursor,omitempty"`
        PrevCursor string `json:"prev_cursor,omitempty"`
    } `json:"pagination"`
}

func (api *UserAPI) listUsers(w http.ResponseWriter, r *http.Request) {
    // Parse pagination parameters
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    
    // Parse filter parameters
    filters := FilterParams{
        Name:    r.URL.Query().Get("name"),
        Email:   r.URL.Query().Get("email"),
        Status:  r.URL.Query().Get("status"),
        SortBy:  r.URL.Query().Get("sort_by"),
        SortDir: r.URL.Query().Get("sort_dir"),
    }
    
    // Get users with pagination
    users, total, err := api.userService.ListUsers(r.Context(), PaginationParams{
        Page:     page,
        PageSize: pageSize,
    }, filters)
    
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
        return
    }
    
    // Build response
    response := PaginatedResponse{
        Data: users,
    }
    
    response.Pagination.Page = page
    response.Pagination.PageSize = pageSize
    response.Pagination.Total = total
    response.Pagination.TotalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
    
    json.NewEncoder(w).Encode(response)
}

// Cursor-based pagination for large datasets
func (api *UserAPI) listUsersWithCursor(w http.ResponseWriter, r *http.Request) {
    cursor := r.URL.Query().Get("cursor")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    users, nextCursor, err := api.userService.ListUsersWithCursor(r.Context(), cursor, limit)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
        return
    }
    
    response := struct {
        Data       []User `json:"data"`
        NextCursor string `json:"next_cursor,omitempty"`
        HasMore    bool   `json:"has_more"`
    }{
        Data:       users,
        NextCursor: nextCursor,
        HasMore:    nextCursor != "",
    }
    
    json.NewEncoder(w).Encode(response)
}
```

## 🚀 GraphQL Implementation
```go
type GraphQLResolver struct {
    userService UserService
    postService PostService
}

// Schema definition
const schema = `
    type User {
        id: ID!
        name: String!
        email: String!
        posts: [Post!]!
    }
    
    type Post {
        id: ID!
        title: String!
        content: String!
        author: User!
        createdAt: String!
    }
    
    type Query {
        user(id: ID!): User
        users(first: Int, after: String): UserConnection
        post(id: ID!): Post
    }
    
    type Mutation {
        createUser(input: CreateUserInput!): User!
        updateUser(id: ID!, input: UpdateUserInput!): User!
        createPost(input: CreatePostInput!): Post!
    }
    
    input CreateUserInput {
        name: String!
        email: String!
    }
`

func (r *GraphQLResolver) User(ctx context.Context, args struct{ ID string }) (*User, error) {
    return r.userService.GetUser(ctx, args.ID)
}

func (r *GraphQLResolver) Users(ctx context.Context, args struct {
    First *int32
    After *string
}) (*UserConnection, error) {
    limit := 20
    if args.First != nil {
        limit = int(*args.First)
    }
    
    cursor := ""
    if args.After != nil {
        cursor = *args.After
    }
    
    users, nextCursor, err := r.userService.ListUsersWithCursor(ctx, cursor, limit)
    if err != nil {
        return nil, err
    }
    
    return &UserConnection{
        Edges: users,
        PageInfo: PageInfo{
            HasNextPage: nextCursor != "",
            EndCursor:   nextCursor,
        },
    }, nil
}

// Resolver for nested fields (N+1 problem solution with DataLoader)
func (u *User) Posts(ctx context.Context) ([]*Post, error) {
    // Use DataLoader to batch requests
    loader := ctx.Value("postLoader").(*PostLoader)
    return loader.LoadByUserID(u.ID)
}

// DataLoader implementation to solve N+1 problem
type PostLoader struct {
    postService PostService
    cache       map[string][]*Post
    batch       []string
    mutex       sync.Mutex
}

func (pl *PostLoader) LoadByUserID(userID string) ([]*Post, error) {
    pl.mutex.Lock()
    defer pl.mutex.Unlock()
    
    // Check cache first
    if posts, found := pl.cache[userID]; found {
        return posts, nil
    }
    
    // Add to batch
    pl.batch = append(pl.batch, userID)
    
    // If batch is full or timeout, execute batch
    if len(pl.batch) >= 10 {
        return pl.executeBatch()
    }
    
    // Set timer for batch execution
    time.AfterFunc(10*time.Millisecond, func() {
        pl.executeBatch()
    })
    
    return pl.cache[userID], nil
}

func (pl *PostLoader) executeBatch() ([]*Post, error) {
    if len(pl.batch) == 0 {
        return nil, nil
    }
    
    // Fetch all posts for batched user IDs
    postsByUser, err := pl.postService.GetPostsByUserIDs(pl.batch)
    if err != nil {
        return nil, err
    }
    
    // Cache results
    for userID, posts := range postsByUser {
        pl.cache[userID] = posts
    }
    
    // Clear batch
    pl.batch = pl.batch[:0]
    
    return nil, nil
}
```

## 🔄 gRPC Implementation
```protobuf
// user.proto
syntax = "proto3";

package user.v1;

option go_package = "github.com/yourorg/userservice/proto/user/v1";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    
    // Streaming RPCs
    rpc StreamUsers(StreamUsersRequest) returns (stream User);
    rpc BulkCreateUsers(stream CreateUserRequest) returns (BulkCreateUsersResponse);
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    google.protobuf.Timestamp created_at = 4;
    Profile profile = 5;
}

message Profile {
    string bio = 1;
    string avatar_url = 2;
    repeated string interests = 3;
}

message GetUserRequest {
    string id = 1;
}

message GetUserResponse {
    User user = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
    string filter = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
}
```

```go
// gRPC server implementation
type UserServiceServer struct {
    userService UserService
    pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Validate request
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "user ID is required")
    }
    
    // Get user from service
    user, err := s.userService.GetUser(ctx, req.Id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "internal server error")
    }
    
    // Convert to protobuf
    pbUser := &pb.User{
        Id:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        CreatedAt: timestamppb.New(user.CreatedAt),
        Profile: &pb.Profile{
            Bio:       user.Profile.Bio,
            AvatarUrl: user.Profile.AvatarURL,
            Interests: user.Profile.Interests,
        },
    }
    
    return &pb.GetUserResponse{User: pbUser}, nil
}

// Server streaming
func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    users, err := s.userService.ListAllUsers(stream.Context())
    if err != nil {
        return status.Error(codes.Internal, "failed to fetch users")
    }
    
    for _, user := range users {
        pbUser := convertUserToProto(user)
        
        if err := stream.Send(pbUser); err != nil {
            return err
        }
        
        // Rate limiting
        time.Sleep(10 * time.Millisecond)
    }
    
    return nil
}

// Client streaming
func (s *UserServiceServer) BulkCreateUsers(stream pb.UserService_BulkCreateUsersServer) error {
    var createdCount int32
    
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            // Client finished sending
            return stream.SendAndClose(&pb.BulkCreateUsersResponse{
                CreatedCount: createdCount,
            })
        }
        if err != nil {
            return err
        }
        
        // Create user
        _, err = s.userService.CreateUser(stream.Context(), CreateUserRequest{
            Name:  req.Name,
            Email: req.Email,
        })
        
        if err == nil {
            createdCount++
        }
    }
}

// gRPC client with connection pooling
type UserServiceClient struct {
    conn   *grpc.ClientConn
    client pb.UserServiceClient
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
    conn, err := grpc.Dial(address,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             time.Second,
            PermitWithoutStream: true,
        }),
    )
    if err != nil {
        return nil, err
    }
    
    client := pb.NewUserServiceClient(conn)
    
    return &UserServiceClient{
        conn:   conn,
        client: client,
    }, nil
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*User, error) {
    resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
        Id: userID,
    })
    if err != nil {
        return nil, err
    }
    
    return convertProtoToUser(resp.User), nil
}
```

## 🔐 API Security & Rate Limiting

### JWT Authentication Middleware
```go
type JWTMiddleware struct {
    secretKey []byte
}

func (jm *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            writeErrorResponse(w, http.StatusUnauthorized, "Bearer token required")
            return
        }
        
        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jm.secretKey, nil
        })
        
        if err != nil || !token.Valid {
            writeErrorResponse(w, http.StatusUnauthorized, "Invalid token")
            return
        }
        
        // Extract claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID := claims["user_id"].(string)
            userRole := claims["role"].(string)
            
            // Add to context
            ctx := context.WithValue(r.Context(), "user_id", userID)
            ctx = context.WithValue(ctx, "user_role", userRole)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            writeErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
        }
    })
}

// Role-based authorization
func (jm *JWTMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := r.Context().Value("user_role").(string)
            
            if userRole != role && userRole != "admin" {
                writeErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Rate Limiting
```go
type RateLimiter struct {
    clients map[string]*ClientLimiter
    mutex   sync.RWMutex
    limit   int
    window  time.Duration
}

type ClientLimiter struct {
    requests []time.Time
    mutex    sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        clients: make(map[string]*ClientLimiter),
        limit:   limit,
        window:  window,
    }
    
    // Cleanup goroutine
    go rl.cleanup()
    
    return rl
}

func (rl *RateLimiter) Allow(clientID string) bool {
    rl.mutex.RLock()
    client, exists := rl.clients[clientID]
    rl.mutex.RUnlock()
    
    if !exists {
        rl.mutex.Lock()
        client = &ClientLimiter{
            requests: make([]time.Time, 0),
        }
        rl.clients[clientID] = client
        rl.mutex.Unlock()
    }
    
    client.mutex.Lock()
    defer client.mutex.Unlock()
    
    now := time.Now()
    
    // Remove old requests outside the window
    cutoff := now.Add(-rl.window)
    validRequests := 0
    for i, reqTime := range client.requests {
        if reqTime.After(cutoff) {
            client.requests = client.requests[i:]
            validRequests = len(client.requests)
            break
        }
    }
    
    if validRequests == 0 {
        client.requests = client.requests[:0]
    }
    
    // Check if limit exceeded
    if len(client.requests) >= rl.limit {
        return false
    }
    
    // Add current request
    client.requests = append(client.requests, now)
    return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Use IP address as client ID (or extract from JWT)
        clientID := r.RemoteAddr
        if userID := r.Context().Value("user_id"); userID != nil {
            clientID = userID.(string)
        }
        
        if !rl.Allow(clientID) {
            w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
            w.Header().Set("X-RateLimit-Window", rl.window.String())
            writeErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

## 🔗 Next Steps
- Study microservices architecture patterns
- Learn about event-driven architecture
- Explore security and authentication in distributed systems
