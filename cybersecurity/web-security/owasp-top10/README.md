# 🛡️ OWASP Top 10 (2021)

The OWASP Top 10 represents the most critical security risks to web applications.

## 📋 Table of Contents

- [A01: Broken Access Control](#a01-broken-access-control)
- [A02: Cryptographic Failures](#a02-cryptographic-failures)
- [A03: Injection](#a03-injection)
- [A04: Insecure Design](#a04-insecure-design)
- [A05: Security Misconfiguration](#a05-security-misconfiguration)
- [A06: Vulnerable Components](#a06-vulnerable-components)
- [A07: Authentication Failures](#a07-authentication-failures)
- [A08: Software Integrity Failures](#a08-software-integrity-failures)
- [A09: Logging Failures](#a09-logging-failures)
- [A10: SSRF](#a10-ssrf)

## A01: Broken Access Control

### Vulnerability Examples
- Accessing other users' data by modifying URLs
- Viewing/editing without proper authorization checks
- Metadata manipulation (JWT, cookies, hidden fields)
- API endpoints without access controls

### Prevention

```python
from functools import wraps
from flask import request, jsonify, g

def require_permission(permission):
    """Decorator for access control"""
    def decorator(f):
        @wraps(f)
        def decorated_function(*args, **kwargs):
            user = g.current_user
            
            if not user:
                return jsonify({"error": "Unauthorized"}), 401
            
            if permission not in user.permissions:
                return jsonify({"error": "Forbidden"}), 403
            
            # Resource ownership check
            resource_id = kwargs.get('id') or request.view_args.get('id')
            if resource_id and not user.owns_resource(resource_id):
                return jsonify({"error": "Access denied"}), 403
            
            return f(*args, **kwargs)
        return decorated_function
    return decorator

# Usage
@app.route('/api/users/<int:id>')
@require_permission('read_user')
def get_user(id):
    user = User.query.get_or_404(id)
    return jsonify(user.to_dict())

@app.route('/api/users/<int:id>', methods=['DELETE'])
@require_permission('delete_user')
def delete_user(id):
    user = User.query.get_or_404(id)
    db.session.delete(user)
    db.session.commit()
    return '', 204
```

### Access Control Checklist
- [ ] Deny by default
- [ ] Implement proper authorization checks
- [ ] Disable directory listing
- [ ] Log access control failures
- [ ] Rate limit API access
- [ ] Invalidate sessions on logout

## A02: Cryptographic Failures

### Vulnerability Examples
- Storing passwords in plaintext
- Using weak algorithms (MD5, SHA1)
- Not using TLS
- Hardcoded encryption keys

### Prevention

```python
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import bcrypt
import os

# Password hashing
def hash_password(password: str) -> str:
    """Secure password hashing with bcrypt"""
    salt = bcrypt.gensalt(rounds=12)
    return bcrypt.hashpw(password.encode(), salt).decode()

def verify_password(password: str, hashed: str) -> bool:
    """Verify password against hash"""
    return bcrypt.checkpw(password.encode(), hashed.encode())

# Data encryption
class DataEncryption:
    def __init__(self, key: bytes = None):
        self.key = key or Fernet.generate_key()
        self.cipher = Fernet(self.key)
    
    def encrypt(self, data: str) -> bytes:
        return self.cipher.encrypt(data.encode())
    
    def decrypt(self, encrypted: bytes) -> str:
        return self.cipher.decrypt(encrypted).decode()
    
    @staticmethod
    def derive_key(password: str, salt: bytes) -> bytes:
        """Derive encryption key from password"""
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=480000,
        )
        return kdf.derive(password.encode())

# Usage
# Passwords
hashed = hash_password("MySecretPassword123")
print(f"Verified: {verify_password('MySecretPassword123', hashed)}")

# Data encryption
enc = DataEncryption()
encrypted = enc.encrypt("Sensitive data like SSN")
print(f"Decrypted: {enc.decrypt(encrypted)}")
```

## A03: Injection

### Types
- SQL Injection
- NoSQL Injection
- OS Command Injection
- LDAP Injection

### Prevention

```python
from typing import List, Any
import sqlite3
import psycopg2
from psycopg2 import sql

# BAD: Vulnerable to SQL injection
def get_user_unsafe(user_id: str):
    query = f"SELECT * FROM users WHERE id = {user_id}"  # DANGEROUS!
    cursor.execute(query)

# GOOD: Parameterized queries
def get_user_safe(user_id: int):
    query = "SELECT * FROM users WHERE id = %s"
    cursor.execute(query, (user_id,))
    return cursor.fetchone()

# GOOD: Using ORM
def get_user_orm(user_id: int):
    return User.query.filter_by(id=user_id).first()

# Input validation
import re
from html import escape

class InputValidator:
    @staticmethod
    def sanitize_string(value: str, max_length: int = 255) -> str:
        """Sanitize string input"""
        if not isinstance(value, str):
            raise ValueError("Expected string")
        if len(value) > max_length:
            raise ValueError(f"Max length {max_length} exceeded")
        # Remove dangerous characters
        return escape(value.strip())
    
    @staticmethod
    def validate_email(email: str) -> bool:
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        return bool(re.match(pattern, email))
    
    @staticmethod
    def validate_integer(value: str) -> int:
        try:
            return int(value)
        except ValueError:
            raise ValueError("Expected integer")

# Command injection prevention
import subprocess

def run_command_safe(filename: str):
    # BAD
    # subprocess.call(f"cat {filename}", shell=True)  # DANGEROUS!
    
    # GOOD
    # Validate input
    if not re.match(r'^[\w\-\.]+$', filename):
        raise ValueError("Invalid filename")
    
    # Use list form (no shell)
    result = subprocess.run(
        ['cat', filename],
        capture_output=True,
        text=True,
        timeout=30
    )
    return result.stdout
```

## A04: Insecure Design

### Prevention
- Threat modeling during design
- Secure design patterns
- Reference architectures

```python
# Secure design pattern: Circuit Breaker
class CircuitBreaker:
    def __init__(self, failure_threshold=5, timeout=60):
        self.failures = 0
        self.failure_threshold = failure_threshold
        self.timeout = timeout
        self.last_failure = None
        self.state = "CLOSED"  # CLOSED, OPEN, HALF_OPEN
    
    def call(self, func, *args):
        if self.state == "OPEN":
            if self._should_attempt_reset():
                self.state = "HALF_OPEN"
            else:
                raise Exception("Circuit breaker is OPEN")
        
        try:
            result = func(*args)
            self._on_success()
            return result
        except Exception as e:
            self._on_failure()
            raise e
    
    def _on_success(self):
        self.failures = 0
        self.state = "CLOSED"
    
    def _on_failure(self):
        self.failures += 1
        self.last_failure = time.time()
        if self.failures >= self.failure_threshold:
            self.state = "OPEN"
    
    def _should_attempt_reset(self):
        return time.time() - self.last_failure > self.timeout
```

## A05: Security Misconfiguration

### Hardening Checklist

```python
# Security headers middleware
from flask import Flask, Response

app = Flask(__name__)

@app.after_request
def add_security_headers(response: Response):
    # Prevent clickjacking
    response.headers['X-Frame-Options'] = 'DENY'
    
    # Prevent MIME sniffing
    response.headers['X-Content-Type-Options'] = 'nosniff'
    
    # XSS protection (legacy browsers)
    response.headers['X-XSS-Protection'] = '1; mode=block'
    
    # HTTPS only
    response.headers['Strict-Transport-Security'] = 'max-age=31536000; includeSubDomains'
    
    # Content Security Policy
    response.headers['Content-Security-Policy'] = (
        "default-src 'self'; "
        "script-src 'self' 'unsafe-inline'; "
        "style-src 'self' 'unsafe-inline'; "
        "img-src 'self' data:; "
        "font-src 'self';"
    )
    
    # Referrer policy
    response.headers['Referrer-Policy'] = 'strict-origin-when-cross-origin'
    
    # Permissions policy
    response.headers['Permissions-Policy'] = (
        "geolocation=(), microphone=(), camera=()"
    )
    
    return response

# Disable debug mode in production
app.config['DEBUG'] = False
app.config['TESTING'] = False

# Remove unnecessary headers
@app.after_request
def remove_server_header(response):
    response.headers.pop('Server', None)
    response.headers.pop('X-Powered-By', None)
    return response
```

## A06: Vulnerable Components

### Dependency Scanning

```python
# requirements.txt security check
import subprocess
import json

def check_vulnerabilities():
    """Check dependencies for known vulnerabilities"""
    try:
        result = subprocess.run(
            ['pip-audit', '--json'],
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            vulnerabilities = json.loads(result.stdout)
            for vuln in vulnerabilities:
                print(f"VULNERABILITY: {vuln['name']} {vuln['version']}")
                print(f"  CVE: {vuln['id']}")
                print(f"  Severity: {vuln['severity']}")
                print(f"  Fix: Upgrade to {vuln['fix_versions']}")
    except FileNotFoundError:
        print("pip-audit not installed")

# Safe dependency versions
SAFE_VERSIONS = {
    'django': '>=4.2',
    'flask': '>=2.3',
    'requests': '>=2.31',
    'cryptography': '>=41.0',
}

def verify_versions():
    import pkg_resources
    for package, version_spec in SAFE_VERSIONS.items():
        try:
            dist = pkg_resources.get_distribution(package)
            # Check if version satisfies requirement
            req = pkg_resources.Requirement.parse(f"{package}{version_spec}")
            if dist.version not in req:
                print(f"WARNING: {package} {dist.version} does not meet {version_spec}")
        except pkg_resources.DistributionNotFound:
            print(f"NOT INSTALLED: {package}")
```

## A07: Authentication Failures

### Secure Authentication Implementation

```python
import secrets
import time
from dataclasses import dataclass
from typing import Optional

@dataclass
class Session:
    user_id: int
    token: str
    expires_at: float
    ip_address: str

class AuthenticationManager:
    def __init__(self):
        self.sessions = {}
        self.failed_attempts = {}
        self.lockouts = {}
    
    def login(self, username: str, password: str, ip: str) -> Optional[Session]:
        # Check lockout
        if self._is_locked_out(username, ip):
            raise Exception("Account temporarily locked")
        
        user = self._authenticate_user(username, password)
        
        if not user:
            self._record_failed_attempt(username, ip)
            return None
        
        # Clear failed attempts on success
        self._clear_failed_attempts(username, ip)
        
        # Create session
        session = Session(
            user_id=user.id,
            token=secrets.token_urlsafe(32),
            expires_at=time.time() + 3600,  # 1 hour
            ip_address=ip
        )
        self.sessions[session.token] = session
        
        return session
    
    def _authenticate_user(self, username: str, password: str):
        user = User.query.filter_by(username=username).first()
        if user and verify_password(password, user.password_hash):
            return user
        return None
    
    def _is_locked_out(self, username: str, ip: str) -> bool:
        key = f"{username}:{ip}"
        lockout = self.lockouts.get(key, 0)
        return time.time() < lockout
    
    def _record_failed_attempt(self, username: str, ip: str):
        key = f"{username}:{ip}"
        self.failed_attempts[key] = self.failed_attempts.get(key, 0) + 1
        
        # Lock after 5 failures for 15 minutes
        if self.failed_attempts[key] >= 5:
            self.lockouts[key] = time.time() + 900
    
    def _clear_failed_attempts(self, username: str, ip: str):
        key = f"{username}:{ip}"
        self.failed_attempts.pop(key, None)
        self.lockouts.pop(key, None)
    
    def validate_session(self, token: str, ip: str) -> Optional[Session]:
        session = self.sessions.get(token)
        
        if not session:
            return None
        
        if time.time() > session.expires_at:
            del self.sessions[token]
            return None
        
        # Optional: Validate IP hasn't changed
        # if session.ip_address != ip:
        #     return None
        
        return session
    
    def logout(self, token: str):
        self.sessions.pop(token, None)
```

## A08: Software Integrity Failures

### Code Signing Verification

```python
import hashlib
import hmac

class IntegrityChecker:
    def __init__(self, secret_key: bytes):
        self.secret_key = secret_key
    
    def sign_data(self, data: bytes) -> str:
        """Create HMAC signature"""
        return hmac.new(self.secret_key, data, hashlib.sha256).hexdigest()
    
    def verify_signature(self, data: bytes, signature: str) -> bool:
        """Verify HMAC signature"""
        expected = self.sign_data(data)
        return hmac.compare_digest(expected, signature)
    
    def checksum_file(self, filepath: str) -> str:
        """Generate SHA-256 checksum"""
        sha256 = hashlib.sha256()
        with open(filepath, 'rb') as f:
            for chunk in iter(lambda: f.read(4096), b''):
                sha256.update(chunk)
        return sha256.hexdigest()
```

## A09: Logging Failures

### Secure Logging Implementation

```python
import logging
import json
from datetime import datetime

class SecurityLogger:
    def __init__(self, name: str = "security"):
        self.logger = logging.getLogger(name)
        self.logger.setLevel(logging.INFO)
        
        # File handler
        handler = logging.FileHandler('security.log')
        handler.setFormatter(logging.Formatter(
            '%(asctime)s - %(levelname)s - %(message)s'
        ))
        self.logger.addHandler(handler)
    
    def log_auth_event(self, event_type: str, username: str, 
                       ip: str, success: bool, details: str = ""):
        event = {
            "timestamp": datetime.utcnow().isoformat(),
            "event_type": event_type,
            "username": username,
            "ip_address": ip,
            "success": success,
            "details": details
        }
        
        if success:
            self.logger.info(json.dumps(event))
        else:
            self.logger.warning(json.dumps(event))
    
    def log_access_control(self, user: str, resource: str, 
                          action: str, allowed: bool):
        event = {
            "timestamp": datetime.utcnow().isoformat(),
            "event_type": "access_control",
            "user": user,
            "resource": resource,
            "action": action,
            "allowed": allowed
        }
        
        if allowed:
            self.logger.info(json.dumps(event))
        else:
            self.logger.warning(json.dumps(event))
    
    def log_security_alert(self, alert_type: str, severity: str, 
                          description: str, ip: str = ""):
        event = {
            "timestamp": datetime.utcnow().isoformat(),
            "event_type": "security_alert",
            "alert_type": alert_type,
            "severity": severity,
            "description": description,
            "ip_address": ip
        }
        self.logger.critical(json.dumps(event))

# Usage
sec_log = SecurityLogger()
sec_log.log_auth_event("login", "john_doe", "192.168.1.100", True)
sec_log.log_access_control("john_doe", "/admin/users", "read", False)
sec_log.log_security_alert("brute_force", "high", "Multiple failed logins", "10.0.0.50")
```

## A10: SSRF (Server-Side Request Forgery)

### Prevention

```python
import ipaddress
import socket
from urllib.parse import urlparse

ALLOWED_DOMAINS = {'api.example.com', 'cdn.example.com'}
BLOCKED_IPS = {'127.0.0.1', '0.0.0.0', '169.254.169.254'}  # AWS metadata

def validate_url(url: str) -> bool:
    """Validate URL to prevent SSRF"""
    try:
        parsed = urlparse(url)
        
        # Only allow http/https
        if parsed.scheme not in ['http', 'https']:
            return False
        
        # Check domain allowlist
        if parsed.hostname not in ALLOWED_DOMAINS:
            return False
        
        # Resolve and check IP
        ip = socket.gethostbyname(parsed.hostname)
        
        # Block private/internal IPs
        ip_obj = ipaddress.ip_address(ip)
        if ip_obj.is_private or ip_obj.is_loopback or ip_obj.is_link_local:
            return False
        
        if ip in BLOCKED_IPS:
            return False
        
        return True
    except:
        return False

def safe_request(url: str):
    """Make safe HTTP request"""
    if not validate_url(url):
        raise ValueError("URL validation failed")
    
    # Use with timeout and size limits
    response = requests.get(url, timeout=10, allow_redirects=False)
    return response
```

## Quick Reference

| Risk | Prevention Priority |
|------|---------------------|
| Injection | Parameterized queries, input validation |
| Auth Failures | MFA, secure session management |
| Sensitive Data | Encrypt at rest and in transit |
| XXE | Disable DTD, external entities |
| Access Control | Deny by default, validate on server |
| Misconfig | Harden configs, remove defaults |
| XSS | Output encoding, CSP headers |
| Insecure Deserialization | Validate, sanitize all inputs |
| Components | Update dependencies, track CVEs |
| Logging | Log security events, monitor alerts |