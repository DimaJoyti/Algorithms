# 🔍 Security Analyst Interview Questions

Comprehensive interview preparation for Security Analyst positions.

## 📋 Table of Contents

- [Technical Fundamentals](#technical-fundamentals)
- [Threat Detection](#threat-detection)
- [Log Analysis](#log-analysis)
- [Incident Response](#incident-response)
- [Tools and Technologies](#tools-and-technologies)
- [Scenario-Based Questions](#scenario-based-questions)

## 🎯 Technical Fundamentals

### Basic Questions

**Q1: What is the CIA Triad?**

**A**: The CIA Triad represents three core principles:
- **Confidentiality**: Ensuring data is only accessible to authorized users
- **Integrity**: Ensuring data accuracy and preventing unauthorized modification
- **Availability**: Ensuring systems and data are accessible when needed

**Q2: Explain the difference between vulnerability, threat, and risk.**

**A**:
- **Vulnerability**: A weakness in a system (e.g., unpatched software)
- **Threat**: A potential danger that could exploit a vulnerability (e.g., attacker)
- **Risk**: The likelihood and impact of a threat exploiting a vulnerability (Risk = Threat × Vulnerability × Asset Value)

**Q3: What is defense in depth?**

**A**: A layered security approach using multiple controls at different levels:
1. Physical security
2. Network security (firewalls, IDS/IPS)
3. Host security (antivirus, hardening)
4. Application security (input validation)
5. Data security (encryption)

## 🔍 Threat Detection

### Intermediate Questions

**Q4: How do you differentiate between a false positive and a true positive security alert?**

**A**: 
- **True Positive**: Alert correctly identifies malicious activity
- **False Positive**: Alert incorrectly flags legitimate activity as malicious
- Investigation steps:
  1. Check the source context (IP reputation, user behavior)
  2. Correlate with other logs and events
  3. Verify against known good baselines
  4. Consider the business context

**Q5: What indicators of compromise (IOCs) would you look for in a suspected breach?**

**A**:
- **Network**: Unusual outbound connections, communication with known C2 servers, data exfiltration patterns
- **Host**: Unknown processes, unusual scheduled tasks, modified system files
- **User**: Anomalous login patterns, privilege escalation, unusual file access
- **File**: New/modified executables, scripts in unusual locations
- **Log**: Cleared logs, authentication failures, service restarts

**Q6: Explain the MITRE ATT&CK framework and how you use it.**

**A**: MITRE ATT&CK is a knowledge base of adversary tactics, techniques, and procedures (TTPs). Used for:
1. Mapping detected threats to known techniques
2. Understanding attack progression
3. Identifying detection gaps
4. Building detection rules
5. Threat intelligence correlation

## 📊 Log Analysis

### Practical Questions

**Q7: Analyze this authentication log. What do you see?**

```
2024-01-15 08:00:12 192.168.1.50 - admin LOGIN_SUCCESS
2024-01-15 08:15:23 10.0.0.100 - admin LOGIN_FAILURE
2024-01-15 08:15:25 10.0.0.100 - admin LOGIN_FAILURE
2024-01-15 08:15:27 10.0.0.100 - admin LOGIN_FAILURE
2024-01-15 08:15:29 10.0.0.100 - admin LOGIN_SUCCESS
2024-01-15 08:16:00 10.0.0.100 - admin CREATE_USER hacker
```

**A**: Indicators of compromise:
1. Brute force attempt (rapid failures followed by success)
2. Different IP from previous legitimate access (192.168.1.50 vs 10.0.0.100)
3. Immediate creation of new user account after compromise
4. Suggests account takeover and persistence establishment

**Q8: What SIEM queries would you write to detect lateral movement?**

**A**:
```sql
-- Detect RDP connections to multiple internal hosts
SELECT source_ip, COUNT(DISTINCT dest_ip) as unique_targets
FROM connection_logs 
WHERE dest_port = 3389 
  AND timestamp > NOW() - INTERVAL 1 HOUR
GROUP BY source_ip
HAVING unique_targets > 3

-- Detect unusual service account usage
SELECT username, source_ip, COUNT(*) as auth_count
FROM auth_logs
WHERE account_type = 'service'
  AND timestamp > NOW() - INTERVAL 24 HOUR
  AND login_type != 'batch'
GROUP BY username, source_ip
HAVING auth_count > 10
```

## 🚨 Incident Response

### Scenario Questions

**Q9: Walk through your response to a suspected malware infection on an employee workstation.**

**A**:
1. **Identification**: Confirm the alert, gather initial information
2. **Containment**:
   - Isolate the workstation from the network
   - Preserve volatile memory if possible
   - Block known malicious indicators
3. **Investigation**:
   - Analyze memory and disk for malware
   - Review logs for infection vector
   - Determine scope of compromise
4. **Eradication**: Remove malware, patch vulnerability
5. **Recovery**: Rebuild or clean system, restore from backup
6. **Lessons Learned**: Document findings, improve detection

**Q10: How would you handle a DDoS attack affecting your company's web services?**

**A**:
1. **Immediate**: Activate DDoS mitigation service, enable rate limiting
2. **Communication**: Notify stakeholders, update status page
3. **Mitigation**: 
   - Block attack source IPs
   - Enable CDN/CDN protection
   - Scale up infrastructure
   - Switch to static content delivery
4. **Monitoring**: Watch for attack resumption
5. **Post-incident**: Review defenses, implement long-term improvements

## 🛠️ Tools and Technologies

### Practical Questions

**Q11: Compare IDS vs IPS. When would you use each?**

**A**:
| Feature | IDS | IPS |
|---------|-----|-----|
| Position | Out-of-band | In-line |
| Action | Alert only | Alert + Block |
| Impact | No latency | Potential latency |
| Use Case | Monitoring, compliance | Active protection |

Use IDS when you need monitoring without affecting traffic. Use IPS when you need active blocking of threats.

**Q12: What is your approach to tuning SIEM rules to reduce false positives?**

**A**:
1. Analyze false positive patterns
2. Add context (asset criticality, user role)
3. Implement whitelists for known good behavior
4. Adjust thresholds based on environment
5. Correlate with multiple data sources
6. Test rules in detection-only mode first
7. Regular rule review and refinement

## 🎭 Scenario-Based Questions

**Q13: You receive an alert for unusual data transfer from a database server. Walk through your investigation.**

**A**:
1. **Validate**: Is this expected behavior? Check maintenance windows, backup schedules
2. **Quantify**: How much data? Where is it going?
3. **Timeline**: When did it start? Correlate with other events
4. **Access**: Who has access? Were there privilege changes?
5. **Network**: External or internal destination? Known or unknown?
6. **Process**: What process initiated the transfer?
7. **Escalate**: If confirmed exfiltration, initiate incident response

**Q14: A user reports receiving a suspicious email with an attachment. What do you do?**

**A**:
1. **Isolate**: Don't open the attachment on a production system
2. **Collect**: Obtain email headers and attachment
3. **Analyze**:
   - Check sender reputation
   - Analyze headers for spoofing
   - Sandbox the attachment
   - Check for known IOCs
4. **Search**: Look for similar emails to other users
5. **Respond**:
   - Block sender/domain
   - Update email filters
   - Notify potentially affected users
6. **Document**: Record findings for future reference

## 📈 Behavioral Questions

**Q15: Describe a time you identified a significant security issue.**

**A** Structure: STAR method
- **Situation**: Describe the context
- **Task**: What was your responsibility
- **Action**: What you did to investigate and resolve
- **Result**: The outcome and impact

Example:
"At my previous role, I noticed unusual authentication patterns during routine log review. After investigation, I discovered an attacker had compromised a service account and was attempting lateral movement. I immediately isolated affected systems, reset credentials, and worked with the team to close the attack vector. This led to implementing additional monitoring for service account activity."

## 💡 Tips for Success

1. **Know your tools**: Be familiar with SIEM, EDR, vulnerability scanners
2. **Practice log analysis**: Work with sample logs and challenges
3. **Stay current**: Follow security news, threats, and techniques
4. **Think out loud**: Explain your thought process during scenarios
5. **Ask clarifying questions**: Understand the full context before answering