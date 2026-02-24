# ✅ Input Validation and Output Encoding

Critical security controls for preventing injection attacks and data validation.

## 📋 Table of Contents

- [Input Validation Strategies](#input-validation-strategies)
- [Output Encoding](#output-encoding)
- [XSS Prevention](#xss-prevention)
- [SQL Injection Prevention](#sql-injection-prevention)
- [Validation Framework](#validation-framework)
- [Interview Questions](#interview-questions)

## 🔍 Input Validation Strategies

### Validation Layers

```
┌─────────────────────────────────────────┐
│         Client-side Validation          │ ← UX (easily bypassed)
├─────────────────────────────────────────┤
│         Server-side Validation          │ ← Security
│  ┌─────────────────────────────────┐    │
│  │     Syntax Validation           │    │
│  │  (type, length, format)         │    │
│  ├─────────────────────────────────┤    │
│  │     Semantic Validation         │    │
│  │  (business rules, ranges)       │    │
│  └─────────────────────────────────┘    │
├─────────────────────────────────────────┤
│         Database Parameterization       │ ← Last line of defense
└─────────────────────────────────────────┘
```

### Python Validation Framework

```python
import re
from dataclasses import dataclass
from typing import Any, List, Optional, Callable
from html import escape

@dataclass
class ValidationResult:
    valid: bool
    value: Any
    errors: List[str]

class Validator:
    def __init__(self, value: Any):
        self.value = value
        self.errors = []
        self._valid = True
    
    def required(self) -> 'Validator':
        """Check value is not empty"""
        if self.value is None or self.value == '':
            self.errors.append("Field is required")
            self._valid = False
        return self
    
    def string(self) -> 'Validator':
        """Check value is string"""
        if not isinstance(self.value, str):
            self.errors.append("Must be a string")
            self._valid = False
        return self
    
    def min_length(self, length: int) -> 'Validator':
        """Check minimum length"""
        if isinstance(self.value, str) and len(self.value) < length:
            self.errors.append(f"Must be at least {length} characters")
            self._valid = False
        return self
    
    def max_length(self, length: int) -> 'Validator':
        """Check maximum length"""
        if isinstance(self.value, str) and len(self.value) > length:
            self.errors.append(f"Must be at most {length} characters")
            self._valid = False
        return self
    
    def email(self) -> 'Validator':
        """Validate email format"""
        pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        if not re.match(pattern, str(self.value)):
            self.errors.append("Invalid email format")
            self._valid = False
        return self
    
    def regex(self, pattern: str, message: str = "Invalid format") -> 'Validator':
        """Validate against regex pattern"""
        if not re.match(pattern, str(self.value)):
            self.errors.append(message)
            self._valid = False
        return self
    
    def integer(self) -> 'Validator':
        """Check value is integer"""
        try:
            int(self.value)
        except (ValueError, TypeError):
            self.errors.append("Must be an integer")
            self._valid = False
        return self
    
    def range(self, min_val: int, max_val: int) -> 'Validator':
        """Check value is within range"""
        try:
            val = int(self.value)
            if not min_val <= val <= max_val:
                self.errors.append(f"Must be between {min_val} and {max_val}")
                self._valid = False
        except:
            pass
        return self
    
    def in_list(self, allowed: List[str]) -> 'Validator':
        """Check value is in allowed list"""
        if str(self.value) not in allowed:
            self.errors.append(f"Must be one of: {', '.join(allowed)}")
            self._valid = False
        return self
    
    def custom(self, validator: Callable[[Any], bool], 
               message: str = "Validation failed") -> 'Validator':
        """Custom validation function"""
        if not validator(self.value):
            self.errors.append(message)
            self._valid = False
        return self
    
    def sanitize(self) -> 'Validator':
        """Sanitize string value"""
        if isinstance(self.value, str):
            self.value = escape(self.value.strip())
        return self
    
    def result(self) -> ValidationResult:
        return ValidationResult(
            valid=self._valid,
            value=self.value,
            errors=self.errors
        )

# Usage
def validate_user_input(data: dict) -> ValidationResult:
    username = Validator(data.get('username')) \
        .required().string().min_length(3).max_length(50) \
        .regex(r'^[a-zA-Z0-9_]+$', "Only letters, numbers, underscore") \
        .result()
    
    email = Validator(data.get('email')) \
        .required().string().email().result()
    
    age = Validator(data.get('age')) \
        .required().integer().range(18, 120).result()
    
    errors = username.errors + email.errors + age.errors
    
    return ValidationResult(
        valid=len(errors) == 0,
        value={'username': username.value, 'email': email.value, 'age': age.value},
        errors=errors
    )
```

## 📤 Output Encoding

### Context-Specific Encoding

```python
import html
import json
import urllib.parse

class OutputEncoder:
    """Context-specific output encoding"""
    
    @staticmethod
    def html_encode(value: str) -> str:
        """Encode for HTML context"""
        return html.escape(str(value), quote=True)
    
    @staticmethod
    def html_attribute_encode(value: str) -> str:
        """Encode for HTML attribute context"""
        # More aggressive encoding for attributes
        encoded = ""
        for char in str(value):
            encoded += f"&#{ord(char)};" if ord(char) > 127 else html.escape(char)
        return encoded
    
    @staticmethod
    def javascript_encode(value: str) -> str:
        """Encode for JavaScript context"""
        return json.dumps(str(value))
    
    @staticmethod
    def url_encode(value: str) -> str:
        """Encode for URL context"""
        return urllib.parse.quote(str(value), safe='')
    
    @staticmethod
    def css_encode(value: str) -> str:
        """Encode for CSS context"""
        encoded = ""
        for char in str(value):
            encoded += f"\\{hex(ord(char))[1:]}" if not char.isalnum() else char
        return encoded
    
    @staticmethod
    def json_encode(value: Any) -> str:
        """Encode for JSON context"""
        return json.dumps(value)

# Template rendering with auto-encoding
class SecureTemplate:
    def __init__(self, template: str):
        self.template = template
    
    def render(self, context: dict) -> str:
        """Render template with auto-encoding"""
        result = self.template
        
        for key, value in context.items():
            # HTML encode by default
            encoded = OutputEncoder.html_encode(str(value))
            result = result.replace(f"{{{{{key}}}}}", encoded)
        
        return result
```

## 🛡️ XSS Prevention

### Content Security Policy

```python
from flask import Flask, Response, request

app = Flask(__name__)

@app.after_request
def add_csp(response: Response):
    """Add Content Security Policy header"""
    csp = (
        "default-src 'self'; "
        "script-src 'self' 'nonce-{nonce}'; "
        "style-src 'self' 'unsafe-inline'; "
        "img-src 'self' data: https:; "
        "font-src 'self'; "
        "connect-src 'self'; "
        "frame-ancestors 'none'; "
        "base-uri 'self'; "
        "form-action 'self';"
    )
    
    nonce = generate_nonce()
    csp = csp.replace('{nonce}', nonce)
    
    response.headers['Content-Security-Policy'] = csp
    response.headers['X-Content-Type-Options'] = 'nosniff'
    response.headers['X-Frame-Options'] = 'DENY'
    
    return response

def generate_nonce():
    import secrets
    return secrets.token_urlsafe(16)

# XSS sanitization for rich text
import bleach

ALLOWED_TAGS = ['b', 'i', 'u', 'strong', 'em', 'p', 'br', 'ul', 'ol', 'li', 'a']
ALLOWED_ATTRIBUTES = {'a': ['href', 'title']}
ALLOWED_PROTOCOLS = ['http', 'https', 'mailto']

def sanitize_html(value: str) -> str:
    """Sanitize HTML for safe display"""
    return bleach.clean(
        value,
        tags=ALLOWED_TAGS,
        attributes=ALLOWED_ATTRIBUTES,
        protocols=ALLOWED_PROTOCOLS,
        strip=True
    )
```

### React/JavaScript XSS Prevention

```javascript
class XSSPrevention {
    // React automatically escapes, but for raw HTML:
    static dangerouslySetInnerHTML(content) {
        return {
            __html: this.sanitizeHTML(content)
        };
    }
    
    static sanitizeHTML(html) {
        const allowedTags = ['b', 'i', 'u', 'strong', 'em', 'p', 'a'];
        const doc = new DOMParser().parseFromString(html, 'text/html');
        
        const walk = (node) => {
            for (let child of Array.from(node.childNodes)) {
                if (child.nodeType === Node.ELEMENT_NODE) {
                    if (!allowedTags.includes(child.tagName.toLowerCase())) {
                        child.replaceWith(...child.childNodes);
                        continue;
                    }
                    
                    // Remove dangerous attributes
                    for (let attr of Array.from(child.attributes)) {
                        if (attr.name.startsWith('on') || attr.name === 'href' && attr.value.startsWith('javascript:')) {
                            child.removeAttribute(attr.name);
                        }
                    }
                }
                walk(child);
            }
        };
        
        walk(doc.body);
        return doc.body.innerHTML;
    }
    
    static escapeHTML(str) {
        const escapeMap = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#x27;',
            '/': '&#x2F;'
        };
        
        return String(str).replace(/[&<>"'/]/g, c => escapeMap[c]);
    }
}
```

## 💉 SQL Injection Prevention

### Parameterized Queries

```python
import sqlite3
from typing import List, Any, Optional

class SafeDatabase:
    def __init__(self, db_path: str):
        self.conn = sqlite3.connect(db_path)
        self.conn.row_factory = sqlite3.Row
    
    def execute(self, query: str, params: tuple = ()) -> List[dict]:
        """Execute parameterized query"""
        cursor = self.conn.cursor()
        cursor.execute(query, params)
        return [dict(row) for row in cursor.fetchall()]
    
    def find_user(self, user_id: int) -> Optional[dict]:
        """Safe: Parameterized query"""
        query = "SELECT * FROM users WHERE id = ?"
        results = self.execute(query, (user_id,))
        return results[0] if results else None
    
    def find_users_by_email(self, email: str) -> List[dict]:
        """Safe: Parameterized query with LIKE"""
        query = "SELECT * FROM users WHERE email LIKE ?"
        results = self.execute(query, (f"%{email}%",))
        return results
    
    def insert_user(self, username: str, email: str, password_hash: str) -> int:
        """Safe: Parameterized INSERT"""
        query = """
            INSERT INTO users (username, email, password_hash)
            VALUES (?, ?, ?)
        """
        cursor = self.conn.cursor()
        cursor.execute(query, (username, email, password_hash))
        self.conn.commit()
        return cursor.lastrowid
    
    def dynamic_order_by(self, column: str, direction: str = 'ASC') -> List[dict]:
        """Safe: Whitelist for dynamic SQL"""
        allowed_columns = ['id', 'username', 'email', 'created_at']
        allowed_directions = ['ASC', 'DESC']
        
        if column not in allowed_columns:
            raise ValueError(f"Invalid column: {column}")
        if direction.upper() not in allowed_directions:
            raise ValueError(f"Invalid direction: {direction}")
        
        # Safe to interpolate because we whitelisted
        query = f"SELECT * FROM users ORDER BY {column} {direction}"
        return self.execute(query)

# ORM usage (safest)
from sqlalchemy.orm import Session
from models import User

def get_user_safe(db: Session, user_id: int) -> Optional[User]:
    """Using ORM - automatic parameterization"""
    return db.query(User).filter(User.id == user_id).first()

def search_users(db: Session, search_term: str) -> List[User]:
    """Using ORM with LIKE"""
    return db.query(User).filter(
        User.username.ilike(f"%{search_term}%")
    ).all()
```

## 📋 Validation Framework

### Comprehensive Validation

```python
from dataclasses import dataclass
from typing import Any, Dict, Type, get_type_hints
import re

@dataclass
class FieldRule:
    required: bool = True
    min_length: int = None
    max_length: int = None
    pattern: str = None
    min_value: int = None
    max_value: int = None
    allowed_values: list = None
    custom_validator: callable = None

class SchemaValidator:
    def __init__(self, schema: Dict[str, FieldRule]):
        self.schema = schema
    
    def validate(self, data: Dict[str, Any]) -> tuple[bool, Dict, List[str]]:
        """Validate data against schema"""
        errors = []
        validated = {}
        
        for field, rules in self.schema.items():
            value = data.get(field)
            
            # Required check
            if rules.required and value is None:
                errors.append(f"{field}: Required field missing")
                continue
            
            if value is None:
                continue
            
            # String validations
            if isinstance(value, str):
                if rules.min_length and len(value) < rules.min_length:
                    errors.append(f"{field}: Minimum {rules.min_length} characters")
                
                if rules.max_length and len(value) > rules.max_length:
                    errors.append(f"{field}: Maximum {rules.max_length} characters")
                
                if rules.pattern and not re.match(rules.pattern, value):
                    errors.append(f"{field}: Invalid format")
            
            # Numeric validations
            if isinstance(value, (int, float)):
                if rules.min_value is not None and value < rules.min_value:
                    errors.append(f"{field}: Minimum value is {rules.min_value}")
                
                if rules.max_value is not None and value > rules.max_value:
                    errors.append(f"{field}: Maximum value is {rules.max_value}")
            
            # Allowed values
            if rules.allowed_values and value not in rules.allowed_values:
                errors.append(f"{field}: Must be one of {rules.allowed_values}")
            
            # Custom validator
            if rules.custom_validator:
                try:
                    rules.custom_validator(value)
                except ValueError as e:
                    errors.append(f"{field}: {str(e)}")
            
            validated[field] = value
        
        # Check for unexpected fields
        for field in data:
            if field not in self.schema:
                errors.append(f"{field}: Unexpected field")
        
        return len(errors) == 0, validated, errors

# Define schemas
USER_SCHEMA = {
    'username': FieldRule(
        required=True,
        min_length=3,
        max_length=50,
        pattern=r'^[a-zA-Z0-9_]+$'
    ),
    'email': FieldRule(
        required=True,
        pattern=r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    ),
    'age': FieldRule(
        required=False,
        min_value=13,
        max_value=120
    ),
    'role': FieldRule(
        required=True,
        allowed_values=['user', 'admin', 'moderator']
    )
}

validator = SchemaValidator(USER_SCHEMA)
valid, data, errors = validator.validate({
    'username': 'john_doe',
    'email': 'john@example.com',
    'age': 25,
    'role': 'user'
})
```

## ❓ Interview Questions

### Basic
**Q: What's the difference between input validation and output encoding?**
A: Input validation checks data meets requirements before processing. Output encoding transforms data for safe display in specific contexts (HTML, JS, URL).

**Q: Why is client-side validation not sufficient?**
A: Easily bypassed - attackers can modify requests directly. Always validate server-side.

### Intermediate
**Q: Explain the principle of "validate on input, encode on output".**
A: Validate when receiving data (reject invalid). Encode when displaying data (prevent XSS in context). Different outputs need different encoding.

**Q: What is Content Security Policy and how does it help?**
A: HTTP header that controls which resources can be loaded. Mitigates XSS by restricting script sources, preventing inline scripts, controlling form destinations.

### Advanced
**Q: How do you handle validation for user-generated content (rich text)?**
A: (1) Use allowlist of permitted tags/attributes, (2) Use library like bleach, (3) Apply CSP, (4) Consider separate domain for user content, (5) Sanitize on input AND encode on output.