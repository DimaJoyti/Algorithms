# 🎯 Threat Modeling - Systematic Security Analysis

Threat modeling is a structured approach to identifying, quantifying, and addressing security threats in systems and applications.

## 📋 Table of Contents

- [Overview](#overview)
- [Threat Modeling Process](#threat-modeling-process)
- [STRIDE Methodology](#stride-methodology)
- [PASTA Framework](#pasta-framework)
- [Attack Trees](#attack-trees)
- [Practical Examples](#practical-examples)
- [Tools and Templates](#tools-and-templates)
- [Interview Questions](#interview-questions)

## 🎯 Overview

Threat modeling helps organizations:
- **Identify threats** early in the development lifecycle
- **Prioritize security efforts** based on risk assessment
- **Design effective countermeasures** for identified threats
- **Communicate security risks** to stakeholders

### When to Perform Threat Modeling
- **Design Phase**: Architecture and system design
- **Implementation**: Code reviews and security testing
- **Deployment**: Infrastructure and configuration review
- **Maintenance**: Regular security assessments

## 🔄 Threat Modeling Process

### 1. Define Security Objectives
- Identify assets to protect
- Determine security requirements
- Establish risk tolerance levels

### 2. Create System Overview
- Document system architecture
- Identify data flows
- Map trust boundaries
- List external dependencies

### 3. Identify Threats
- Use structured methodologies (STRIDE, PASTA)
- Consider different attack vectors
- Analyze threat actors and motivations

### 4. Assess Risk
- Evaluate likelihood and impact
- Prioritize threats based on risk level
- Consider existing security controls

### 5. Design Countermeasures
- Select appropriate security controls
- Implement defense-in-depth strategy
- Plan monitoring and detection

### 6. Validate and Monitor
- Test security controls effectiveness
- Monitor for new threats
- Update threat model regularly

## ⚔️ STRIDE Methodology

STRIDE categorizes threats into six types:

### **S**poofing Identity
- **Definition**: Pretending to be someone or something else
- **Examples**: Email spoofing, IP spoofing, credential theft
- **Countermeasures**: Authentication, digital certificates, MFA

### **T**ampering with Data
- **Definition**: Malicious modification of data
- **Examples**: SQL injection, man-in-the-middle attacks
- **Countermeasures**: Input validation, encryption, digital signatures

### **R**epudiation
- **Definition**: Denying actions or transactions
- **Examples**: Claiming "I didn't send that email"
- **Countermeasures**: Audit logs, digital signatures, timestamps

### **I**nformation Disclosure
- **Definition**: Exposing information to unauthorized parties
- **Examples**: Data breaches, eavesdropping, information leakage
- **Countermeasures**: Encryption, access controls, data classification

### **D**enial of Service
- **Definition**: Making systems unavailable to legitimate users
- **Examples**: DDoS attacks, resource exhaustion
- **Countermeasures**: Rate limiting, load balancing, redundancy

### **E**levation of Privilege
- **Definition**: Gaining unauthorized access or permissions
- **Examples**: Buffer overflows, privilege escalation
- **Countermeasures**: Least privilege, input validation, sandboxing

### STRIDE Analysis Template

```python
class STRIDEAnalysis:
    def __init__(self, component_name):
        self.component = component_name
        self.threats = {
            'Spoofing': [],
            'Tampering': [],
            'Repudiation': [],
            'Information_Disclosure': [],
            'Denial_of_Service': [],
            'Elevation_of_Privilege': []
        }
    
    def add_threat(self, category, description, likelihood, impact, mitigation):
        threat = {
            'description': description,
            'likelihood': likelihood,  # 1-5 scale
            'impact': impact,         # 1-5 scale
            'risk_score': likelihood * impact,
            'mitigation': mitigation
        }
        self.threats[category].append(threat)
    
    def get_high_risk_threats(self, threshold=15):
        high_risk = []
        for category, threats in self.threats.items():
            for threat in threats:
                if threat['risk_score'] >= threshold:
                    high_risk.append({
                        'category': category,
                        'threat': threat
                    })
        return sorted(high_risk, key=lambda x: x['threat']['risk_score'], reverse=True)

# Example usage
web_app = STRIDEAnalysis("Web Application Login")

web_app.add_threat(
    'Spoofing',
    'Attacker uses stolen credentials to impersonate legitimate user',
    4,  # High likelihood
    5,  # High impact
    'Implement MFA and account lockout policies'
)

web_app.add_threat(
    'Tampering',
    'SQL injection attack modifies database records',
    3,  # Medium likelihood
    5,  # High impact
    'Use parameterized queries and input validation'
)

high_risks = web_app.get_high_risk_threats()
for risk in high_risks:
    print(f"{risk['category']}: {risk['threat']['description']}")
    print(f"Risk Score: {risk['threat']['risk_score']}")
    print(f"Mitigation: {risk['threat']['mitigation']}\n")
```

## 🍝 PASTA Framework

Process for Attack Simulation and Threat Analysis (PASTA) is a seven-stage methodology:

### Stage 1: Define Objectives
- Business objectives
- Security and compliance requirements
- Risk tolerance levels

### Stage 2: Define Technical Scope
- Application boundaries
- Infrastructure components
- Data flows and trust boundaries

### Stage 3: Application Decomposition
- Identify application components
- Map data flows
- Document trust boundaries

### Stage 4: Threat Analysis
- Threat intelligence gathering
- Attack vector analysis
- Threat actor profiling

### Stage 5: Weakness and Vulnerability Analysis
- Code review findings
- Infrastructure vulnerabilities
- Configuration weaknesses

### Stage 6: Attack Modeling
- Attack tree construction
- Attack scenario development
- Exploit analysis

### Stage 7: Risk and Impact Analysis
- Risk scoring and prioritization
- Business impact assessment
- Countermeasure recommendations

### PASTA Implementation Example

```go
package main

import (
    "fmt"
    "sort"
)

type ThreatActor struct {
    Name        string
    Motivation  string
    Capability  int // 1-5 scale
    Resources   int // 1-5 scale
}

type Vulnerability struct {
    ID          string
    Description string
    CVSS        float64
    Component   string
}

type AttackPath struct {
    ID          string
    Description string
    Steps       []string
    Likelihood  int // 1-5 scale
    Impact      int // 1-5 scale
    RiskScore   int
}

type PASTAModel struct {
    Objectives      []string
    TechnicalScope  []string
    Components      []string
    ThreatActors    []ThreatActor
    Vulnerabilities []Vulnerability
    AttackPaths     []AttackPath
}

func (p *PASTAModel) CalculateRisk() {
    for i := range p.AttackPaths {
        p.AttackPaths[i].RiskScore = p.AttackPaths[i].Likelihood * p.AttackPaths[i].Impact
    }
    
    // Sort by risk score (highest first)
    sort.Slice(p.AttackPaths, func(i, j int) bool {
        return p.AttackPaths[i].RiskScore > p.AttackPaths[j].RiskScore
    })
}

func (p *PASTAModel) GetHighRiskAttacks(threshold int) []AttackPath {
    var highRisk []AttackPath
    for _, attack := range p.AttackPaths {
        if attack.RiskScore >= threshold {
            highRisk = append(highRisk, attack)
        }
    }
    return highRisk
}

// Example usage
func main() {
    model := PASTAModel{
        Objectives: []string{
            "Protect customer financial data",
            "Ensure system availability",
            "Maintain regulatory compliance",
        },
        TechnicalScope: []string{
            "Web application",
            "Database server",
            "Payment processing API",
        },
    }
    
    // Add threat actors
    model.ThreatActors = append(model.ThreatActors, ThreatActor{
        Name:       "External Hacker",
        Motivation: "Financial gain",
        Capability: 4,
        Resources:  3,
    })
    
    // Add attack paths
    model.AttackPaths = append(model.AttackPaths, AttackPath{
        ID:          "ATK-001",
        Description: "SQL Injection leading to data breach",
        Steps: []string{
            "Identify vulnerable input field",
            "Craft malicious SQL payload",
            "Extract sensitive data",
            "Exfiltrate customer records",
        },
        Likelihood: 3,
        Impact:     5,
    })
    
    model.CalculateRisk()
    
    highRiskAttacks := model.GetHighRiskAttacks(12)
    fmt.Printf("High Risk Attacks (%d found):\n", len(highRiskAttacks))
    for _, attack := range highRiskAttacks {
        fmt.Printf("- %s (Risk: %d)\n", attack.Description, attack.RiskScore)
    }
}
```

## 🌳 Attack Trees

Attack trees provide a visual representation of how attacks can be carried out against a system.

### Attack Tree Structure
- **Root Node**: Primary attack goal
- **Leaf Nodes**: Basic attack steps
- **AND Gates**: All child nodes must succeed
- **OR Gates**: Any child node can succeed

### Example: ATM Attack Tree

```
                    Steal Money from ATM
                           |
                    ┌──────┴──────┐
                    │             │
              Physical Attack  Logical Attack
                    │             │
            ┌───────┼───────┐     ├─────────────┐
            │       │       │     │             │
        Break    Cut     Explosive  Card      Network
        Into     Power   Device     Skimming  Attack
        ATM      Cable              │         │
                                   │         │
                              ┌────┴───┐    │
                              │        │    │
                          Install   Steal   Man-in-
                          Skimmer   PINs    Middle
```

### Attack Tree Implementation

```javascript
class AttackTreeNode {
    constructor(name, type = 'OR', cost = 0, probability = 0) {
        this.name = name;
        this.type = type; // 'OR', 'AND', 'LEAF'
        this.cost = cost;
        this.probability = probability;
        this.children = [];
    }
    
    addChild(child) {
        this.children.push(child);
    }
    
    calculateProbability() {
        if (this.type === 'LEAF') {
            return this.probability;
        }
        
        if (this.type === 'OR') {
            // At least one child succeeds
            let failureProbability = 1;
            for (let child of this.children) {
                failureProbability *= (1 - child.calculateProbability());
            }
            return 1 - failureProbability;
        }
        
        if (this.type === 'AND') {
            // All children must succeed
            let successProbability = 1;
            for (let child of this.children) {
                successProbability *= child.calculateProbability();
            }
            return successProbability;
        }
        
        return 0;
    }
    
    calculateMinCost() {
        if (this.type === 'LEAF') {
            return this.cost;
        }
        
        if (this.type === 'OR') {
            // Minimum cost among children
            return Math.min(...this.children.map(child => child.calculateMinCost()));
        }
        
        if (this.type === 'AND') {
            // Sum of all children costs
            return this.children.reduce((sum, child) => sum + child.calculateMinCost(), 0);
        }
        
        return Infinity;
    }
}

// Example: Web Application Attack Tree
const webAppAttack = new AttackTreeNode('Compromise Web Application', 'OR');

const sqlInjection = new AttackTreeNode('SQL Injection', 'AND');
sqlInjection.addChild(new AttackTreeNode('Find Vulnerable Input', 'LEAF', 10, 0.8));
sqlInjection.addChild(new AttackTreeNode('Craft Payload', 'LEAF', 5, 0.9));
sqlInjection.addChild(new AttackTreeNode('Execute Attack', 'LEAF', 2, 0.7));

const xssAttack = new AttackTreeNode('XSS Attack', 'AND');
xssAttack.addChild(new AttackTreeNode('Identify XSS Vector', 'LEAF', 8, 0.6));
xssAttack.addChild(new AttackTreeNode('Create Malicious Script', 'LEAF', 3, 0.9));
xssAttack.addChild(new AttackTreeNode('Social Engineer Victim', 'LEAF', 15, 0.4));

webAppAttack.addChild(sqlInjection);
webAppAttack.addChild(xssAttack);

console.log(`Attack Success Probability: ${webAppAttack.calculateProbability().toFixed(3)}`);
console.log(`Minimum Attack Cost: ${webAppAttack.calculateMinCost()}`);
```

## 🛠️ Tools and Templates

### Popular Threat Modeling Tools
- **Microsoft Threat Modeling Tool**: Free, STRIDE-based
- **OWASP Threat Dragon**: Open-source, web-based
- **IriusRisk**: Commercial, automated threat modeling
- **ThreatModeler**: Enterprise threat modeling platform

### Threat Modeling Template

```markdown
# Threat Model: [System Name]

## 1. System Overview
- **Purpose**: [What does the system do?]
- **Users**: [Who uses the system?]
- **Assets**: [What needs protection?]
- **Trust Boundaries**: [Where do trust levels change?]

## 2. Data Flow Diagram
[Include DFD showing data flows, processes, and trust boundaries]

## 3. Threat Analysis
| Threat ID | Category | Description | Likelihood | Impact | Risk Score | Mitigation |
|-----------|----------|-------------|------------|--------|------------|------------|
| T001      | Spoofing | ...         | 3          | 4      | 12         | ...        |

## 4. Security Controls
| Control ID | Type | Description | Threats Addressed |
|------------|------|-------------|-------------------|
| C001       | Preventive | MFA | T001, T003 |

## 5. Residual Risks
[List risks that remain after implementing controls]

## 6. Recommendations
[Prioritized list of security improvements]
```

## ❓ Interview Questions

### Basic Questions
1. **Q**: What is threat modeling and why is it important?
   **A**: Threat modeling is a structured process to identify, analyze, and mitigate security threats in systems. It's important because it helps find security issues early, prioritizes security efforts, and ensures comprehensive security coverage.

2. **Q**: Explain the STRIDE methodology.
   **A**: STRIDE categorizes threats into six types: Spoofing (identity), Tampering (data), Repudiation (denying actions), Information disclosure, Denial of service, and Elevation of privilege. Each category helps systematically identify different types of security threats.

### Intermediate Questions
3. **Q**: How would you threat model a mobile banking application?
   **A**: I'd start by identifying assets (user credentials, financial data), map data flows (login, transactions, account queries), identify trust boundaries (client-server, internal services), then apply STRIDE to each component, focusing on high-risk areas like authentication and transaction processing.

4. **Q**: What's the difference between threat modeling and vulnerability assessment?
   **A**: Threat modeling is proactive and design-focused, identifying potential threats during development. Vulnerability assessment is reactive, testing existing systems for known vulnerabilities. Threat modeling asks "what could go wrong?" while vulnerability assessment asks "what is currently wrong?"

### Advanced Questions
5. **Q**: How do you prioritize threats in a complex system with hundreds of potential threats?
   **A**: Use risk-based prioritization: calculate risk scores (likelihood × impact), consider business context, evaluate existing controls, focus on high-value assets, and consider threat actor capabilities. Also factor in regulatory requirements and business continuity needs.

---

**🎯 Key Takeaway**: Threat modeling is most effective when integrated into the development lifecycle, not treated as a one-time activity. Regular updates and stakeholder involvement are crucial for success.
