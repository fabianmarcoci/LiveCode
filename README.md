# LiveCode

A modern desktop application for real-time code collaboration and file transfer, built as a WinSCP alternative with enterprise-level authentication.

## Project Structure

This is a **monorepo** containing two main projects:

- **[desktop-app/](desktop-app/)** - Tauri desktop application (React + TypeScript + Rust)
- **[backend-api/](backend-api/)** - REST API server (Go + Gin + PostgreSQL)

## Tech Stack

### Desktop App
- **Frontend:** React 18 + TypeScript + Vite
- **Backend:** Tauri v2 (Rust)
- **UI:** Custom components with modern design
- **Auth:** JWT tokens stored in OS Keyring (Windows Credential Manager)

### Backend API
- **Language:** Go 1.22
- **Framework:** Gin
- **Database:** PostgreSQL
- **Auth:** JWT (access tokens 15min) + Refresh tokens (30 days)

## Getting Started

### Prerequisites
- Node.js 18+
- Rust (latest stable)
- Go 1.22+
- PostgreSQL 14+

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd LiveCode
   ```

2. **Backend API Setup**
   ```bash
   cd backend-api
   cp .env.example .env
   # Edit .env with your database credentials
   go mod download
   go run main.go
   ```

3. **Desktop App Setup**
   ```bash
   cd desktop-app
   npm install
   npm run tauri dev
   ```

## Development

Open the workspace in VS Code:
```bash
code livecode.code-workspace
```

### Recommended IDE Extensions
- [Tauri](https://marketplace.visualstudio.com/items?itemName=tauri-apps.tauri-vscode)
- [rust-analyzer](https://marketplace.visualstudio.com/items?itemName=rust-lang.rust-analyzer)
- [Go](https://marketplace.visualstudio.com/items?itemName=golang.go)

## Architecture

```
Desktop App (Tauri) → Backend API (Go) → PostgreSQL
       ↓
  OS Keyring (JWT storage)
```

## License

This project is part of a bachelor's thesis and is intended for educational purposes.
