# 🔐 Zero Trust Architecture

Never trust, always verify - security model for modern distributed systems.

## 📋 Table of Contents

- [Core Principles](#core-principles)
- [Implementation Components](#implementation-components)
- [Network Architecture](#network-architecture)
- [Identity and Access](#identity-and-access)
- [Micro-segmentation](#micro-segmentation)
- [Interview Questions](#interview-questions)

## 🎯 Core Principles

```
┌─────────────────────────────────────────────────────────┐
│                 ZERO TRUST PRINCIPLES                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Never Trust, Always Verify                         │
│     → Authenticate and authorize every request         │
│                                                         │
│  2. Least Privilege Access                             │
│     → Minimum permissions necessary                    │
│                                                         │
│  3. Assume Breach                                      │
│     → Design for compromise, limit blast radius        │
│                                                         │
│  4. Verify Explicitly                                  │
│     → Use all available data points                    │
│                                                         │
│  5. Continuous Monitoring                              │
│     → Real-time security analytics                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Traditional vs Zero Trust

| Aspect | Traditional | Zero Trust |
|--------|-------------|------------|
| Trust Model | Trust inside perimeter | Trust nothing |
| Access | Network-based | Identity-based |
| Segmentation | Coarse (DMZ, internal) | Fine-grained |
| Verification | One-time | Continuous |
| Visibility | Limited | Full visibility |

## 🏗️ Implementation Components

### Policy Decision Point (PDP)

```python
from dataclasses import dataclass
from typing import List, Dict, Optional
from enum import Enum
import time

class TrustLevel(Enum):
    LOW = 1
    MEDIUM = 2
    HIGH = 3
    VERY_HIGH = 4

@dataclass
class AccessRequest:
    user_id: str
    device_id: str
    resource: str
    action: str
    context: Dict
    timestamp: float

@dataclass
class TrustScore:
    user_trust: float
    device_trust: float
    network_trust: float
    behavior_trust: float
    overall: float

class ZeroTrustEngine:
    def __init__(self):
        self.user_risk_engine = UserRiskEngine()
        self.device_trust_engine = DeviceTrustEngine()
        self.network_analyzer = NetworkAnalyzer()
        self.behavior_analyzer = BehaviorAnalyzer()
        self.policy_store = PolicyStore()
    
    def evaluate_access(self, request: AccessRequest) -> tuple[bool, TrustScore]:
        """Evaluate access request against zero trust policies"""
        
        # Calculate trust scores
        user_trust = self.user_risk_engine.calculate_trust(request.user_id)
        device_trust = self.device_trust_engine.calculate_trust(request.device_id)
        network_trust = self.network_analyzer.calculate_trust(request.context)
        behavior_trust = self.behavior_analyzer.calculate_trust(request)
        
        # Weighted overall trust score
        overall = (
            user_trust * 0.35 +
            device_trust * 0.25 +
            network_trust * 0.20 +
            behavior_trust * 0.20
        )
        
        trust_score = TrustScore(
            user_trust=user_trust,
            device_trust=device_trust,
            network_trust=network_trust,
            behavior_trust=behavior_trust,
            overall=overall
        )
        
        # Get required trust level for resource
        required_level = self.policy_store.get_required_trust(
            request.resource, request.action
        )
        
        # Determine access
        granted = overall >= required_level.value / 4.0
        
        # Apply additional conditions
        if granted:
            granted = self._apply_conditions(request, trust_score)
        
        return granted, trust_score
    
    def _apply_conditions(self, request: AccessRequest, 
                         trust_score: TrustScore) -> bool:
        """Apply additional access conditions"""
        # MFA required for sensitive resources
        if self._is_sensitive(request.resource):
            if not request.context.get('mfa_verified'):
                return False
        
        # Time-based access
        if not self._is_business_hours(request.context.get('time')):
            if trust_score.overall < 0.8:
                return False
        
        # Location-based access
        if self._is_restricted_location(request.context.get('location')):
            return False
        
        return True

class UserRiskEngine:
    def calculate_trust(self, user_id: str) -> float:
        """Calculate user trust score"""
        score = 1.0  # Start with full trust
        
        # Check recent failed logins
        failed_logins = self._get_recent_failed_logins(user_id)
        score -= failed_logins * 0.1
        
        # Check password age
        password_age = self._get_password_age(user_id)
        if password_age > 90:
            score -= 0.1
        
        # Check MFA enabled
        if not self._has_mfa(user_id):
            score -= 0.2
        
        # Check privileged access
        if self._is_privileged(user_id):
            score -= 0.1  # Higher scrutiny
        
        return max(0, min(1, score))

class DeviceTrustEngine:
    def calculate_trust(self, device_id: str) -> float:
        """Calculate device trust score"""
        score = 1.0
        
        device_info = self._get_device_info(device_id)
        
        # Check device compliance
        if not device_info.get('compliant'):
            score -= 0.3
        
        # Check OS patch level
        patch_age = device_info.get('patch_age_days', 0)
        if patch_age > 30:
            score -= 0.2
        elif patch_age > 60:
            score -= 0.4
        
        # Check encryption
        if not device_info.get('encrypted'):
            score -= 0.3
        
        # Check antivirus
        if not device_info.get('antivirus_active'):
            score -= 0.2
        
        # Known device
        if not device_info.get('registered'):
            score -= 0.2
        
        return max(0, min(1, score))

class BehaviorAnalyzer:
    def calculate_trust(self, request: AccessRequest) -> float:
        """Calculate behavioral trust score"""
        score = 1.0
        
        # Check access time anomaly
        if self._is_unusual_time(request.user_id, request.timestamp):
            score -= 0.2
        
        # Check access location anomaly
        if self._is_unusual_location(request.user_id, request.context.get('location')):
            score -= 0.3
        
        # Check resource access pattern
        if self._is_unusual_resource(request.user_id, request.resource):
            score -= 0.2
        
        # Check data volume anomaly
        if self._is_unusual_data_volume(request.user_id):
            score -= 0.2
        
        return max(0, min(1, score))
```

## 🌐 Network Architecture

### Micro-perimeter Design

```
                    ┌─────────────────────────────────────┐
                    │         Identity Provider           │
                    │         (SSO, MFA, IAM)            │
                    └──────────────┬──────────────────────┘
                                   │
                                   ▼
                    ┌─────────────────────────────────────┐
                    │      Policy Decision Point          │
                    │      (Access Control Engine)        │
                    └──────────────┬──────────────────────┘
                                   │
            ┌──────────────────────┼──────────────────────┐
            │                      │                      │
            ▼                      ▼                      ▼
    ┌───────────────┐    ┌───────────────┐    ┌───────────────┐
    │  Policy       │    │  Policy       │    │  Policy       │
    │  Enforcement  │    │  Enforcement  │    │  Enforcement  │
    │  Point        │    │  Point        │    │  Point        │
    └───────┬───────┘    └───────┬───────┘    └───────┬───────┘
            │                    │                    │
            ▼                    ▼                    ▼
    ┌───────────────┐    ┌───────────────┐    ┌───────────────┐
    │   Service A   │    │   Service B   │    │   Service C   │
    │   (Micro-     │    │   (Micro-     │    │   (Micro-     │
    │   segment)    │    │   segment)    │    │   segment)    │
    └───────────────┘    └───────────────┘    └───────────────┘
```

### Service Mesh with mTLS

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: production
spec:
  mtls:
    mode: STRICT

---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: service-a-policy
  namespace: production
spec:
  selector:
    matchLabels:
      app: service-a
  rules:
    - from:
        - source:
            principals: ["cluster.local/ns/production/sa/frontend"]
      to:
        - operation:
            methods: ["GET", "POST"]
            paths: ["/api/*"]
      when:
        - key: request.auth.claims[role]
          values: ["user", "admin"]
```

## 🔑 Identity and Access

### Continuous Authentication

```go
package main

import (
	"context"
	"time"
)

type SessionManager struct {
	sessions      map[string]*Session
	riskEngine    *RiskEngine
	policyEngine  *PolicyEngine
}

type Session struct {
	ID           string
	UserID       string
	DeviceID     string
	CreatedAt    time.Time
	LastActivity time.Time
	TrustScore   float64
	ReauthNeeded bool
}

func (sm *SessionManager) ValidateSession(ctx context.Context, sessionID, action string) (*AuthDecision, error) {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return &AuthDecision{Allowed: false, Reason: "invalid_session"}, nil
	}
	
	// Check session timeout
	if time.Since(session.LastActivity) > 30*time.Minute {
		return &AuthDecision{Allowed: false, Reason: "session_expired"}, nil
	}
	
	// Re-evaluate trust score
	currentTrust := sm.riskEngine.EvaluateTrust(session.UserID, session.DeviceID)
	
	// Check if trust degraded
	if currentTrust < session.TrustScore-0.2 {
		session.ReauthNeeded = true
	}
	
	// Step-up authentication for sensitive actions
	if sm.policyEngine.RequiresStepUp(action) {
		if !sm.hasRecentMFA(session.ID) {
			return &AuthDecision{
				Allowed:     false,
				Reason:      "step_up_required",
				RequireMFA:  true,
			}, nil
		}
	}
	
	// Update activity
	session.LastActivity = time.Now()
	session.TrustScore = currentTrust
	
	return &AuthDecision{
		Allowed:    true,
		TrustScore: currentTrust,
	}, nil
}

type AuthDecision struct {
	Allowed     bool
	Reason      string
	TrustScore  float64
	RequireMFA  bool
}
```

## 🔒 Micro-segmentation

### Network Segmentation Policy

```python
from dataclasses import dataclass
from typing import List, Set

@dataclass
class Segment:
    id: str
    name: str
    services: Set[str]
    allowed_flows: List['FlowRule']

@dataclass
class FlowRule:
    source: str
    destination: str
    port: int
    protocol: str
    action: str  # 'allow' or 'deny'

class MicroSegmentationEngine:
    def __init__(self):
        self.segments: Dict[str, Segment] = {}
        self.service_to_segment: Dict[str, str] = {}
    
    def create_segment(self, name: str, services: List[str]) -> Segment:
        segment = Segment(
            id=f"seg_{name}",
            name=name,
            services=set(services),
            allowed_flows=[]
        )
        
        for service in services:
            self.service_to_segment[service] = segment.id
        
        self.segments[segment.id] = segment
        return segment
    
    def allow_flow(self, source_segment: str, dest_segment: str, 
                   port: int, protocol: str = "tcp"):
        """Allow traffic between segments"""
        source = self.segments[source_segment]
        source.allowed_flows.append(FlowRule(
            source=source_segment,
            destination=dest_segment,
            port=port,
            protocol=protocol,
            action="allow"
        ))
    
    def evaluate_flow(self, source_service: str, dest_service: str, 
                      port: int) -> bool:
        """Evaluate if flow is allowed"""
        source_seg = self.service_to_segment.get(source_service)
        dest_seg = self.service_to_segment.get(dest_service)
        
        if not source_seg or not dest_seg:
            return False
        
        segment = self.segments[source_seg]
        
        for flow in segment.allowed_flows:
            if (flow.destination == dest_seg and 
                flow.port == port and 
                flow.action == "allow"):
                return True
        
        return False  # Default deny
    
    def generate_firewall_rules(self) -> List[str]:
        """Generate firewall rules from segmentation policy"""
        rules = []
        
        for segment_id, segment in self.segments.items():
            for flow in segment.allowed_flows:
                for src_service in self.segments[flow.source].services:
                    for dst_service in self.segments[flow.destination].services:
                        rules.append(
                            f"ALLOW {src_service} -> {dst_service}:{flow.port}/{flow.protocol}"
                        )
        
        rules.append("DENY ALL")  # Default deny
        return rules

# Usage
engine = MicroSegmentationEngine()

# Create segments
web = engine.create_segment("web", ["frontend", "nginx"])
app = engine.create_segment("app", ["api-server", "worker"])
data = engine.create_segment("data", ["postgres", "redis"])

# Define allowed flows
engine.allow_flow("web", "app", 8080)
engine.allow_flow("app", "data", 5432)  # API to postgres
engine.allow_flow("app", "data", 6379)  # API to redis

# Evaluate
print(engine.evaluate_flow("frontend", "api-server", 8080))  # True
print(engine.evaluate_flow("frontend", "postgres", 5432))    # False
```

## ❓ Interview Questions

### Basic
**Q: What is Zero Trust?**
A: Security model that assumes breach and requires verification for every access request, regardless of location. No implicit trust based on network location.

**Q: How does Zero Trust differ from traditional security?**
A: Traditional trusts inside perimeter. Zero Trust trusts nothing, verifies everything continuously.

### Intermediate
**Q: What are the key components of Zero Trust?**
A: Identity provider, policy decision point, policy enforcement points, continuous monitoring, micro-segmentation, strong authentication.

**Q: How do you implement least privilege in Zero Trust?**
A: Just-in-time access, role-based permissions, attribute-based access control, continuous re-evaluation, time-limited access grants.

### Advanced
**Q: How would you implement Zero Trust for a hybrid cloud environment?**
A: (1) Unified identity across cloud and on-prem, (2) Software-defined perimeter, (3) Service mesh with mTLS, (4) Cloud-native PDP/PEP, (5) Unified policy engine, (6) Continuous monitoring across all environments.