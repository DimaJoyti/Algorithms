# 🖥️ OS Hardening

Secure configuration of operating systems to reduce attack surface.

## 📋 Table of Contents

- [Linux Hardening](#linux-hardening)
- [User Management](#user-management)
- [File Permissions](#file-permissions)
- [Service Hardening](#service-hardening)
- [Logging and Auditing](#logging-and-auditing)
- [Interview Questions](#interview-questions)

## 🐧 Linux Hardening

### Initial Hardening Script

```bash
#!/bin/bash
# Linux Hardening Script

# Update system
apt update && apt upgrade -y

# 1. User and Access Control
# Remove unnecessary users
for user in games news lp sync; do
    userdel $user 2>/dev/null
done

# Set strong password policy
cat > /etc/login.defs << 'EOF'
PASS_MAX_DAYS   90
PASS_MIN_DAYS   7
PASS_WARN_AGE   14
EOF

# 2. SSH Hardening
cat > /etc/ssh/sshd_config << 'EOF'
Port 2222
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
AuthorizedKeysFile .ssh/authorized_keys
MaxAuthTries 3
MaxSessions 2
ClientAliveInterval 300
ClientAliveCountMax 2
AllowUsers admin
X11Forwarding no
AllowTcpForwarding no
PermitTunnel no
EOF

# 3. Disable unnecessary services
for service in telnet ftp rsh rlogin; do
    systemctl disable $service 2>/dev/null
    systemctl stop $service 2>/dev/null
done

# 4. Firewall setup
apt install -y ufw
ufw default deny incoming
ufw default allow outgoing
ufw allow 2222/tcp comment 'SSH'
ufw allow 80/tcp comment 'HTTP'
ufw allow 443/tcp comment 'HTTPS'
ufw --force enable

# 5. Kernel hardening
cat > /etc/sysctl.d/99-security.conf << 'EOF'
# Disable IP forwarding
net.ipv4.ip_forward = 0

# Disable send redirects
net.ipv4.conf.all.send_redirects = 0

# Disable source routing
net.ipv4.conf.all.accept_source_route = 0
net.ipv6.conf.all.accept_source_route = 0

# Enable reverse path filtering
net.ipv4.conf.all.rp_filter = 1

# Ignore ICMP redirects
net.ipv4.conf.all.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0

# Disable ICMP broadcast responses
net.ipv4.icmp_echo_ignore_broadcasts = 1

# Enable bad error message protection
net.ipv4.icmp_ignore_bogus_error_responses = 1

# Enable TCP SYN cookies
net.ipv4.tcp_syncookies = 1

# Disable IPv6 if not needed
net.ipv6.conf.all.disable_ipv6 = 1
EOF

sysctl -p /etc/sysctl.d/99-security.conf

# 6. File permissions
chmod 600 /etc/shadow
chmod 600 /etc/gshadow
chmod 644 /etc/passwd
chmod 644 /etc/group
chmod 700 /root
chmod 600 /etc/ssh/sshd_config

# 7. Install security tools
apt install -y fail2ban aide rkhunter

# Configure fail2ban
cat > /etc/fail2ban/jail.local << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = 2222
EOF

systemctl enable fail2ban
systemctl start fail2ban

echo "Hardening complete!"
```

### Automated Compliance Check

```python
import subprocess
import os
from dataclasses import dataclass
from typing import List

@dataclass
class SecurityCheck:
    name: str
    passed: bool
    message: str
    severity: str

class LinuxHardeningChecker:
    def __init__(self):
        self.results: List[SecurityCheck] = []
    
    def run_all_checks(self) -> List[SecurityCheck]:
        """Run all security checks"""
        self.check_ssh_config()
        self.check_password_policy()
        self.check_file_permissions()
        self.check_services()
        self.check_firewall()
        self.check_users()
        return self.results
    
    def check_ssh_config(self):
        """Check SSH configuration"""
        config_path = "/etc/ssh/sshd_config"
        
        if not os.path.exists(config_path):
            self.results.append(SecurityCheck(
                "SSH Config", False, "SSH config not found", "HIGH"
            ))
            return
        
        with open(config_path, 'r') as f:
            content = f.read()
        
        # Check root login
        if "PermitRootLogin no" in content:
            self.results.append(SecurityCheck(
                "SSH Root Login", True, "Root login disabled", "PASS"
            ))
        else:
            self.results.append(SecurityCheck(
                "SSH Root Login", False, "Root login should be disabled", "HIGH"
            ))
        
        # Check password auth
        if "PasswordAuthentication no" in content:
            self.results.append(SecurityCheck(
                "SSH Password Auth", True, "Password auth disabled", "PASS"
            ))
        else:
            self.results.append(SecurityCheck(
                "SSH Password Auth", False, "Password auth should be disabled", "HIGH"
            ))
    
    def check_password_policy(self):
        """Check password policy"""
        login_defs = "/etc/login.defs"
        
        if os.path.exists(login_defs):
            with open(login_defs, 'r') as f:
                content = f.read()
            
            if "PASS_MAX_DAYS" in content:
                for line in content.split('\n'):
                    if line.startswith('PASS_MAX_DAYS'):
                        days = int(line.split()[1])
                        if days <= 90:
                            self.results.append(SecurityCheck(
                                "Password Expiry", True, 
                                f"Max password age: {days} days", "PASS"
                            ))
                        else:
                            self.results.append(SecurityCheck(
                                "Password Expiry", False,
                                f"Max password age too long: {days} days", "MEDIUM"
                            ))
    
    def check_file_permissions(self):
        """Check critical file permissions"""
        checks = {
            "/etc/shadow": (0o600, "Shadow file"),
            "/etc/passwd": (0o644, "Passwd file"),
            "/etc/ssh/sshd_config": (0o600, "SSH config"),
            "/root": (0o700, "Root directory"),
        }
        
        for path, (expected_perm, name) in checks.items():
            if os.path.exists(path):
                actual_perm = os.stat(path).st_mode & 0o777
                if actual_perm == expected_perm:
                    self.results.append(SecurityCheck(
                        f"Permissions: {name}", True,
                        f"Correct: {oct(actual_perm)}", "PASS"
                    ))
                else:
                    self.results.append(SecurityCheck(
                        f"Permissions: {name}", False,
                        f"Expected {oct(expected_perm)}, got {oct(actual_perm)}", "MEDIUM"
                    ))
    
    def check_services(self):
        """Check running services"""
        dangerous_services = ['telnet', 'ftp', 'rsh', 'rlogin']
        
        for service in dangerous_services:
            result = subprocess.run(
                ['systemctl', 'is-active', service],
                capture_output=True, text=True
            )
            
            if result.returncode == 0:
                self.results.append(SecurityCheck(
                    f"Service: {service}", False,
                    f"Dangerous service {service} is running", "HIGH"
                ))
            else:
                self.results.append(SecurityCheck(
                    f"Service: {service}", True,
                    f"Service {service} not running", "PASS"
                ))
    
    def check_firewall(self):
        """Check firewall status"""
        result = subprocess.run(
            ['ufw', 'status'],
            capture_output=True, text=True
        )
        
        if 'active' in result.stdout.lower():
            self.results.append(SecurityCheck(
                "Firewall", True, "UFW firewall is active", "PASS"
            ))
        else:
            self.results.append(SecurityCheck(
                "Firewall", False, "UFW firewall is not active", "HIGH"
            ))
    
    def check_users(self):
        """Check for security issues with users"""
        # Check for users with empty passwords
        result = subprocess.run(
            ['awk', '-F:', '($2 == "") {print $1}', '/etc/shadow'],
            capture_output=True, text=True
        )
        
        empty_pass_users = result.stdout.strip().split('\n')
        empty_pass_users = [u for u in empty_pass_users if u]
        
        if empty_pass_users:
            self.results.append(SecurityCheck(
                "Empty Passwords", False,
                f"Users with empty passwords: {', '.join(empty_pass_users)}", "HIGH"
            ))
        else:
            self.results.append(SecurityCheck(
                "Empty Passwords", True, "No users with empty passwords", "PASS"
            ))

# Generate report
def generate_report(results: List[SecurityCheck]):
    print("\n" + "="*60)
    print("LINUX HARDENING SECURITY REPORT")
    print("="*60)
    
    passed = sum(1 for r in results if r.passed)
    failed = len(results) - passed
    
    print(f"\nTotal Checks: {len(results)}")
    print(f"Passed: {passed}")
    print(f"Failed: {failed}")
    
    print("\n" + "-"*60)
    print("DETAILED RESULTS:")
    print("-"*60)
    
    for result in results:
        status = "✓ PASS" if result.passed else "✗ FAIL"
        print(f"\n[{status}] {result.name}")
        print(f"  Severity: {result.severity}")
        print(f"  Details: {result.message}")

if __name__ == "__main__":
    checker = LinuxHardeningChecker()
    results = checker.run_all_checks()
    generate_report(results)
```

## 👤 User Management

### Secure User Creation

```bash
#!/bin/bash
# Secure user creation script

create_secure_user() {
    local username=$1
    local shell=$2
    
    # Create user with specific shell
    useradd -m -s $shell $username
    
    # Expire password immediately (force change on first login)
    chage -d 0 $username
    
    # Set password policy
    chage -M 90 -m 7 -W 14 $username
    
    # Lock account until password is set
    passwd -l $username
    
    # Create SSH directory
    mkdir -p /home/$username/.ssh
    chmod 700 /home/$username/.ssh
    chown $username:$username /home/$username/.ssh
    
    echo "User $username created. Unlock with: passwd -u $username"
}

# Add user to sudoers with specific commands only
add_limited_sudo() {
    local username=$1
    
    cat > /etc/sudoers.d/$username << EOF
$username ALL=(ALL) NOPASSWD: /usr/bin/systemctl restart nginx
$username ALL=(ALL) NOPASSWD: /usr/bin/systemctl status *
$username ALL=(ALL) !/usr/bin/su
$username ALL=(ALL) !/usr/bin/bash
EOF
    
    chmod 440 /etc/sudoers.d/$username
}
```

## 📁 File Permissions

### Permission Auditing

```python
import os
import stat
from pathlib import Path
from typing import List, Tuple

class PermissionAuditor:
    def __init__(self):
        self.issues: List[Tuple[str, str]] = []
    
    def audit_suid_sgid(self, root: str = "/") -> List[str]:
        """Find SUID/SGID files"""
        suid_files = []
        
        for path in Path(root).rglob('*'):
            try:
                mode = path.stat().st_mode
                if mode & stat.S_ISUID or mode & stat.S_ISGID:
                    suid_files.append(str(path))
            except (PermissionError, FileNotFoundError):
                pass
        
        return suid_files
    
    def audit_world_writable(self, root: str = "/") -> List[str]:
        """Find world-writable files"""
        writable = []
        
        for path in Path(root).rglob('*'):
            try:
                mode = path.stat().st_mode
                if mode & stat.S_IWOTH:
                    writable.append(str(path))
            except (PermissionError, FileNotFoundError):
                pass
        
        return writable
    
    def audit_unowned_files(self, root: str = "/") -> List[str]:
        """Find files with no owner"""
        unowned = []
        
        for path in Path(root).rglob('*'):
            try:
                stat_info = path.stat()
                try:
                    import pwd
                    pwd.getpwuid(stat_info.st_uid)
                except KeyError:
                    unowned.append(str(path))
            except (PermissionError, FileNotFoundError):
                pass
        
        return unowned

# Critical file permissions
CRITICAL_FILES = {
    "/etc/passwd": {"mode": 0o644, "owner": "root", "group": "root"},
    "/etc/shadow": {"mode": 0o600, "owner": "root", "group": "root"},
    "/etc/gshadow": {"mode": 0o600, "owner": "root", "group": "root"},
    "/etc/group": {"mode": 0o644, "owner": "root", "group": "root"},
    "/etc/sudoers": {"mode": 0o440, "owner": "root", "group": "root"},
    "/etc/ssh/sshd_config": {"mode": 0o600, "owner": "root", "group": "root"},
}
```

## 🔧 Service Hardening

### Systemd Service Security

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Secure Application
After=network.target

[Service]
Type=simple
User=myapp
Group=myapp

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadOnlyPaths=/etc /usr
ReadWritePaths=/var/lib/myapp /var/log/myapp

# Network isolation
PrivateNetwork=false
RestrictAddressFamilies=AF_INET AF_INET6

# Resource limits
LimitNOFILE=1024
MemoryMax=512M
CPUQuota=50%

# Capabilities
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
AmbientCapabilities=

# System call filtering
SystemCallFilter=@system-service
SystemCallArchitectures=native

# Environment
Environment=NODE_ENV=production
ExecStart=/usr/bin/node /opt/myapp/server.js
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 📝 Logging and Auditing

### Auditd Configuration

```bash
# /etc/audit/auditd.conf
local_events = yes
write_logs = yes
log_format = RAW
log_group = root
priority_boost = 4
flush = INCREMENTAL_ASYNC
freq = 50
max_log_file = 8
num_logs = 5
name_format = NONE
max_log_file_action = ROTATE
space_left = 75
space_left_action = SYSLOG
admin_space_left = 50
admin_space_left_action = SUSPEND
disk_full_action = SUSPEND
disk_error_action = SUSPEND

# /etc/audit/rules.d/audit.rules

# Monitor file changes
-w /etc/passwd -p wa -k identity
-w /etc/shadow -p wa -k identity
-w /etc/group -p wa -k identity
-w /etc/sudoers -p wa -k privilege

# Monitor login/logout
-w /var/log/wtmp -p wa -k session
-w /var/run/utmp -p wa -k session
-w /var/log/btmp -p wa -k session

# Monitor SSH
-w /etc/ssh/sshd_config -p wa -k ssh
-w /var/log/auth.log -p wa -k auth

# Monitor time changes
-a always,exit -F arch=b64 -S adjtimex -S settimeofday -k time-change
-a always,exit -F arch=b32 -S adjtimex -S settimeofday -S stime -k time-change
-a always,exit -F arch=b64 -S clock_settime -k time-change

# Monitor user/group changes
-w /usr/sbin/useradd -p x -k user_modification
-w /usr/sbin/usermod -p x -k user_modification
-w /usr/sbin/userdel -p x -k user_modification
-w /usr/sbin/groupadd -p x -k group_modification
-w /usr/sbin/groupmod -p x -k group_modification
-w /usr/sbin/groupdel -p x -k group_modification
```

## ❓ Interview Questions

### Basic
**Q: What's the first thing to do when hardening a Linux server?**
A: Update all packages, change default passwords, disable root SSH login, set up firewall.

**Q: Why disable unused services?**
A: Reduces attack surface, each service is potential vulnerability, saves resources.

### Intermediate
**Q: Explain the principle of least privilege in Linux.**
A: Users and services run with minimum permissions needed. Use non-root users, sudo for specific commands, restrict file permissions.

**Q: How do you secure SSH?**
A: Disable root login, key-based auth only, non-default port, fail2ban, limit users, strong ciphers, session timeout.

### Advanced
**Q: How would you implement a comprehensive logging strategy?**
A: Central log server (syslog/ELK), auditd for file changes, application logs, retention policy, log analysis with alerting, tamper-proof storage.