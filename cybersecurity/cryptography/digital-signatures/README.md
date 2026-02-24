# ✍️ Digital Signatures

Digital signatures provide authentication, integrity, and non-repudiation for digital messages.

## 📋 Table of Contents

- [Overview](#overview)
- [How Digital Signatures Work](#how-digital-signatures-work)
- [Signature Algorithms](#signature-algorithms)
- [Certificate Chains](#certificate-chains)
- [Implementation Examples](#implementation-examples)
- [Interview Questions](#interview-questions)

## 🎯 Overview

### Security Properties

| Property | Description |
|----------|-------------|
| **Authentication** | Verifies sender identity |
| **Integrity** | Detects message tampering |
| **Non-repudiation** | Sender cannot deny signing |
| **Unforgeable** | Only private key holder can sign |

### Process Flow

```
SIGNING:
┌──────────┐     ┌───────┐     ┌───────────┐     ┌──────────┐
│ Message  ├────►│  Hash ├────►│ Encrypt   ├────►│Signature │
└──────────┘     └───────┘     │with PrivK │     └──────────┘
                               └───────────┘

VERIFICATION:
┌──────────┐     ┌───────┐
│ Message  ├────►│  Hash ├────┐
└──────────┘     └───────┘    │        ┌──────────┐
                              ├───────►│  Match?  │
┌──────────┐     ┌───────────┐│        └──────────┘
│Signature ├────►│ Decrypt   ├┘
└──────────┘     │with PubK  │
                 └───────────┘
```

## 🔐 Signature Algorithms

### RSA Signatures

```python
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization
import hashlib

class RSASignature:
    def __init__(self):
        self.private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=2048
        )
        self.public_key = self.private_key.public_key()
    
    def sign(self, message: bytes) -> bytes:
        """Create signature"""
        signature = self.private_key.sign(
            message,
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()),
                salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256()
        )
        return signature
    
    def verify(self, message: bytes, signature: bytes) -> bool:
        """Verify signature"""
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
        except Exception:
            return False
    
    def export_public_key(self) -> bytes:
        return self.public_key.public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo
        )
```

### ECDSA (Elliptic Curve Digital Signature Algorithm)

```go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"
)

type ECDSASignature struct {
	R, S *big.Int
}

type Signer struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewSigner() (*Signer, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Signer{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

func (s *Signer) Sign(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	r, sig, err := ecdsa.Sign(rand.Reader, s.privateKey, hash[:])
	if err != nil {
		return nil, err
	}
	
	signature := ECDSASignature{R: r, S: sig}
	return asn1.Marshal(signature)
}

func (s *Signer) Verify(message, signatureBytes []byte) bool {
	var signature ECDSASignature
	if _, err := asn1.Unmarshal(signatureBytes, &signature); err != nil {
		return false
	}
	
	hash := sha256.Sum256(message)
	return ecdsa.Verify(s.publicKey, hash[:], signature.R, signature.S)
}

func (s *Signer) GetPublicKey() (x, y *big.Int) {
	return s.publicKey.X, s.publicKey.Y
}

func main() {
	signer, _ := NewSigner()
	
	message := []byte("Important document to sign")
	signature, _ := signer.Sign(message)
	
	valid := signer.Verify(message, signature)
	fmt.Printf("Signature valid: %v\n", valid)
	
	// Tampered message
	tampered := []byte("Important document to sign!")
	invalid := signer.Verify(tampered, signature)
	fmt.Printf("Tampered valid: %v\n", invalid)
}
```

### Ed25519 (Recommended)

```javascript
const crypto = require('crypto');

class Ed25519Signer {
    constructor() {
        const { privateKey, publicKey } = crypto.generateKeyPairSync('ed25519');
        this.privateKey = privateKey;
        this.publicKey = publicKey;
    }
    
    sign(message) {
        return crypto.sign(null, Buffer.from(message), this.privateKey);
    }
    
    verify(message, signature) {
        return crypto.verify(null, Buffer.from(message), this.publicKey, signature);
    }
    
    exportPublicKey() {
        return this.publicKey.export({ type: 'spki', format: 'pem' });
    }
}

// Usage
const signer = new Ed25519Signer();
const message = 'Critical transaction data';
const signature = signer.sign(message);

console.log('Valid signature:', signer.verify(message, signature));
console.log('Invalid message:', signer.verify('tampered', signature));
```

## 🔗 Certificate Chains

### X.509 Certificate Structure

```
Root CA (Self-signed)
    │
    ├── Intermediate CA
    │       │
    │       ├── Server Certificate
    │       └── Client Certificate
    │
    └── Another Intermediate CA
```

### Certificate Verification

```python
from cryptography import x509
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID
import datetime

class CertificateAuthority:
    def __init__(self, name: str, is_root: bool = True):
        self.private_key = rsa.generate_private_key(65537, 2048)
        self.public_key = self.private_key.public_key()
        self.name = name
        self.is_root = is_root
        self.certificate = None
        
        if is_root:
            self.certificate = self._create_self_signed()
    
    def _create_self_signed(self) -> x509.Certificate:
        subject = issuer = x509.Name([
            x509.NameAttribute(NameOID.COMMON_NAME, self.name),
        ])
        
        return x509.CertificateBuilder().subject_name(
            subject
        ).issuer_name(
            issuer
        ).public_key(
            self.public_key
        ).serial_number(
            x509.random_serial_number()
        ).not_valid_before(
            datetime.datetime.utcnow()
        ).not_valid_after(
            datetime.datetime.utcnow() + datetime.timedelta(days=365)
        ).sign(self.private_key, hashes.SHA256())
    
    def issue_certificate(self, subject_name: str, subject_public_key) -> x509.Certificate:
        subject = x509.Name([
            x509.NameAttribute(NameOID.COMMON_NAME, subject_name),
        ])
        
        return x509.CertificateBuilder().subject_name(
            subject
        ).issuer_name(
            self.certificate.subject
        ).public_key(
            subject_public_key
        ).serial_number(
            x509.random_serial_number()
        ).not_valid_before(
            datetime.datetime.utcnow()
        ).not_valid_after(
            datetime.datetime.utcnow() + datetime.timedelta(days=90)
        ).sign(self.private_key, hashes.SHA256())
    
    def verify_certificate(self, cert: x509.Certificate) -> bool:
        try:
            self.public_key.verify(
                cert.signature,
                cert.tbs_certificate_bytes,
                cert.signature_algorithm_parameters
            )
            return True
        except:
            return False

# Usage
root_ca = CertificateAuthority("Root CA")
intermediate = CertificateAuthority("Intermediate CA", is_root=False)
intermediate.certificate = root_ca.issue_certificate("Intermediate CA", intermediate.public_key)

server_key = rsa.generate_private_key(65537, 2048)
server_cert = intermediate.issue_certificate("server.example.com", server_key.public_key())

print(f"Server cert valid: {intermediate.verify_certificate(server_cert)}")
```

## ❓ Interview Questions

### Basic
**Q: What's the difference between encryption and digital signatures?**
A: Encryption provides confidentiality (only intended recipient can read). Signatures provide authentication and integrity (proves who sent it and that it wasn't modified).

**Q: Can you decrypt a digital signature?**
A: You "decrypt" with public key to verify, but you're not recovering secret data - you're comparing hash values.

### Intermediate
**Q: Why do we hash before signing?**
A: Signing algorithms work on fixed-size data. Hashing: (1) enables signing arbitrary-length messages, (2) is faster than signing large data, (3) provides additional security.

**Q: What's the difference between PSS and PKCS#1 v1.5 padding?**
A: PSS is probabilistic (includes random salt), provably secure. PKCS#1 v1.5 is deterministic, has known vulnerabilities (bleichenbacher attacks). Use PSS.

### Advanced
**Q: How do you implement non-repudiation in practice?**
A: Use trusted timestamp authorities, secure key storage (HSM), certificate revocation checking, and maintain audit logs. Consider legal requirements for digital signature validity.
