# LiveCode Project Context

## Vision
A minimalist WinSCP alternative with real-time collaboration features.
Users can manage SSH/SFTP connections, browse remote files, and see when others are editing files (visual locking).

**Goal**: Build an enterprise-level product with modern architecture and scalable infrastructure.

## What is LiveCode?
- **Type**: Desktop application for remote file management
- **Target Users**: Developers managing remote servers/projects
- **Key Differentiator**: Real-time collaboration (show locked files when others edit)
- **Alternative to**: WinSCP, FileZilla (but minimalist + collaborative)
- **Ambition**: Production-ready, enterprise-grade software

## Tech Stack

### Backend
- **Framework**: Tauri v2
- **Language**: Rust
- **Database**: PostgreSQL with SQLx (compile-time verified queries)
- **Async Runtime**: Tokio
- **Migrations**: SQLx migrations (auto-run on startup)
- **Env Config**: dotenvy (.env files)

### Frontend
- **Framework**: React 19
- **Language**: TypeScript 5.8 (strict mode)
- **Build Tool**: Vite 7
- **Styling**: Pure CSS with CSS variables (no framework)
- **State Management**: React hooks (no Redux/Zustand yet)

### Future Stack (Planned)
- **Containerization**: Docker
- **Orchestration**: Kubernetes (if scaling to cloud/multi-instance)
- **Password Hashing**: bcrypt or argon2
- **SSH/SFTP**: Rust crate (TBD)
- **Real-time**: WebSockets or similar
- **Cloud Storage**: AWS S3, Google Cloud Storage, etc.
- **Monitoring**: Prometheus, Grafana (enterprise-level observability)
- **CI/CD**: GitHub Actions, automated deployments

## What's Implemented
- ✅ Custom frameless window with titlebar
- ✅ Light/Dark theme switching
- ✅ Registration form UI (complete with validation)
- ✅ PostgreSQL database with users table
- ✅ SQLx migrations system
- ✅ Auth backend handlers (written but not connected)

## Learning Goals
- **Primary**: Deep Rust knowledge (async, error handling, system integration)
- **Secondary**: Docker, Kubernetes, SQL migrations, modern frontend, AI-assisted development
- **Interest areas**: Cybersecurity patterns, networking protocols, secure system design
- **Outcome**: Enterprise-level portfolio project to showcase to employers

## Architecture Notes

### Database Schema
- `users` table: UUID primary key, email/username unique, OAuth support, timestamps
- Future: `projects`, `connections`, `sessions`, `file_locks` tables

### Auth Strategy
- Dual auth: Email/password + OAuth (Google, GitHub)
- Session-based (JWT or similar)
- Password hashing with bcrypt/argon2

### File Management Strategy
- SSH/SFTP connections stored per user
- Real-time file lock tracking (WebSocket or polling)
- Visual indicators (red = locked by other user)

## Important Context
- **Owner**: Marco (learning-focused developer)
- **Development Style**: Incremental, learn-as-you-go
- **Code Quality**: Enterprise-level architecture for portfolio
- **Git**: Version controlled, main branch
- **Deployment**: Docker + potential Kubernetes orchestration

## Technologies to Explore
- Kubernetes
- gRPC (if microservices needed)
- Redis (caching layer)
- Message queues (RabbitMQ, Kafka)
- Observability tools (Prometheus, Grafana, OpenTelemetry)
- Advanced security (Vault for secrets management)