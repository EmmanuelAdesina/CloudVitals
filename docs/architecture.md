# CloudVitals Architecture

**Concurrent security scanner. Go orchestrates. Python inspects. YAML configures.**

---

## Layers
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Terminal   │  │    JSON     │  │    SARIF    │
│  Renderer   │  │  Renderer   │  │  Renderer   │
└─────────────┘  └─────────────┘  └─────────────┘
└────────────────┬─────────────────┘
▼
┌─────────────────────────────────────────────┐
│  Runner │ Scorer │ Registry (checks.yaml)   │
│  (goroutines)     (0-100)    (no recompile) │
└─────────────────────────────────────────────┘
▼
┌─────────────────────────────────────────────┐
│         Provider Interface (Go)               │
│  RunCheck() → spawns Python → JSON back     │
└─────────────────────────────────────────────┘
▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│    AWS      │  │    GCP      │  │   Azure     │
│  (Python)   │  │  (Python)   │  │  (Python)   │
│ s3_public   │  │ storage_pub │  │  blob_pub   │
│ sg_open     │  │ firewall    │  │   nsg       │
│ ebs_encrypt │  │  ...        │  │   ...       │
└─────────────┘  └─────────────┘  └─────────────┘
plain
Copy

---

## Data Flow

| Step | What Happens | Time |
|---|---|---|
| 1. Load | `checks.yaml` → `[]CheckConfig` | 1ms |
| 2. Run | 5 goroutines spawn 5 Python scripts | ~3s |
| 3. Score | Severity weights deducted from 100 | 1ms |
| 4. Render | Terminal table + exact fix commands | 1ms |

---

## Why Go + Python?

| Layer | Language | Reason |
|---|---|---|
| CLI / Concurrency | Go | Single binary, goroutines, fast |
| Cloud APIs | Python | boto3, google-cloud, azure-identity are mature |
| Config | YAML | Human-editable, no recompile |

---

## Add a Check in 20 Minutes

1. Write `internal/providers/aws/checks/YOUR_CHECK.py`
2. Add one entry to `config/checks.yaml`
3. **Zero Go changes.** Registry discovers it at runtime.

---

## Add a Provider in 1 Day

1. Create `internal/providers/gcp/gcp.go` — implement `Provider`
2. Add GCP checks in `internal/providers/gcp/checks/`
3. `runner`, `scorer`, `renderer` **do not change.**

---

## Security

- **Local only.** No SaaS. Credentials never leave your machine.
- **Read-only checks.** No `eval`, no `exec`, no shell injection.
- **Binary pass/fail.** No heuristic noise. `BlockPublicPolicy=true` is either set or not.

---

## One Sentence

Go runs Python scripts in parallel, scores the results, and renders a security score with exact fix commands.
