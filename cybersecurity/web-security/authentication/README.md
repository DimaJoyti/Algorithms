# 🔐 Authentication and Authorization

Secure implementation of identity verification and access control.

## 📋 Table of Contents

- [Authentication Methods](#authentication-methods)
- [OAuth 2.0](#oauth-20)
- [JWT Implementation](#jwt-implementation)
- [Session Management](#session-management)
- [Authorization Patterns](#authorization-patterns)
- [Interview Questions](#interview-questions)

## 🔑 Authentication Methods

### Multi-Factor Authentication (MFA)

```python
import pyotp
import qrcode
from io import BytesIO
import base64

class MFAHandler:
    def __init__(self, issuer: str = "MyApp"):
        self.issuer = issuer
    
    def generate_secret(self) -> str:
        """Generate TOTP secret for user"""
        return pyotp.random_base32()
    
    def get_qr_code(self, email: str, secret: str) -> str:
        """Generate QR code for authenticator app"""
        uri = pyotp.totp.TOTP(secret).provisioning_uri(
            name=email,
            issuer_name=self.issuer
        )
        
        qr = qrcode.make(uri)
        buffer = BytesIO()
        qr.save(buffer, format="PNG")
        return base64.b64encode(buffer.getvalue()).decode()
    
    def verify_code(self, secret: str, code: str) -> bool:
        """Verify TOTP code"""
        totp = pyotp.TOTP(secret)
        return totp.verify(code, valid_window=1)  # Allow 1 step drift

class BackupCodes:
    @staticmethod
    def generate_codes(count: int = 10) -> list:
        """Generate backup codes"""
        import secrets
        return [secrets.token_hex(4).upper() for _ in range(count)]
    
    @staticmethod
    def hash_codes(codes: list) -> list:
        """Hash backup codes for storage"""
        import bcrypt
        return [bcrypt.hashpw(code.encode(), bcrypt.gensalt()).decode() 
                for code in codes]
    
    @staticmethod
    def verify_code(code: str, hashed_codes: list) -> bool:
        """Verify backup code"""
        import bcrypt
        for hashed in hashed_codes:
            if bcrypt.checkpw(code.encode(), hashed.encode()):
                return True
        return False
```

### Password Policies

```python
import re
from dataclasses import dataclass

@dataclass
class PasswordPolicy:
    min_length: int = 12
    require_uppercase: bool = True
    require_lowercase: bool = True
    require_digits: bool = True
    require_special: bool = True
    max_age_days: int = 90
    history_count: int = 10

class PasswordValidator:
    def __init__(self, policy: PasswordPolicy = None):
        self.policy = policy or PasswordPolicy()
        self.common_passwords = self._load_common_passwords()
    
    def _load_common_passwords(self) -> set:
        """Load list of common/breached passwords"""
        return {
            'password', '123456', 'qwerty', 'letmein',
            'admin', 'welcome', 'monkey', 'dragon'
        }
    
    def validate(self, password: str) -> tuple[bool, list]:
        """Validate password against policy"""
        errors = []
        
        if len(password) < self.policy.min_length:
            errors.append(f"Must be at least {self.policy.min_length} characters")
        
        if self.policy.require_uppercase and not re.search(r'[A-Z]', password):
            errors.append("Must contain uppercase letter")
        
        if self.policy.require_lowercase and not re.search(r'[a-z]', password):
            errors.append("Must contain lowercase letter")
        
        if self.policy.require_digits and not re.search(r'\d', password):
            errors.append("Must contain digit")
        
        if self.policy.require_special and not re.search(r'[!@#$%^&*(),.?":{}|<>]', password):
            errors.append("Must contain special character")
        
        if password.lower() in self.common_passwords:
            errors.append("Password is too common")
        
        return len(errors) == 0, errors
    
    def check_strength(self, password: str) -> int:
        """Calculate password strength score (0-100)"""
        score = 0
        
        # Length
        score += min(len(password) * 4, 40)
        
        # Character variety
        if re.search(r'[a-z]', password):
            score += 10
        if re.search(r'[A-Z]', password):
            score += 10
        if re.search(r'\d', password):
            score += 10
        if re.search(r'[!@#$%^&*(),.?":{}|<>]', password):
            score += 15
        
        # Entropy bonus
        unique_chars = len(set(password))
        score += min(unique_chars * 2, 15)
        
        return min(score, 100)
```

## 🔄 OAuth 2.0

### Authorization Code Flow

```python
from dataclasses import dataclass
from typing import Optional
import requests
import secrets

@dataclass
class OAuthConfig:
    client_id: str
    client_secret: str
    authorize_url: str
    token_url: str
    redirect_uri: str
    scope: str

class OAuthClient:
    def __init__(self, config: OAuthConfig):
        self.config = config
        self.state_store = {}
    
    def get_authorization_url(self, state: str = None) -> str:
        """Generate authorization URL"""
        if not state:
            state = secrets.token_urlsafe(16)
        
        self.state_store[state] = True
        
        params = {
            'client_id': self.config.client_id,
            'redirect_uri': self.config.redirect_uri,
            'response_type': 'code',
            'scope': self.config.scope,
            'state': state
        }
        
        query = '&'.join(f"{k}={v}" for k, v in params.items())
        return f"{self.config.authorize_url}?{query}"
    
    def validate_state(self, state: str) -> bool:
        """Validate OAuth state parameter"""
        return self.state_store.pop(state, False)
    
    def exchange_code(self, code: str) -> dict:
        """Exchange authorization code for tokens"""
        data = {
            'grant_type': 'authorization_code',
            'code': code,
            'redirect_uri': self.config.redirect_uri,
            'client_id': self.config.client_id,
            'client_secret': self.config.client_secret
        }
        
        response = requests.post(
            self.config.token_url,
            data=data,
            headers={'Accept': 'application/json'}
        )
        
        return response.json()
    
    def refresh_token(self, refresh_token: str) -> dict:
        """Refresh access token"""
        data = {
            'grant_type': 'refresh_token',
            'refresh_token': refresh_token,
            'client_id': self.config.client_id,
            'client_secret': self.config.client_secret
        }
        
        response = requests.post(
            self.config.token_url,
            data=data,
            headers={'Accept': 'application/json'}
        )
        
        return response.json()
```

## 🎫 JWT Implementation

### Secure JWT Handler

```python
import jwt
from datetime import datetime, timedelta
from typing import Optional, Any
import secrets

class JWTHandler:
    def __init__(self, secret_key: str, algorithm: str = "HS256"):
        self.secret_key = secret_key
        self.algorithm = algorithm
        self.access_token_expire = 15  # minutes
        self.refresh_token_expire = 7  # days
    
    def create_access_token(self, user_id: int, 
                           additional_claims: dict = None) -> str:
        """Create JWT access token"""
        now = datetime.utcnow()
        payload = {
            'sub': str(user_id),
            'type': 'access',
            'iat': now,
            'exp': now + timedelta(minutes=self.access_token_expire),
            'jti': secrets.token_urlsafe(16)  # Unique token ID
        }
        
        if additional_claims:
            payload.update(additional_claims)
        
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)
    
    def create_refresh_token(self, user_id: int) -> str:
        """Create refresh token"""
        now = datetime.utcnow()
        payload = {
            'sub': str(user_id),
            'type': 'refresh',
            'iat': now,
            'exp': now + timedelta(days=self.refresh_token_expire),
            'jti': secrets.token_urlsafe(16)
        }
        
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)
    
    def decode_token(self, token: str) -> Optional[dict]:
        """Decode and validate JWT"""
        try:
            payload = jwt.decode(
                token, 
                self.secret_key, 
                algorithms=[self.algorithm]
            )
            return payload
        except jwt.ExpiredSignatureError:
            raise Exception("Token has expired")
        except jwt.InvalidTokenError:
            raise Exception("Invalid token")
    
    def create_token_pair(self, user_id: int) -> dict:
        """Create access and refresh token pair"""
        return {
            'access_token': self.create_access_token(user_id),
            'refresh_token': self.create_refresh_token(user_id),
            'token_type': 'Bearer',
            'expires_in': self.access_token_expire * 60
        }

# Token blacklist for revocation
class TokenBlacklist:
    def __init__(self, redis_client):
        self.redis = redis_client
        self.prefix = "blacklist:"
    
    def add(self, jti: str, expires_in: int):
        """Add token to blacklist"""
        self.redis.setex(f"{self.prefix}{jti}", expires_in, "1")
    
    def is_blacklisted(self, jti: str) -> bool:
        """Check if token is blacklisted"""
        return bool(self.redis.get(f"{self.prefix}{jti}"))
```

## 📦 Session Management

### Secure Session Handler

```python
import secrets
from dataclasses import dataclass
from datetime import datetime, timedelta
from typing import Optional

@dataclass
class SessionData:
    session_id: str
    user_id: int
    created_at: datetime
    last_activity: datetime
    ip_address: str
    user_agent: str

class SessionManager:
    def __init__(self, storage, timeout_minutes: int = 30):
        self.storage = storage
        self.timeout = timedelta(minutes=timeout_minutes)
    
    def create_session(self, user_id: int, ip: str, 
                       user_agent: str) -> SessionData:
        """Create new session"""
        session = SessionData(
            session_id=secrets.token_urlsafe(32),
            user_id=user_id,
            created_at=datetime.utcnow(),
            last_activity=datetime.utcnow(),
            ip_address=ip,
            user_agent=user_agent
        )
        
        self.storage.set(
            f"session:{session.session_id}",
            session.__dict__,
            ex=3600  # 1 hour TTL
        )
        
        return session
    
    def get_session(self, session_id: str) -> Optional[SessionData]:
        """Get and validate session"""
        data = self.storage.get(f"session:{session_id}")
        
        if not data:
            return None
        
        session = SessionData(**data)
        
        # Check timeout
        if datetime.utcnow() - session.last_activity > self.timeout:
            self.delete_session(session_id)
            return None
        
        # Update last activity
        session.last_activity = datetime.utcnow()
        self.storage.set(
            f"session:{session_id}",
            session.__dict__,
            ex=3600
        )
        
        return session
    
    def delete_session(self, session_id: str):
        """Delete session"""
        self.storage.delete(f"session:{session_id}")
    
    def delete_user_sessions(self, user_id: int):
        """Delete all sessions for user"""
        for key in self.storage.scan_iter(f"session:*"):
            data = self.storage.get(key)
            if data and data.get('user_id') == user_id:
                self.storage.delete(key)
```

## 🔒 Authorization Patterns

### Role-Based Access Control (RBAC)

```python
from enum import Enum
from functools import wraps
from typing import List, Set

class Role(Enum):
    ADMIN = "admin"
    MANAGER = "manager"
    USER = "user"
    GUEST = "guest"

class Permission(Enum):
    READ_USERS = "read:users"
    WRITE_USERS = "write:users"
    DELETE_USERS = "delete:users"
    READ_REPORTS = "read:reports"
    WRITE_REPORTS = "write:reports"
    ADMIN_ACCESS = "admin:access"

ROLE_PERMISSIONS: dict[Role, Set[Permission]] = {
    Role.ADMIN: {
        Permission.READ_USERS, Permission.WRITE_USERS, Permission.DELETE_USERS,
        Permission.READ_REPORTS, Permission.WRITE_REPORTS, Permission.ADMIN_ACCESS
    },
    Role.MANAGER: {
        Permission.READ_USERS, Permission.WRITE_USERS,
        Permission.READ_REPORTS, Permission.WRITE_REPORTS
    },
    Role.USER: {
        Permission.READ_REPORTS
    },
    Role.GUEST: set()
}

@dataclass
class User:
    id: int
    username: str
    roles: List[Role]
    
    def has_permission(self, permission: Permission) -> bool:
        """Check if user has permission through any role"""
        for role in self.roles:
            if permission in ROLE_PERMISSIONS.get(role, set()):
                return True
        return False
    
    def has_role(self, role: Role) -> bool:
        """Check if user has specific role"""
        return role in self.roles

def require_permission(permission: Permission):
    """Decorator for permission-based access control"""
    def decorator(f):
        @wraps(f)
        def decorated(*args, **kwargs):
            user = get_current_user()
            
            if not user:
                raise Unauthorized("Authentication required")
            
            if not user.has_permission(permission):
                raise Forbidden(f"Permission '{permission.value}' required")
            
            return f(*args, **kwargs)
        return decorated
    return decorator

def require_role(role: Role):
    """Decorator for role-based access control"""
    def decorator(f):
        @wraps(f)
        def decorated(*args, **kwargs):
            user = get_current_user()
            
            if not user:
                raise Unauthorized("Authentication required")
            
            if not user.has_role(role):
                raise Forbidden(f"Role '{role.value}' required")
            
            return f(*args, **kwargs)
        return decorated
    return decorator
```

### Attribute-Based Access Control (ABAC)

```python
from dataclasses import dataclass
from typing import Any, Dict

@dataclass
class AccessRequest:
    subject: Dict[str, Any]  # User attributes
    resource: Dict[str, Any]  # Resource attributes
    action: str  # Action to perform
    environment: Dict[str, Any]  # Context (time, IP, etc.)

class ABACEngine:
    def __init__(self):
        self.policies = []
    
    def add_policy(self, condition: callable, effect: str):
        """Add access policy"""
        self.policies.append({
            'condition': condition,
            'effect': effect  # 'allow' or 'deny'
        })
    
    def evaluate(self, request: AccessRequest) -> bool:
        """Evaluate access request against policies"""
        for policy in self.policies:
            if policy['condition'](request):
                return policy['effect'] == 'allow'
        
        return False  # Default deny

# Example policies
engine = ABACEngine()

# Policy: Users can edit their own documents
engine.add_policy(
    lambda r: (r.action == 'edit' and 
               r.subject['id'] == r.resource['owner_id']),
    'allow'
)

# Policy: Managers can edit documents in their department
engine.add_policy(
    lambda r: (r.action == 'edit' and
               r.subject['role'] == 'manager' and
               r.subject['department'] == r.resource['department']),
    'allow'
)

# Policy: No access outside business hours
engine.add_policy(
    lambda r: (r.environment['hour'] < 9 or 
               r.environment['hour'] > 17),
    'deny'
)
```

## ❓ Interview Questions

### Basic
**Q: What's the difference between authentication and authorization?**
A: Authentication verifies who you are (identity). Authorization determines what you can do (permissions).

**Q: Why use short-lived access tokens with refresh tokens?**
A: Reduces exposure window if token is compromised. Refresh tokens are stored securely and only used to get new access tokens.

### Intermediate
**Q: Explain OAuth 2.0 authorization code flow.**
A: (1) User redirected to auth server, (2) User authenticates and consents, (3) Auth server redirects with code, (4) Client exchanges code for tokens (server-side).

**Q: When would you use RBAC vs ABAC?**
A: RBAC when permissions are primarily role-based and static. ABAC when access depends on dynamic attributes (user, resource, environment context).

### Advanced
**Q: How do you implement secure session management?**
A: (1) Generate cryptographically random session IDs, (2) Store server-side, (3) Set appropriate timeouts, (4) Regenerate on authentication, (5) Validate IP/User-Agent, (6) Secure cookie flags (HttpOnly, Secure, SameSite).