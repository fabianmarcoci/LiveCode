# Technology Decisions & Tradeoffs

This document captures architectural decisions, alternatives considered, and tradeoffs for the LiveCode project.

---

## Logging Stack

### Decision: Zap (Uber) for Structured Logging

**Alternatives Considered:**

1. **slog** (Go 1.21+ built-in)
2. **Zap** (Uber library) ← **SELECTED**
3. **Logrus** (older, still popular)
4. **Zerolog** (fastest, minimal allocations)

**Comparison:**

| Feature               | slog         | Zap              | Logrus             | Zerolog              |
| --------------------- | ------------ | ---------------- | ------------------ | -------------------- |
| **Performance**       | ~500 ns/op   | ~300 ns/op       | ~3000 ns/op        | ~100 ns/op           |
| **Zero allocation**   | No           | Yes              | No                 | Yes                  |
| **JSON output**       | Built-in     | Built-in         | Built-in           | Built-in             |
| **Structured fields** | Yes          | Yes              | Yes                | Yes                  |
| **Sampling**          | No           | Yes              | No                 | No                   |
| **Hooks**             | No           | Yes              | Yes                | No                   |
| **Context support**   | Native       | With()           | WithFields()       | With()               |
| **Production use**    | New (2023+)  | Uber, Cloudflare | Kubernetes, Docker | RabbitMQ, Mattermost |
| **Maintenance**       | Go core team | Uber (active)    | Community (slower) | Community (active)   |
| **Learning curve**    | Easy         | Medium           | Easy               | Medium               |

**Why Zap:**

- ✅ **Performance:** 60% faster than slog, good enough (zerolog's 200ns gain irrelevant at our scale)
- ✅ **Sampling:** Can reduce log volume at high traffic (log 1/100 requests when > 1000 req/s)
- ✅ **Hooks:** Can send ERROR logs to Sentry/PagerDuty automatically
- ✅ **Production-proven:** Battle-tested at Uber scale (millions of req/s)
- ✅ **Zero allocation:** Reduces GC pressure (important for high-throughput APIs)

**Why NOT slog:**

- ❌ Newer (less battle-tested in production)
- ❌ No sampling (would need manual implementation)
- ❌ No hooks (can't auto-alert on errors)

**Why NOT Logrus:**

- ❌ 10x slower than Zap
- ❌ Allocates on every log call (GC pressure)
- ❌ Maintenance slowing down

**Why NOT Zerolog:**

- ❌ Harder to read (chaining API: `log.Info().Str("key", "val").Msg("...")`)
- ❌ Less features (no sampling, no hooks)
- ✅ Fastest (but overkill for our 10-100 req/s scale)

**Tradeoff Accepted:**

- Zap has more boilerplate than slog (`zap.String("key", "val")` vs `slog.String("key", "val")`)
- But gain: production features (sampling, hooks) + proven reliability

---

## Timestamp Format

### Decision: ISO8601 (`2025-12-11T03:00:00Z`)

**Alternatives Considered:**

1. **ISO8601** ← **SELECTED**
2. **RFC3339** (`2025-12-11T03:00:00+02:00`)
3. **Unix Epoch** (`1702260000`)
4. **Custom** (`2025-12-11 03:00:00`)

**Comparison:**

| Format     | Example                     | Timezone Info  | Human Readable | Machine Parseable | Size (bytes) |
| ---------- | --------------------------- | -------------- | -------------- | ----------------- | ------------ |
| ISO8601    | `2025-12-11T03:00:00Z`      | UTC only       | ✅ Yes         | ✅ Yes            | 20           |
| RFC3339    | `2025-12-11T03:00:00+02:00` | ✅ Full offset | ✅ Yes         | ✅ Yes            | 25           |
| Unix Epoch | `1702260000`                | Implicit UTC   | ❌ No          | ✅ Yes            | 10           |
| Custom     | `2025-12-11 03:00:00`       | ❌ Ambiguous   | ✅ Yes         | ⚠️ Fragile        | 19           |

**Why ISO8601:**

- ✅ **Industry standard** (used by AWS CloudWatch, Grafana Loki, Datadog)
- ✅ **Human readable** (developers can read logs without conversion)
- ✅ **UTC explicit** (`Z` suffix = Zulu time = UTC)
- ✅ **Sortable** (lexicographic sort = chronological sort)
- ✅ **Compact** (shorter than RFC3339)

**Why NOT RFC3339:**

- ❌ Longer (includes timezone offset like `+02:00`)
- ❌ Unnecessary (we always log in UTC, no need for offset)

**Why NOT Unix Epoch:**

- ❌ Not human readable (need calculator to convert)
- ❌ Harder debugging (can't scan logs visually)
- ✅ Smallest size (good for high-volume logs, but we're not there yet)

**Why NOT Custom:**

- ❌ Ambiguous timezone (is it local time? UTC? server time?)
- ❌ Parsing fragility (different languages parse differently)

**Tradeoff Accepted:**

- 10 bytes larger than Unix Epoch
- But gain: human readability + no ambiguity

---

## Log Storage

### Decision: Loki + Grafana (for production)

**Alternatives Considered:**

1. **stdout (console)** ← **Development**
2. **File (`logs/app.log`)** ← **VPS/Bare Metal**
3. **Loki + Grafana** ← **SELECTED for Production**
4. **CloudWatch** (AWS)
5. **Datadog** (SaaS)
6. **ELK Stack** (Elasticsearch, Logstash, Kibana)

**Comparison:**

| Solution       | Cost             | Setup Complexity | Query Speed  | Scalability | Vendor Lock-in | Alerting |
| -------------- | ---------------- | ---------------- | ------------ | ----------- | -------------- | -------- |
| **stdout**     | Free             | Trivial          | N/A          | ❌ No       | No             | No       |
| **File**       | Free             | Easy             | Slow (grep)  | ❌ No       | No             | No       |
| **Loki**       | Free (self-host) | Medium           | Fast (LogQL) | ✅ Yes      | No             | ✅ Yes   |
| **CloudWatch** | $$$ ($0.50/GB)   | Easy (AWS only)  | Fast         | ✅ Yes      | ✅ AWS         | ✅ Yes   |
| **Datadog**    | $$$$ ($1.27/GB)  | Easy             | Fast         | ✅ Yes      | ✅ Datadog     | ✅ Yes   |
| **ELK**        | Free/$$$         | Hard             | Very Fast    | ✅ Yes      | No             | ✅ Yes   |

**Why Loki + Grafana:**

- ✅ **Open-source** (no vendor lock-in, can self-host anywhere)
- ✅ **Cost-effective** (free if self-hosted, cheap if using Grafana Cloud)
- ✅ **Industry standard** (Netflix, Grafana Labs, CNCF project)
- ✅ **Kubernetes-native** (integrates seamlessly with K8s deployments)
- ✅ **Unified observability** (Loki logs + Prometheus metrics in same Grafana dashboard)
- ✅ **LogQL query language** (like PromQL, easy to learn)

**Why NOT CloudWatch:**

- ❌ AWS-only (can't use on GCP/Azure/DigitalOcean)
- ❌ Expensive at scale ($0.50/GB ingested + $0.03/GB stored)
- ❌ Vendor lock-in (logs stuck in AWS)

**Why NOT Datadog:**

- ❌ Most expensive ($1.27/GB + per-host fees)
- ❌ Vendor lock-in (can't export logs easily)
- ✅ Best-in-class UI and features (but overkill for us)

**Why NOT ELK Stack:**

- ❌ Complex setup (3 components: Elasticsearch, Logstash, Kibana)
- ❌ Heavy resource usage (Elasticsearch needs 4-8GB RAM minimum)
- ❌ Older architecture (pre-cloud-native era)
- ✅ Very powerful search (full-text indexing), but we don't need it

**Tradeoff Accepted:**

- Self-hosting requires managing Loki/Grafana infrastructure
- But gain: no vendor lock-in + no per-GB costs

**Deployment Strategy:**

- **Development:** stdout (console logs, easy to see)
- **Production (VPS):** Loki + Grafana (self-hosted)
- **Production (Cloud):** Grafana Cloud (managed Loki, $50/month for 100GB)

---

## VPS vs Bare Metal vs Cloud

**Definitions:**

### Bare Metal

- **What:** Physical server hardware dedicated to you
- **Example:** Dell PowerEdge R740 in a datacenter (OVH, Hetzner)
- **Specs:** 64GB RAM, 12 CPU cores, 1TB SSD
- **Cost:** $100-500/month
- **Use case:** High-performance databases, compliance (HIPAA, PCI-DSS), no "noisy neighbors"

**Pros:**

- ✅ Maximum performance (no virtualization overhead)
- ✅ Dedicated resources (CPU, RAM, network)
- ✅ Compliance-friendly (physical isolation)

**Cons:**

- ❌ Expensive (paying for full server even if using 20%)
- ❌ No auto-scaling (can't add RAM on-demand)
- ❌ Long provisioning (hours to days)

---

### VPS (Virtual Private Server)

- **What:** Virtual machine on shared physical hardware
- **Example:** DigitalOcean Droplet ($12/month for 2GB RAM, 1 CPU)
- **Provider:** DigitalOcean, Linode, Vultr, Hetzner Cloud
- **Use case:** 90% of web applications

**Pros:**

- ✅ Affordable ($5-50/month for small-medium apps)
- ✅ Fast provisioning (minutes)
- ✅ Easy to upgrade (resize VM in UI)
- ✅ Predictable pricing (fixed monthly cost)

**Cons:**

- ❌ Shared hardware ("noisy neighbor" problem if another VM hammers CPU)
- ❌ Virtualization overhead (~5-10% performance loss)
- ❌ Manual scaling (need to resize or add VMs manually)

---

### Cloud (AWS/GCP/Azure)

- **What:** Managed services on top of VPS infrastructure
- **Example:** AWS EC2 + RDS + S3 + CloudWatch
- **Use case:** Enterprise applications, auto-scaling needs

**Pros:**

- ✅ Auto-scaling (scale to 1000 VMs automatically)
- ✅ Managed services (RDS = managed PostgreSQL, no DB admin needed)
- ✅ Global presence (deploy in 20+ regions worldwide)
- ✅ Pay-per-use (only pay for what you use)

**Cons:**

- ❌ Complex pricing (hundreds of line items, surprise bills)
- ❌ Vendor lock-in (hard to migrate from AWS to GCP)
- ❌ Overkill for small apps (VPS is simpler)
- ❌ More expensive than VPS (AWS EC2 t3.small = $15/month vs DigitalOcean $12/month with better specs)

---

**Our Choice:**

- **Development:** Local machine (free)
- **Production (Phase 1):** VPS (DigitalOcean $12/month)
- **Production (Scale):** Cloud (AWS/GCP when > 1000 users)

---

## Production Config vs Development Config

### Zap Config Modes

**Development:**
