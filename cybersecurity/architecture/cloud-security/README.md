# ☁️ Cloud Security

Securing cloud infrastructure, platforms, and services.

## 📋 Table of Contents

- [Shared Responsibility Model](#shared-responsibility-model)
- [AWS Security](#aws-security)
- [IAM Best Practices](#iam-best-practices)
- [Data Protection](#data-protection)
- [Security Monitoring](#security-monitoring)
- [Interview Questions](#interview-questions)

## 🔄 Shared Responsibility Model

```
┌─────────────────────────────────────────────────────────┐
│                    RESPONSIBILITY                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │           Customer Responsibility               │   │
│  │  • Data in transit/at rest                      │   │
│  │  • Identity and access management               │   │
│  │  • Security group configuration                 │   │
│  │  • OS patching (EC2)                           │   │
│  │  • Application security                         │   │
│  └─────────────────────────────────────────────────┘   │
│                        │                               │
│  ┌─────────────────────────────────────────────────┐   │
│  │           AWS Responsibility                    │   │
│  │  • Physical security                            │   │
│  │  • Hypervisor                                   │   │
│  │  • Network infrastructure                       │   │
│  │  • Managed services (RDS, S3)                   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘

IaaS (EC2)      → More customer responsibility
PaaS (RDS)      → Shared
SaaS (S3)       → More provider responsibility
```

## 🟠 AWS Security

### Security Hub Configuration

```python
import boto3
from dataclasses import dataclass
from typing import List

@dataclass
class SecurityFinding:
    id: str
    title: str
    severity: str
    description: str
    remediation: str
    resources: List[str]

class AWSSecurityHub:
    def __init__(self, region: str = 'us-east-1'):
        self.client = boto3.client('securityhub', region_name=region)
    
    def get_findings(self, severity: List[str] = None) -> List[SecurityFinding]:
        """Get security findings from Security Hub"""
        filters = {
            'RecordState': [{'Value': 'ACTIVE', 'Comparison': 'EQUALS'}]
        }
        
        if severity:
            filters['SeverityLabel'] = [
                {'Value': s.upper(), 'Comparison': 'EQUALS'} 
                for s in severity
            ]
        
        response = self.client.get_findings(Filters=filters)
        
        findings = []
        for item in response.get('Findings', []):
            findings.append(SecurityFinding(
                id=item['Id'],
                title=item['Title'],
                severity=item['Severity']['Label'],
                description=item['Description'],
                remediation=item.get('Remediation', {}).get('Recommendation', {}).get('Text', ''),
                resources=[r['Id'] for r in item.get('Resources', [])]
            ))
        
        return findings

class EC2SecurityChecker:
    def __init__(self):
        self.ec2 = boto3.client('ec2')
    
    def check_security_groups(self) -> List[dict]:
        """Check for insecure security group rules"""
        issues = []
        
        groups = self.ec2.describe_security_groups()
        
        for sg in groups['SecurityGroups']:
            for rule in sg['IpPermissions']:
                for ip_range in rule.get('IpRanges', []):
                    cidr = ip_range.get('CidrIp', '')
                    
                    # Check for open to world on sensitive ports
                    if cidr == '0.0.0.0/0':
                        port = rule.get('FromPort', 'all')
                        sensitive_ports = [22, 3389, 3306, 5432, 1433, 27017]
                        
                        if port in sensitive_ports or port == 'all':
                            issues.append({
                                'type': 'open_security_group',
                                'group_id': sg['GroupId'],
                                'group_name': sg['GroupName'],
                                'port': port,
                                'severity': 'HIGH',
                                'message': f"Security group {sg['GroupName']} allows access from anywhere on port {port}"
                            })
        
        return issues
    
    def check_unencrypted_volumes(self) -> List[dict]:
        """Check for unencrypted EBS volumes"""
        issues = []
        
        volumes = self.ec2.describe_volumes()
        
        for volume in volumes['Volumes']:
            if not volume['Encrypted']:
                issues.append({
                    'type': 'unencrypted_volume',
                    'volume_id': volume['VolumeId'],
                    'severity': 'MEDIUM',
                    'message': f"Volume {volume['VolumeId']} is not encrypted"
                })
        
        return issues

class S3SecurityChecker:
    def __init__(self):
        self.s3 = boto3.client('s3')
    
    def check_bucket_security(self) -> List[dict]:
        """Check S3 bucket security configurations"""
        issues = []
        
        buckets = self.s3.list_buckets()
        
        for bucket in buckets['Buckets']:
            bucket_name = bucket['Name']
            
            # Check public access block
            try:
                public_access = self.s3.get_public_access_block(Bucket=bucket_name)
                config = public_access['PublicAccessBlockConfiguration']
                
                if not all([
                    config['BlockPublicAcls'],
                    config['BlockPublicPolicy'],
                    config['IgnorePublicAcls'],
                    config['RestrictPublicBuckets']
                ]):
                    issues.append({
                        'type': 'incomplete_public_access_block',
                        'bucket': bucket_name,
                        'severity': 'HIGH',
                        'message': f"Bucket {bucket_name} has incomplete public access block"
                    })
            except:
                issues.append({
                    'type': 'no_public_access_block',
                    'bucket': bucket_name,
                    'severity': 'HIGH',
                    'message': f"Bucket {bucket_name} has no public access block"
                })
            
            # Check encryption
            try:
                encryption = self.s3.get_bucket_encryption(Bucket=bucket_name)
            except:
                issues.append({
                    'type': 'no_encryption',
                    'bucket': bucket_name,
                    'severity': 'MEDIUM',
                    'message': f"Bucket {bucket_name} has no default encryption"
                })
            
            # Check versioning
            versioning = self.s3.get_bucket_versioning(Bucket=bucket_name)
            if versioning.get('Status') != 'Enabled':
                issues.append({
                    'type': 'versioning_disabled',
                    'bucket': bucket_name,
                    'severity': 'LOW',
                    'message': f"Bucket {bucket_name} has versioning disabled"
                })
        
        return issues
```

### CloudFormation Security Template

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Description: Secure infrastructure template

Parameters:
  VpcCidr:
    Type: String
    Default: 10.0.0.0/16

Resources:
  # VPC with flow logs
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: !Ref VpcCidr
      EnableDnsHostnames: true
      EnableDnsSupport: true
      Tags:
        - Key: Name
          Value: SecureVPC

  VPCFlowLog:
    Type: AWS::EC2::FlowLog
    Properties:
      DeliverLogsPermissionArn: !GetAtt FlowLogRole.Arn
      LogGroupName: !Ref FlowLogGroup
      ResourceId: !Ref VPC
      ResourceType: VPC
      TrafficType: ALL

  FlowLogRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: vpc-flow-logs.amazonaws.com
            Action: sts:AssumeRole
      Policies:
        - PolicyName: FlowLogPolicy
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action:
                  - logs:CreateLogGroup
                  - logs:CreateLogStream
                  - logs:PutLogEvents
                Resource: '*'

  # Security Groups
  WebSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Web tier security group
      VpcId: !Ref VPC
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 443
          ToPort: 443
          CidrIp: 0.0.0.0/0
        - IpProtocol: tcp
          FromPort: 80
          ToPort: 80
          CidrIp: 0.0.0.0/0
      SecurityGroupEgress:
        - IpProtocol: tcp
          FromPort: 8080
          ToPort: 8080
          SourceSecurityGroupId: !Ref AppSecurityGroup

  AppSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: App tier security group
      VpcId: !Ref VPC
      SecurityGroupIngress:
        - IpProtocol: tcp
          FromPort: 8080
          ToPort: 8080
          SourceSecurityGroupId: !Ref WebSecurityGroup

  # KMS Key for encryption
  EncryptionKey:
    Type: AWS::KMS::Key
    Properties:
      Description: Encryption key for sensitive data
      KeyPolicy:
        Version: '2012-10-17'
        Statement:
          - Sid: Enable IAM User Permissions
            Effect: Allow
            Principal:
              AWS: !Sub 'arn:aws:iam::${AWS::AccountId}:root'
            Action: 'kms:*'
            Resource: '*'
          - Sid: Allow use of the key
            Effect: Allow
            Principal:
              AWS: !GetAtt AppRole.Arn
            Action:
              - kms:Encrypt
              - kms:Decrypt
              - kms:ReEncrypt*
              - kms:GenerateDataKey*
              - kms:DescribeKey
            Resource: '*'

  # S3 Bucket with security
  SecureBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketEncryption:
        ServerSideEncryptionConfiguration:
          - ServerSideEncryptionByDefault:
              SSEAlgorithm: aws:kms
              KMSMasterKeyID: !GetAtt EncryptionKey.Arn
      PublicAccessBlockConfiguration:
        BlockPublicAcls: true
        BlockPublicPolicy: true
        IgnorePublicAcls: true
        RestrictPublicBuckets: true
      VersioningConfiguration:
        Status: Enabled
      LoggingConfiguration:
        DestinationBucketName: !Ref LogBucket
        LogFilePrefix: access-logs/

  # IAM Role with least privilege
  AppRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: ec2.amazonaws.com
            Action: sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess
      Policies:
        - PolicyName: AppSpecificAccess
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action:
                  - s3:PutObject
                  - s3:GetObject
                Resource: !Sub '${SecureBucket.Arn}/*'
              - Effect: Allow
                Action:
                  - kms:Encrypt
                  - kms:Decrypt
                Resource: !GetAtt EncryptionKey.Arn
```

## 🔐 IAM Best Practices

### Least Privilege Policy Generator

```python
import json
from typing import List, Dict

class IAMPolicyGenerator:
    def __init__(self):
        self.statements = []
    
    def add_read_only(self, resource: str, services: List[str]):
        """Add read-only permissions for services"""
        for service in services:
            self.statements.append({
                "Effect": "Allow",
                "Action": [
                    f"{service}:Get*",
                    f"{service}:List*",
                    f"{service}:Describe*"
                ],
                "Resource": resource
            })
        return self
    
    def add_s3_bucket_access(self, bucket_name: str, read: bool = True, write: bool = False):
        """Add S3 bucket access"""
        actions = []
        if read:
            actions.extend(["s3:GetObject", "s3:ListBucket"])
        if write:
            actions.extend(["s3:PutObject", "s3:DeleteObject"])
        
        self.statements.append({
            "Effect": "Allow",
            "Action": actions,
            "Resource": [
                f"arn:aws:s3:::{bucket_name}",
                f"arn:aws:s3:::{bucket_name}/*"
            ]
        })
        return self
    
    def add_secrets_access(self, secret_names: List[str]):
        """Add Secrets Manager access"""
        resources = [f"arn:aws:secretsmanager:*:*:secret:{name}*" for name in secret_names]
        
        self.statements.append({
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret"
            ],
            "Resource": resources
        })
        return self
    
    def add_condition(self, condition: Dict):
        """Add condition to last statement"""
        if self.statements:
            self.statements[-1]["Condition"] = condition
        return self
    
    def generate(self) -> Dict:
        """Generate IAM policy document"""
        return {
            "Version": "2012-10-17",
            "Statement": self.statements
        }

# Usage
policy = IAMPolicyGenerator() \
    .add_read_only("*", ["ec2", "rds"]) \
    .add_s3_bucket_access("my-app-bucket", read=True, write=True) \
    .add_secrets_access(["database-password", "api-key"]) \
    .add_condition({
        "IpAddress": {
            "aws:SourceIp": ["10.0.0.0/8"]
        }
    }) \
    .generate()

print(json.dumps(policy, indent=2))
```

## 🔒 Data Protection

### Encryption Strategy

```python
import boto3
from cryptography.fernet import Fernet
import base64

class CloudEncryptionManager:
    def __init__(self, kms_key_id: str, region: str = 'us-east-1'):
        self.kms = boto3.client('kms', region_name=region)
        self.kms_key_id = kms_key_id
        self._data_key = None
    
    def generate_data_key(self) -> tuple:
        """Generate data encryption key using KMS"""
        response = self.kms.generate_data_key(
            KeyId=self.kms_key_id,
            KeySpec='AES_256'
        )
        
        plaintext_key = response['Plaintext']
        encrypted_key = response['CiphertextBlob']
        
        return plaintext_key, encrypted_key
    
    def encrypt_data(self, plaintext: bytes) -> dict:
        """Encrypt data using envelope encryption"""
        # Generate data key
        data_key, encrypted_key = self.generate_data_key()
        
        # Use data key for encryption
        fernet = Fernet(base64.urlsafe_b64encode(data_key))
        ciphertext = fernet.encrypt(plaintext)
        
        return {
            'ciphertext': ciphertext,
            'encrypted_data_key': base64.b64encode(encrypted_key).decode(),
            'algorithm': 'AES256'
        }
    
    def decrypt_data(self, ciphertext: bytes, encrypted_key: str) -> bytes:
        """Decrypt data using KMS"""
        # Decrypt data key
        encrypted_key_bytes = base64.b64decode(encrypted_key)
        
        response = self.kms.decrypt(CiphertextBlob=encrypted_key_bytes)
        data_key = response['Plaintext']
        
        # Decrypt data
        fernet = Fernet(base64.urlsafe_b64encode(data_key))
        return fernet.decrypt(ciphertext)

# S3 encryption
class S3EncryptionHelper:
    def __init__(self, bucket: str, kms_key_id: str):
        self.s3 = boto3.client('s3')
        self.bucket = bucket
        self.kms_key_id = kms_key_id
    
    def upload_encrypted(self, key: str, data: bytes):
        """Upload data with server-side encryption"""
        self.s3.put_object(
            Bucket=self.bucket,
            Key=key,
            Body=data,
            ServerSideEncryption='aws:kms',
            SSEKMSKeyId=self.kms_key_id
        )
    
    def download(self, key: str) -> bytes:
        """Download and decrypt data"""
        response = self.s3.get_object(Bucket=self.bucket, Key=key)
        return response['Body'].read()
```

## 📊 Security Monitoring

### CloudTrail and GuardDuty

```python
import boto3
from datetime import datetime, timedelta

class SecurityMonitor:
    def __init__(self):
        self.cloudtrail = boto3.client('cloudtrail')
        self.guardduty = boto3.client('guardduty')
        self.cloudwatch = boto3.client('cloudwatch')
    
    def get_suspicious_activity(self, hours: int = 24) -> List[dict]:
        """Get suspicious CloudTrail events"""
        start_time = datetime.utcnow() - timedelta(hours=hours)
        
        suspicious_events = []
        
        # Look for sensitive API calls
        sensitive_calls = [
            'DeleteBucket',
            'DeleteSecurityGroup', 
            'DetachRolePolicy',
            'DeleteAccessKey',
            'ConsoleLogin',  # Check for failed logins
        ]
        
        for event_name in sensitive_calls:
            response = self.cloudtrail.lookup_events(
                LookupAttributes=[
                    {'AttributeKey': 'EventName', 'AttributeValue': event_name}
                ],
                StartTime=start_time
            )
            
            for event in response.get('Events', []):
                suspicious_events.append({
                    'event_name': event['EventName'],
                    'username': event['Username'],
                    'time': event['EventTime'].isoformat(),
                    'source_ip': event.get('SourceIpAddress', 'unknown'),
                    'resources': [r.get('ResourceName') for r in event.get('Resources', [])]
                })
        
        return suspicious_events
    
    def get_guardduty_findings(self) -> List[dict]:
        """Get GuardDuty findings"""
        detector_ids = self.guardduty.list_detectors().get('DetectorIds', [])
        
        findings = []
        for detector_id in detector_ids:
            response = self.guardduty.list_findings(
                DetectorId=detector_id,
                FindingCriteria={
                    'Criterion': {
                        'service.archived': {'Eq': ['false']},
                        'severity': {'Gte': 4}  # Medium and above
                    }
                }
            )
            
            for finding_id in response.get('FindingIds', []):
                detail = self.guardduty.get_findings(
                    DetectorId=detector_id,
                    FindingIds=[finding_id]
                )
                findings.append(detail['Findings'][0])
        
        return findings
    
    def put_metric_alarm(self, metric_name: str, threshold: int = 10):
        """Create CloudWatch alarm for security metrics"""
        self.cloudwatch.put_metric_alarm(
            AlarmName=f'SecurityAlert-{metric_name}',
            AlarmDescription=f'Alert when {metric_name} exceeds threshold',
            MetricName=metric_name,
            Namespace='Security/Monitoring',
            Statistic='Sum',
            Period=300,
            EvaluationPeriods=1,
            Threshold=threshold,
            ComparisonOperator='GreaterThanThreshold',
            AlarmActions=['arn:aws:sns:us-east-1:123456789:security-alerts']
        )
```

## ❓ Interview Questions

### Basic
**Q: What is the shared responsibility model?**
A: AWS secures the cloud (infrastructure), customer secures what's in the cloud (data, apps, IAM). Division depends on service type.

**Q: Why enable CloudTrail?**
A: Logs all API calls for auditing, compliance, incident response, and security analysis.

### Intermediate
**Q: How do you secure an S3 bucket?**
A: Enable bucket encryption, block public access, enable versioning, use bucket policies, enable logging, use IAM policies for access control.

**Q: Explain least privilege in IAM.**
A: Grant only permissions needed for task. Use resource-level permissions, conditions, and regularly review permissions with Access Analyzer.

### Advanced
**Q: How would you architect a secure multi-account AWS environment?**
A: (1) AWS Organizations with SCPs, (2) Separate accounts for security, logging, workloads, (3) Central IAM with SSO, (4) CloudTrail in all accounts to central bucket, (5) Security Hub aggregation, (6) VPC peering/Transit Gateway with strict routing.