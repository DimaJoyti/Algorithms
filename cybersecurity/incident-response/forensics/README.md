# 🔬 Digital Forensics

Techniques for investigating security incidents and preserving evidence.

## 📋 Table of Contents

- [Forensic Methodology](#forensic-methodology)
- [Evidence Collection](#evidence-collection)
- [Log Analysis](#log-analysis)
- [Memory Forensics](#memory-forensics)
- [Interview Questions](#interview-questions)

## 🔍 Forensic Methodology

### Chain of Custody

```
┌────────────────────────────────────────────────────────┐
│                CHAIN OF CUSTODY                        │
├────────────────────────────────────────────────────────┤
│                                                        │
│  1. IDENTIFICATION   → Identify potential evidence    │
│                                                        │
│  2. PRESERVATION     → Protect from alteration        │
│                                                        │
│  3. COLLECTION       → Documented acquisition         │
│                                                        │
│  4. EXAMINATION      → Extract data                   │
│                                                        │
│  5. ANALYSIS         → Interpret findings             │
│                                                        │
│  6. PRESENTATION     → Report findings                │
│                                                        │
└────────────────────────────────────────────────────────┘
```

### Evidence Documentation

```python
from dataclasses import dataclass
from datetime import datetime
from enum import Enum
import hashlib

class EvidenceType(Enum):
    DISK_IMAGE = "disk_image"
    MEMORY_DUMP = "memory_dump"
    LOG_FILE = "log_file"
    NETWORK_CAPTURE = "network_capture"
    SCREENSHOT = "screenshot"
    DOCUMENT = "document"

@dataclass
class Evidence:
    id: str
    type: EvidenceType
    source: str
    collected_by: str
    collected_at: datetime
    original_location: str
    hash_md5: str
    hash_sha256: str
    chain_of_custody: list
    notes: str = ""

class EvidenceManager:
    def __init__(self):
        self.evidence: dict[str, Evidence] = {}
    
    def collect_evidence(self, filepath: str, evidence_type: EvidenceType,
                        collector: str, notes: str = "") -> Evidence:
        """Collect and document evidence"""
        # Calculate hashes
        md5 = self._calculate_hash(filepath, 'md5')
        sha256 = self._calculate_hash(filepath, 'sha256')
        
        evidence = Evidence(
            id=self._generate_id(),
            type=evidence_type,
            source=filepath,
            collected_by=collector,
            collected_at=datetime.utcnow(),
            original_location=filepath,
            hash_md5=md5,
            hash_sha256=sha256,
            chain_of_custody=[{
                'timestamp': datetime.utcnow(),
                'action': 'collected',
                'by': collector,
                'notes': notes
            }],
            notes=notes
        )
        
        self.evidence[evidence.id] = evidence
        return evidence
    
    def _calculate_hash(self, filepath: str, algorithm: str) -> str:
        """Calculate file hash"""
        h = hashlib.new(algorithm)
        with open(filepath, 'rb') as f:
            for chunk in iter(lambda: f.read(4096), b''):
                h.update(chunk)
        return h.hexdigest()
    
    def verify_integrity(self, evidence_id: str) -> bool:
        """Verify evidence hasn't been tampered with"""
        evidence = self.evidence.get(evidence_id)
        if not evidence:
            return False
        
        current_hash = self._calculate_hash(evidence.source, 'sha256')
        return current_hash == evidence.hash_sha256
    
    def transfer_custody(self, evidence_id: str, from_person: str, 
                        to_person: str, reason: str):
        """Transfer custody of evidence"""
        evidence = self.evidence.get(evidence_id)
        if not evidence:
            return
        
        evidence.chain_of_custody.append({
            'timestamp': datetime.utcnow(),
            'action': 'custody_transfer',
            'from': from_person,
            'to': to_person,
            'reason': reason
        })
    
    def generate_chain_report(self, evidence_id: str) -> str:
        """Generate chain of custody report"""
        evidence = self.evidence.get(evidence_id)
        if not evidence:
            return "Evidence not found"
        
        report = [
            f"CHAIN OF CUSTODY REPORT - {evidence_id}",
            "=" * 50,
            f"Evidence Type: {evidence.type.value}",
            f"Original Location: {evidence.original_location}",
            f"MD5: {evidence.hash_md5}",
            f"SHA256: {evidence.hash_sha256}",
            "",
            "CUSTODY HISTORY:",
            "-" * 50
        ]
        
        for entry in evidence.chain_of_custody:
            report.append(
                f"{entry['timestamp']}: {entry['action']} by {entry.get('by', entry.get('to'))}"
            )
        
        return "\n".join(report)
```

## 📦 Evidence Collection

### Disk Imaging

```python
import subprocess
from pathlib import Path

class DiskImager:
    def __init__(self, output_dir: str):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
    
    def create_image(self, device: str, case_id: str) -> dict:
        """Create forensic disk image"""
        output_file = self.output_dir / f"{case_id}_disk.img"
        log_file = self.output_dir / f"{case_id}_imaging.log"
        
        # Use dd or dc3dd for imaging
        result = subprocess.run([
            'dc3dd',
            'if=' + device,
            'of=' + str(output_file),
            'log=' + str(log_file),
            'hash=md5,sha256',
            'verify=on'
        ], capture_output=True, text=True)
        
        return {
            'image_file': str(output_file),
            'log_file': str(log_file),
            'success': result.returncode == 0,
            'size': output_file.stat().st_size if output_file.exists() else 0
        }
    
    def mount_readonly(self, image_path: str, mount_point: str):
        """Mount disk image read-only"""
        subprocess.run([
            'mount', '-o', 'ro,loop,noexec',
            image_path, mount_point
        ])
    
    def unmount(self, mount_point: str):
        """Unmount disk image"""
        subprocess.run(['umount', mount_point])

class ArtifactCollector:
    """Collect common forensic artifacts"""
    
    ARTIFACTS = {
        'windows': [
            ('Registry', r'%SystemRoot%\System32\config\*'),
            ('Event Logs', r'%SystemRoot%\System32\winevt\Logs\*'),
            ('Prefetch', r'%SystemRoot%\Prefetch\*'),
            ('Browser History', r'%AppData%\Local\*\User Data\Default\History'),
            ('Recent Documents', r'%AppData%\Microsoft\Windows\Recent\*'),
            ('Recycle Bin', r'$Recycle.Bin\*\*'),
        ],
        'linux': [
            ('Auth Logs', '/var/log/auth.log*'),
            ('Syslog', '/var/log/syslog*'),
            ('Bash History', '/home/*/.bash_history'),
            ('Cron Logs', '/var/log/cron*'),
            ('Secure Logs', '/var/log/secure*'),
        ]
    }
    
    def collect_artifacts(self, system_type: str, output_dir: str) -> list:
        """Collect system artifacts"""
        artifacts = self.ARTIFACTS.get(system_type, [])
        collected = []
        
        for name, pattern in artifacts:
            matches = list(Path('/').glob(pattern))
            for match in matches:
                dest = Path(output_dir) / name / match.name
                dest.parent.mkdir(parents=True, exist_ok=True)
                shutil.copy2(match, dest)
                collected.append({
                    'name': name,
                    'source': str(match),
                    'destination': str(dest)
                })
        
        return collected
```

## 📊 Log Analysis

### Log Parser

```python
import re
from datetime import datetime
from collections import Counter
from dataclasses import dataclass

@dataclass
class LogEntry:
    timestamp: datetime
    source: str
    level: str
    message: str
    raw: str

class LogAnalyzer:
    def __init__(self):
        self.patterns = {
            'auth_failure': re.compile(
                r'Failed password for .* from (\d+\.\d+\.\d+\.\d+)'
            ),
            'auth_success': re.compile(
                r'Accepted password for .* from (\d+\.\d+\.\d+\.\d+)'
            ),
            'sudo_usage': re.compile(
                r'sudo:\s+(\w+)\s+:.*COMMAND=(.+)'
            ),
            'ssh_connection': re.compile(
                r'sshd\[\d+\]: Connection from (\d+\.\d+\.\d+\.\d+)'
            )
        }
    
    def parse_syslog(self, filepath: str) -> list[LogEntry]:
        """Parse syslog format"""
        entries = []
        
        with open(filepath, 'r') as f:
            for line in f:
                entry = self._parse_line(line.strip())
                if entry:
                    entries.append(entry)
        
        return entries
    
    def _parse_line(self, line: str) -> LogEntry:
        """Parse single log line"""
        # Syslog format: Mon DD HH:MM:SS hostname process: message
        match = re.match(
            r'(\w{3}\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+(\S+):\s+(.+)',
            line
        )
        
        if match:
            timestamp_str, host, process, message = match.groups()
            timestamp = datetime.strptime(
                f"{datetime.now().year} {timestamp_str}",
                "%Y %b %d %H:%M:%S"
            )
            
            return LogEntry(
                timestamp=timestamp,
                source=host,
                level='info',
                message=f"{process}: {message}",
                raw=line
            )
        
        return None
    
    def find_brute_force(self, entries: list[LogEntry], 
                        threshold: int = 5, window_minutes: int = 5) -> list:
        """Detect brute force attempts"""
        failures_by_ip = Counter()
        
        for entry in entries:
            match = self.patterns['auth_failure'].search(entry.raw)
            if match:
                failures_by_ip[match.group(1)] += 1
        
        return [ip for ip, count in failures_by_ip.items() if count >= threshold]
    
    def find_lateral_movement(self, entries: list[LogEntry]) -> list:
        """Detect potential lateral movement"""
        connections = []
        
        for entry in entries:
            match = self.patterns['ssh_connection'].search(entry.raw)
            if match:
                connections.append({
                    'timestamp': entry.timestamp,
                    'source_ip': match.group(1)
                })
        
        # Look for rapid connections from same IP
        # (potential lateral movement after compromise)
        
        return connections
```

## 🧠 Memory Forensics

### Volatility Integration

```python
import subprocess
import json

class MemoryForensics:
    def __init__(self, memory_dump: str, volatility_path: str = 'vol.py'):
        self.memory_dump = memory_dump
        self.volatility = volatility_path
        self.profile = None
    
    def identify_profile(self) -> str:
        """Identify OS profile from memory dump"""
        result = subprocess.run([
            'python', self.volatility, '-f', self.memory_dump,
            'imageinfo'
        ], capture_output=True, text=True)
        
        # Extract suggested profile
        match = re.search(r'Suggested Profile\(s\): (\S+)', result.stdout)
        if match:
            self.profile = match.group(1).split(',')[0]
        
        return self.profile
    
    def list_processes(self) -> list:
        """List running processes"""
        result = subprocess.run([
            'python', self.volatility, '-f', self.memory_dump,
            '--profile', self.profile, 'pslist'
        ], capture_output=True, text=True)
        
        processes = []
        for line in result.stdout.split('\n')[3:]:  # Skip header
            if line.strip():
                parts = line.split()
                if len(parts) >= 8:
                    processes.append({
                        'offset': parts[0],
                        'name': parts[1],
                        'pid': int(parts[2]),
                        'ppid': int(parts[3]),
                        'threads': int(parts[4]),
                        'handles': int(parts[5]),
                        'session': parts[6],
                        'create_time': ' '.join(parts[7:])
                    })
        
        return processes
    
    def find_malware_indicators(self) -> list:
        """Find potential malware indicators"""
        indicators = []
        
        # Check for suspicious processes
        processes = self.list_processes()
        suspicious_names = ['cmd.exe', 'powershell.exe', 'rundll32.exe', 
                           'svchost.exe', 'lsass.exe']
        
        for proc in processes:
            if proc['name'].lower() in suspicious_names:
                # Check for anomalies
                if proc['ppid'] == -1:  # No parent
                    indicators.append({
                        'type': 'orphan_process',
                        'details': proc
                    })
        
        # Check for hidden processes
        hidden = self._find_hidden_processes()
        if hidden:
            indicators.append({
                'type': 'hidden_processes',
                'details': hidden
            })
        
        return indicators
    
    def extract_strings(self, min_length: int = 8) -> list:
        """Extract strings from memory"""
        result = subprocess.run([
            'strings', '-n', str(min_length), self.memory_dump
        ], capture_output=True, text=True)
        
        return result.stdout.split('\n')
    
    def dump_process_memory(self, pid: int, output_dir: str):
        """Dump memory of specific process"""
        subprocess.run([
            'python', self.volatility, '-f', self.memory_dump,
            '--profile', self.profile, 'procdump',
            '-p', str(pid), '-D', output_dir
        ])
```

## ❓ Interview Questions

### Basic
**Q: What is chain of custody?**
A: Documentation of who handled evidence, when, and what was done. Ensures evidence is admissible in court.

**Q: Why create a disk image instead of examining the original?**
A: Preserve original evidence, allow multiple examiners, prevent accidental modification, maintain integrity.

### Intermediate
**Q: What are important Windows forensic artifacts?**
A: Registry hives, event logs, prefetch files, browser history, MFT, USN journal, recycle bin, LNK files.

**Q: How do you verify evidence integrity?**
A: Calculate and compare cryptographic hashes (MD5, SHA256) at collection and throughout investigation.

### Advanced
**Q: What can memory forensics reveal that disk forensics cannot?**
A: Running processes, loaded DLLs, network connections, encryption keys, passwords in memory, malware in memory (never written to disk).