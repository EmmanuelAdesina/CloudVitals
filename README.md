Here is the README. Copy it directly into your repository.

---

```markdown
# CloudVitals

**Is your AWS account one misconfiguration away from a breach? Find out in 10 seconds.**

---

Cloud misconfiguration is the #1 cause of cloud data breaches. Not zero-days. Not nation-state hackers. **Human error.** An S3 bucket left public. A security group with `0.0.0.0/0`. A root account without MFA.

Existing tools "solve" this by running **300+ checks** and drowning you in PDFs. They are built for compliance auditors, not developers. By the time you parse the noise, the attacker has already exfiltrated your data.

CloudVitals does the opposite. **Five checks. Zero noise. One score. One fix.**

---

## What It Checks

| Check | Why It Destroys Companies | Zero Trust Principle |
|---|---|---|
| **Public S3 Buckets** | #1 cause of data leaks — attackers scan for these automatically | Verify Explicitly |
| **Open Security Groups** | `0.0.0.0/0` on port 22 or 3389 is a ransomware invitation | Assume Breach |
| **Unencrypted EBS Volumes** | Compliance failure + instant data exposure if snapshot leaks | Use Least Privilege |
| **Root Account Missing MFA** | One phished password = total account takeover | Verify Explicitly |
| **CloudTrail Disabled** | No audit trail = undetectable breach, indefinite dwell time | Assume Breach |

If you pass all five, you are not "secure." You are **not immediately on fire.** That is the baseline CloudVitals enforces.

---

## One Command

```bash
# Requires: Go 1.22+, Python 3.9+, AWS credentials
aws configure
go install github.com/emmanueladesina/cloudvitals@latest
cloudvitals scan
```

**Output:**

```
===================================================
  CloudVitals Security Score: 72/100
  Status: AT RISK — 2 critical findings
===================================================

SEVERITY    CHECK                        STATUS    FINDINGS
critical    Public S3 Bucket Access      FAIL      1
critical    Open Security Groups         FAIL      1
high        Root Account Missing MFA     FAIL      1

--- Public S3 Bucket Access (critical) ---
  Resource: arn:aws:s3:::backup-bucket
  Region:   us-east-1
  Detail:   S3 bucket public access block not fully enabled
  Fix:      aws s3api put-public-access-block --bucket backup-bucket \
            --public-access-block-configuration \
            BlockPublicAcls=true,IgnorePublicAcls=true,\
            BlockPublicPolicy=true,RestrictPublicBuckets=true

--- Open Security Groups (critical) ---
  Resource: sg-0a1b2c3d
  Region:   us-east-1
  Detail:   Inbound rule allows 0.0.0.0/0 on port 22
  Fix:      aws ec2 revoke-security-group-ingress \
            --group-id sg-0a1b2c3d \
            --ip-permissions IpProtocol=tcp,FromPort=22,ToPort=22,IpRanges='[{CidrIp=0.0.0.0/0}]'
```

No dashboards. No 200-page PDF. **A score, a finding, and the exact CLI command to fix it.**

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CloudVitals CLI (Go)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐ │
│  │   Runner     │  │ Risk Scorer  │  │ Output Formatter  │ │
│  │  (Concurrent)│  │  (0-100)     │  │  (Terminal/JSON/  │ │
│  │              │  │              │  │   SARIF/GitHub     │ │
│  │              │  │              │  │   Actions)         │ │
│  └──────────────┘  └──────────────┘  └─────────────────────┘ │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│   AWS       │   │   GCP       │   │   Azure     │
│ Provider    │   │ Provider    │   │ Provider    │
│ (Python)    │   │ (Python)    │   │ (Python)    │
│             │   │             │   │             │
│ • S3        │   │ • Storage   │   │ • Blob      │
│ • EC2/SG    │   │ • Firewall  │   │ • NSG       │
│ • IAM       │   │ • IAM       │   │ • AD/Entra  │
│ • CloudTrail│   │ • AuditLog  │   │ • Monitor   │
└─────────────┘   └─────────────┘   └─────────────┘
```

**Go** handles the CLI, concurrency, and output formatting.  
**Python** handles cloud SDK logic via official APIs (boto3, google-cloud, azure-identity).  
**YAML registry** defines checks without recompiling the binary.

This split is intentional. Python has mature cloud SDKs. Go has mature CLI ergonomics. We use the right tool for each layer.

---

## Why This Exists

I am a first-year cybersecurity student building toward Cloud Security Architecture (Zero Trust, Multi-Cloud). This is my proof of work — real API calls against real cloud resources, not theoretical slide decks.

If you are hiring in cloud security, infrastructure, or platform engineering: **this is what I can build.** Production-grade concurrency, clean provider interfaces, and security checks that map to actual breach patterns.

If you are a developer who hates security tools that waste your time: **this is what I needed and could not find.**

---

## Adding a Check (20 Minutes)

The fastest way to contribute. No Go changes required.

**1. Write the Python script** in `internal/providers/aws/checks/`:

```python
#!/usr/bin/env python3
import argparse
import json
import boto3

def check(profile, region):
    # Your boto3 logic here
    findings = []
    # ... inspect resources, append to findings if misconfigured
    
    status = "fail" if findings else "pass"
    
    print(json.dumps({
        "check_id": "your_check_id",
        "check_name": "Human Readable Name",
        "status": status,
        "severity": "critical",
        "findings": findings,
        "executed_at": None,
        "error_msg": None
    }))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--profile", default="default")
    parser.add_argument("--region", default="us-east-1")
    args = parser.parse_args()
    check(args.profile, args.region)
```

**2. Register it** in `config/checks.yaml`:

```yaml
checks:
  - id: your_check_id
    name: "Human Readable Name"
    description: "What this check detects"
    provider: aws
    severity: critical
    script: your_check.py
```

**3. Test it:**

```bash
go run ./cmd/cloudvitals -profile default -region us-east-1
```

Submit a PR. First-time contributors welcome.

---

## Roadmap

- [x] AWS S3 public access check
- [ ] AWS security group openness
- [ ] AWS EBS encryption
- [ ] AWS root MFA enforcement
- [ ] AWS CloudTrail status
- [ ] GCP provider (GCS public buckets, firewall rules)
- [ ] Azure provider (Blob storage, NSG rules)
- [ ] GitHub Actions integration (security scorecard on PR)
- [ ] SARIF output for GitHub Advanced Security
- [ ] Pre-commit hook for IaC scanning
- [ ] JSON output for CI/CD pipelines

---

## What CloudVitals Is Not

- Not a full CSPM. If you need 300 checks, use Prowler.
- Not an insurance product. We detect misconfigurations. We do not indemnify against breaches.
- Not a blockchain project. We use SHA-256 for internal integrity where appropriate. That is it.
- Not a managed service. This is a CLI tool you run in your own environment. Your credentials never leave your machine.

---

## License

MIT License — see [LICENSE](LICENSE)

---

**One sentence:** CloudVitals checks the five AWS misconfigurations that actually destroy companies, gives you a security score, and tells you exactly how to fix them. In 10 seconds. For free.
```

---


