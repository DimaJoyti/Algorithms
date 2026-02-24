# 🌐 Network Security Protocols

Essential protocols for securing network communications.

## 📋 Table of Contents

- [TLS/SSL](#tlsssl)
- [VPN Protocols](#vpn-protocols)
- [SSH Protocol](#ssh-protocol)
- [DNS Security](#dns-security)
- [Implementation Examples](#implementation-examples)
- [Interview Questions](#interview-questions)

## 🔒 TLS/SSL

### TLS Handshake Process

```
Client                                          Server
  │                                               │
  │  1. ClientHello (version, ciphers, random)   │
  │─────────────────────────────────────────────►│
  │                                               │
  │  2. ServerHello (version, cipher, random)    │
  │◄─────────────────────────────────────────────│
  │  3. Certificate                               │
  │◄─────────────────────────────────────────────│
  │  4. ServerHelloDone                          │
  │◄─────────────────────────────────────────────│
  │                                               │
  │  5. ClientKeyExchange (pre-master secret)    │
  │─────────────────────────────────────────────►│
  │  6. ChangeCipherSpec                         │
  │─────────────────────────────────────────────►│
  │  7. Finished                                 │
  │─────────────────────────────────────────────►│
  │                                               │
  │  8. ChangeCipherSpec                         │
  │◄─────────────────────────────────────────────│
  │  9. Finished                                 │
  │◄─────────────────────────────────────────────│
  │                                               │
  │  ============ Encrypted Data =================│
```

### TLS Configuration Best Practices

```python
import ssl
import socket

def create_secure_context():
    """Create TLS context with security best practices"""
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    
    # Minimum TLS 1.2
    context.minimum_version = ssl.TLSVersion.TLSv1_2
    
    # Strong cipher suites only
    context.set_ciphers('ECDHE+AESGCM:DHE+AESGCM')
    
    # Require certificate verification
    context.verify_mode = ssl.CERT_REQUIRED
    
    # Load certificates
    context.load_cert_chain('server.crt', 'server.key')
    
    # Enable OCSP stapling
    context.check_hostname = True
    
    return context

def secure_client_connection(hostname, port):
    """Create secure client connection"""
    context = ssl.create_default_context()
    
    # Force TLS 1.2+
    context.minimum_version = ssl.TLSVersion.TLSv1_2
    
    # Verify certificate
    context.verify_mode = ssl.CERT_REQUIRED
    context.check_hostname = True
    
    with socket.create_connection((hostname, port)) as sock:
        with context.wrap_socket(sock, server_hostname=hostname) as secure_sock:
            print(f"TLS Version: {secure_sock.version()}")
            print(f"Cipher: {secure_sock.cipher()}")
            return secure_sock
```

### Certificate Pinning

```javascript
const https = require('crypto');
const crypto = require('crypto');

class CertificatePinning {
    constructor(expectedHashes) {
        this.expectedHashes = expectedHashes; // SHA-256 hashes
    }
    
    verifyCertificate(cert) {
        const certHash = crypto
            .createHash('sha256')
            .update(cert.raw)
            .digest('base64');
        
        if (!this.expectedHashes.includes(certHash)) {
            throw new Error('Certificate pin verification failed');
        }
        
        return true;
    }
    
    createPinnedAgent() {
        const agent = new https.Agent({
            checkServerIdentity: (host, cert) => {
                this.verifyCertificate(cert);
            }
        });
        return agent;
    }
}
```

## 🔐 VPN Protocols

### Protocol Comparison

| Protocol | Speed | Security | Use Case |
|----------|-------|----------|----------|
| OpenVPN | Medium | High | General purpose |
| WireGuard | Fast | High | Modern VPN |
| IPSec | Medium | High | Site-to-site |
| IKEv2 | Fast | High | Mobile devices |

### WireGuard Configuration

```go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type WireGuardKey struct {
	PrivateKey string
	PublicKey  string
}

func GenerateWireGuardKeys() (*WireGuardKey, error) {
	privateKey := make([]byte, 32)
	if _, err := rand.Read(privateKey); err != nil {
		return nil, err
	}
	
	// WireGuard private key tweak
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64
	
	// In real implementation, derive public key using Curve25519
	// publicKey := curve25519.X25519(privateKey, curve25519.Basepoint)
	
	return &WireGuardKey{
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
		PublicKey:  "derived_from_private_key", // Actual derivation needed
	}, nil
}

func GenerateWireGuardConfig(name string, privateKey, peerPublicKey, endpoint string) string {
	return fmt.Sprintf(`
[Interface]
PrivateKey = %s
Address = 10.0.0.2/24
DNS = 1.1.1.1

[Peer]
PublicKey = %s
Endpoint = %s
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
`, privateKey, peerPublicKey, endpoint)
}
```

## 🖥️ SSH Protocol

### SSH Hardening

```bash
# /etc/ssh/sshd_config hardening

# Disable root login
PermitRootLogin no

# Key-based authentication only
PubkeyAuthentication yes
PasswordAuthentication no
PermitEmptyPasswords no

# Specific algorithms
KexAlgorithms curve25519-sha256@libssh.org,diffie-hellman-group-exchange-sha256
Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com

# Limit users
AllowUsers user1 user2

# Disable forwarding
X11Forwarding no
AllowTcpForwarding no

# Port (non-default)
Port 2222

# Fail2Ban integration
MaxAuthTries 3
```

### SSH Key Management

```python
import subprocess
import os
from pathlib import Path

class SSHKeyManager:
    def __init__(self, key_dir='~/.ssh'):
        self.key_dir = Path(key_dir).expanduser()
        self.key_dir.mkdir(mode=0o700, exist_ok=True)
    
    def generate_keypair(self, name: str, key_type: str = 'ed25519', 
                         comment: str = '', passphrase: str = ''):
        """Generate SSH keypair"""
        key_path = self.key_dir / name
        
        cmd = ['ssh-keygen', '-t', key_type, '-f', str(key_path), '-N', passphrase]
        if comment:
            cmd.extend(['-C', comment])
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            raise Exception(f"Key generation failed: {result.stderr}")
        
        return {
            'private_key': key_path,
            'public_key': f"{key_path}.pub",
            'fingerprint': self._get_fingerprint(key_path)
        }
    
    def _get_fingerprint(self, key_path):
        result = subprocess.run(
            ['ssh-keygen', '-lf', str(key_path)],
            capture_output=True, text=True
        )
        return result.stdout.strip()
    
    def add_to_agent(self, key_path, passphrase=''):
        """Add key to SSH agent"""
        # In production, use ssh-add with proper handling
        pass
    
    def authorize_key(self, public_key: str, user: str = None):
        """Add public key to authorized_keys"""
        auth_file = self.key_dir / 'authorized_keys'
        auth_file.touch(mode=0o600)
        
        with open(auth_file, 'a') as f:
            f.write(public_key + '\n')
```

## 🔍 DNS Security

### DNSSEC

```python
import dns.resolver
import dns.dnssec

class DNSSECValidator:
    def __init__(self):
        self.resolver = dns.resolver.Resolver()
        self.resolver.use_edns = True
    
    def verify_dnssec(self, domain: str, record_type: str = 'A'):
        """Verify DNSSEC for a domain"""
        try:
            # Get the answer
            answer = self.resolver.resolve(domain, record_type)
            
            # Get DNSKEY
            dnskey = self.resolver.resolve(domain, 'DNSKEY')
            
            # Get RRSIG
            rrsig = self.resolver.resolve(domain, 'RRSIG')
            
            # Verify signature
            # In production, use proper DNSSEC validation chain
            return {
                'domain': domain,
                'dnssec_enabled': True,
                'dnskey_count': len(dnskey.rrset),
                'validated': True
            }
        except Exception as e:
            return {
                'domain': domain,
                'dnssec_enabled': False,
                'error': str(e)
            }
```

### DNS over HTTPS (DoH)

```javascript
const https = require('https');

class DNSOverHTTPS {
    constructor(server = 'https://cloudflare-dns.com/dns-query') {
        this.server = server;
    }
    
    async resolve(domain, type = 'A') {
        const url = `${this.server}?name=${domain}&type=${type}`;
        
        return new Promise((resolve, reject) => {
            https.get(url, {
                headers: { 'Accept': 'application/dns-json' }
            }, (res) => {
                let data = '';
                res.on('data', chunk => data += chunk);
                res.on('end', () => {
                    try {
                        resolve(JSON.parse(data));
                    } catch (e) {
                        reject(e);
                    }
                });
            }).on('error', reject);
        });
    }
}

// Usage
const doh = new DNSOverHTTPS();
doh.resolve('example.com').then(result => {
    console.log('DNS Response:', result.Answer);
});
```

## ❓ Interview Questions

### Basic
**Q: What's the difference between TLS 1.2 and TLS 1.3?**
A: TLS 1.3: (1) Faster - fewer round trips, (2) Removed weak algorithms, (3) Perfect forward secrecy mandatory, (4) 0-RTT mode for resumed connections.

**Q: Why disable SSH root login?**
A: Reduces attack surface - root is a known target, limits brute force attempts, enforces principle of least privilege.

### Intermediate
**Q: Explain certificate pinning and when to use it.**
A: Pinning hardcodes expected certificate/public key in app, preventing MITM even with compromised CA. Use for high-security apps, APIs with known certificates.

**Q: What's the difference between DNSSEC and DoH?**
A: DNSSEC provides authenticity (verifies DNS responses), DoH provides confidentiality (encrypts DNS queries). Both can be used together.

### Advanced
**Q: How would you implement mutual TLS (mTLS)?**
A: (1) Both client and server have certificates, (2) Server verifies client cert during handshake, (3) Client verifies server cert, (4) Use for service-to-service auth, zero-trust networks.
