# 🔒 Security Patterns

Reusable solutions to common security problems in software architecture.

## 📋 Table of Contents

- [Authentication Patterns](#authentication-patterns)
- [Authorization Patterns](#authorization-patterns)
- [Data Protection Patterns](#data-protection-patterns)
- [Network Security Patterns](#network-security-patterns)
- [Logging and Monitoring Patterns](#logging-and-monitoring-patterns)
- [Interview Questions](#interview-questions)

## 🔑 Authentication Patterns

### 1. Token Validation Pattern

```python
from functools import wraps
from flask import request, g
import jwt

class TokenValidator:
    def __init__(self, secret_key: str, algorithm: str = "HS256"):
        self.secret_key = secret_key
        self.algorithm = algorithm
        self.blacklist = TokenBlacklist()
    
    def validate_token(self, token: str) -> dict:
        """Validate JWT token"""
        try:
            payload = jwt.decode(token, self.secret_key, algorithms=[self.algorithm])
            
            # Check blacklist
            if self.blacklist.is_blacklisted(payload.get('jti')):
                raise InvalidTokenError("Token has been revoked")
            
            # Check expiration
            if payload.get('exp', 0) < time.time():
                raise InvalidTokenError("Token has expired")
            
            return payload
            
        except jwt.InvalidTokenError as e:
            raise InvalidTokenError(str(e))
    
    def require_auth(self, f):
        """Decorator for authentication"""
        @wraps(f)
        def decorated(*args, **kwargs):
            token = self._extract_token()
            
            if not token:
                return {"error": "Authentication required"}, 401
            
            try:
                payload = self.validate_token(token)
                g.current_user = payload
            except InvalidTokenError as e:
                return {"error": str(e)}, 401
            
            return f(*args, **kwargs)
        return decorated
    
    def _extract_token(self) -> str:
        """Extract token from request"""
        auth_header = request.headers.get('Authorization', '')
        if auth_header.startswith('Bearer '):
            return auth_header[7:]
        return None
```

### 2. API Key Pattern

```python
import hashlib
import secrets
from datetime import datetime

class APIKeyManager:
    def __init__(self, db):
        self.db = db
    
    def generate_api_key(self, user_id: int, name: str, scopes: list) -> dict:
        """Generate API key"""
        # Generate random key
        raw_key = secrets.token_urlsafe(32)
        
        # Store hashed version
        key_hash = hashlib.sha256(raw_key.encode()).hexdigest()
        key_prefix = raw_key[:8]
        
        self.db.execute("""
            INSERT INTO api_keys (user_id, name, key_hash, key_prefix, scopes, created_at)
            VALUES (?, ?, ?, ?, ?, ?)
        """, (user_id, name, key_hash, key_prefix, json.dumps(scopes), datetime.utcnow()))
        
        return {
            'key': raw_key,  # Only shown once!
            'prefix': key_prefix,
            'name': name
        }
    
    def validate_key(self, api_key: str, required_scope: str = None) -> dict:
        """Validate API key"""
        key_hash = hashlib.sha256(api_key.encode()).hexdigest()
        
        result = self.db.execute("""
            SELECT user_id, scopes, active, rate_limit
            FROM api_keys
            WHERE key_hash = ?
        """, (key_hash,)).fetchone()
        
        if not result:
            raise InvalidAPIKeyError("Invalid API key")
        
        if not result['active']:
            raise InvalidAPIKeyError("API key has been revoked")
        
        scopes = json.loads(result['scopes'])
        
        if required_scope and required_scope not in scopes:
            raise InsufficientScopeError(f"Scope '{required_scope}' required")
        
        return {
            'user_id': result['user_id'],
            'scopes': scopes,
            'rate_limit': result['rate_limit']
        }
```

## 🛡️ Authorization Patterns

### 1. Role-Based Access Control (RBAC)

```python
from enum import Enum
from functools import wraps

class Role(Enum):
    ADMIN = "admin"
    MANAGER = "manager"
    USER = "user"
    GUEST = "guest"

ROLE_HIERARCHY = {
    Role.ADMIN: 4,
    Role.MANAGER: 3,
    Role.USER: 2,
    Role.GUEST: 1
}

ROLE_PERMISSIONS = {
    Role.ADMIN: ['*'],
    Role.MAGER: ['read:*', 'write:own', 'delete:own'],
    Role.USER: ['read:own', 'write:own'],
    Role.GUEST: ['read:public']
}

class RBACAuthorizer:
    def has_permission(self, user_roles: list, required_permission: str) -> bool:
        """Check if user has required permission"""
        for role in user_roles:
            permissions = ROLE_PERMISSIONS.get(role, [])
            
            # Check wildcard
            if '*' in permissions:
                return True
            
            # Check exact match
            if required_permission in permissions:
                return True
            
            # Check pattern match
            for perm in permissions:
                if self._match_permission(perm, required_permission):
                    return True
        
        return False
    
    def _match_permission(self, pattern: str, permission: str) -> bool:
        """Match permission pattern"""
        if pattern.endswith(':*'):
            prefix = pattern[:-2]
            return permission.startswith(prefix + ':')
        return False
    
    def require_role(self, min_role: Role):
        """Decorator for role requirement"""
        def decorator(f):
            @wraps(f)
            def decorated(*args, **kwargs):
                user = g.current_user
                
                user_level = max(ROLE_HIERARCHY.get(r, 0) for r in user.get('roles', []))
                required_level = ROLE_HIERARCHY[min_role]
                
                if user_level < required_level:
                    return {"error": "Insufficient role"}, 403
                
                return f(*args, **kwargs)
            return decorated
        return decorator
```

### 2. Attribute-Based Access Control (ABAC)

```python
class ABACEngine:
    def __init__(self):
        self.policies = []
    
    def add_policy(self, condition: callable, effect: str):
        """Add ABAC policy"""
        self.policies.append({
            'condition': condition,
            'effect': effect  # 'allow' or 'deny'
        })
    
    def evaluate(self, subject: dict, resource: dict, 
                 action: str, environment: dict) -> bool:
        """Evaluate access request"""
        for policy in self.policies:
            if policy['condition'](subject, resource, action, environment):
                return policy['effect'] == 'allow'
        
        return False  # Default deny

# Define policies
engine = ABACEngine()

# Time-based access
engine.add_policy(
    lambda s, r, a, e: (
        a == 'access' and 
        9 <= e.get('hour', 0) <= 17
    ),
    'allow'
)

# Owner access
engine.add_policy(
    lambda s, r, a, e: (
        a in ['read', 'update', 'delete'] and 
        s.get('id') == r.get('owner_id')
    ),
    'allow'
)

# Department access
engine.add_policy(
    lambda s, r, a, e: (
        a == 'read' and 
        s.get('department') == r.get('department')
    ),
    'allow'
)

# Deny sensitive resources outside work location
engine.add_policy(
    lambda s, r, a, e: (
        r.get('classification') == 'confidential' and 
        e.get('location') not in ['office', 'vpn']
    ),
    'deny'
)
```

## 📦 Data Protection Patterns

### 1. Encryption at Rest

```python
from cryptography.fernet import Fernet
import os

class DataEncryption:
    def __init__(self, key: bytes = None):
        self.key = key or os.environ.get('ENCRYPTION_KEY').encode()
        self.cipher = Fernet(self.key)
    
    def encrypt_field(self, value: str) -> str:
        """Encrypt single field"""
        if not value:
            return value
        return self.cipher.encrypt(value.encode()).decode()
    
    def decrypt_field(self, encrypted: str) -> str:
        """Decrypt single field"""
        if not encrypted:
            return encrypted
        return self.cipher.decrypt(encrypted.encode()).decode()

class EncryptedModel:
    """Base model with encrypted fields"""
    encrypted_fields = []
    
    def __init__(self, encryption: DataEncryption):
        self.encryption = encryption
    
    def save(self):
        """Encrypt fields before save"""
        data = self.__dict__.copy()
        
        for field in self.encrypted_fields:
            if field in data:
                data[field] = self.encryption.encrypt_field(data[field])
        
        return self._save_to_db(data)
    
    def load(self, data: dict):
        """Decrypt fields after load"""
        for field in self.encrypted_fields:
            if field in data:
                data[field] = self.encryption.decrypt_field(data[field])
        
        self.__dict__.update(data)
        return self

# Usage
class User(EncryptedModel):
    encrypted_fields = ['ssn', 'credit_card', 'medical_info']
    
    def __init__(self, encryption):
        super().__init__(encryption)
        self.ssn = None
        self.credit_card = None
        self.medical_info = None
```

### 2. Data Masking

```python
import re

class DataMasker:
    @staticmethod
    def mask_email(email: str) -> str:
        """Mask email: j***@example.com"""
        if '@' not in email:
            return email
        
        local, domain = email.split('@', 1)
        masked_local = local[0] + '*' * (len(local) - 1)
        return f"{masked_local}@{domain}"
    
    @staticmethod
    def mask_phone(phone: str) -> str:
        """Mask phone: ***-***-1234"""
        digits = re.sub(r'\D', '', phone)
        if len(digits) >= 4:
            return '*' * (len(digits) - 4) + digits[-4:]
        return phone
    
    @staticmethod
    def mask_credit_card(card: str) -> str:
        """Mask CC: ****-****-****-1234"""
        digits = re.sub(r'\D', '', card)
        if len(digits) >= 4:
            return '*' * (len(digits) - 4) + digits[-4:]
        return card
    
    @staticmethod
    def mask_ssn(ssn: str) -> str:
        """Mask SSN: ***-**-1234"""
        digits = re.sub(r'\D', '', ssn)
        if len(digits) >= 4:
            return '***-**-' + digits[-4:]
        return ssn

class AuditLogger:
    """Logger that automatically masks sensitive data"""
    
    SENSITIVE_FIELDS = {
        'password': '***',
        'ssn': lambda v: DataMasker.mask_ssn(v),
        'credit_card': lambda v: DataMasker.mask_credit_card(v),
        'email': lambda v: DataMasker.mask_email(v),
        'phone': lambda v: DataMasker.mask_phone(v),
    }
    
    def log(self, event: str, data: dict):
        """Log with automatic masking"""
        masked_data = data.copy()
        
        for field, masker in self.SENSITIVE_FIELDS.items():
            if field in masked_data:
                if callable(masker):
                    masked_data[field] = masker(masked_data[field])
                else:
                    masked_data[field] = masker
        
        print(f"[AUDIT] {event}: {json.dumps(masked_data)}")
```

## 🌐 Network Security Patterns

### 1. Rate Limiting Pattern

```python
import time
from collections import defaultdict
from functools import wraps

class RateLimiter:
    def __init__(self, requests: int, window: int):
        self.requests = requests
        self.window = window  # seconds
        self.clients = defaultdict(list)
    
    def is_allowed(self, client_id: str) -> tuple[bool, dict]:
        """Check if request is allowed"""
        now = time.time()
        cutoff = now - self.window
        
        # Clean old requests
        self.clients[client_id] = [
            t for t in self.clients[client_id] if t > cutoff
        ]
        
        current_count = len(self.clients[client_id])
        
        if current_count >= self.requests:
            return False, {
                'retry_after': int(self.clients[client_id][0] - cutoff),
                'limit': self.requests,
                'remaining': 0
            }
        
        # Record this request
        self.clients[client_id].append(now)
        
        return True, {
            'limit': self.requests,
            'remaining': self.requests - current_count - 1,
            'reset': int(cutoff + self.window)
        }
    
    def limit(self, key_func=None):
        """Decorator for rate limiting"""
        def decorator(f):
            @wraps(f)
            def decorated(*args, **kwargs):
                # Get client identifier
                if key_func:
                    client_id = key_func()
                else:
                    client_id = request.remote_addr
                
                allowed, info = self.is_allowed(client_id)
                
                # Add rate limit headers
                headers = {
                    'X-RateLimit-Limit': str(info['limit']),
                    'X-RateLimit-Remaining': str(info['remaining']),
                }
                
                if not allowed:
                    headers['Retry-After'] = str(info['retry_after'])
                    return {'error': 'Rate limit exceeded'}, 429, headers
                
                return f(*args, **kwargs)
            return decorated
        return decorator

# Usage
limiter = RateLimiter(requests=100, window=60)

@app.route('/api/data')
@limiter.limit()
def get_data():
    return {'data': 'response'}
```

### 2. Circuit Breaker Pattern

```python
from enum import Enum
from datetime import datetime, timedelta

class CircuitState(Enum):
    CLOSED = "closed"      # Normal operation
    OPEN = "open"          # Failing, reject requests
    HALF_OPEN = "half_open"  # Testing if recovered

class CircuitBreaker:
    def __init__(self, failure_threshold: int = 5, 
                 recovery_timeout: int = 30):
        self.failure_threshold = failure_threshold
        self.recovery_timeout = timedelta(seconds=recovery_timeout)
        
        self.state = CircuitState.CLOSED
        self.failure_count = 0
        self.last_failure_time = None
    
    def execute(self, func, *args, **kwargs):
        """Execute function with circuit breaker protection"""
        if self.state == CircuitState.OPEN:
            if self._should_attempt_recovery():
                self.state = CircuitState.HALF_OPEN
            else:
                raise CircuitOpenError("Circuit breaker is open")
        
        try:
            result = func(*args, **kwargs)
            self._on_success()
            return result
        except Exception as e:
            self._on_failure()
            raise e
    
    def _should_attempt_recovery(self) -> bool:
        """Check if enough time has passed to try recovery"""
        if not self.last_failure_time:
            return True
        return datetime.now() - self.last_failure_time > self.recovery_timeout
    
    def _on_success(self):
        """Handle successful execution"""
        self.failure_count = 0
        self.state = CircuitState.CLOSED
    
    def _on_failure(self):
        """Handle failed execution"""
        self.failure_count += 1
        self.last_failure_time = datetime.now()
        
        if self.failure_count >= self.failure_threshold:
            self.state = CircuitState.OPEN

# Usage
breaker = CircuitBreaker(failure_threshold=5, recovery_timeout=30)

def call_external_api():
    return breaker.execute(requests.get, 'https://api.example.com/data')
```

## 📊 Logging and Monitoring Patterns

### 1. Security Event Logging

```python
from datetime import datetime
from enum import Enum
import json

class EventType(Enum):
    AUTH_SUCCESS = "auth.success"
    AUTH_FAILURE = "auth.failure"
    ACCESS_GRANTED = "access.granted"
    ACCESS_DENIED = "access.denied"
    DATA_ACCESS = "data.access"
    DATA_MODIFICATION = "data.modification"
    SECURITY_ALERT = "security.alert"

class SecurityEventLogger:
    def __init__(self, sinks: list):
        self.sinks = sinks
    
    def log(self, event_type: EventType, details: dict, 
            severity: str = "info"):
        """Log security event"""
        event = {
            'timestamp': datetime.utcnow().isoformat(),
            'event_type': event_type.value,
            'severity': severity,
            'details': details,
            'context': {
                'user_id': g.get('current_user', {}).get('id'),
                'ip_address': request.remote_addr,
                'user_agent': request.headers.get('User-Agent'),
                'request_id': g.get('request_id')
            }
        }
        
        for sink in self.sinks:
            sink.write(event)
    
    def log_auth_success(self, user_id: str, method: str):
        self.log(EventType.AUTH_SUCCESS, {
            'user_id': user_id,
            'auth_method': method
        })
    
    def log_auth_failure(self, username: str, reason: str):
        self.log(EventType.AUTH_FAILURE, {
            'username': username,
            'reason': reason
        }, severity='warning')
    
    def log_access_denied(self, resource: str, action: str, reason: str):
        self.log(EventType.ACCESS_DENIED, {
            'resource': resource,
            'action': action,
            'reason': reason
        }, severity='warning')

# Sinks
class LogSink:
    def write(self, event: dict):
        raise NotImplementedError

class FileLogSink(LogSink):
    def __init__(self, filepath: str):
        self.filepath = filepath
    
    def write(self, event: dict):
        with open(self.filepath, 'a') as f:
            f.write(json.dumps(event) + '\n')

class CloudWatchSink(LogSink):
    def __init__(self, log_group: str):
        self.client = boto3.client('logs')
        self.log_group = log_group
    
    def write(self, event: dict):
        self.client.put_log_events(
            logGroupName=self.log_group,
            logStreamName=event['event_type'],
            logEvents=[{
                'timestamp': int(datetime.now().timestamp() * 1000),
                'message': json.dumps(event)
            }]
        )
```

## ❓ Interview Questions

### Basic
**Q: What is RBAC?**
A: Role-Based Access Control - permissions assigned to roles, users assigned to roles. Simplifies permission management.

**Q: Why use API keys?**
A: Service-to-service auth, machine access, simpler than OAuth for server apps. Should be rotated and scoped.

### Intermediate
**Q: When would you use ABAC over RBAC?**
A: When access depends on dynamic attributes (user dept, resource owner, time, location). More flexible but more complex.

**Q: Explain the circuit breaker pattern.**
A: Prevents cascading failures. Opens after threshold failures, prevents requests during recovery window, tests recovery in half-open state.

### Advanced
**Q: How do you implement secure audit logging?**
A: (1) Log all security events, (2) Mask sensitive data, (3) Include context (who, what, when, where), (4) Immutable storage, (5) Centralized aggregation, (6) Alert on anomalies.