# Technology Decisions & Tradeoffs

Architectural decisions for LiveCode authentication system.

---

## Logging Stack

### Selected: Zap (Uber)

**Alternatives:**

| Library | Performance | Zero-Alloc | Sampling | Hooks | Maintenance |
|---------|-------------|------------|----------|-------|-------------|
| **Zap** | 300 ns/op | Yes | Yes | Yes | Uber (active) |
| slog | 500 ns/op | No | No | No | Go core team |
| Logrus | 3000 ns/op | No | No | Yes | Community (slow) |
| Zerolog | 100 ns/op | Yes | No | No | Community (active) |

**Why Zap:**
- 60% faster than slog
- Sampling reduces log volume at high traffic
- Hooks for error alerting (Sentry integration)
- Battle-tested at Uber scale (millions req/s)
- Zero allocation (reduces GC pressure)

**Tradeoff:** More boilerplate than slog, but gain production features.

---

## Timestamp Format

### Selected: ISO8601 (`2025-12-11T03:00:00Z`)

**Alternatives:**

| Format | Example | Human Readable | Size | Parseable |
|--------|---------|----------------|------|-----------|
| **ISO8601** | `2025-12-11T03:00:00Z` | Yes | 20 bytes | Yes |
| RFC3339 | `2025-12-11T03:00:00+02:00` | Yes | 25 bytes | Yes |
| Unix Epoch | `1702260000` | No | 10 bytes | Yes |

**Why ISO8601:**
- Industry standard (AWS CloudWatch, Grafana Loki, Datadog)
- Human readable + UTC explicit
- Sortable (lexicographic = chronological)

**Tradeoff:** 10 bytes larger than Unix Epoch, but human debuggable.

---

## Log Storage

### Selected: Loki + Grafana (Production)

**Alternatives:**

| Solution | Cost | Setup | Scalability | Vendor Lock |
|----------|------|-------|-------------|-------------|
| **Loki** | Free (self-host) | Medium | Yes | No |
| CloudWatch | $0.50/GB | Easy | Yes | AWS only |
| Datadog | $1.27/GB | Easy | Yes | Datadog only |
| ELK Stack | Free/$$$ | Hard | Yes | No |

**Why Loki:**
- Open-source (no vendor lock-in)
- Kubernetes-native (Netflix, Grafana Labs)
- Unified observability (Loki logs + Prometheus metrics in Grafana)
- LogQL query language (similar to PromQL)
- Cost-effective (free self-hosted, cheap on Grafana Cloud)

**Deployment:**
- Development: stdout (console)
- Production: Loki + Grafana (self-hosted VPS or Grafana Cloud)

---

## Infrastructure

### VPS vs Bare Metal vs Cloud

**VPS (Virtual Private Server):**
- Example: DigitalOcean Droplet ($12/month, 2GB RAM, 1 CPU)
- Use: 90% of web applications
- Pros: Affordable, fast provisioning, easy upgrades
- Cons: Shared hardware ("noisy neighbor"), virtualization overhead (~5-10%)

**Bare Metal:**
- Example: Dedicated Dell R740 ($100-500/month, 64GB RAM, 12 cores)
- Use: High-performance databases, compliance (HIPAA, PCI-DSS)
- Pros: Maximum performance, dedicated resources, physical isolation
- Cons: Expensive, no auto-scaling, slow provisioning (hours to days)

**Cloud (AWS/GCP/Azure):**
- Example: AWS EC2 + RDS + CloudWatch
- Use: Enterprise auto-scaling, global deployment
- Pros: Managed services, auto-scaling, global presence (20+ regions)
- Cons: Complex pricing, vendor lock-in, expensive vs VPS

**Our Choice:**
- Development: Local machine
- Production Phase 1: VPS (DigitalOcean $12/month)
- Production Scale: Cloud (AWS/GCP when > 1000 concurrent users)

---

## Zap Configuration

### Production vs Development

**Production Config:**
```go
zap.NewProductionConfig()
```
- JSON output (machine-parseable for Loki)
- INFO level and above (no DEBUG)
- Stack traces only on ERROR
- Sampling enabled (reduces volume at high traffic)

**Development Config:**
```go
zap.NewDevelopmentConfig()
```
- Console output (colored, human-readable)
- DEBUG level (verbose logging)
- Stack traces on WARN and above
- No sampling

**Our Choice:** Production config from day 1 (JSON format ready for Loki, easy to switch later).

---

## Log Levels

| Level | Use Case | Example |
|-------|----------|---------|
| DEBUG | Development verbose info | Query parameters, intermediate values |
| INFO | Normal operations | User registered, request completed |
| WARN | Unexpected but recoverable | Rate limit hit, repeated validation fail |
| ERROR | Functionality broken | DB query failed, JWT generation failed |
| FATAL | Unrecoverable (app crash) | DB connection failed at startup |

**Production:** INFO and above (no DEBUG to reduce noise and improve performance).

---

## Caller Info & Stack Traces

**Configuration:**
```go
zap.AddCaller()                        // Always (cheap: ~50ns, very useful)
zap.AddStacktrace(zapcore.ErrorLevel)  // Only on ERROR (expensive, verbose)
```

**Why:**
- Caller adds `"caller": "file.go:42"` to every log (know exact line that logged)
- Stack traces only on ERROR (verbose output, only when needed for debugging)

**Output Example:**
```json
{
  "level": "error",
  "timestamp": "2025-12-11T03:00:00Z",
  "caller": "handlers/register.go:46",
  "msg": "database_query_failed",
  "operation": "email_check",
  "error": "pq: connection timeout",
  "stacktrace": "main.RegisterUserInternal\n\t/path/handlers/register.go:46\n..."
}
```

---

## Correlation IDs

**What:** Unique UUID attached to every request, included in all logs for that request.

**Why:**
- Trace full request lifecycle across layers (middleware → handler → database)
- Debug issues: "Show all logs for correlation_id=abc-123"
- Essential for distributed systems

**Implementation:**
```go
correlationID := uuid.New().String()
c.Set("correlation_id", correlationID)
```

**Example Usage in Grafana:**
```
{correlation_id="abc-123"}
```
Shows ALL logs (incoming request, DB query, errors, response) for that specific request.

---

## Summary

**Key Decisions:**
1. **Logging:** Zap (performance + production features over simplicity)
2. **Timestamps:** ISO8601 (readability + industry standard over size)
3. **Storage:** Loki/Grafana (open-source + scalability over vendor SaaS)
4. **Deployment:** VPS initially (cost + simplicity over cloud complexity)

**Principle:** Choose tools that scale from 10 req/s to 10,000 req/s without forced rewrites.
