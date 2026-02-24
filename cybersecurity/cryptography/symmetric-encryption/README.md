# 🔐 Symmetric Encryption

Symmetric encryption uses the same key for both encryption and decryption. It's fast and efficient for large data volumes.

## 📋 Table of Contents

- [Overview](#overview)
- [Common Algorithms](#common-algorithms)
- [Key Management](#key-management)
- [Modes of Operation](#modes-of-operation)
- [Implementation Examples](#implementation-examples)
- [Interview Questions](#interview-questions)

## 🎯 Overview

```
┌─────────┐    Key    ┌─────────┐    Key    ┌─────────┐
│Plaintext├──────────►│ Encrypt ├──────────►│Ciphertext│
└─────────┘           └─────────┘           └──────────┘
                                              │
                      ┌─────────┐             ▼
                      │Decrypt  ├──────────► Plaintext
                      └─────────┘
                           ▲
                           │
                          Key
```

### Characteristics
- **Speed**: 100-1000x faster than asymmetric
- **Key Size**: Typically 128-256 bits
- **Key Distribution**: Major challenge - same key must be shared securely
- **Use Cases**: Data at rest, bulk encryption, session keys

## 🔧 Common Algorithms

### AES (Advanced Encryption Standard)
- **Key Sizes**: 128, 192, 256 bits
- **Block Size**: 128 bits
- **Status**: Current standard, NIST approved

### ChaCha20
- **Key Size**: 256 bits
- **Stream Cipher**: No block size limitation
- **Status**: Preferred for mobile, TLS 1.3

### 3DES (Triple DES)
- **Key Size**: 168 bits (effective 112 bits)
- **Status**: Deprecated, being phased out

## 🔑 Key Management

```python
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import os
import base64

class SymmetricEncryption:
    def __init__(self):
        self.key = None
        self.cipher = None
    
    def generate_key(self):
        """Generate a random encryption key"""
        self.key = Fernet.generate_key()
        self.cipher = Fernet(self.key)
        return self.key
    
    def derive_key_from_password(self, password: str, salt: bytes = None):
        """Derive encryption key from password"""
        if salt is None:
            salt = os.urandom(16)
        
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=480000,
        )
        key = base64.urlsafe_b64encode(kdf.derive(password.encode()))
        self.cipher = Fernet(key)
        return salt
    
    def encrypt(self, plaintext: str) -> bytes:
        """Encrypt plaintext"""
        return self.cipher.encrypt(plaintext.encode())
    
    def decrypt(self, ciphertext: bytes) -> str:
        """Decrypt ciphertext"""
        return self.cipher.decrypt(ciphertext).decode()
    
    def encrypt_file(self, filepath: str) -> bytes:
        """Encrypt file contents"""
        with open(filepath, 'rb') as f:
            data = f.read()
        return self.cipher.encrypt(data)
    
    def decrypt_file(self, ciphertext: bytes, output_path: str):
        """Decrypt to file"""
        data = self.cipher.decrypt(ciphertext)
        with open(output_path, 'wb') as f:
            f.write(data)

# Usage
enc = SymmetricEncryption()
key = enc.generate_key()
encrypted = enc.encrypt("Sensitive data")
decrypted = enc.decrypt(encrypted)
```

## 📊 Modes of Operation

| Mode | Description | Use Case | Parallelizable |
|------|-------------|----------|----------------|
| ECB | Electronic Codebook | Never use | Yes |
| CBC | Cipher Block Chaining | File encryption | Decryption only |
| CTR | Counter | Streaming | Yes |
| GCM | Galois/Counter | Network, AEAD | Yes |
| XTS | XEX-based | Disk encryption | Yes |

### GCM Mode (Recommended)

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type AESGCM struct {
	key   []byte
	aead  cipher.AEAD
}

func NewAESGCM(key []byte) (*AESGCM, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	return &AESGCM{key: key, aead: aead}, nil
}

func (a *AESGCM) Encrypt(plaintext []byte) (ciphertext, nonce []byte, err error) {
	nonce = make([]byte, a.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	
	ciphertext = a.aead.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func (a *AESGCM) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	return a.aead.Open(nil, nonce, ciphertext, nil)
}

func main() {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	
	aesgcm, err := NewAESGCM(key)
	if err != nil {
		panic(err)
	}
	
	plaintext := []byte("Secret message")
	ciphertext, nonce, err := aesgcm.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Encrypted: %s\n", hex.EncodeToString(ciphertext))
	
	decrypted, err := aesgcm.Decrypt(ciphertext, nonce)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Decrypted: %s\n", decrypted)
}
```

## ❓ Interview Questions

### Basic
**Q: What's the main challenge with symmetric encryption?**
A: Key distribution - how to securely share the same key between parties.

**Q: Why is AES preferred over DES?**
A: AES has larger key sizes (128-256 vs 56 bits), is faster, and has no known practical attacks.

### Intermediate
**Q: When would you use CBC vs GCM mode?**
A: GCM provides authenticated encryption (confidentiality + integrity) and is parallelizable. Use GCM for most cases. CBC only provides confidentiality and requires padding.

**Q: How do you securely store encryption keys?**
A: Use a key management system (KMS), hardware security module (HSM), or derive from password with strong KDF (PBKDF2, Argon2).

### Advanced
**Q: Explain the difference between block ciphers and stream ciphers.**
A: Block ciphers encrypt fixed-size blocks (e.g., 128 bits for AES) and need modes for data larger than block size. Stream ciphers encrypt one bit/byte at a time, suitable for streaming data.
