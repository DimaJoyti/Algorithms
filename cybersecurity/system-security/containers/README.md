# 🐳 Container Security

Securing containerized applications and container orchestration platforms.

## 📋 Table of Contents

- [Docker Security](#docker-security)
- [Container Hardening](#container-hardening)
- [Kubernetes Security](#kubernetes-security)
- [Image Security](#image-security)
- [Interview Questions](#interview-questions)

## 🐋 Docker Security

### Secure Dockerfile

```dockerfile
# BAD: Insecure Dockerfile
# FROM node:latest
# USER root
# COPY . /app
# RUN npm install
# CMD ["node", "server.js"]

# GOOD: Secure Dockerfile
FROM node:20-alpine AS builder

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy dependency files first (better caching)
COPY package*.json ./
RUN npm ci --only=production && \
    npm cache clean --force

# Copy application code
COPY --chown=appuser:appgroup . .

# Switch to non-root user
USER appuser

# Security labels
LABEL maintainer="security@example.com"
LABEL version="1.0"
LABEL description="Secure Node.js application"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Read-only filesystem
VOLUME ["/tmp"]

EXPOSE 3000

CMD ["node", "server.js"]
```

### Docker Security Configuration

```yaml
# docker-compose.yml with security hardening
version: '3.8'

services:
  app:
    build: .
    security_opt:
      - no-new-privileges:true
      - seccomp:seccomp-profile.json
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=64m
    environment:
      - NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
    secrets:
      - db_password
    networks:
      - frontend
      - backend
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 3s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
    read_only: true
    tmpfs:
      - /var/cache/nginx
      - /var/run

secrets:
  db_password:
    file: ./secrets/db_password.txt

networks:
  frontend:
  backend:
    internal: true
```

### Container Security Scanner

```python
import subprocess
import json
from dataclasses import dataclass
from typing import List

@dataclass
class Vulnerability:
    id: str
    severity: str
    package: str
    installed_version: str
    fixed_version: str
    description: str

class ContainerSecurityScanner:
    def __init__(self, image: str):
        self.image = image
    
    def scan_image(self) -> List[Vulnerability]:
        """Scan container image for vulnerabilities"""
        result = subprocess.run(
            ['trivy', 'image', '--format', 'json', self.image],
            capture_output=True, text=True
        )
        
        vulnerabilities = []
        
        if result.returncode == 0:
            data = json.loads(result.stdout)
            
            for result in data.get('Results', []):
                for vuln in result.get('Vulnerabilities', []):
                    vulnerabilities.append(Vulnerability(
                        id=vuln.get('VulnerabilityID', ''),
                        severity=vuln.get('Severity', 'UNKNOWN'),
                        package=vuln.get('PkgName', ''),
                        installed_version=vuln.get('InstalledVersion', ''),
                        fixed_version=vuln.get('FixedVersion', ''),
                        description=vuln.get('Description', '')
                    ))
        
        return vulnerabilities
    
    def check_security_issues(self) -> dict:
        """Check container security configuration"""
        issues = {
            'running_as_root': False,
            'privileged': False,
            'sensitive_mounts': [],
            'exposed_ports': [],
            'secrets_in_env': []
        }
        
        # Inspect image
        result = subprocess.run(
            ['docker', 'image', 'inspect', self.image],
            capture_output=True, text=True
        )
        
        if result.returncode == 0:
            data = json.loads(result.stdout)[0]
            config = data.get('Config', {})
            
            # Check user
            if not config.get('User') or config.get('User') == 'root':
                issues['running_as_root'] = True
            
            # Check for secrets in environment
            env_vars = config.get('Env', [])
            sensitive_keywords = ['PASSWORD', 'SECRET', 'KEY', 'TOKEN', 'API_KEY']
            
            for env in env_vars:
                for keyword in sensitive_keywords:
                    if keyword in env.upper():
                        issues['secrets_in_env'].append(env.split('=')[0])
        
        return issues
    
    def generate_report(self) -> str:
        """Generate security report"""
        vulns = self.scan_image()
        issues = self.check_security_issues()
        
        report = []
        report.append(f"CONTAINER SECURITY REPORT: {self.image}")
        report.append("=" * 60)
        
        # Vulnerability summary
        severity_counts = {}
        for v in vulns:
            severity_counts[v.severity] = severity_counts.get(v.severity, 0) + 1
        
        report.append("\nVULNERABILITY SUMMARY:")
        for severity in ['CRITICAL', 'HIGH', 'MEDIUM', 'LOW']:
            count = severity_counts.get(severity, 0)
            if count > 0:
                report.append(f"  {severity}: {count}")
        
        # Configuration issues
        report.append("\nCONFIGURATION ISSUES:")
        if issues['running_as_root']:
            report.append("  ✗ Container runs as root")
        if issues['privileged']:
            report.append("  ✗ Container runs in privileged mode")
        if issues['secrets_in_env']:
            report.append(f"  ✗ Potential secrets in environment: {issues['secrets_in_env']}")
        
        return "\n".join(report)
```

## 🔒 Container Hardening

### Security Profiles

```json
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": ["SCMP_ARCH_X86_64"],
  "syscalls": [
    {
      "names": [
        "read", "write", "close", "fstat", "mmap", "mprotect",
        "munmap", "brk", "ioctl", "access", "pipe", "dup2",
        "getpid", "socket", "connect", "sendto", "recvfrom",
        "exit_group", "arch_prctl", "gettid", "futex", "set_tid_address"
      ],
      "action": "SCMP_ACT_ALLOW"
    }
  ]
}
```

### AppArmor Profile

```bash
#include <tunables/global>

profile docker-apparmor flags=(attach_disconnected,mediate_deleted) {
  #include <abstractions/base>

  network inet tcp,
  network inet udp,
  network inet6 tcp,
  network inet6 udp,

  # Deny access to sensitive paths
  deny /etc/shadow r,
  deny /etc/gshadow r,
  deny /etc/ssh/* r,
  deny /root/** rwx,
  deny /home/** rwx,

  # Allow specific capabilities
  capability net_bind_service,

  # Deny dangerous capabilities
  deny capability sys_admin,
  deny capability sys_ptrace,
  deny capability sys_module,
  deny capability net_raw,
  deny capability mknod,

  # File access
  /app/** r,
  /tmp/** rw,
  /var/log/app/** rw,

  # Process execution
  /usr/bin/node rix,
}
```

## ☸️ Kubernetes Security

### Pod Security Standards

```yaml
# Restricted Pod Security Policy
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: restricted
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    - 'persistentVolumeClaim'
  hostNetwork: false
  hostIPC: false
  hostPID: false
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  supplementalGroups:
    rule: 'MustRunAs'
    ranges:
      - min: 1
        max: 65535
  fsGroup:
    rule: 'MustRunAs'
    ranges:
      - min: 1
        max: 65535
  readOnlyRootFilesystem: true
```

### Secure Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: secure-app
  labels:
    app: secure-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: secure-app
  template:
    metadata:
      labels:
        app: secure-app
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000
        seccompProfile:
          type: RuntimeDefault
      containers:
        - name: app
          image: myapp:1.0
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
          resources:
            limits:
              cpu: "500m"
              memory: "512Mi"
            requests:
              cpu: "250m"
              memory: "256Mi"
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 3000
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: tmp
              mountPath: /tmp
            - name: config
              mountPath: /app/config
              readOnly: true
      volumes:
        - name: tmp
          emptyDir: {}
        - name: config
          configMap:
            name: app-config
      imagePullSecrets:
        - name: registry-credentials
```

### Network Policy

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: app-network-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: secure-app
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              role: frontend
      ports:
        - protocol: TCP
          port: 3000
  egress:
    - to:
        - podSelector:
            matchLabels:
              role: database
      ports:
        - protocol: TCP
          port: 5432
    - to:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              k8s-app: kube-dns
      ports:
        - protocol: UDP
          port: 53
```

## 🖼️ Image Security

### Image Scanning Pipeline

```yaml
# .gitlab-ci.yml
stages:
  - build
  - scan
  - deploy

build-image:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

scan-image:
  stage: scan
  script:
    - trivy image --exit-code 1 --severity HIGH,CRITICAL $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - grype $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA --fail-on high
  allow_failure: false

scan-dockerfile:
  stage: scan
  script:
    - hadolint Dockerfile
  allow_failure: true

deploy:
  stage: deploy
  script:
    - kubectl set image deployment/app app=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main
```

### Base Image Management

```python
import subprocess
from datetime import datetime
from typing import List, Dict

class ImageManager:
    def __init__(self, registry: str):
        self.registry = registry
    
    def list_images(self) -> List[Dict]:
        """List all images in registry"""
        result = subprocess.run(
            ['skopeo', 'list', 'tags', f'docker://{self.registry}'],
            capture_output=True, text=True
        )
        return json.loads(result.stdout).get('Tags', [])
    
    def get_image_age(self, image: str) -> int:
        """Get age of image in days"""
        result = subprocess.run(
            ['skopeo', 'inspect', f'docker://{image}'],
            capture_output=True, text=True
        )
        
        if result.returncode == 0:
            data = json.loads(result.stdout)
            created = data.get('Created', '')
            if created:
                created_date = datetime.fromisoformat(created.replace('Z', '+00:00'))
                return (datetime.now(created_date.tzinfo) - created_date).days
        return -1
    
    def cleanup_old_images(self, keep_last: int = 10):
        """Remove images older than threshold"""
        images = self.list_images()
        
        for image in images[keep_last:]:
            print(f"Would remove: {image}")
            # subprocess.run(['skopeo', 'delete', f'docker://{image}'])
```

## ❓ Interview Questions

### Basic
**Q: Why run containers as non-root?**
A: Reduces privilege escalation risk, limits container break-out impact, follows least privilege.

**Q: What's the purpose of image scanning?**
A: Detect known vulnerabilities in base image and dependencies before deployment.

### Intermediate
**Q: Explain the difference between capabilities and privileged mode.**
A: Privileged gives all capabilities + host access. Capabilities allow fine-grained permissions (e.g., NET_BIND_SERVICE only).

**Q: How do you secure secrets in containers?**
A: Use secrets management (K8s secrets, Docker secrets, Vault), mount as files, never in environment variables or image layers.

### Advanced
**Q: Design a secure container deployment pipeline.**
A: (1) Base image scanning, (2) Dockerfile linting, (3) Build with non-root user, (4) Image scanning after build, (5) Sign images, (6) Deploy with security context, (7) Network policies, (8) Runtime monitoring with Falco.