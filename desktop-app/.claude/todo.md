# LiveCode Authentication - Professional Implementation Plan

## FASE 1: FUNDAȚIE (Current Focus)

### 1.1 API Versionare ✅

- [x] Adaugă `/api/v1` prefix în main.go
- [x] Update config.rs cu auth_url() și api_url()
- [x] Test versionare cu curl

### 1.2 Refresh Token Logic

- [ ] Backend: Creează endpoint POST /api/v1/auth/refresh
- [ ] Backend: Implementează verificare refresh token
- [ ] Backend: Generare nou access token
- [ ] Rust: Adaugă refresh_token command în auth.rs
- [ ] Frontend: Auto-refresh logic (retry request la 401)
- [ ] Test: Token expirat → refresh → success

### 1.3 Role-Based Access Control (RBAC)

- [ ] Database: Adaugă coloană `role` în users table (DEFAULT 'user')
- [ ] Backend: Include `role` în JWT claims la login
- [ ] Backend: Creează middleware/rbac.go cu RequireRole()
- [ ] Backend: Protejează /health cu RequireRole("admin")
- [ ] Test: User normal → 403 Forbidden
- [ ] Test: Admin user → 200 OK

### 1.4 Protected Routes (Exemplu: Submit Code)

- [ ] Backend: Creează routes/submission.go
- [ ] Backend: POST /api/v1/submissions (create)
- [ ] Backend: GET /api/v1/submissions/:id (read)
- [ ] Database: Creează tabel `submissions`
- [ ] Rust: Adaugă submit_code command
- [ ] Frontend: Creează SubmissionForm.tsx
- [ ] Test: Submit code cu token valid

## FAZA 2: SECURITATE & BEST PRACTICES

### 2.1 Rate Limiting

- [ ] Backend: Adaugă middleware rate limiter
- [ ] Config: Max 100 requests/minute per IP
- [ ] Test: Depășire limit → 429 Too Many Requests

### 2.2 Request Validation

- [ ] Backend: Validare max length pentru toate inputs
- [ ] Backend: Sanitizare input (prevent XSS)
- [ ] Backend: Validare email format
- [ ] Backend: Validare password strength

### 2.3 Error Handling

- [ ] Backend: Centralizare error responses
- [ ] Backend: NU expune internal errors în production
- [ ] Backend: Logging pentru debugging
- [ ] Frontend: User-friendly error messages

### 2.4 CORS Configuration

- [ ] Backend: Configurează CORS middleware
- [ ] Whitelist doar origin-uri known
- [ ] Test: Cross-origin requests

## FAZA 3: SCALABILITATE & DEPLOYMENT

### 3.1 Dockerization

- [ ] Creează backend-api/Dockerfile
- [ ] Creează docker-compose.yml
- [ ] Multi-stage build pentru size optimization
- [ ] Test: docker-compose up

### 3.2 Nginx Load Balancer

- [ ] Creează nginx.conf
- [ ] Configurează round-robin între 3 instanțe API
- [ ] Health checks (ping /health)
- [ ] Test: Load distribution

### 3.3 Monitoring

- [ ] Prometheus metrics endpoint
- [ ] Grafana dashboards
- [ ] Alert rules (high error rate, latency)

### 3.4 Database Optimizations

- [ ] Index pe email și username
- [ ] Connection pooling tuning
- [ ] Query performance analysis
- [ ] Migrations versioning

## FAZA 4: ADVANCED FEATURES

### 4.1 Token Refresh Strategy (Hybrid)

- [ ] Database: refresh_tokens table cu jti
- [ ] Backend: Whitelist/blacklist pentru revocation
- [ ] Redis cache pentru fast lookup
- [ ] Automatic cleanup expired tokens

### 4.2 Multi-Device Support (Enterprise Session Management)

**Database Schema:**
- [ ] CREATE TABLE sessions (id, user_id, refresh_token_hash, device_name, ip_address, user_agent, last_active, expires_at)
- [ ] Add indexes on user_id and expires_at

**Backend Implementation:**
- [ ] Salvează refresh token în sessions (hashed cu SHA-256)
- [ ] La login: INSERT session în DB, returnează session_id
- [ ] La refresh: Verifică session în DB, generează access token nou
- [ ] GET /api/v1/sessions: Lista toate dispozitivele utilizatorului
- [ ] DELETE /api/v1/sessions/:id: Logout de pe un dispozitiv specific
- [ ] DELETE /api/v1/sessions/all: Logout all devices (revoke toate sessions)

**Frontend/Rust:**
- [ ] Salvează session_id în Tauri storage (în loc de refresh_token)
- [ ] UI: Account Settings → Active Devices
- [ ] UI: Afișează device_name, last_active, current device indicator
- [ ] UI: Buton "Logout" per device
- [ ] UI: Buton "Logout all other devices"

**Security Benefits:**
- ✅ Revocație instantanee (ștergi session din DB → invalid imediat)
- ✅ User vede toate dispozitivele logat (transparency)
- ✅ Refresh token NU e pe client (doar session_id)
- ✅ Audit trail complet (când, unde, ce device)

### 4.3 OAuth Integration

- [ ] Google OAuth provider
- [ ] GitHub OAuth provider
- [ ] Link multiple providers per user

### 4.4 Audit Logging

- [ ] Log toate login attempts
- [ ] Log toate changes la resources
- [ ] Admin dashboard pentru audit trail

## BEST PRACTICES CHECKLIST

### Security

- [x] JWT signatures verified
- [x] Passwords hashed cu Argon2
- [x] Tokens stored encrypted (Tauri secure storage)
- [ ] Rate limiting implemented
- [ ] HTTPS enforced (production)
- [ ] CORS configured
- [ ] SQL injection prevented (parameterized queries)
- [ ] XSS prevented (input sanitization)

### Code Quality

- [x] API versionat (/api/v1)
- [x] Middleware pattern folosit
- [x] Error handling consistent
- [ ] Unit tests (>80% coverage)
- [ ] Integration tests
- [ ] Code documentation
- [ ] Git commit messages descriptive

### System Design

- [x] Stateless API (JWT-based)
- [ ] Horizontal scaling ready (Docker + LB)
- [ ] Database indexed
- [ ] Caching strategy (Redis)
- [ ] Monitoring & alerting
- [ ] Graceful shutdown
- [ ] Health checks

## NEXT IMMEDIATE STEPS

1. **Refresh Token Endpoint** (2-3 ore)
   - Implementează backend logic
   - Testează cu curl
   - Adaugă Rust command
2. **RBAC Implementation** (1-2 ore)

   - Adaugă role în database
   - Creează middleware
   - Testează cu admin/user

3. **Protected Route Exemplu** (2 ore)
   - Submissions API
   - Test end-to-end

**Target: Fundație completă în 1-2 zile**
