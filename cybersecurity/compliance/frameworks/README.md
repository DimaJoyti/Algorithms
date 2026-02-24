# 📜 Security Frameworks and Compliance

Understanding and implementing security frameworks for compliance and risk management.

## 📋 Table of Contents

- [NIST Cybersecurity Framework](#nist-cybersecurity-framework)
- [ISO 27001](#iso-27001)
- [SOC 2](#soc-2)
- [GDPR Compliance](#gdpr-compliance)
- [Interview Questions](#interview-questions)

## 🏛️ NIST Cybersecurity Framework

### Core Functions

```
┌─────────────────────────────────────────────────────────────┐
│           NIST CSF CORE FUNCTIONS                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  IDENTIFY    → Understand assets, risks, capabilities      │
│  PROTECT     → Implement safeguards                        │
│  DETECT      → Identify security events                    │
│  RESPOND     → Take action on detected events              │
│  RECOVER     → Restore capabilities                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### NIST Controls Mapping

```python
from enum import Enum
from dataclasses import dataclass
from typing import List

class NISTFunction(Enum):
    IDENTIFY = "ID"
    PROTECT = "PR"
    DETECT = "DE"
    RESPOND = "RS"
    RECOVER = "RC"

@dataclass
class NISTControl:
    id: str
    function: NISTFunction
    category: str
    description: str
    implementation: str

NIST_CONTROLS = {
    'ID.AM-1': NISTControl(
        id='ID.AM-1',
        function=NISTFunction.IDENTIFY,
        category='Asset Management',
        description='Physical devices and systems within the organization are inventoried',
        implementation='Maintain CMDB with all hardware and software assets'
    ),
    'PR.AC-1': NISTControl(
        id='PR.AC-1',
        function=NISTFunction.PROTECT,
        category='Access Control',
        description='Identities and credentials are issued, managed, verified, revoked, and audited',
        implementation='Implement IAM solution with lifecycle management'
    ),
    'PR.AC-4': NISTControl(
        id='PR.AC-4',
        function=NISTFunction.PROTECT,
        category='Access Control',
        description='Access permissions and authorizations are managed using principles of least privilege',
        implementation='RBAC/ABAC implementation, regular access reviews'
    ),
    'DE.CM-1': NISTControl(
        id='DE.CM-1',
        function=NISTFunction.DETECT,
        category='Security Monitoring',
        description='The network is monitored for potential cybersecurity events',
        implementation='Deploy IDS/IPS, SIEM, network monitoring tools'
    ),
    'RS.AN-1': NISTControl(
        id='RS.AN-1',
        function=NISTFunction.RESPOND,
        category='Analysis',
        description='Notifications from detection systems are investigated',
        implementation='Establish SOC, incident response procedures'
    ),
}

class NISTComplianceChecker:
    def __init__(self):
        self.implemented_controls = {}
        self.gaps = []
    
    def assess_control(self, control_id: str, implemented: bool, 
                       evidence: str = None):
        """Assess control implementation"""
        control = NIST_CONTROLS.get(control_id)
        if not control:
            return
        
        self.implemented_controls[control_id] = {
            'control': control,
            'implemented': implemented,
            'evidence': evidence,
            'assessed_at': datetime.utcnow()
        }
        
        if not implemented:
            self.gaps.append(control_id)
    
    def generate_gap_report(self) -> dict:
        """Generate compliance gap report"""
        by_function = {}
        
        for control_id, data in self.implemented_controls.items():
            control = data['control']
            func = control.function.value
            
            if func not in by_function:
                by_function[func] = {'total': 0, 'implemented': 0, 'gaps': []}
            
            by_function[func]['total'] += 1
            if data['implemented']:
                by_function[func]['implemented'] += 1
            else:
                by_function[func]['gaps'].append(control_id)
        
        return by_function
    
    def calculate_maturity(self) -> float:
        """Calculate compliance maturity score"""
        if not self.implemented_controls:
            return 0.0
        
        implemented = sum(1 for c in self.implemented_controls.values() 
                         if c['implemented'])
        return implemented / len(self.implemented_controls)
```

## 🌐 ISO 27001

### Annex A Controls

```python
from dataclasses import dataclass
from typing import List, Dict

@dataclass
class ISOControl:
    id: str
    domain: str
    description: str
    control_type: str  # preventive, detective, corrective

ISO27001_CONTROLS = {
    'A.5.1.1': ISOControl(
        id='A.5.1.1',
        domain='Information Security Policies',
        description='An information security policy shall be defined and approved',
        control_type='preventive'
    ),
    'A.6.1.2': ISOControl(
        id='A.6.1.2',
        domain='Organization of Information Security',
        description='Information security roles and responsibilities shall be defined',
        control_type='preventive'
    ),
    'A.7.1.1': ISOControl(
        id='A.7.1.1',
        domain='Human Resource Security',
        description='Background verification checks shall be carried out',
        control_type='preventive'
    ),
    'A.8.2.1': ISOControl(
        id='A.8.2.1',
        domain='Asset Management',
        description='Assets shall be classified in terms of legal requirements and criticality',
        control_type='preventive'
    ),
    'A.9.1.1': ISOControl(
        id='A.9.1.1',
        domain='Access Control',
        description='An access control policy shall be established',
        control_type='preventive'
    ),
    'A.12.4.1': ISOControl(
        id='A.12.4.1',
        domain='Operations Security',
        description='Event logs shall be produced and protected',
        control_type='detective'
    ),
    'A.16.1.1': ISOControl(
        id='A.16.1.1',
        domain='Information Security Incident Management',
        description='Responsibilities and procedures shall be defined',
        control_type='corrective'
    ),
}

class ISMSManager:
    """Information Security Management System Manager"""
    
    def __init__(self):
        self.controls = {}
        self.risk_assessments = []
        self.policies = []
    
    def create_statement_of_applicability(self) -> List[str]:
        """Generate Statement of Applicability"""
        soa = ["STATEMENT OF APPLICABILITY"]
        soa.append("=" * 50)
        soa.append("")
        
        for control_id, control in ISO27001_CONTROLS.items():
            applicability = self._determine_applicability(control)
            if applicability['applicable']:
                soa.append(f"Control {control_id}: {control.description}")
                soa.append(f"  Status: {applicability['status']}")
                soa.append(f"  Justification: {applicability['justification']}")
                soa.append("")
        
        return "\n".join(soa)
    
    def _determine_applicability(self, control: ISOControl) -> dict:
        """Determine if control is applicable"""
        # Implementation would assess based on organization
        return {
            'applicable': True,
            'status': 'implemented',
            'justification': 'Applicable to all information assets'
        }
    
    def perform_risk_assessment(self, assets: List[dict], 
                                threats: List[dict]) -> List[dict]:
        """Perform risk assessment"""
        risks = []
        
        for asset in assets:
            for threat in threats:
                likelihood = self._assess_likelihood(asset, threat)
                impact = self._assess_impact(asset, threat)
                risk_score = likelihood * impact
                
                risks.append({
                    'asset': asset['name'],
                    'threat': threat['name'],
                    'likelihood': likelihood,
                    'impact': impact,
                    'risk_score': risk_score,
                    'treatment': self._determine_treatment(risk_score)
                })
        
        self.risk_assessments = risks
        return risks
    
    def _assess_likelihood(self, asset: dict, threat: dict) -> int:
        """Assess likelihood (1-5)"""
        return 3  # Simplified
    
    def _assess_impact(self, asset: dict, threat: dict) -> int:
        """Assess impact (1-5)"""
        return asset.get('criticality', 3)
    
    def _determine_treatment(self, risk_score: int) -> str:
        """Determine risk treatment"""
        if risk_score >= 15:
            return 'mitigate'
        elif risk_score >= 10:
            return 'mitigate'  # or transfer
        elif risk_score >= 5:
            return 'accept'
        else:
            return 'accept'
```

## ✅ SOC 2

### Trust Service Criteria

```python
from enum import Enum
from dataclasses import dataclass

class TrustServiceCriteria(Enum):
    SECURITY = "CC"  # Common Criteria
    AVAILABILITY = "A"
    PROCESSING_INTEGRITY = "PI"
    CONFIDENTIALITY = "C"
    PRIVACY = "P"

@dataclass
class SOC2Control:
    id: str
    criteria: TrustServiceCriteria
    description: str
    points_of_focus: List[str]

SOC2_CONTROLS = {
    'CC6.1': SOC2Control(
        id='CC6.1',
        criteria=TrustServiceCriteria.SECURITY,
        description='Logical access security',
        points_of_focus=[
            'Access to software, data, infrastructure components',
            'Access restricted based on job requirements',
            'Access removed on termination'
        ]
    ),
    'CC6.6': SOC2Control(
        id='CC6.6',
        criteria=TrustServiceCriteria.SECURITY,
        description='Transmission and disposal',
        points_of_focus=[
            'Data protected during transmission',
            'Data disposed of securely'
        ]
    ),
    'CC7.1': SOC2Control(
        id='CC7.1',
        criteria=TrustServiceCriteria.SECURITY,
        description='Detection of anomalies',
        points_of_focus=[
            'System monitored for anomalies',
            'Alerts analyzed and acted upon'
        ]
    ),
}

class SOC2Auditor:
    def __init__(self):
        self.evidence = {}
        self.findings = []
    
    def collect_evidence(self, control_id: str, evidence: dict):
        """Collect evidence for control"""
        if control_id not in self.evidence:
            self.evidence[control_id] = []
        
        self.evidence[control_id].append({
            'collected_at': datetime.utcnow(),
            'evidence': evidence,
            'collector': 'auditor'
        })
    
    def test_control(self, control_id: str, test_procedures: List[str]) -> dict:
        """Test control effectiveness"""
        results = {
            'control_id': control_id,
            'tests': [],
            'effective': True
        }
        
        for procedure in test_procedures:
            # Execute test procedure
            passed = self._execute_test(control_id, procedure)
            
            results['tests'].append({
                'procedure': procedure,
                'passed': passed
            })
            
            if not passed:
                results['effective'] = False
        
        return results
    
    def _execute_test(self, control_id: str, procedure: str) -> bool:
        """Execute specific test procedure"""
        # Implementation would perform actual testing
        return True
    
    def generate_audit_report(self) -> str:
        """Generate SOC 2 audit report"""
        report = [
            "SOC 2 TYPE II AUDIT REPORT",
            "=" * 50,
            "",
            "SCOPE:",
            "[Describe systems in scope]",
            "",
            "PERIOD:",
            "[Audit period dates]",
            "",
            "CONTROLS TESTED:",
            "-" * 50
        ]
        
        for control_id, control in SOC2_CONTROLS.items():
            evidence_list = self.evidence.get(control_id, [])
            report.append(f"\n{control_id}: {control.description}")
            report.append(f"  Evidence collected: {len(evidence_list)} items")
        
        return "\n".join(report)
```

## 🇪🇺 GDPR Compliance

### Key Requirements

```python
from dataclasses import dataclass
from enum import Enum
from typing import List, Optional
from datetime import datetime

class LawfulBasis(Enum):
    CONSENT = "consent"
    CONTRACT = "contract"
    LEGAL_OBLIGATION = "legal_obligation"
    VITAL_INTERESTS = "vital_interests"
    PUBLIC_TASK = "public_task"
    LEGITIMATE_INTERESTS = "legitimate_interests"

@dataclass
class DataProcessingActivity:
    name: str
    purpose: str
    lawful_basis: LawfulBasis
    data_categories: List[str]
    recipients: List[str]
    retention_period: str
    safeguards: List[str]

class GDPRCompliance:
    def __init__(self):
        self.processing_activities = []
        self.data_subjects = {}
        self.dsr_requests = []
    
    def register_processing_activity(self, activity: DataProcessingActivity):
        """Register data processing activity (Article 30)"""
        self.processing_activities.append(activity)
    
    def handle_dsr_request(self, request_type: str, data_subject_id: str,
                          details: dict) -> dict:
        """Handle Data Subject Rights request"""
        request = {
            'id': self._generate_id(),
            'type': request_type,  # access, rectification, erasure, portability
            'data_subject_id': data_subject_id,
            'received_at': datetime.utcnow(),
            'deadline': datetime.utcnow() + timedelta(days=30),
            'status': 'pending',
            'details': details
        }
        
        self.dsr_requests.append(request)
        return request
    
    def process_erasure_request(self, request_id: str) -> dict:
        """Process right to erasure (Article 17)"""
        request = next((r for r in self.dsr_requests if r['id'] == request_id), None)
        
        if not request:
            return {'success': False, 'reason': 'Request not found'}
        
        # Verify grounds for erasure
        can_erase = self._verify_erasure_grounds(request)
        
        if can_erase:
            # Delete data
            deleted = self._delete_data(request['data_subject_id'])
            request['status'] = 'completed'
            return {'success': True, 'deleted': deleted}
        else:
            request['status'] = 'rejected'
            return {'success': False, 'reason': 'Legal grounds to retain'}
    
    def check_breach_notification_required(self, breach: dict) -> dict:
        """Check if breach notification required (Article 33/34)"""
        # Assess risk to rights and freedoms
        risk_level = self._assess_breach_risk(breach)
        
        result = {
            'risk_level': risk_level,
            'supervisory_authority_notification': False,
            'data_subject_notification': False,
            'deadline': None
        }
        
        if risk_level in ['high', 'critical']:
            result['supervisory_authority_notification'] = True
            result['deadline'] = datetime.utcnow() + timedelta(hours=72)
        
        if risk_level == 'critical':
            result['data_subject_notification'] = True
        
        return result
    
    def _assess_breach_risk(self, breach: dict) -> str:
        """Assess breach risk level"""
        score = 0
        
        # Data sensitivity
        sensitive_data = ['health', 'racial', 'political', 'religious', 
                         'genetic', 'biometric', 'sexual']
        if any(d in breach.get('data_types', []) for d in sensitive_data):
            score += 3
        
        # Volume
        if breach.get('records_affected', 0) > 10000:
            score += 2
        
        # Vulnerability of data subjects
        if breach.get('vulnerable_subjects'):
            score += 2
        
        if score >= 5:
            return 'critical'
        elif score >= 3:
            return 'high'
        elif score >= 1:
            return 'medium'
        else:
            return 'low'

class ConsentManager:
    """GDPR Consent Management"""
    
    def __init__(self):
        self.consents = {}
    
    def record_consent(self, data_subject_id: str, purpose: str,
                      consent_text: str, method: str) -> dict:
        """Record consent (Article 7)"""
        consent = {
            'id': self._generate_id(),
            'data_subject_id': data_subject_id,
            'purpose': purpose,
            'consent_text': consent_text,
            'method': method,  # checkbox, signature, etc.
            'granted_at': datetime.utcnow(),
            'ip_address': self._get_ip(),
            'active': True
        }
        
        if data_subject_id not in self.consents:
            self.consents[data_subject_id] = []
        
        self.consents[data_subject_id].append(consent)
        return consent
    
    def withdraw_consent(self, data_subject_id: str, purpose: str):
        """Withdraw consent"""
        if data_subject_id in self.consents:
            for consent in self.consents[data_subject_id]:
                if consent['purpose'] == purpose:
                    consent['active'] = False
                    consent['withdrawn_at'] = datetime.utcnow()
    
    def verify_consent(self, data_subject_id: str, purpose: str) -> bool:
        """Verify valid consent exists"""
        if data_subject_id not in self.consents:
            return False
        
        for consent in self.consents[data_subject_id]:
            if (consent['purpose'] == purpose and 
                consent['active'] and
                not self._is_expired(consent)):
                return True
        
        return False
```

## ❓ Interview Questions

### Basic
**Q: What is NIST CSF?**
A: Framework of standards, guidelines, and best practices for managing cybersecurity risk. Five core functions: Identify, Protect, Detect, Respond, Recover.

**Q: What is SOC 2?**
A: Audit report on controls relevant to security, availability, processing integrity, confidentiality, or privacy. Based on Trust Service Criteria.

### Intermediate
**Q: What's the difference between SOC 2 Type I and Type II?**
A: Type I evaluates control design at a point in time. Type II evaluates operating effectiveness over a period (6-12 months).

**Q: What are GDPR data subject rights?**
A: Right to access, rectification, erasure, restriction, portability, object, and not be subject to automated decision-making.

### Advanced
**Q: How would you prepare for a SOC 2 audit?**
A: (1) Define scope, (2) Map controls to criteria, (3) Document policies and procedures, (4) Implement controls, (5) Collect evidence for 6+ months, (6) Perform readiness assessment, (7) Remediate gaps.