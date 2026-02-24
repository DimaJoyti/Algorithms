# 🔐 Hash Functions - Data Integrity and Authentication

Hash functions are mathematical algorithms that transform input data into fixed-size strings, providing data integrity verification and authentication mechanisms.

## 📋 Table of Contents

- [Overview](#overview)
- [Properties of Cryptographic Hash Functions](#properties)
- [Common Hash Algorithms](#common-hash-algorithms)
- [Implementation Examples](#implementation-examples)
- [Security Considerations](#security-considerations)
- [Practical Applications](#practical-applications)
- [Interview Questions](#interview-questions)

## 🎯 Overview

Hash functions are fundamental building blocks in cryptography, used for:
- **Data Integrity**: Verify data hasn't been modified
- **Password Storage**: Securely store user passwords
- **Digital Signatures**: Create message digests for signing
- **Proof of Work**: Blockchain and cryptocurrency mining
- **Data Deduplication**: Identify duplicate files

### Hash Function Properties

A cryptographic hash function must have these properties:
1. **Deterministic**: Same input always produces same output
2. **Fixed Output Size**: Output length is constant regardless of input size
3. **Fast Computation**: Efficient to compute hash value
4. **Avalanche Effect**: Small input change causes large output change
5. **Pre-image Resistance**: Hard to find input from hash output
6. **Second Pre-image Resistance**: Hard to find different input with same hash
7. **Collision Resistance**: Hard to find two inputs with same hash

## 🔢 Common Hash Algorithms

### SHA-256 (Secure Hash Algorithm)
- **Output Size**: 256 bits (32 bytes)
- **Security**: Currently secure, widely used
- **Applications**: Bitcoin, TLS certificates, digital signatures

### SHA-3 (Keccak)
- **Output Size**: Variable (224, 256, 384, 512 bits)
- **Security**: Latest NIST standard, different construction than SHA-2
- **Applications**: New systems requiring latest standards

### MD5 (Message Digest 5)
- **Output Size**: 128 bits (16 bytes)
- **Security**: ⚠️ BROKEN - vulnerable to collision attacks
- **Applications**: Legacy systems, file integrity (non-security)

### bcrypt
- **Purpose**: Password hashing with built-in salt
- **Security**: Adaptive cost, resistant to rainbow table attacks
- **Applications**: User password storage

## 💻 Implementation Examples

### Python Hash Functions

```python
import hashlib
import hmac
import bcrypt
import secrets
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import base64

class HashFunctions:
    @staticmethod
    def sha256_hash(data):
        """Compute SHA-256 hash of data"""
        if isinstance(data, str):
            data = data.encode('utf-8')
        return hashlib.sha256(data).hexdigest()
    
    @staticmethod
    def sha3_hash(data, output_size=256):
        """Compute SHA-3 hash with specified output size"""
        if isinstance(data, str):
            data = data.encode('utf-8')
        
        if output_size == 224:
            return hashlib.sha3_224(data).hexdigest()
        elif output_size == 256:
            return hashlib.sha3_256(data).hexdigest()
        elif output_size == 384:
            return hashlib.sha3_384(data).hexdigest()
        elif output_size == 512:
            return hashlib.sha3_512(data).hexdigest()
        else:
            raise ValueError("Invalid output size")
    
    @staticmethod
    def hmac_hash(key, message, algorithm='sha256'):
        """Compute HMAC (Hash-based Message Authentication Code)"""
        if isinstance(key, str):
            key = key.encode('utf-8')
        if isinstance(message, str):
            message = message.encode('utf-8')
        
        if algorithm == 'sha256':
            return hmac.new(key, message, hashlib.sha256).hexdigest()
        elif algorithm == 'sha512':
            return hmac.new(key, message, hashlib.sha512).hexdigest()
        else:
            raise ValueError("Unsupported algorithm")
    
    @staticmethod
    def bcrypt_hash(password, rounds=12):
        """Hash password using bcrypt with salt"""
        if isinstance(password, str):
            password = password.encode('utf-8')
        
        salt = bcrypt.gensalt(rounds=rounds)
        return bcrypt.hashpw(password, salt).decode('utf-8')
    
    @staticmethod
    def verify_bcrypt(password, hashed):
        """Verify password against bcrypt hash"""
        if isinstance(password, str):
            password = password.encode('utf-8')
        if isinstance(hashed, str):
            hashed = hashed.encode('utf-8')
        
        return bcrypt.checkpw(password, hashed)
    
    @staticmethod
    def pbkdf2_hash(password, salt=None, iterations=100000, key_length=32):
        """Derive key using PBKDF2"""
        if isinstance(password, str):
            password = password.encode('utf-8')
        
        if salt is None:
            salt = secrets.token_bytes(16)
        elif isinstance(salt, str):
            salt = salt.encode('utf-8')
        
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=key_length,
            salt=salt,
            iterations=iterations,
        )
        
        key = kdf.derive(password)
        return {
            'key': base64.b64encode(key).decode('utf-8'),
            'salt': base64.b64encode(salt).decode('utf-8'),
            'iterations': iterations
        }

# Example usage and testing
def demonstrate_hash_functions():
    hash_func = HashFunctions()
    
    # Test data
    message = "Hello, Cybersecurity Interview!"
    password = "SecurePassword123!"
    
    print("=== Hash Function Demonstrations ===\n")
    
    # SHA-256
    sha256_result = hash_func.sha256_hash(message)
    print(f"SHA-256: {sha256_result}")
    
    # SHA-3
    sha3_result = hash_func.sha3_hash(message, 256)
    print(f"SHA-3-256: {sha3_result}")
    
    # HMAC
    key = "secret_key_for_hmac"
    hmac_result = hash_func.hmac_hash(key, message)
    print(f"HMAC-SHA256: {hmac_result}")
    
    # bcrypt password hashing
    bcrypt_hash = hash_func.bcrypt_hash(password)
    print(f"bcrypt hash: {bcrypt_hash}")
    
    # Verify bcrypt
    is_valid = hash_func.verify_bcrypt(password, bcrypt_hash)
    print(f"bcrypt verification: {is_valid}")
    
    # PBKDF2
    pbkdf2_result = hash_func.pbkdf2_hash(password)
    print(f"PBKDF2 key: {pbkdf2_result['key']}")
    print(f"PBKDF2 salt: {pbkdf2_result['salt']}")
    
    # Demonstrate avalanche effect
    print("\n=== Avalanche Effect Demonstration ===")
    original = "Hello World"
    modified = "Hello World!"  # Just added one character
    
    hash1 = hash_func.sha256_hash(original)
    hash2 = hash_func.sha256_hash(modified)
    
    print(f"Original: '{original}'")
    print(f"Hash: {hash1}")
    print(f"Modified: '{modified}'")
    print(f"Hash: {hash2}")
    
    # Calculate bit differences
    diff_bits = sum(c1 != c2 for c1, c2 in zip(hash1, hash2))
    print(f"Different characters: {diff_bits} out of {len(hash1)}")

if __name__ == "__main__":
    demonstrate_hash_functions()
```

### Go Hash Implementation

```go
package main

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/rand"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/pbkdf2"
    "golang.org/x/crypto/sha3"
)

type HashService struct{}

func (hs *HashService) SHA256Hash(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}

func (hs *HashService) SHA3Hash(data []byte) string {
    hash := sha3.Sum256(data)
    return hex.EncodeToString(hash[:])
}

func (hs *HashService) HMACHash(key, message []byte) string {
    h := hmac.New(sha256.New, key)
    h.Write(message)
    return hex.EncodeToString(h.Sum(nil))
}

func (hs *HashService) BcryptHash(password string, cost int) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
    return string(bytes), err
}

func (hs *HashService) VerifyBcrypt(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (hs *HashService) PBKDF2Hash(password, salt []byte, iterations, keyLength int) []byte {
    return pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
}

func (hs *HashService) GenerateSalt(length int) ([]byte, error) {
    salt := make([]byte, length)
    _, err := rand.Read(salt)
    return salt, err
}

// File integrity checker
func (hs *HashService) FileIntegrityChecker(files map[string][]byte) map[string]string {
    manifest := make(map[string]string)
    for filename, content := range files {
        manifest[filename] = hs.SHA256Hash(content)
    }
    return manifest
}

func (hs *HashService) VerifyFileIntegrity(files map[string][]byte, manifest map[string]string) map[string]bool {
    results := make(map[string]bool)
    for filename, content := range files {
        expectedHash, exists := manifest[filename]
        if !exists {
            results[filename] = false
            continue
        }
        actualHash := hs.SHA256Hash(content)
        results[filename] = actualHash == expectedHash
    }
    return results
}

func main() {
    hs := &HashService{}
    
    message := []byte("Hello, Cybersecurity Interview!")
    password := "SecurePassword123!"
    
    fmt.Println("=== Go Hash Function Demonstrations ===\n")
    
    // SHA-256
    sha256Result := hs.SHA256Hash(message)
    fmt.Printf("SHA-256: %s\n", sha256Result)
    
    // SHA-3
    sha3Result := hs.SHA3Hash(message)
    fmt.Printf("SHA-3-256: %s\n", sha3Result)
    
    // HMAC
    key := []byte("secret_key_for_hmac")
    hmacResult := hs.HMACHash(key, message)
    fmt.Printf("HMAC-SHA256: %s\n", hmacResult)
    
    // bcrypt
    bcryptHash, err := hs.BcryptHash(password, 12)
    if err != nil {
        fmt.Printf("bcrypt error: %v\n", err)
    } else {
        fmt.Printf("bcrypt hash: %s\n", bcryptHash)
        
        // Verify bcrypt
        isValid := hs.VerifyBcrypt(password, bcryptHash)
        fmt.Printf("bcrypt verification: %t\n", isValid)
    }
    
    // PBKDF2
    salt, _ := hs.GenerateSalt(16)
    pbkdf2Key := hs.PBKDF2Hash([]byte(password), salt, 100000, 32)
    fmt.Printf("PBKDF2 key: %s\n", hex.EncodeToString(pbkdf2Key))
    fmt.Printf("PBKDF2 salt: %s\n", hex.EncodeToString(salt))
    
    // File integrity example
    fmt.Println("\n=== File Integrity Checking ===")
    files := map[string][]byte{
        "config.txt":  []byte("server_port=8080\ndb_host=localhost"),
        "secrets.env": []byte("API_KEY=secret123\nDB_PASSWORD=pass456"),
    }
    
    manifest := hs.FileIntegrityChecker(files)
    fmt.Println("File manifest created:")
    for filename, hash := range manifest {
        fmt.Printf("  %s: %s\n", filename, hash)
    }
    
    // Simulate file modification
    files["config.txt"] = []byte("server_port=9090\ndb_host=localhost")
    
    results := hs.VerifyFileIntegrity(files, manifest)
    fmt.Println("\nIntegrity check results:")
    for filename, isValid := range results {
        status := "VALID"
        if !isValid {
            status = "MODIFIED"
        }
        fmt.Printf("  %s: %s\n", filename, status)
    }
}
```

### JavaScript Hash Implementation

```javascript
const crypto = require('crypto');
const bcrypt = require('bcrypt');

class HashFunctions {
    static sha256Hash(data) {
        return crypto.createHash('sha256').update(data).digest('hex');
    }
    
    static sha512Hash(data) {
        return crypto.createHash('sha512').update(data).digest('hex');
    }
    
    static hmacHash(key, message, algorithm = 'sha256') {
        return crypto.createHmac(algorithm, key).update(message).digest('hex');
    }
    
    static async bcryptHash(password, rounds = 12) {
        return await bcrypt.hash(password, rounds);
    }
    
    static async verifyBcrypt(password, hash) {
        return await bcrypt.compare(password, hash);
    }
    
    static pbkdf2Hash(password, salt, iterations = 100000, keyLength = 32) {
        return new Promise((resolve, reject) => {
            crypto.pbkdf2(password, salt, iterations, keyLength, 'sha256', (err, derivedKey) => {
                if (err) reject(err);
                else resolve({
                    key: derivedKey.toString('hex'),
                    salt: salt.toString('hex'),
                    iterations: iterations
                });
            });
        });
    }
    
    static generateSalt(length = 16) {
        return crypto.randomBytes(length);
    }
    
    // Hash-based rate limiting
    static createRateLimitKey(ip, endpoint, timeWindow) {
        const timestamp = Math.floor(Date.now() / (timeWindow * 1000));
        const data = `${ip}:${endpoint}:${timestamp}`;
        return this.sha256Hash(data);
    }
    
    // Merkle tree implementation for data integrity
    static buildMerkleTree(data) {
        if (data.length === 0) return null;
        if (data.length === 1) return this.sha256Hash(data[0]);
        
        const hashes = data.map(item => this.sha256Hash(item));
        
        while (hashes.length > 1) {
            const newLevel = [];
            for (let i = 0; i < hashes.length; i += 2) {
                const left = hashes[i];
                const right = hashes[i + 1] || left; // Handle odd number of nodes
                const combined = this.sha256Hash(left + right);
                newLevel.push(combined);
            }
            hashes.splice(0, hashes.length, ...newLevel);
        }
        
        return hashes[0];
    }
}

// Password strength checker using hash analysis
class PasswordAnalyzer {
    static async analyzePassword(password) {
        const analysis = {
            length: password.length,
            hasUppercase: /[A-Z]/.test(password),
            hasLowercase: /[a-z]/.test(password),
            hasNumbers: /\d/.test(password),
            hasSpecialChars: /[!@#$%^&*(),.?":{}|<>]/.test(password),
            entropy: this.calculateEntropy(password),
            commonPassword: await this.checkCommonPassword(password)
        };
        
        analysis.score = this.calculateScore(analysis);
        analysis.strength = this.getStrengthLevel(analysis.score);
        
        return analysis;
    }
    
    static calculateEntropy(password) {
        const charFreq = {};
        for (let char of password) {
            charFreq[char] = (charFreq[char] || 0) + 1;
        }
        
        let entropy = 0;
        const length = password.length;
        
        for (let freq of Object.values(charFreq)) {
            const probability = freq / length;
            entropy -= probability * Math.log2(probability);
        }
        
        return entropy;
    }
    
    static async checkCommonPassword(password) {
        // In real implementation, check against common password database
        const commonPasswords = [
            'password', '123456', 'password123', 'admin', 'qwerty',
            'letmein', 'welcome', 'monkey', '1234567890'
        ];
        
        const passwordHash = HashFunctions.sha256Hash(password.toLowerCase());
        const commonHashes = commonPasswords.map(p => HashFunctions.sha256Hash(p));
        
        return commonHashes.includes(passwordHash);
    }
    
    static calculateScore(analysis) {
        let score = 0;
        
        // Length scoring
        if (analysis.length >= 12) score += 25;
        else if (analysis.length >= 8) score += 15;
        else if (analysis.length >= 6) score += 5;
        
        // Character variety
        if (analysis.hasUppercase) score += 15;
        if (analysis.hasLowercase) score += 15;
        if (analysis.hasNumbers) score += 15;
        if (analysis.hasSpecialChars) score += 20;
        
        // Entropy bonus
        if (analysis.entropy > 3.5) score += 10;
        
        // Common password penalty
        if (analysis.commonPassword) score -= 50;
        
        return Math.max(0, Math.min(100, score));
    }
    
    static getStrengthLevel(score) {
        if (score >= 80) return 'Very Strong';
        if (score >= 60) return 'Strong';
        if (score >= 40) return 'Moderate';
        if (score >= 20) return 'Weak';
        return 'Very Weak';
    }
}

// Demonstration
async function demonstrateHashFunctions() {
    console.log('=== JavaScript Hash Function Demonstrations ===\n');
    
    const message = 'Hello, Cybersecurity Interview!';
    const password = 'SecurePassword123!';
    
    // Basic hash functions
    console.log(`SHA-256: ${HashFunctions.sha256Hash(message)}`);
    console.log(`SHA-512: ${HashFunctions.sha512Hash(message)}`);
    
    // HMAC
    const key = 'secret_key_for_hmac';
    console.log(`HMAC-SHA256: ${HashFunctions.hmacHash(key, message)}`);
    
    // bcrypt
    const bcryptHash = await HashFunctions.bcryptHash(password);
    console.log(`bcrypt hash: ${bcryptHash}`);
    
    const isValid = await HashFunctions.verifyBcrypt(password, bcryptHash);
    console.log(`bcrypt verification: ${isValid}`);
    
    // PBKDF2
    const salt = HashFunctions.generateSalt();
    const pbkdf2Result = await HashFunctions.pbkdf2Hash(password, salt);
    console.log(`PBKDF2 key: ${pbkdf2Result.key}`);
    console.log(`PBKDF2 salt: ${pbkdf2Result.salt}`);
    
    // Merkle tree
    const data = ['file1.txt', 'file2.txt', 'file3.txt', 'file4.txt'];
    const merkleRoot = HashFunctions.buildMerkleTree(data);
    console.log(`Merkle root: ${merkleRoot}`);
    
    // Password analysis
    console.log('\n=== Password Analysis ===');
    const analysis = await PasswordAnalyzer.analyzePassword(password);
    console.log(`Password: ${password}`);
    console.log(`Strength: ${analysis.strength} (Score: ${analysis.score})`);
    console.log(`Entropy: ${analysis.entropy.toFixed(2)} bits`);
    console.log(`Common password: ${analysis.commonPassword}`);
}

// Run demonstration
demonstrateHashFunctions().catch(console.error);
```

## 🔒 Security Considerations

### Hash Function Selection
- **Use SHA-256 or SHA-3** for new applications
- **Avoid MD5 and SHA-1** for security purposes
- **Use bcrypt or Argon2** for password hashing
- **Consider PBKDF2** for key derivation

### Common Vulnerabilities
- **Rainbow Table Attacks**: Use salts with password hashing
- **Length Extension Attacks**: Use HMAC instead of simple concatenation
- **Timing Attacks**: Use constant-time comparison functions
- **Collision Attacks**: Choose collision-resistant algorithms

### Best Practices
- **Always use salt** for password hashing
- **Use appropriate work factors** (bcrypt rounds, PBKDF2 iterations)
- **Implement proper key management** for HMAC keys
- **Regular security updates** for cryptographic libraries

## ❓ Interview Questions

### Basic Questions
1. **Q**: What is the difference between encryption and hashing?
   **A**: Encryption is reversible with the correct key, while hashing is a one-way function. Encryption provides confidentiality, while hashing provides integrity verification.

2. **Q**: Why should you never use MD5 for security purposes?
   **A**: MD5 is vulnerable to collision attacks, meaning attackers can create two different inputs that produce the same hash. This breaks the collision resistance property required for security.

### Intermediate Questions
3. **Q**: Explain why salts are important in password hashing.
   **A**: Salts prevent rainbow table attacks by ensuring each password has a unique hash, even if multiple users have the same password. They add randomness that makes precomputed attack tables ineffective.

4. **Q**: What is HMAC and when would you use it?
   **A**: HMAC (Hash-based Message Authentication Code) provides both data integrity and authenticity using a secret key. Use it when you need to verify that data hasn't been tampered with and comes from a trusted source.

---

**🔑 Key Takeaway**: Hash functions are fundamental to cybersecurity, but choosing the right algorithm and implementing it correctly is crucial for security. Always use current best practices and avoid deprecated algorithms.
