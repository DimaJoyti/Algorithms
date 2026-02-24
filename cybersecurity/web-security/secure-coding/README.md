# 💻 Secure Coding Practices

Essential security practices for writing secure software.

## 📋 Table of Contents

- [Principles of Secure Coding](#principles-of-secure-coding)
- [Secure Development Lifecycle](#secure-development-lifecycle)
- [Code Review Checklist](#code-review-checklist)
- [Common Vulnerabilities and Fixes](#common-vulnerabilities-and-fixes)
- [Secure Error Handling](#secure-error-handling)
- [Interview Questions](#interview-questions)

## 🎯 Principles of Secure Coding

### Core Principles

```
1. Defense in Depth    → Multiple layers of security controls
2. Least Privilege     → Minimum necessary permissions
3. Fail Secure         → Default to secure state on failure
4. Input Validation    → Never trust user input
5. Output Encoding     → Encode for context
6. Separation of Concerns → Isolate security-critical code
```

### Secure Defaults

```python
# BAD: Insecure defaults
class Config:
    DEBUG = True  # Exposes sensitive info
    SECRET_KEY = "secret"  # Hardcoded, weak
    DATABASE_URL = "postgres://localhost/db"  # No auth

# GOOD: Secure defaults
import os

class Config:
    DEBUG = False  # Secure default
    SECRET_KEY = os.environ.get('SECRET_KEY') or os.urandom(32)
    DATABASE_URL = os.environ.get('DATABASE_URL')
    
    # Validate required settings
    @classmethod
    def validate(cls):
        required = ['SECRET_KEY', 'DATABASE_URL']
        for setting in required:
            if not getattr(cls, setting):
                raise ValueError(f"{setting} must be set")
```

## 🔄 Secure Development Lifecycle

### Security Gates

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Planning   │───►│  Design     │───►│  Coding     │
│             │    │             │    │             │
│ Threat      │    │ Threat      │    │ Secure      │
│ Modeling    │    │ Modeling    │    │ Coding      │
└─────────────┘    └─────────────┘    └──────┬──────┘
                                             │
┌─────────────┐    ┌─────────────┐    ┌──────▼──────┐
│  Deployment │◄───│  Testing    │◄───│  Code       │
│             │    │             │    │  Review     │
│ Security    │    │ SAST, DAST  │    │             │
│ Config      │    │ Pen Test    │    │ Security    │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Pre-commit Security Hooks

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/PyCQA/bandit
    rev: '1.7.5'
    hooks:
      - id: bandit
        args: ['-r', '-ll']
  
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.4.0
    hooks:
      - id: detect-secrets
        args: ['--baseline', '.secrets.baseline']
  
  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.0
    hooks:
      - id: gitleaks

  - repo: local
    hooks:
      - id: check-hardcoded-passwords
        name: Check for hardcoded passwords
        entry: scripts/check_passwords.sh
        language: script
```

## ✅ Code Review Checklist

### Security Checklist

```markdown
## Authentication & Authorization
- [ ] All endpoints have proper authentication
- [ ] Authorization checks before sensitive operations
- [ ] Session management is secure
- [ ] Password storage uses strong hashing

## Input Validation
- [ ] All user input is validated
- [ ] Parameterized queries for database access
- [ ] File uploads are validated and restricted
- [ ] Input is sanitized before use

## Data Protection
- [ ] Sensitive data encrypted at rest
- [ ] TLS for data in transit
- [ ] Secrets not hardcoded
- [ ] Proper key management

## Error Handling
- [ ] Errors don't leak sensitive info
- [ ] Proper logging without sensitive data
- [ ] Graceful failure handling

## Dependencies
- [ ] No known vulnerable dependencies
- [ ] Dependencies are up to date
- [ ] Minimal dependency footprint
```

## 🔧 Common Vulnerabilities and Fixes

### 1. Hardcoded Secrets

```python
# BAD
API_KEY = "sk-1234567890abcdef"
DB_PASSWORD = "password123"

# GOOD
import os
API_KEY = os.environ.get('API_KEY')
DB_PASSWORD = os.environ.get('DB_PASSWORD')

# Or use secrets manager
from aws_secretsmanager import get_secret
api_key = get_secret('myapp/api-key')
```

### 2. Insecure Deserialization

```python
# BAD
import pickle
data = pickle.loads(user_input)  # Remote code execution risk

# GOOD
import json
data = json.loads(user_input)  # Safe parsing

# For complex objects
def safe_deserialize(data: dict):
    """Validate structure before creating object"""
    required_fields = ['name', 'email']
    for field in required_fields:
        if field not in data:
            raise ValueError(f"Missing field: {field}")
    
    if not re.match(r'^[^@]+@[^@]+\.[^@]+$', data['email']):
        raise ValueError("Invalid email")
    
    return User(name=data['name'], email=data['email'])
```

### 3. Path Traversal

```python
# BAD
def read_file(filename):
    with open(f"/uploads/{filename}") as f:
        return f.read()  # Can access any file with ../../../etc/passwd

# GOOD
import os
from pathlib import Path

UPLOAD_DIR = Path("/uploads").resolve()

def read_file(filename: str) -> str:
    # Sanitize filename
    filename = os.path.basename(filename)
    
    # Resolve and check path
    filepath = (UPLOAD_DIR / filename).resolve()
    
    # Ensure file is within allowed directory
    if not str(filepath).startswith(str(UPLOAD_DIR)):
        raise SecurityError("Path traversal attempt")
    
    # Check file exists and is a file
    if not filepath.is_file():
        raise FileNotFoundError("File not found")
    
    with open(filepath) as f:
        return f.read()
```

### 4. Command Injection

```python
# BAD
import os
def ping_host(host):
    return os.system(f"ping -c 1 {host}")  # Command injection!

# GOOD
import subprocess
import shlex

def ping_host(host: str) -> str:
    # Validate input
    if not re.match(r'^[a-zA-Z0-9.-]+$', host):
        raise ValueError("Invalid hostname")
    
    # Use list form (no shell interpretation)
    result = subprocess.run(
        ['ping', '-c', '1', host],
        capture_output=True,
        text=True,
        timeout=10
    )
    return result.stdout
```

### 5. Insecure Random Numbers

```python
# BAD
import random
token = random.randint(0, 999999)  # Predictable!

# GOOD
import secrets
token = secrets.randbelow(1000000)
api_key = secrets.token_urlsafe(32)
```

## 🚨 Secure Error Handling

### Error Handling Best Practices

```python
import logging
from functools import wraps
from flask import jsonify

logger = logging.getLogger(__name__)

# Custom exceptions
class AppError(Exception):
    """Base application error"""
    def __init__(self, message: str, code: str = "INTERNAL_ERROR"):
        self.message = message
        self.code = code
        super().__init__(message)

class ValidationError(AppError):
    def __init__(self, message: str):
        super().__init__(message, "VALIDATION_ERROR")

class AuthenticationError(AppError):
    def __init__(self, message: str = "Authentication required"):
        super().__init__(message, "AUTH_ERROR")

class AuthorizationError(AppError):
    def __init__(self, message: str = "Access denied"):
        super().__init__(message, "FORBIDDEN")

# Error handler
def handle_errors(f):
    """Decorator for consistent error handling"""
    @wraps(f)
    def wrapper(*args, **kwargs):
        try:
            return f(*args, **kwargs)
        except ValidationError as e:
            logger.warning(f"Validation error: {e.message}")
            return jsonify({
                "error": e.code,
                "message": e.message
            }), 400
        except AuthenticationError as e:
            logger.info(f"Auth error: {e.message}")
            return jsonify({
                "error": e.code,
                "message": e.message
            }), 401
        except AuthorizationError as e:
            logger.warning(f"Authz error: {e.message}")
            return jsonify({
                "error": e.code,
                "message": e.message
            }), 403
        except Exception as e:
            # Log full error internally
            logger.exception("Unexpected error")
            # Return generic message to client
            return jsonify({
                "error": "INTERNAL_ERROR",
                "message": "An unexpected error occurred"
            }), 500
    return wrapper

# Usage
@app.route('/api/users', methods=['POST'])
@handle_errors
def create_user():
    data = request.get_json()
    
    # Validate input
    if not data.get('email'):
        raise ValidationError("Email is required")
    
    # Create user
    user = create_user_in_db(data)
    
    return jsonify({"id": user.id}), 201
```

### Secure Logging

```python
import logging
import re

class SecureFormatter(logging.Formatter):
    """Formatter that redacts sensitive data"""
    
    SENSITIVE_PATTERNS = [
        (r'password["\']?\s*[:=]\s*["\']?([^"\',\s]+)', 'password=***REDACTED***'),
        (r'token["\']?\s*[:=]\s*["\']?([^"\',\s]+)', 'token=***REDACTED***'),
        (r'api_key["\']?\s*[:=]\s*["\']?([^"\',\s]+)', 'api_key=***REDACTED***'),
        (r'Authorization:\s*Bearer\s+\S+', 'Authorization: Bearer ***REDACTED***'),
    ]
    
    def format(self, record):
        message = super().format(record)
        
        for pattern, replacement in self.SENSITIVE_PATTERNS:
            message = re.sub(pattern, replacement, message, flags=re.IGNORECASE)
        
        return message

# Setup
handler = logging.StreamHandler()
handler.setFormatter(SecureFormatter(
    '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
))

logger = logging.getLogger(__name__)
logger.addHandler(handler)
logger.setLevel(logging.INFO)
```

## ❓ Interview Questions

### Basic
**Q: What are the OWASP Top 10?**
A: Standard awareness document for developers: Injection, Broken Auth, Sensitive Data, XXE, Access Control, Misconfiguration, XSS, Deserialization, Components, Logging.

**Q: Why should you never trust user input?**
A: Users may intentionally or unintentionally provide malicious input that can cause injection attacks, XSS, or other vulnerabilities.

### Intermediate
**Q: Explain defense in depth.**
A: Multiple independent security controls so if one fails, others still protect. Example: Input validation + parameterized queries + least privilege DB user.

**Q: How do you handle secrets in code?**
A: Use environment variables or secret managers. Never commit secrets to version control. Rotate secrets regularly. Use pre-commit hooks to detect leaks.

### Advanced
**Q: How would you implement secure coding practices in a development team?**
A: (1) Training and awareness, (2) Secure coding standards, (3) Code review process with security checklist, (4) Automated SAST in CI/CD, (5) Dependency scanning, (6) Security champion program, (7) Regular security audits.