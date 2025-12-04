# LiveCode TODO List

## 🎯 Current Sprint: AI Setup + Auth Completion

### Priority 1: AI & Development Setup
- [x] Create `.claude/` folder structure
- [x] Write project-context.md
- [x] Write rules.md
- [x] Write todo.md
- [x] Write architecture.md
- [x] Add CRITICAL rules for TODO/docs updates
- [ ] Test context persistence (new conversation)
- [ ] Setup thesis documentation (LaTeX) - waiting for PDF guidelines

### Priority 2: Complete Authentication System
- [ ] **Password Security**
  - [ ] Research bcrypt vs argon2 for Rust
  - [ ] Add password hashing dependency
  - [ ] Implement hash_password() function
  - [ ] Implement verify_password() function

- [ ] **Backend - Register**
  - [ ] Uncomment password hashing in register handler
  - [ ] Implement user insertion into DB
  - [ ] Add proper error handling
  - [ ] Register auth commands in main.rs
  - [ ] Test with Postman/curl

- [ ] **Frontend - Register**
  - [ ] Replace mock setTimeout with actual Tauri invoke
  - [ ] Handle backend errors properly
  - [ ] Add network error handling
  - [ ] Test full registration flow

- [ ] **Login Implementation**
  - [ ] Design login UI (LoginForm.tsx)
  - [ ] Implement login backend handler
  - [ ] Add session management (decide: JWT, cookies, etc.)
  - [ ] Connect frontend to backend
  - [ ] Test login flow

- [ ] **Session Management**
  - [ ] Research session strategies (JWT vs server-side)
  - [ ] Implement chosen strategy
  - [ ] Add logout functionality
  - [ ] Persist auth state across app restarts

### Priority 3: Project Architecture Improvements
- [ ] **Database**
  - [ ] Review schema design with Claude
  - [ ] Plan future tables (projects, connections, sessions)
  - [ ] Add indexes for performance
  - [ ] Implement updated_at auto-trigger

- [ ] **Code Structure**
  - [ ] Review Rust module organization
  - [ ] Separate concerns (handlers, services, models)
  - [ ] Add proper error types (custom Error enum)
  - [ ] Implement logging (tracing crate?)

- [ ] **Security**
  - [ ] Move .env to .env.example (gitignore .env)
  - [ ] Add input sanitization
  - [ ] Review SQL injection protection
  - [ ] Add rate limiting (future)

### Priority 4: Core Features (After Auth)
- [ ] **Project/Connection Management**
  - [ ] Design DB schema for projects/connections
  - [ ] Design UI for project list
  - [ ] Implement CRUD operations
  - [ ] Add encryption for stored credentials

- [ ] **SSH/SFTP Integration**
  - [ ] Research Rust SSH/SFTP crates
  - [ ] Implement connection testing
  - [ ] Build file browser backend
  - [ ] Build file browser UI

- [ ] **Real-time Collaboration**
  - [ ] Design file locking mechanism
  - [ ] Choose real-time tech (WebSocket, polling)
  - [ ] Implement backend locking system
  - [ ] Add visual indicators (red for locked files)

### Priority 5: DevOps & Deployment
- [ ] **Docker**
  - [ ] Learn Docker basics
  - [ ] Create Dockerfile for app
  - [ ] Create docker-compose.yml (app + PostgreSQL)
  - [ ] Test containerized deployment

- [ ] **Kubernetes** (Future - Enterprise Level)
  - [ ] Learn Kubernetes fundamentals
  - [ ] Create k8s manifests
  - [ ] Deploy to local cluster (minikube)
  - [ ] Explore Helm charts

- [ ] **Testing**
  - [ ] Add Rust unit tests
  - [ ] Add integration tests
  - [ ] Add frontend tests (Vitest?)

- [ ] **CI/CD**
  - [ ] GitHub Actions for tests
  - [ ] Automated builds
  - [ ] Release workflow

## 📚 Learning Goals
- [ ] Understand Rust async/await deeply
- [ ] Master SQLx compile-time queries
- [ ] Learn Docker containerization
- [ ] Explore Kubernetes orchestration
- [ ] Understand Tauri architecture
- [ ] Study real-time collaboration patterns
- [ ] Enterprise-level architecture patterns

## 💡 Ideas / Future Features
- [ ] OAuth integration (Google, GitHub)
- [ ] Cloud storage integration (S3, GCS)
- [ ] File editing in-app (code editor)
- [ ] Terminal integration (SSH shell)
- [ ] File transfer progress bars
- [ ] Keyboard shortcuts
- [ ] Dark/light theme preferences saved
- [ ] Redis caching layer
- [ ] Monitoring with Prometheus/Grafana
- [ ] gRPC for internal services (if microservices)

## 🐛 Known Issues
- [ ] Auth handlers not registered (main.rs)
- [ ] Password hashing not implemented
- [ ] updated_at field not auto-updating (needs trigger)
- [ ] LoginForm.tsx empty
- [ ] Frontend using mock data (not calling backend)

## ✅ Completed
- [x] Tauri v2 project setup
- [x] PostgreSQL database connection
- [x] SQLx migrations system
- [x] Users table schema
- [x] Custom window with titlebar
- [x] Light/Dark theme toggle
- [x] Registration form UI with validation
- [x] Auth backend handlers (code written)