# 🔒 Security Engineer Interview Questions

Comprehensive collection of interview questions for Security Engineer positions, organized by difficulty and topic area.

## 📋 Table of Contents

- [Technical Fundamentals](#technical-fundamentals)
- [Application Security](#application-security)
- [Network Security](#network-security)
- [Cryptography](#cryptography)
- [Incident Response](#incident-response)
- [System Design](#system-design)
- [Scenario-Based Questions](#scenario-based-questions)
- [Behavioral Questions](#behavioral-questions)

## 🎯 Technical Fundamentals

### Basic Questions

**Q1: Explain the CIA Triad and provide examples of how each principle can be compromised.**

**Expected Answer:**
- **Confidentiality**: Ensuring data is only accessible to authorized parties
  - Compromise: Data breaches, eavesdropping, insider threats
  - Example: Customer database exposed due to misconfigured S3 bucket
- **Integrity**: Ensuring data accuracy and preventing unauthorized modification
  - Compromise: Data tampering, man-in-the-middle attacks, malware
  - Example: Financial transaction modified during transmission
- **Availability**: Ensuring systems and data are accessible when needed
  - Compromise: DDoS attacks, hardware failures, natural disasters
  - Example: Website taken offline by distributed denial of service attack

**Follow-up**: How would you prioritize these principles for a banking application?

---

**Q2: What is the difference between vulnerability, threat, and risk?**

**Expected Answer:**
- **Vulnerability**: A weakness in a system that can be exploited
  - Example: Unpatched software, weak passwords, misconfiguration
- **Threat**: A potential danger that could exploit a vulnerability
  - Example: Malicious hacker, malware, natural disaster
- **Risk**: The likelihood and impact of a threat exploiting a vulnerability
  - Formula: Risk = Threat × Vulnerability × Asset Value

**Follow-up**: How would you calculate the risk score for a SQL injection vulnerability?

---

**Q3: Describe the principle of least privilege and how you would implement it.**

**Expected Answer:**
- **Definition**: Users should only have the minimum access necessary to perform their job functions
- **Implementation strategies**:
  - Role-Based Access Control (RBAC)
  - Just-in-time access provisioning
  - Regular access reviews and deprovisioning
  - Separation of duties
  - Privileged access management (PAM)
- **Benefits**: Reduces attack surface, limits blast radius of breaches

## 🌐 Application Security

### Intermediate Questions

**Q4: Walk me through how you would secure a REST API.**

**Expected Answer:**
```
1. Authentication & Authorization
   - Implement OAuth 2.0 or JWT tokens
   - Use API keys for service-to-service communication
   - Implement proper session management

2. Input Validation
   - Validate all inputs (type, length, format)
   - Use parameterized queries to prevent SQL injection
   - Implement rate limiting

3. Data Protection
   - Use HTTPS/TLS for all communications
   - Encrypt sensitive data at rest
   - Implement proper error handling (don't leak information)

4. Security Headers
   - Content-Security-Policy
   - X-Frame-Options
   - X-Content-Type-Options
   - Strict-Transport-Security

5. Monitoring & Logging
   - Log all API calls and authentication attempts
   - Implement anomaly detection
   - Set up alerting for suspicious activities
```

**Follow-up**: How would you handle API versioning from a security perspective?

---

**Q5: Explain OWASP Top 10 and how you would mitigate the top 3 risks.**

**Expected Answer:**
**Top 3 OWASP Risks (2021):**

1. **Broken Access Control**
   - Mitigation: Implement proper authorization checks, use RBAC, deny by default
   
2. **Cryptographic Failures**
   - Mitigation: Use strong encryption (AES-256), proper key management, TLS 1.2+
   
3. **Injection**
   - Mitigation: Parameterized queries, input validation, use ORM frameworks

**Implementation Example:**
```python
# SQL Injection Prevention
def get_user_by_id(user_id):
    # BAD: Vulnerable to SQL injection
    # query = f"SELECT * FROM users WHERE id = {user_id}"
    
    # GOOD: Parameterized query
    query = "SELECT * FROM users WHERE id = %s"
    cursor.execute(query, (user_id,))
    return cursor.fetchone()
```

---

**Q6: How would you implement secure session management?**

**Expected Answer:**
```
1. Session Creation
   - Generate cryptographically secure session IDs
   - Use sufficient entropy (128+ bits)
   - Regenerate session ID after authentication

2. Session Storage
   - Store sessions server-side (not in cookies)
   - Use secure session storage (Redis, database)
   - Encrypt session data

3. Session Transmission
   - Use Secure and HttpOnly cookie flags
   - Implement SameSite attribute
   - Use HTTPS only

4. Session Lifecycle
   - Implement session timeout (idle and absolute)
   - Provide secure logout functionality
   - Clean up expired sessions

5. Session Validation
   - Validate session on each request
   - Check for session fixation attacks
   - Implement concurrent session limits
```

## 🔐 Cryptography

### Advanced Questions

**Q7: Explain the difference between symmetric and asymmetric encryption. When would you use each?**

**Expected Answer:**

**Symmetric Encryption:**
- Same key for encryption and decryption
- Fast and efficient for large data
- Key distribution challenge
- Examples: AES, ChaCha20
- Use cases: Data at rest, bulk data encryption

**Asymmetric Encryption:**
- Different keys (public/private key pair)
- Slower than symmetric encryption
- Solves key distribution problem
- Examples: RSA, ECC
- Use cases: Key exchange, digital signatures, small data

**Hybrid Approach:**
```
1. Generate symmetric key (AES)
2. Encrypt data with symmetric key
3. Encrypt symmetric key with recipient's public key
4. Send both encrypted data and encrypted key
```

**Follow-up**: How would you implement perfect forward secrecy?

---

**Q8: Describe how digital signatures work and their security properties.**

**Expected Answer:**
```
Digital Signature Process:
1. Create hash of message (SHA-256)
2. Encrypt hash with sender's private key
3. Attach signature to message
4. Recipient decrypts signature with sender's public key
5. Compare decrypted hash with computed hash of message

Security Properties:
- Authentication: Verifies sender identity
- Integrity: Detects message tampering
- Non-repudiation: Sender cannot deny signing
- Unforgeable: Only private key holder can create valid signature
```

**Implementation Considerations:**
- Use strong hash functions (SHA-256+)
- Proper key management and storage
- Certificate validation and revocation
- Timestamp signatures to prevent replay attacks

## 🚨 Incident Response

### Scenario-Based Questions

**Q9: You discover that your company's customer database has been breached. Walk me through your incident response process.**

**Expected Answer:**
```
1. IMMEDIATE RESPONSE (0-1 hour)
   - Confirm and contain the breach
   - Activate incident response team
   - Preserve evidence and logs
   - Document timeline and actions

2. ASSESSMENT (1-4 hours)
   - Determine scope and impact
   - Identify compromised systems and data
   - Assess ongoing threats
   - Notify key stakeholders

3. CONTAINMENT (4-24 hours)
   - Isolate affected systems
   - Change compromised credentials
   - Apply emergency patches
   - Monitor for lateral movement

4. ERADICATION (1-3 days)
   - Remove malware and backdoors
   - Close attack vectors
   - Strengthen security controls
   - Validate system integrity

5. RECOVERY (3-7 days)
   - Restore systems from clean backups
   - Implement additional monitoring
   - Gradual service restoration
   - Continuous monitoring

6. LESSONS LEARNED (1-2 weeks)
   - Post-incident review
   - Update procedures and controls
   - Staff training and awareness
   - Regulatory reporting if required
```

**Follow-up**: How would you handle media and customer communications?

---

**Q10: A developer reports that they accidentally committed AWS credentials to a public GitHub repository. What do you do?**

**Expected Answer:**
```
IMMEDIATE ACTIONS (0-15 minutes):
1. Rotate the exposed credentials immediately
2. Check AWS CloudTrail for unauthorized usage
3. Review all resources created with those credentials
4. Remove credentials from GitHub (though history remains)

ASSESSMENT (15-60 minutes):
1. Determine credential permissions and scope
2. Check for any unauthorized resource creation
3. Review billing for unexpected charges
4. Scan for data exfiltration attempts

REMEDIATION:
1. Delete any unauthorized resources
2. Implement GitHub secret scanning
3. Set up pre-commit hooks to prevent future exposure
4. Review and update IAM policies (least privilege)
5. Implement credential rotation policies

PREVENTION:
1. Use AWS IAM roles instead of long-term credentials
2. Implement secrets management (AWS Secrets Manager)
3. Regular security training for developers
4. Automated secret scanning in CI/CD pipeline
```

## 🏗️ System Design

### Architecture Questions

**Q11: Design a secure authentication system for a microservices architecture.**

**Expected Answer:**
```
ARCHITECTURE COMPONENTS:

1. Authentication Service
   - Centralized authentication (OAuth 2.0/OpenID Connect)
   - JWT token generation and validation
   - Multi-factor authentication support
   - Rate limiting and brute force protection

2. API Gateway
   - Token validation and routing
   - Rate limiting per user/service
   - Request/response logging
   - SSL termination

3. Service-to-Service Authentication
   - mTLS for internal communication
   - Service mesh (Istio) for policy enforcement
   - Short-lived tokens with automatic rotation

4. Token Management
   - Short-lived access tokens (15 minutes)
   - Refresh tokens with rotation
   - Token revocation capability
   - Secure token storage

SECURITY CONSIDERATIONS:
- Zero-trust network model
- Principle of least privilege
- Comprehensive audit logging
- Secrets management (HashiCorp Vault)
- Regular security assessments
```

**Follow-up**: How would you handle token revocation across all services?

---

**Q12: How would you secure a CI/CD pipeline?**

**Expected Answer:**
```
PIPELINE SECURITY LAYERS:

1. Source Code Security
   - Branch protection rules
   - Code review requirements
   - Secret scanning (GitLeaks, TruffleHog)
   - Static Application Security Testing (SAST)

2. Build Security
   - Secure build environments (containers)
   - Dependency scanning (Snyk, OWASP Dependency Check)
   - Container image scanning
   - Build artifact signing

3. Deployment Security
   - Infrastructure as Code (IaC) scanning
   - Dynamic Application Security Testing (DAST)
   - Deployment approval workflows
   - Blue-green deployments for rollback capability

4. Runtime Security
   - Runtime Application Self-Protection (RASP)
   - Container runtime security (Falco)
   - Continuous monitoring and alerting
   - Automated incident response

ACCESS CONTROLS:
- Role-based access to pipeline stages
- Service accounts with minimal permissions
- Audit logging of all pipeline activities
- Secrets management integration
```

## 🎭 Behavioral Questions

**Q13: Describe a time when you had to balance security requirements with business needs. How did you handle it?**

**Expected Answer Structure:**
```
SITUATION:
- Describe the specific scenario
- Explain the conflicting requirements

TASK:
- What was your role and responsibility?
- What needed to be accomplished?

ACTION:
- How did you analyze the risk vs. business impact?
- What alternatives did you propose?
- How did you communicate with stakeholders?

RESULT:
- What was the outcome?
- How did you measure success?
- What did you learn?
```

**Example Response:**
"The business wanted to launch a new feature quickly, but it required integrating with a third-party API that didn't meet our security standards. I conducted a risk assessment, proposed a phased approach with temporary compensating controls, and worked with the vendor to improve their security posture for the long-term solution."

---

**Q14: How do you stay current with cybersecurity threats and trends?**

**Expected Answer:**
```
INFORMATION SOURCES:
- Security blogs and publications (Krebs on Security, Dark Reading)
- Threat intelligence feeds (MITRE ATT&CK, NIST)
- Security conferences (Black Hat, DEF CON, RSA)
- Professional communities (OWASP, SANS)
- Vendor security advisories

PRACTICAL APPLICATION:
- Home lab for testing new tools and techniques
- Capture The Flag (CTF) competitions
- Bug bounty programs participation
- Open source security tool contributions
- Security certifications and training

KNOWLEDGE SHARING:
- Internal security awareness presentations
- Mentoring junior team members
- Writing technical blog posts
- Speaking at local security meetups
```

## 📊 Evaluation Criteria

### Technical Competency
- **Depth of Knowledge**: Understanding of security principles and technologies
- **Practical Experience**: Ability to apply concepts to real-world scenarios
- **Problem-Solving**: Analytical thinking and systematic approach
- **Communication**: Ability to explain complex concepts clearly

### Security Mindset
- **Risk Assessment**: Ability to evaluate and prioritize security risks
- **Defense in Depth**: Understanding of layered security approaches
- **Threat Modeling**: Systematic approach to identifying threats
- **Continuous Learning**: Commitment to staying current with evolving threats

### Collaboration Skills
- **Stakeholder Management**: Working with business and technical teams
- **Incident Response**: Ability to work under pressure during security incidents
- **Training and Awareness**: Helping others understand security concepts
- **Documentation**: Clear and comprehensive security documentation

---

**💡 Interview Tips:**
- Prepare specific examples from your experience
- Practice explaining technical concepts to non-technical audiences
- Stay current with recent security incidents and trends
- Be ready to discuss both successes and failures
- Ask thoughtful questions about the company's security posture
