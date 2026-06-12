# Go-CMS

> Personal portfolio & blog — Go monolith with all frontend assets embedded into a single binary.

**Stack:** Go · Chi · GORM · PostgreSQL · Redis · html/template · Docker

---

## Features

- 📝 **Blog** — listing & detail pages with server-side Markdown rendering (goldmark) and reading duration details.
- 💼 **Portfolio** — project showcase with tech stack badges and live/repository links.
- 🔐 **Admin Dashboard** — private CMS for blog & portfolio CRUD, plus integrated image upload.
- 🌗 **Dark Mode** — Light/Dark toggle using CSS Variables + localStorage persistence.
- 📦 **Single Binary** — all templates & static assets embedded via `go:embed`.
- 🔒 **Security** — Bcrypt password hashing, server-side sessions, CSRF protection, and security headers.

---

## Quick Start

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- Atlas CLI (for database migrations)
- [`air`](https://github.com/air-verse/air) for hot-reload (optional)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/danangamw/go-cms.git
cd go-cms

# Copy and configure environment variables
cp .env.example .env

# Start services (PostgreSQL + Redis + MinIO) via Docker Compose
docker compose up -d

# Run database migrations (Atlas)
make migrate-apply

# Seed the admin user and mock data
make seed

# Start the server with hot-reload (air)
make run

# Or without hot-reload
go run ./cmd/app
```

Open **http://localhost:8080** in your browser.  
Admin dashboard: navigate manually to **http://localhost:8080/login**.

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values:

| Variable | Description | Default |
|---|---|---|
| `APP_ENV` | `development` or `production` | `development` |
| `APP_PORT` | Server port | `8080` |
| `APP_SECRET_KEY` | Secret key, minimum 32 characters | — (required) |
| `DATABASE_URL` | PostgreSQL connection string | — (required) |
| `REDIS_URL` | Redis URL (optional, falls back to Postgres sessions) | — |
| `UPLOAD_STORAGE` | `local` or `s3` | `local` |
| `UPLOAD_DIR` | Local upload directory | `./uploads` |
| `S3_BUCKET` | S3/MinIO bucket name | `go-cms` |
| `S3_REGION` | S3/MinIO region | `us-east-1` |
| `S3_ENDPOINT` | Internal S3/MinIO endpoint (uploads) | `http://localhost:9000` |
| `S3_PUBLIC_ENDPOINT` | Public subdomain reverse proxy endpoint (generated URLs) | `https://media.danangamw.com` |
| `S3_ACCESS_KEY_ID` | S3/MinIO access key | — |
| `S3_SECRET_ACCESS_KEY` | S3/MinIO secret access key | — |
| `ADMIN_USERNAME` | Admin username for seeding | `admin` |
| `ADMIN_PASSWORD` | Admin password for seeding | — |

---

## Project Structure

```
go-cms/
├── cmd/
│   ├── app/main.go        ← Server entry point
│   ├── seed/main.go       ← Admin user & sample data seeder
│   └── atlas-loader/      ← Atlas schema bridge loader
│
├── internal/
│   ├── config/            ← Environment variable loader
│   ├── database/          ← GORM + PostgreSQL connection & health check
│   ├── model/             ← Domain entities (User, Blog, Portfolio)
│   ├── repository/        ← Database repositories
│   ├── service/           ← Business logic (Auth service, Upload service)
│   ├── handler/           ← HTTP handlers (Public & Admin controllers)
│   ├── middleware/        ← Auth guard, CSRF validation, telemetry
│   └── session/           ← Session store (Redis / Postgres fallback)
│
├── pkg/
│   ├── logger/            ← Structured slog logger
│   └── storage/           ← S3/MinIO storage package
│
├── web/
│   ├── templates/         ← HTML templates (embedded into binary)
│   └── static/            ← CSS, JS, images (embedded into binary)
│
├── migrations/            ← Atlas database versioned migration SQL files
└── uploads/               ← Local upload directory (ignored in git)
```

---

## Commands

```bash
# Development
make run            # Start server with hot-reload (using air)
make build          # Build production binary to bin/go-cms
make test           # Run all unit tests
make seed           # Seed admin user and mock contents

# Database Migrations (Atlas)
make migrate-status # Show current migration status
make migrate-apply  # Apply pending schema migrations
make migrate-diff   # Generate a new migration file (usage: make migrate-diff name=description)

# Docker
make docker-run     # Build and start all services via Docker Compose
make docker-down    # Stop all Docker Compose services
```

---

## Architecture

```
HTTP Request
    ↓
[Middleware]  — auth guard, CSRF validation, telemetry, security headers
    ↓
[Handler]     — input validation, form/URL parsing, template rendering
    ↓
[Service]     — business logic (slug generation, Upload management)
    ↓
[Repository]  — GORM queries to PostgreSQL
```

All frontend assets (`web/`) are embedded into the binary via `//go:embed`, so deployment requires only a single executable file.

---

## Deployment

```bash
# Build Docker image
docker build -t go-cms .

# Start production stack
docker compose up -d
```

Health check endpoint: `GET /health`

---

## License

MIT
