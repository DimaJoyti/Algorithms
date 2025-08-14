# 🔐 Security & Authentication in Distributed Systems

## 🎯 Authentication vs Authorization

### Core Concepts
```go
// Authentication: "Who are you?"
type AuthenticationResult struct {
    UserID      string
    Email       string
    Roles       []string
    Permissions []string
    TokenType   string
    ExpiresAt   time.Time
}

// Authorization: "What can you do?"
type AuthorizationContext struct {
    UserID      string
    Roles       []string
    Permissions []string
    Resource    string
    Action      string
    Context     map[string]interface{}
}

func (ac *AuthorizationContext) IsAuthorized() bool {
    // Role-based access control (RBAC)
    for _, role := range ac.Roles {
        if ac.hasRolePermission(role, ac.Resource, ac.Action) {
            return true
        }
    }
    
    // Permission-based access control
    requiredPermission := fmt.Sprintf("%s:%s", ac.Resource, ac.Action)
    for _, permission := range ac.Permissions {
        if permission == requiredPermission || permission == "*" {
            return true
        }
    }
    
    return false
}
```

## 🎫 JWT (JSON Web Tokens)

### JWT Implementation
```go
type JWTService struct {
    secretKey     []byte
    issuer        string
    expiration    time.Duration
    refreshExpiry time.Duration
}

type Claims struct {
    UserID      string   `json:"user_id"`
    Email       string   `json:"email"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
    TokenType   string   `json:"token_type"` // "access" or "refresh"
    jwt.RegisteredClaims
}

func (js *JWTService) GenerateTokenPair(userID, email string, roles []string) (*TokenPair, error) {
    // Access token (short-lived)
    accessClaims := &Claims{
        UserID:      userID,
        Email:       email,
        Roles:       roles,
        Permissions: js.getRolePermissions(roles),
        TokenType:   "access",
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    js.issuer,
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(js.secretKey)
    if err != nil {
        return nil, err
    }
    
    // Refresh token (long-lived)
    refreshClaims := &Claims{
        UserID:    userID,
        TokenType: "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    js.issuer,
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.refreshExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(js.secretKey)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int(js.expiration.Seconds()),
        TokenType:    "Bearer",
    }, nil
}

func (js *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return js.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, ErrInvalidToken
}

// JWT Middleware
func (js *JWTService) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Bearer token required", http.StatusUnauthorized)
            return
        }
        
        claims, err := js.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        if claims.TokenType != "access" {
            http.Error(w, "Access token required", http.StatusUnauthorized)
            return
        }
        
        // Add claims to context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "user_email", claims.Email)
        ctx = context.WithValue(ctx, "user_roles", claims.Roles)
        ctx = context.WithValue(ctx, "user_permissions", claims.Permissions)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Token Refresh Flow
```go
func (js *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
    claims, err := js.ValidateToken(refreshTokenString)
    if err != nil {
        return nil, err
    }
    
    if claims.TokenType != "refresh" {
        return nil, ErrInvalidTokenType
    }
    
    // Get fresh user data
    user, err := js.userService.GetUser(claims.UserID)
    if err != nil {
        return nil, err
    }
    
    if !user.IsActive {
        return nil, ErrUserInactive
    }
    
    // Generate new token pair
    return js.GenerateTokenPair(user.ID, user.Email, user.Roles)
}

// Token blacklist for logout
type TokenBlacklist struct {
    blacklistedTokens map[string]time.Time
    mutex             sync.RWMutex
}

func (tb *TokenBlacklist) BlacklistToken(tokenString string, expiresAt time.Time) {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    
    tb.blacklistedTokens[tokenString] = expiresAt
}

func (tb *TokenBlacklist) IsBlacklisted(tokenString string) bool {
    tb.mutex.RLock()
    defer tb.mutex.RUnlock()
    
    expiresAt, exists := tb.blacklistedTokens[tokenString]
    if !exists {
        return false
    }
    
    // Clean up expired tokens
    if time.Now().After(expiresAt) {
        delete(tb.blacklistedTokens, tokenString)
        return false
    }
    
    return true
}
```

## 🔑 OAuth 2.0 Implementation

### OAuth 2.0 Authorization Server
```go
type OAuthServer struct {
    clients       map[string]*OAuthClient
    authCodes     map[string]*AuthorizationCode
    accessTokens  map[string]*AccessToken
    refreshTokens map[string]*RefreshToken
    userService   UserService
    mutex         sync.RWMutex
}

type OAuthClient struct {
    ID           string
    Secret       string
    RedirectURIs []string
    Scopes       []string
    GrantTypes   []string
}

type AuthorizationCode struct {
    Code        string
    ClientID    string
    UserID      string
    RedirectURI string
    Scopes      []string
    ExpiresAt   time.Time
}

// Authorization endpoint
func (os *OAuthServer) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
    clientID := r.URL.Query().Get("client_id")
    redirectURI := r.URL.Query().Get("redirect_uri")
    responseType := r.URL.Query().Get("response_type")
    scopes := strings.Split(r.URL.Query().Get("scope"), " ")
    state := r.URL.Query().Get("state")
    
    // Validate client
    client, exists := os.clients[clientID]
    if !exists {
        http.Error(w, "Invalid client", http.StatusBadRequest)
        return
    }
    
    // Validate redirect URI
    if !contains(client.RedirectURIs, redirectURI) {
        http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
        return
    }
    
    // Only support authorization code flow
    if responseType != "code" {
        http.Redirect(w, r, fmt.Sprintf("%s?error=unsupported_response_type&state=%s", redirectURI, state), http.StatusFound)
        return
    }
    
    // Get authenticated user (from session)
    userID := os.getAuthenticatedUser(r)
    if userID == "" {
        // Redirect to login
        loginURL := fmt.Sprintf("/login?redirect=%s", url.QueryEscape(r.URL.String()))
        http.Redirect(w, r, loginURL, http.StatusFound)
        return
    }
    
    // Generate authorization code
    code := &AuthorizationCode{
        Code:        generateRandomString(32),
        ClientID:    clientID,
        UserID:      userID,
        RedirectURI: redirectURI,
        Scopes:      scopes,
        ExpiresAt:   time.Now().Add(10 * time.Minute),
    }
    
    os.mutex.Lock()
    os.authCodes[code.Code] = code
    os.mutex.Unlock()
    
    // Redirect with authorization code
    redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code.Code, state)
    http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Token endpoint
func (os *OAuthServer) TokenHandler(w http.ResponseWriter, r *http.Request) {
    grantType := r.FormValue("grant_type")
    
    switch grantType {
    case "authorization_code":
        os.handleAuthorizationCodeGrant(w, r)
    case "refresh_token":
        os.handleRefreshTokenGrant(w, r)
    case "client_credentials":
        os.handleClientCredentialsGrant(w, r)
    default:
        os.writeTokenError(w, "unsupported_grant_type", "Grant type not supported")
    }
}

func (os *OAuthServer) handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request) {
    code := r.FormValue("code")
    clientID := r.FormValue("client_id")
    clientSecret := r.FormValue("client_secret")
    redirectURI := r.FormValue("redirect_uri")
    
    // Validate client credentials
    client, exists := os.clients[clientID]
    if !exists || client.Secret != clientSecret {
        os.writeTokenError(w, "invalid_client", "Invalid client credentials")
        return
    }
    
    // Validate authorization code
    os.mutex.RLock()
    authCode, exists := os.authCodes[code]
    os.mutex.RUnlock()
    
    if !exists || authCode.ClientID != clientID || authCode.RedirectURI != redirectURI {
        os.writeTokenError(w, "invalid_grant", "Invalid authorization code")
        return
    }
    
    if time.Now().After(authCode.ExpiresAt) {
        os.writeTokenError(w, "invalid_grant", "Authorization code expired")
        return
    }
    
    // Generate tokens
    accessToken := &AccessToken{
        Token:     generateRandomString(32),
        ClientID:  clientID,
        UserID:    authCode.UserID,
        Scopes:    authCode.Scopes,
        ExpiresAt: time.Now().Add(1 * time.Hour),
    }
    
    refreshToken := &RefreshToken{
        Token:     generateRandomString(32),
        ClientID:  clientID,
        UserID:    authCode.UserID,
        Scopes:    authCode.Scopes,
        ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
    }
    
    // Store tokens
    os.mutex.Lock()
    os.accessTokens[accessToken.Token] = accessToken
    os.refreshTokens[refreshToken.Token] = refreshToken
    delete(os.authCodes, code) // Use authorization code only once
    os.mutex.Unlock()
    
    // Return token response
    response := map[string]interface{}{
        "access_token":  accessToken.Token,
        "token_type":    "Bearer",
        "expires_in":    3600,
        "refresh_token": refreshToken.Token,
        "scope":         strings.Join(authCode.Scopes, " "),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### OAuth 2.0 Resource Server
```go
type ResourceServer struct {
    authServer *OAuthServer
    scopes     map[string][]string // endpoint -> required scopes
}

func (rs *ResourceServer) RequireScope(requiredScopes ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                rs.writeError(w, http.StatusUnauthorized, "access_denied", "Authorization header required")
                return
            }
            
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == authHeader {
                rs.writeError(w, http.StatusUnauthorized, "invalid_token", "Bearer token required")
                return
            }
            
            // Validate access token
            rs.authServer.mutex.RLock()
            accessToken, exists := rs.authServer.accessTokens[tokenString]
            rs.authServer.mutex.RUnlock()
            
            if !exists {
                rs.writeError(w, http.StatusUnauthorized, "invalid_token", "Invalid access token")
                return
            }
            
            if time.Now().After(accessToken.ExpiresAt) {
                rs.writeError(w, http.StatusUnauthorized, "invalid_token", "Access token expired")
                return
            }
            
            // Check scopes
            if !rs.hasRequiredScopes(accessToken.Scopes, requiredScopes) {
                rs.writeError(w, http.StatusForbidden, "insufficient_scope", "Insufficient scope")
                return
            }
            
            // Add token info to context
            ctx := context.WithValue(r.Context(), "user_id", accessToken.UserID)
            ctx = context.WithValue(ctx, "client_id", accessToken.ClientID)
            ctx = context.WithValue(ctx, "scopes", accessToken.Scopes)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func (rs *ResourceServer) hasRequiredScopes(tokenScopes, requiredScopes []string) bool {
    for _, required := range requiredScopes {
        found := false
        for _, tokenScope := range tokenScopes {
            if tokenScope == required {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}
```

## 🛡️ API Security Best Practices

### Rate Limiting by User/IP
```go
type SecurityMiddleware struct {
    rateLimiter   RateLimiter
    ipWhitelist   map[string]bool
    ipBlacklist   map[string]bool
    bruteForceProtection BruteForceProtection
}

func (sm *SecurityMiddleware) SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        next.ServeHTTP(w, r)
    })
}

func (sm *SecurityMiddleware) IPFiltering(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientIP := getClientIP(r)
        
        // Check blacklist
        if sm.ipBlacklist[clientIP] {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        
        // Check whitelist (if configured)
        if len(sm.ipWhitelist) > 0 && !sm.ipWhitelist[clientIP] {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func (sm *SecurityMiddleware) BruteForceProtection(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientIP := getClientIP(r)
        
        if sm.bruteForceProtection.IsBlocked(clientIP) {
            http.Error(w, "Too many failed attempts", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

type BruteForceProtection struct {
    attempts map[string]*AttemptTracker
    mutex    sync.RWMutex
}

type AttemptTracker struct {
    Count     int
    LastAttempt time.Time
    BlockedUntil time.Time
}

func (bfp *BruteForceProtection) RecordFailedAttempt(identifier string) {
    bfp.mutex.Lock()
    defer bfp.mutex.Unlock()
    
    tracker, exists := bfp.attempts[identifier]
    if !exists {
        tracker = &AttemptTracker{}
        bfp.attempts[identifier] = tracker
    }
    
    tracker.Count++
    tracker.LastAttempt = time.Now()
    
    // Progressive blocking
    if tracker.Count >= 5 {
        blockDuration := time.Duration(tracker.Count-4) * 5 * time.Minute
        tracker.BlockedUntil = time.Now().Add(blockDuration)
    }
}

func (bfp *BruteForceProtection) IsBlocked(identifier string) bool {
    bfp.mutex.RLock()
    defer bfp.mutex.RUnlock()
    
    tracker, exists := bfp.attempts[identifier]
    if !exists {
        return false
    }
    
    return time.Now().Before(tracker.BlockedUntil)
}
```

### Input Validation & Sanitization
```go
type Validator struct {
    emailRegex *regexp.Regexp
    phoneRegex *regexp.Regexp
}

func NewValidator() *Validator {
    return &Validator{
        emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
        phoneRegex: regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
    }
}

func (v *Validator) ValidateEmail(email string) error {
    if len(email) > 254 {
        return ErrEmailTooLong
    }
    
    if !v.emailRegex.MatchString(email) {
        return ErrInvalidEmailFormat
    }
    
    return nil
}

func (v *Validator) SanitizeInput(input string) string {
    // Remove potentially dangerous characters
    input = strings.ReplaceAll(input, "<", "&lt;")
    input = strings.ReplaceAll(input, ">", "&gt;")
    input = strings.ReplaceAll(input, "\"", "&quot;")
    input = strings.ReplaceAll(input, "'", "&#x27;")
    input = strings.ReplaceAll(input, "&", "&amp;")
    
    return strings.TrimSpace(input)
}

// SQL Injection Prevention
func (v *Validator) ValidateAndSanitizeSQL(query string, args ...interface{}) error {
    // Use parameterized queries - never concatenate user input
    dangerousPatterns := []string{
        "DROP TABLE",
        "DELETE FROM",
        "UPDATE.*SET",
        "INSERT INTO",
        "UNION SELECT",
        "--",
        "/*",
        "*/",
        "xp_",
        "sp_",
    }
    
    upperQuery := strings.ToUpper(query)
    for _, pattern := range dangerousPatterns {
        if matched, _ := regexp.MatchString(pattern, upperQuery); matched {
            return ErrSuspiciousQuery
        }
    }
    
    return nil
}
```

## 🔒 Encryption & Data Protection

### Data Encryption at Rest
```go
type EncryptionService struct {
    key []byte
}

func NewEncryptionService(key []byte) *EncryptionService {
    return &EncryptionService{key: key}
}

func (es *EncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(es.key)
    if err != nil {
        return nil, err
    }
    
    // Generate random IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    
    // Encrypt using CBC mode
    ciphertext := make([]byte, len(plaintext))
    mode := cipher.NewCBCEncrypter(block, iv)
    
    // Pad plaintext to block size
    paddedPlaintext := pkcs7Pad(plaintext, aes.BlockSize)
    mode.CryptBlocks(ciphertext, paddedPlaintext)
    
    // Prepend IV to ciphertext
    result := make([]byte, len(iv)+len(ciphertext))
    copy(result[:len(iv)], iv)
    copy(result[len(iv):], ciphertext)
    
    return result, nil
}

func (es *EncryptionService) Decrypt(ciphertext []byte) ([]byte, error) {
    if len(ciphertext) < aes.BlockSize {
        return nil, ErrInvalidCiphertext
    }
    
    block, err := aes.NewCipher(es.key)
    if err != nil {
        return nil, err
    }
    
    // Extract IV
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    
    // Decrypt using CBC mode
    mode := cipher.NewCBCDecrypter(block, iv)
    plaintext := make([]byte, len(ciphertext))
    mode.CryptBlocks(plaintext, ciphertext)
    
    // Remove padding
    return pkcs7Unpad(plaintext)
}

// Password hashing
func (es *EncryptionService) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func (es *EncryptionService) VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## 🔗 Next Steps
- Study monitoring and observability with distributed tracing
- Learn about advanced security patterns like zero-trust architecture
- Explore compliance and regulatory requirements (GDPR, HIPAA, etc.)
