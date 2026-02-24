# 🚨 Incident Response Procedures

Structured approach to handling security incidents effectively.

## 📋 Table of Contents

- [Incident Response Lifecycle](#incident-response-lifecycle)
- [Incident Classification](#incident-classification)
- [Response Playbooks](#response-playbooks)
- [Communication Templates](#communication-templates)
- [Post-Incident Activities](#post-incident-activities)
- [Interview Questions](#interview-questions)

## 🔄 Incident Response Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                  INCIDENT RESPONSE LIFECYCLE                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│    ┌──────────┐    ┌──────────┐    ┌──────────┐           │
│    │ 1. Prep  │───►│2. Detect │───►│3. Analyze│           │
│    │          │    │          │    │          │           │
│    └──────────┘    └──────────┘    └────┬─────┘           │
│                                         │                  │
│    ┌──────────┐    ┌──────────┐    ┌────▼─────┐           │
│    │6. Lessons│◄───│ 5. Recover│◄───│4. Contain│           │
│    │          │    │          │    │          │           │
│    └──────────┘    └──────────┘    └──────────┘           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Phase Details

| Phase | Activities | Duration |
|-------|------------|----------|
| Preparation | Team training, tools, playbooks | Ongoing |
| Detection | Monitoring, alerts, user reports | Minutes-hours |
| Analysis | Triage, scope, impact assessment | Hours |
| Containment | Isolate, stop spread | Hours-days |
| Eradication | Remove threat, patch vulnerabilities | Days |
| Recovery | Restore services, verify security | Days-weeks |
| Lessons Learned | Review, update procedures | 1-2 weeks |

## 📊 Incident Classification

### Severity Levels

```python
from enum import Enum
from dataclasses import dataclass
from datetime import datetime

class Severity(Enum):
    CRITICAL = 1  # Active breach, data exfiltration
    HIGH = 2      # Compromised system, malware spread
    MEDIUM = 3    # Attempted breach, suspicious activity
    LOW = 4       # Policy violation, minor security event

class IncidentStatus(Enum):
    NEW = "new"
    IN_PROGRESS = "in_progress"
    CONTAINED = "contained"
    RESOLVED = "resolved"
    CLOSED = "closed"

@dataclass
class Incident:
    id: str
    title: str
    severity: Severity
    status: IncidentStatus
    description: str
    affected_systems: list
    detected_at: datetime
    reporter: str
    assignee: str = None
    timeline: list = None
    
    def calculate_response_deadline(self) -> datetime:
        """Calculate SLA deadline based on severity"""
        hours = {
            Severity.CRITICAL: 1,
            Severity.HIGH: 4,
            Severity.MEDIUM: 24,
            Severity.LOW: 72
        }
        return self.detected_at + timedelta(hours=hours[self.severity])

class IncidentClassifier:
    @staticmethod
    def classify(event: dict) -> Severity:
        """Classify incident severity based on event data"""
        score = 0
        
        # Data sensitivity
        if event.get('data_sensitivity') == 'confidential':
            score += 3
        elif event.get('data_sensitivity') == 'internal':
            score += 2
        
        # System criticality
        if event.get('system_criticality') == 'critical':
            score += 3
        elif event.get('system_criticality') == 'important':
            score += 2
        
        # Active exploitation
        if event.get('active_exploitation'):
            score += 4
        
        # Lateral movement
        if event.get('lateral_movement'):
            score += 3
        
        # Data exfiltration
        if event.get('data_exfiltration'):
            score += 4
        
        # Map to severity
        if score >= 8:
            return Severity.CRITICAL
        elif score >= 5:
            return Severity.HIGH
        elif score >= 3:
            return Severity.MEDIUM
        else:
            return Severity.LOW
```

## 📖 Response Playbooks

### Data Breach Playbook

```python
class DataBreachPlaybook:
    """
    Playbook for handling data breach incidents
    """
    
    STEPS = [
        {
            'phase': 'Detection',
            'actions': [
                'Verify the breach is legitimate',
                'Document initial findings',
                'Activate incident response team',
                'Begin evidence preservation'
            ],
            'duration': '0-1 hour'
        },
        {
            'phase': 'Analysis',
            'actions': [
                'Identify type of data compromised',
                'Determine number of affected records',
                'Identify affected individuals',
                'Assess regulatory implications',
                'Document attack vector'
            ],
            'duration': '1-4 hours'
        },
        {
            'phase': 'Containment',
            'actions': [
                'Isolate affected systems',
                'Revoke compromised credentials',
                'Block malicious IPs/domains',
                'Preserve forensic evidence',
                'Implement emergency monitoring'
            ],
            'duration': '4-24 hours'
        },
        {
            'phase': 'Eradication',
            'actions': [
                'Remove malware/backdoors',
                'Patch vulnerability',
                'Reset all potentially compromised passwords',
                'Review and strengthen access controls'
            ],
            'duration': '24-72 hours'
        },
        {
            'phase': 'Recovery',
            'actions': [
                'Restore systems from clean backups',
                'Verify system integrity',
                'Implement additional monitoring',
                'Gradually restore services'
            ],
            'duration': '72+ hours'
        },
        {
            'phase': 'Notification',
            'actions': [
                'Prepare breach notification',
                'Notify affected individuals',
                'Notify regulatory bodies (GDPR: 72h)',
                'Prepare public statement if needed'
            ],
            'duration': 'As required by regulation'
        }
    ]

class IncidentResponseOrchestrator:
    def __init__(self):
        self.incidents = {}
        self.team = ResponseTeam()
        self.playbooks = {
            'data_breach': DataBreachPlaybook(),
            'malware': MalwarePlaybook(),
            'ddos': DDoSPlaybook(),
            'ransomware': RansomwarePlaybook()
        }
    
    def create_incident(self, event: dict) -> Incident:
        """Create new incident from detection event"""
        incident = Incident(
            id=self._generate_id(),
            title=event.get('title', 'Unnamed Incident'),
            severity=IncidentClassifier.classify(event),
            status=IncidentStatus.NEW,
            description=event.get('description', ''),
            affected_systems=event.get('affected_systems', []),
            detected_at=datetime.utcnow(),
            reporter=event.get('reporter', 'system'),
            timeline=[]
        )
        
        self.incidents[incident.id] = incident
        
        # Auto-assign based on severity
        if incident.severity in [Severity.CRITICAL, Severity.HIGH]:
            incident.assignee = self.team.get_on_call()
            self._send_alert(incident)
        
        return incident
    
    def update_status(self, incident_id: str, status: IncidentStatus, 
                      note: str = None):
        """Update incident status"""
        incident = self.incidents.get(incident_id)
        if not incident:
            return
        
        incident.status = status
        incident.timeline.append({
            'timestamp': datetime.utcnow(),
            'action': f'Status changed to {status.value}',
            'note': note,
            'actor': self.team.current_user
        })
    
    def get_playbook(self, incident_type: str) -> dict:
        """Get response playbook for incident type"""
        return self.playbooks.get(incident_type)

# DDoS Response
class DDoSPlaybook:
    STEPS = [
        {
            'phase': 'Detection',
            'actions': [
                'Monitor traffic anomalies',
                'Identify attack type (volumetric, protocol, application)',
                'Document attack characteristics',
                'Estimate attack volume'
            ]
        },
        {
            'phase': 'Mitigation',
            'actions': [
                'Activate DDoS protection service',
                'Enable rate limiting',
                'Block attack source IPs',
                'Enable SYN cookies',
                'Scale up infrastructure if needed',
                'Switch to static content delivery'
            ]
        },
        {
            'phase': 'Communication',
            'actions': [
                'Notify stakeholders',
                'Update status page',
                'Coordinate with ISP/upstream provider',
                'Engage DDoS mitigation vendor'
            ]
        },
        {
            'phase': 'Recovery',
            'actions': [
                'Monitor for attack resumption',
                'Gradually restore normal operations',
                'Review and improve defenses',
                'Document incident'
            ]
        }
    ]
```

## 📧 Communication Templates

### Incident Notification

```python
class IncidentCommunicator:
    
    @staticmethod
    def generate_internal_notification(incident: Incident) -> str:
        """Generate internal incident notification"""
        return f"""
SECURITY INCIDENT NOTIFICATION

Incident ID: {incident.id}
Severity: {incident.severity.name}
Status: {incident.status.value}
Detected: {incident.detected_at.strftime('%Y-%m-%d %H:%M UTC')}

SUMMARY:
{incident.title}

DESCRIPTION:
{incident.description}

AFFECTED SYSTEMS:
{chr(10).join('- ' + s for s in incident.affected_systems)}

RESPONSE DEADLINE: {incident.calculate_response_deadline().strftime('%Y-%m-%d %H:%M UTC')}

ASSIGNEE: {incident.assignee or 'TBD'}

This is an automated notification. Please do not reply directly.
        """.strip()
    
    @staticmethod
    def generate_breach_notification(incident: Incident, 
                                      affected_count: int) -> str:
        """Generate data breach notification for affected individuals"""
        return f"""
Dear [Customer Name],

We are writing to inform you of a security incident that may have affected 
your personal information.

WHAT HAPPENED:
On {incident.detected_at.strftime('%B %d, %Y')}, we discovered unauthorized 
access to our systems.

WHAT INFORMATION WAS INVOLVED:
[Describe types of personal information affected]

WHAT WE ARE DOING:
We have taken immediate steps to secure our systems and are working with 
cybersecurity experts to investigate this incident.

WHAT YOU CAN DO:
We recommend you:
- Monitor your accounts for suspicious activity
- Change your passwords
- Consider placing a fraud alert on your credit files

FOR MORE INFORMATION:
Contact us at security@example.com or call 1-800-XXX-XXXX

We sincerely apologize for any inconvenience this may cause.

Sincerely,
[Company Name] Security Team
        """.strip()
```

## 📝 Post-Incident Activities

### Incident Report Template

```python
class PostIncidentReport:
    def __init__(self, incident: Incident):
        self.incident = incident
    
    def generate_report(self) -> dict:
        """Generate post-incident report"""
        return {
            'executive_summary': self._executive_summary(),
            'timeline': self._timeline(),
            'root_cause_analysis': self._root_cause(),
            'impact_assessment': self._impact(),
            'lessons_learned': self._lessons_learned(),
            'recommendations': self._recommendations()
        }
    
    def _executive_summary(self) -> str:
        return f"""
On {self.incident.detected_at.strftime('%B %d, %Y')}, {self.incident.title}.
The incident was classified as {self.incident.severity.name} severity and 
was resolved within [X hours/days]. [Brief description of impact].
        """.strip()
    
    def _timeline(self) -> list:
        return self.incident.timeline
    
    def _root_cause(self) -> dict:
        return {
            'root_cause': 'Description of root cause',
            'contributing_factors': [
                'Factor 1',
                'Factor 2'
            ],
            'attack_vector': 'How attacker gained access'
        }
    
    def _impact(self) -> dict:
        return {
            'systems_affected': self.incident.affected_systems,
            'data_compromised': 'Description of data',
            'business_impact': 'Operational/financial impact',
            'customer_impact': 'Impact on customers'
        }
    
    def _lessons_learned(self) -> list:
        return [
            'Lesson 1: What went well',
            'Lesson 2: What could be improved',
            'Lesson 3: Gap identified'
        ]
    
    def _recommendations(self) -> list:
        return [
            {
                'recommendation': 'Description',
                'priority': 'High/Medium/Low',
                'owner': 'Team/Person',
                'deadline': 'Date'
            }
        ]
```

## ❓ Interview Questions

### Basic
**Q: What are the phases of incident response?**
A: Preparation, Detection & Analysis, Containment, Eradication, Recovery, Lessons Learned.

**Q: What's the first thing to do when you detect a breach?**
A: Verify it's legitimate, document findings, preserve evidence, activate IR team.

### Intermediate
**Q: How do you prioritize incidents?**
A: Based on severity (impact, likelihood), business criticality, data sensitivity, active exploitation, regulatory requirements.

**Q: What information should be in an incident report?**
A: Timeline, root cause, impact assessment, affected systems, data compromised, lessons learned, recommendations.

### Advanced
**Q: How do you handle a ransomware attack?**
A: (1) Isolate affected systems, (2) Don't pay ransom, (3) Assess backup integrity, (4) Engage IR team, (5) Report to authorities, (6) Restore from clean backups, (7) Investigate entry point.