# 🔍 Intrusion Detection and Prevention Systems (IDS/IPS)

IDS monitors for malicious activity, IPS can actively block detected threats.

## 📋 Table of Contents

- [Overview](#overview)
- [Detection Methods](#detection-methods)
- [SIEM Integration](#siem-integration)
- [Implementation Examples](#implementation-examples)
- [Interview Questions](#interview-questions)

## 🎯 Overview

### IDS vs IPS

| Feature | IDS | IPS |
|---------|-----|-----|
| Position | Out-of-band | In-line |
| Action | Alert only | Alert + Block |
| Impact | No latency | Can add latency |
| Failure mode | No network impact | Can block traffic |

### Architecture

```
                    ┌─────────┐
    Internet ──────►│  IPS    │──────► Firewall ──────► Internal
                    │ In-line │
                    └─────────┘
                         │
                         ▼
                    ┌─────────┐
                    │   IDS   │ (Out-of-band via SPAN/mirror port)
                    │ Monitor │
                    └─────────┘
                         │
                         ▼
                    ┌─────────┐
                    │  SIEM   │
                    └─────────┘
```

## 🔬 Detection Methods

### Signature-Based Detection

```python
from dataclasses import dataclass
from typing import List, Dict
import re

@dataclass
class Signature:
    id: str
    name: str
    pattern: str
    severity: int  # 1-5
    category: str
    reference: str = ""

class SignatureEngine:
    def __init__(self):
        self.signatures: List[Signature] = []
        self._load_signatures()
    
    def _load_signatures(self):
        """Load detection signatures"""
        self.signatures = [
            Signature(
                id="ET-0001",
                name="SQL Injection Attempt",
                pattern=r"(?i)(\%27|\'|\-\-|union.*select|select.*from)",
                severity=5,
                category="web-attack",
                reference="CVE-2024-XXXX"
            ),
            Signature(
                id="ET-0002", 
                name="Nmap Scan Detection",
                pattern=r"(?i)Nmap",
                severity=3,
                category="reconnaissance"
            ),
            Signature(
                id="ET-0003",
                name="Shellshock Attack",
                pattern=r"\(\)\s*\{[^}]*\}\s*[^a-zA-Z0-9]",
                severity=5,
                category="exploit"
            ),
            Signature(
                id="ET-0004",
                name="Suspicious User Agent",
                pattern=r"(?i)(nikto|sqlmap|dirbuster|gobuster|metasploit)",
                severity=4,
                category="scanner"
            )
        ]
    
    def scan(self, data: str) -> List[Dict]:
        """Scan data against signatures"""
        matches = []
        
        for sig in self.signatures:
            if re.search(sig.pattern, data):
                matches.append({
                    "signature_id": sig.id,
                    "name": sig.name,
                    "severity": sig.severity,
                    "category": sig.category,
                    "reference": sig.reference
                })
        
        return matches
```

### Anomaly-Based Detection

```python
import numpy as np
from collections import defaultdict
from datetime import datetime, timedelta

class AnomalyDetector:
    def __init__(self, baseline_window: int = 7):
        self.baseline_window = baseline_window  # days
        self.traffic_baseline = defaultdict(list)
        self.threshold_multiplier = 3.0  # Standard deviations
    
    def record_traffic(self, source_ip: str, bytes_sent: int, 
                       requests: int, timestamp: datetime):
        """Record traffic metrics for baseline calculation"""
        key = f"{source_ip}"
        self.traffic_baseline[key].append({
            "timestamp": timestamp,
            "bytes": bytes_sent,
            "requests": requests
        })
        
        # Keep only recent data
        cutoff = timestamp - timedelta(days=self.baseline_window)
        self.traffic_baseline[key] = [
            t for t in self.traffic_baseline[key] 
            if t["timestamp"] > cutoff
        ]
    
    def calculate_baseline(self, source_ip: str) -> dict:
        """Calculate statistical baseline for source"""
        data = self.traffic_baseline.get(source_ip, [])
        
        if len(data) < 10:  # Need minimum data points
            return None
        
        bytes_list = [d["bytes"] for d in data]
        requests_list = [d["requests"] for d in data]
        
        return {
            "bytes_mean": np.mean(bytes_list),
            "bytes_std": np.std(bytes_list),
            "requests_mean": np.mean(requests_list),
            "requests_std": np.std(requests_list)
        }
    
    def detect_anomaly(self, source_ip: str, bytes_sent: int, 
                       requests: int) -> List[dict]:
        """Detect anomalies against baseline"""
        baseline = self.calculate_baseline(source_ip)
        
        if not baseline:
            return []  # Insufficient baseline data
        
        anomalies = []
        
        # Check for volume anomaly
        if baseline["bytes_std"] > 0:
            z_score = (bytes_sent - baseline["bytes_mean"]) / baseline["bytes_std"]
            if z_score > self.threshold_multiplier:
                anomalies.append({
                    "type": "volume_anomaly",
                    "z_score": z_score,
                    "description": f"Traffic volume {z_score:.1f}x above normal",
                    "severity": min(5, int(z_score))
                })
        
        # Check for request rate anomaly
        if baseline["requests_std"] > 0:
            z_score = (requests - baseline["requests_mean"]) / baseline["requests_std"]
            if z_score > self.threshold_multiplier:
                anomalies.append({
                    "type": "rate_anomaly",
                    "z_score": z_score,
                    "description": f"Request rate {z_score:.1f}x above normal",
                    "severity": min(5, int(z_score))
                })
        
        return anomalies

class BehavioralAnalysis:
    """Advanced behavioral analysis for threat detection"""
    
    def __init__(self):
        self.user_profiles = defaultdict(lambda: {
            "login_times": [],
            "login_locations": set(),
            "devices": set(),
            "accessed_resources": set()
        })
    
    def record_user_activity(self, user: str, login_time: datetime,
                           location: str, device: str, resource: str):
        """Record user activity for profiling"""
        profile = self.user_profiles[user]
        profile["login_times"].append(login_time.hour)
        profile["login_locations"].add(location)
        profile["devices"].add(device)
        profile["accessed_resources"].add(resource)
    
    def detect_behavioral_anomaly(self, user: str, login_time: datetime,
                                   location: str, device: str) -> List[dict]:
        """Detect behavioral anomalies"""
        profile = self.user_profiles[user]
        anomalies = []
        
        # Check login time (if enough history)
        if len(profile["login_times"]) >= 5:
            times = profile["login_times"]
            mean_time = np.mean(times)
            std_time = np.std(times)
            
            if std_time > 0:
                z_score = abs(login_time.hour - mean_time) / std_time
                if z_score > 2.5:
                    anomalies.append({
                        "type": "unusual_login_time",
                        "severity": 2,
                        "description": f"Login at unusual hour"
                    })
        
        # Check for new location
        if location not in profile["login_locations"] and len(profile["login_locations"]) > 0:
            anomalies.append({
                "type": "new_location",
                "severity": 3,
                "description": f"Login from new location: {location}"
            })
        
        # Check for new device
        if device not in profile["devices"] and len(profile["devices"]) > 0:
            anomalies.append({
                "type": "new_device",
                "severity": 2,
                "description": f"Login from new device"
            })
        
        return anomalies
```

## 📊 SIEM Integration

### Log Collection and Analysis

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type SecurityEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
	EventType   string    `json:"event_type"`
	Severity    int       `json:"severity"`
	SourceIP    string    `json:"source_ip"`
	Destination string    `json:"destination"`
	Message     string    `json:"message"`
	Details     string    `json:"details"`
}

type SIEMCollector struct {
	eventChan  chan SecurityEvent
	alertRules []AlertRule
}

type AlertRule struct {
	Name      string
	Condition func(SecurityEvent) bool
	Severity  int
	Message   string
}

func NewSIEMCollector() *SIEMCollector {
	return &SIEMCollector{
		eventChan: make(chan SecurityEvent, 1000),
		alertRules: []AlertRule{
			{
				Name: "Brute Force Detection",
				Condition: func(e SecurityEvent) bool {
					return e.EventType == "auth_failure" && 
						strings.Contains(e.Message, "Invalid credentials")
				},
				Severity: 4,
				Message:  "Potential brute force attack detected",
			},
			{
				Name: "Privilege Escalation",
				Condition: func(e SecurityEvent) bool {
					return e.EventType == "privilege_change" && 
						e.Severity >= 4
				},
				Severity: 5,
				Message:  "Privilege escalation attempt detected",
			},
			{
				Name: "Malware Signature",
				Condition: func(e SecurityEvent) bool {
					return e.EventType == "malware_detected"
				},
				Severity: 5,
				Message:  "Malware detected on endpoint",
			},
		},
	}
}

func (s *SIEMCollector) CollectSyslog(message string, source string) {
	event := s.parseSyslog(message, source)
	if event != nil {
		s.eventChan <- *event
	}
}

func (s *SIEMCollector) parseSyslog(message, source string) *SecurityEvent {
	// Basic syslog parsing
	event := &SecurityEvent{
		Timestamp: time.Now(),
		Source:    source,
		Message:   message,
	}
	
	// Extract IP addresses
	ipRegex := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}`)
	if ips := ipRegex.FindAllString(message, -1); len(ips) > 0 {
		event.SourceIP = ips[0]
		if len(ips) > 1 {
			event.Destination = ips[1]
		}
	}
	
	// Detect event type
	lowerMsg := strings.ToLower(message)
	switch {
	case strings.Contains(lowerMsg, "failed password") || 
	     strings.Contains(lowerMsg, "authentication failure"):
		event.EventType = "auth_failure"
		event.Severity = 3
	case strings.Contains(lowerMsg, "accepted password") ||
	     strings.Contains(lowerMsg, "session opened"):
		event.EventType = "auth_success"
		event.Severity = 1
	case strings.Contains(lowerMsg, "sudo") || 
	     strings.Contains(lowerMsg, "su:"):
		event.EventType = "privilege_change"
		event.Severity = 2
	case strings.Contains(lowerMsg, "connection from"):
		event.EventType = "connection"
		event.Severity = 1
	default:
		event.EventType = "general"
		event.Severity = 1
	}
	
	return event
}

func (s *SIEMCollector) StartProcessor() {
	for event := range s.eventChan {
		// Process event
		s.processEvent(event)
	}
}

func (s *SIEMCollector) processEvent(event SecurityEvent) {
	// Check alert rules
	for _, rule := range s.alertRules {
		if rule.Condition(event) {
			alert := SecurityEvent{
				Timestamp:   time.Now(),
				Source:      "SIEM",
				EventType:   "alert",
				Severity:    rule.Severity,
				SourceIP:    event.SourceIP,
				Destination: event.Destination,
				Message:     rule.Message,
				Details:     fmt.Sprintf("Rule: %s, Original: %s", rule.Name, event.Message),
			}
			s.sendAlert(alert)
		}
	}
	
	// Log event
	s.logEvent(event)
}

func (s *SIEMCollector) sendAlert(alert SecurityEvent) {
	alertJSON, _ := json.Marshal(alert)
	log.Printf("ALERT: %s", alertJSON)
}

func (s *SIEMCollector) logEvent(event SecurityEvent) {
	eventJSON, _ := json.Marshal(event)
	fmt.Printf("EVENT: %s\n", eventJSON)
}

func (s *SIEMCollector) StartSyslogServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	buf := make([]byte, 65535)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		message := string(buf[:n])
		go s.CollectSyslog(message, addr.IP.String())
	}
}

func main() {
	siem := NewSIEMCollector()
	
	// Start processor
	go siem.StartProcessor()
	
	// Simulate events
	events := []string{
		"Feb 24 10:00:01 server sshd[1234]: Failed password for root from 192.168.1.100",
		"Feb 24 10:00:02 server sshd[1234]: Failed password for root from 192.168.1.100",
		"Feb 24 10:00:03 server sshd[1234]: Accepted password for user from 192.168.1.100",
		"Feb 24 10:00:05 server sudo: user : TTY=pts/0 ; USER=root ; COMMAND=/bin/bash",
	}
	
	for _, event := range events {
		siem.CollectSyslog(event, "localhost")
	}
	
	time.Sleep(1 * time.Second)
}
```

## ❓ Interview Questions

### Basic
**Q: What's the difference between IDS and IPS?**
A: IDS monitors and alerts but doesn't block. IPS is in-line and can actively block detected threats.

**Q: What are false positives and false negatives?**
A: False positive = legitimate activity flagged as malicious. False negative = malicious activity not detected.

### Intermediate
**Q: When would you choose signature-based vs anomaly-based detection?**
A: Signatures for known threats (accurate, low false positives). Anomaly for zero-days and insider threats (can detect unknown threats but higher false positives).

**Q: How do you handle false positives?**
A: Tune signatures, adjust thresholds, whitelist known good behavior, use multiple detection methods for correlation.

### Advanced
**Q: Design an IDS/IPS deployment for a cloud environment.**
A: (1) Cloud-native IDS/IPS at VPC level, (2) Host-based IDS on instances, (3) Centralized SIEM for correlation, (4) Auto-scaling detection rules, (5) Integration with cloud security services.
