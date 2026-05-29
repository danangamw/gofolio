# Go-CMS

> Personal portfolio & blog — Go monolith with all frontend assets embedded into a single binary.

**Stack:** Go · Chi · GORM · PostgreSQL · Redis · html/template · Docker

---

## Features

- 📝 **Blog** — listing & detail pages with server-side Markdown rendering (goldmark), pagination, and read time estimation
- 💼 **Portfolio** — project showcase with tech stack badges, thumbnails, and repo/live demo links
- 🔐 **Admin Dashboard** — private CMS for blog & portfolio CRUD, image uploads
- 🌗 **Dark Mode** — Light/Dark toggle using OKLCH CSS variables + localStorage persistence
- 📦 **Single Binary** — all templates & static assets embedded via `go:embed`
- 🔒 **Security** — Argon2id password hashing, server-side sessions, CSRF protection, rate limiting, security headers

---

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- [`air`](https://github.com/air-verse/air) for hot-reload (optional)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/danangamw/go-cms.git
cd go-cms

# Copy and configure environment variables
cp .env.example .env

# Start PostgreSQL + Redis via Docker
docker compose up -d postgres redis

# Run database migrations
go run ./cmd/migrate

# Seed the first admin user
go run ./cmd/seed

# Start the server with hot-reload
air

# Or without hot-reload
go run ./cmd/api
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
| `ADMIN_USERNAME` | Admin username for seeding | — |
| `ADMIN_PASSWORD` | Admin password for seeding | — |

---

## Project Structure

```
go-cms/
├── cmd/
│   ├── api/main.go        ← Server entry point
│   ├── migrate/main.go    ← GORM migration runner
│   └── seed/main.go       ← Admin user seeder
│
├── internal/
│   ├── config/            ← Environment variable loader
│   ├── database/          ← GORM + PostgreSQL connection & health check
│   ├── model/             ← Domain entities (User, Blog, Portfolio, Session)
│   ├── repository/        ← Database queries
│   ├── service/           ← Business logic
│   ├── handler/           ← HTTP handlers
│   ├── middleware/        ← Auth guard, CSRF, rate limiter, security headers
│   └── session/           ← Session store (Redis / Postgres fallback)
│
├── web/
│   ├── templates/         ← HTML templates (embedded into binary)
│   └── static/            ← CSS, JS, images (embedded into binary)
│
├── migrations/            ← SQL migration files
└── uploads/               ← Upload directory (not embedded, mounted as volume)
```

---

## Commands

```bash
# Development
make run          # Start server with hot-reload (air)
make build        # Build production binary
make test         # Run all unit tests
make migrate      # Run database migrations
make seed         # Seed admin user

# Docker
make docker-up    # Start stack (postgres + redis)
make docker-build # Build and start all services

# Other
make lint         # go vet
make clean        # Remove binary
```

---

## Architecture

```
HTTP Request
    ↓
[Middleware]  — auth guard, CSRF validation, rate limiter, security headers
    ↓
[Handler]     — input validation, form/URL parsing, template rendering
    ↓
[Service]     — business logic (slug generation, Markdown rendering, hashing)
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

Health check endpoint: `GET /healthz`

---

## License

MIT
