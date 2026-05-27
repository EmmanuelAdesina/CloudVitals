# CloudVitals Architecture

## Overview

CloudVitals is a developer-first cloud security scanner focused on identifying the highest-risk AWS misconfigurations with minimal noise.

The system is intentionally lightweight:
- Fast scans
- Opinionated checks
- CLI-native UX
- Actionable remediation guidance

---

# High-Level Architecture

```mermaid
flowchart TD

%%{init: {'theme': 'dark', 'themeVariables': { 'primaryColor': '#1e1e2e', 'primaryTextColor': '#cdd6f4', 'primaryBorderColor': '#89b4fa', 'lineColor': '#89b4fa', 'secondaryColor': '#313244', 'tertiaryColor': '#45475a'}}}%%
flowchart TB
    subgraph CLI["CloudVitals CLI (Go)"]
        direction TB
        M[("main.go")]
        R["Runner<br/>goroutines"]
        S["Scorer<br/>0-100"]
        RE["Renderer<br/>terminal / JSON"]
    end

    subgraph REG["Configuration"]
        direction TB
        Y[("checks.yaml<br/>registry")]
    end

    subgraph PROV["Provider Layer"]
        direction TB
        PI["Provider Interface"]
        AWS["AWS Provider<br/>Go wrapper"]
    end

    subgraph EXEC["Check Execution (Python)"]
        direction LR
        C1["s3_public.py"]
        C2["sg_open.py"]
        C3["ebs_encrypt.py"]
        C4["root_mfa.py"]
        C5["cloudtrail.py"]
    end

    subgraph OUT["Output"]
        direction LR
        T["Terminal Table"]
        J["JSON"]
        SAR["SARIF"]
    end

    M -->|"loads"| Y
    Y -->|"[]CheckConfig"| M
    M -->|"dispatches"| R
    R -->|"RunCheck()"| PI
    PI -->|"spawns"| AWS
    AWS -->|"exec.Command"| C1
    AWS -->|"exec.Command"| C2
    AWS -->|"exec.Command"| C3
    AWS -->|"exec.Command"| C4
    AWS -->|"exec.Command"| C5
    C1 -->|"JSON stdout"| AWS
    C2 -->|"JSON stdout"| AWS
    C3 -->|"JSON stdout"| AWS
    C4 -->|"JSON stdout"| AWS
    C5 -->|"JSON stdout"| AWS
    AWS -->|"CheckResult"| R
    R -->|"[]CheckResult"| S
    S -->|"score + results"| RE
    RE --> T
    RE --> J
    RE --> SAR

    style CLI fill:#1e1e2e,stroke:#89b4fa,stroke-width:2px
    style REG fill:#1e1e2e,stroke:#f9e2af,stroke-width:2px
    style PROV fill:#1e1e2e,stroke:#a6e3a1,stroke-width:2px
    style EXEC fill:#1e1e2e,stroke:#fab387,stroke-width:2px
    style OUT fill:#1e1e2e,stroke:#cba6f7,stroke-width:2px
```

---

# Core Components

## CLI Layer
Responsible for:
- command parsing
- scan orchestration
- output rendering
- configuration loading

Built with:
- Go
- Cobra
- Charm/Lipgloss

---

## AWS Integration Layer
Handles:
- AWS authentication
- SDK session creation
- account metadata retrieval

Built with:
- AWS SDK for Go v2

---

## Check Engine
Each security check:
- runs independently
- returns structured findings
- maps to a Zero Trust principle
- includes remediation guidance

Example finding structure:

```json
{
  "id": "s3_public",
  "severity": "CRITICAL",
  "resource": "backup-bucket",
  "risk": "Public data exposure",
  "fix_command": "aws s3api put-public-access-block ..."
}
```

---

# Security Philosophy

CloudVitals intentionally avoids:
- excessive compliance checks
- alert fatigue
- enterprise dashboard complexity

Instead, it focuses on:
- breach-critical misconfigurations
- developer usability
- rapid remediation

---

# Future Architecture Expansion

Planned future capabilities:
- GitHub Actions integration
- Multi-cloud support (GCP/Azure)
- Drift detection
- CI/CD policy enforcement
- JSON and SARIF exports
- Risk correlation engine

---

# Design Principles

- Fast by default
- Minimal dependencies
- Clear findings
- Human-readable remediation
- Zero Trust aligned
