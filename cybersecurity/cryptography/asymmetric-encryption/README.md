# 🔑 Asymmetric Encryption

Asymmetric encryption uses a pair of keys: public key for encryption, private key for decryption. Solves key distribution problem.

## 📋 Table of Contents

- [Overview](#overview)
- [RSA Algorithm](#rsa-algorithm)
- [Elliptic Curve Cryptography](#elliptic-curve-cryptography)
- [Key Exchange Protocols](#key-exchange-protocols)
- [Implementation Examples](#implementation-examples)
- [Interview Questions](#interview-questions)

## 🎯 Overview

```
                    Public Key (Shared)
                          │
┌─────────┐               ▼          ┌─────────┐
│Plaintext├──────────► Encrypt ─────►│Ciphertext│
└─────────┘                          └──────────┘
                                          │
                    ┌─────────┐           ▼
                    │Decrypt  ├───────► Plaintext
                    └─────────┘
                         ▲
                         │
                   Private Key (Secret)
```

### Characteristics
- **Speed**: Slower than symmetric (100-1000x)
- **Key Size**: Larger (RSA: 2048-4096 bits, ECC: 256-521 bits)
- **Key Distribution**: Public key can be shared openly
- **Use Cases**: Key exchange, digital signatures, small data

## 🔧 RSA Algorithm

Based on difficulty of factoring large prime numbers.

```python
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization

class RSAEncryption:
    def __init__(self, key_size=2048):
        self.private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=key_size
        )
        self.public_key = self.private_key.public_key()
    
    def get_public_key_pem(self):
        """Export public key in PEM format"""
        return self.public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
    
    def encrypt(self, plaintext: bytes) -> bytes:
        """Encrypt with public key"""
        return self.public_key.encrypt(
            plaintext,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
    
    def decrypt(self, ciphertext: bytes) -> bytes:
        """Decrypt with private key"""
        return self.private_key.decrypt(
            ciphertext,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
    
    def sign(self, message: bytes) -> bytes:
        """Create digital signature"""
        return self.private_key.sign(
            message,
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
    
    def verify(self, message: bytes, signature: bytes) -> bool:
        """Verify digital signature"""
        try:
            self.public_key.verify(
                signature,
                message,
                padding.PSS(
                    mgf=padding.MGF1(hashes.SHA256()),
                    salt_length=padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            return True
        except:
            return False

# Usage
rsa_enc = RSAEncryption()
encrypted = rsa_enc.encrypt(b"Secret message")
decrypted = rsa_enc.decrypt(encrypted)
print(f"Decrypted: {decrypted}")
```

## 📐 Elliptic Curve Cryptography (ECC)

Same security with smaller keys. ECDH for key exchange, ECDSA for signatures.

```go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

type ECCKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func GenerateECCKeyPair() (*ECCKeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	
	return &ECCKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (kp *ECCKeyPair) ECDH(peerPublicKey *ecdsa.PublicKey) ([]byte, error) {
	// Simplified ECDH - in production use crypto/ecdh
	curve := elliptic.P256()
	
	x, _ := curve.ScalarMult(peerPublicKey.X, peerPublicKey.Y, kp.PrivateKey.D.Bytes())
	
	sharedSecret := sha256.Sum256(x.Bytes())
	return sharedSecret[:], nil
}

func (kp *ECCKeyPair) Sign(message []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(message)
	return ecdsa.Sign(rand.Reader, kp.PrivateKey, hash[:])
}

func (kp *ECCKeyPair) Verify(message []byte, r, s *big.Int) bool {
	hash := sha256.Sum256(message)
	return ecdsa.Verify(kp.PublicKey, hash[:], r, s)
}

func main() {
	alice, _ := GenerateECCKeyPair()
	bob, _ := GenerateECCKeyPair()
	
	// ECDH Key Exchange
	aliceSecret, _ := alice.ECDH(bob.PublicKey)
	bobSecret, _ := bob.ECDH(alice.PublicKey)
	
	fmt.Printf("Alice secret: %s\n", hex.EncodeToString(aliceSecret)[:16]+"...")
	fmt.Printf("Bob secret:   %s\n", hex.EncodeToString(bobSecret)[:16]+"...")
	fmt.Printf("Shared secret match: %v\n", hex.EncodeToString(aliceSecret) == hex.EncodeToString(bobSecret))
	
	// ECDSA Signature
	message := []byte("Important message")
	r, s, _ := alice.Sign(message)
	valid := alice.Verify(message, r, s)
	fmt.Printf("Signature valid: %v\n", valid)
}
```

## 🔄 Key Exchange Protocols

### Diffie-Hellman Key Exchange

```javascript
const crypto = require('crypto');

class DiffieHellman {
    constructor() {
        this.dh = crypto.createDiffieHellman(2048);
        this.keys = this.dh.generateKeys();
    }
    
    getPublicKey() {
        return this.dh.getPublicKey('hex');
    }
    
    computeSharedSecret(otherPublicKey) {
        return this.dh.computeSecret(otherPublicKey, 'hex', 'hex');
    }
}

// Usage
const alice = new DiffieHellman();
const bob = new DiffieHellman();

const aliceSecret = alice.computeSharedSecret(bob.getPublicKey());
const bobSecret = bob.computeSharedSecret(alice.getPublicKey());

console.log('Secrets match:', aliceSecret === bobSecret);
```

### Hybrid Encryption (Best Practice)

Combine symmetric + asymmetric for efficiency:

```python
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
import os

class HybridEncryption:
    def __init__(self):
        self.rsa_private = rsa.generate_private_key(65537, 2048)
        self.rsa_public = self.rsa_private.public_key()
    
    def encrypt(self, plaintext: bytes) -> dict:
        # 1. Generate random symmetric key
        sym_key = Fernet.generate_key()
        cipher = Fernet(sym_key)
        
        # 2. Encrypt data with symmetric key
        encrypted_data = cipher.encrypt(plaintext)
        
        # 3. Encrypt symmetric key with RSA
        encrypted_key = self.rsa_public.encrypt(
            sym_key,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
        
        return {
            'encrypted_data': encrypted_data,
            'encrypted_key': encrypted_key
        }
    
    def decrypt(self, encrypted_data: bytes, encrypted_key: bytes) -> bytes:
        # 1. Decrypt symmetric key
        sym_key = self.rsa_private.decrypt(
            encrypted_key,
            padding.OAEP(
                mgf=padding.MGF1(algorithm=hashes.SHA256()),
                algorithm=hashes.SHA256(),
                label=None
            )
        )
        
        # 2. Decrypt data
        cipher = Fernet(sym_key)
        return cipher.decrypt(encrypted_data)
```

## ❓ Interview Questions

### Basic
**Q: What's the main advantage of asymmetric encryption?**
A: Solves key distribution problem - public key can be shared openly, private key stays secret.

**Q: Why not use asymmetric encryption for all data?**
A: Too slow - 100-1000x slower than symmetric. Use hybrid approach instead.

### Intermediate
**Q: When would you choose ECC over RSA?**
A: ECC provides same security with smaller keys (256-bit ECC ≈ 3072-bit RSA). Use for mobile, IoT, TLS certificates.

**Q: What is perfect forward secrecy?**
A: Compromise of long-term private key doesn't compromise past session keys. Achieved with ephemeral key exchange (ECDHE).

### Advanced
**Q: Explain the padding oracle attack.**
A: Attacker exploits error messages about padding validity to decrypt ciphertext byte-by-byte. Prevented by using authenticated encryption (GCM) or verifying MAC before padding.
