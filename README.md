# Benah Backend

Go + PostgreSQL REST API for **Benah**, a field service / work order app:

- **Benah Admin** (web, React): admins create work orders, assign technicians, and track status.
- **Benah Mitra** (Android, Kotlin): technicians see their tasks, start work, and complete it with photo proof and customer approval.

This repo is the backend (`/api/v1`). It is also a learning project for working the way a Forward Deployed Engineer (FDE) does. See the [spec](docs/superpowers/specs/rencana-belajar-fullstack-fde.md).

> **Status: Phase 0 (foundation).** The folder structure and docs are in place, but there is **no Go code yet**. Commands marked *(planned)* below will work once their issue is done. Track progress on the [GitHub Project](https://github.com/users/irfanandriansyah1997/projects/1).

---

## Minimum requirements

### Hardware

| | Backend only | Full stack (backend + web + Android emulator + Kibana) |
|---|---|---|
| CPU | Any 64-bit (Apple Silicon or x86_64) | Apple Silicon M1/M2 or equivalent |
| RAM | 8 GB | **16 GB** (run Elasticsearch/Kibana only when needed) |
| Disk | ~5 GB free | ~30 GB free (Android Studio + emulator images) |
| OS | macOS, Linux, or Windows with WSL2 | macOS recommended (the project is developed on macOS) |

### Software

| Tool | Version | Needed for | Install (macOS) |
|---|---|---|---|
| Go | **1.22+** (needed for method/path routing in `net/http`) | Building and running the API | `brew install go` |
| Docker runtime | Docker Compose v2 | PostgreSQL (and optional observability stack) | `brew install orbstack` (lighter than Docker Desktop on Apple Silicon) |
| PostgreSQL | 16 (runs in Docker; no local install needed) | Database | via `docker-compose.yml` |
| GNU Make | any | Task shortcuts | preinstalled on macOS |
| lefthook | latest | Git hooks (pre-commit checks, like husky) | `brew install lefthook` |
| goose | latest | Database migrations | `brew install goose` |
| sqlc | latest | Generating type-safe Go from SQL | `brew install sqlc` |
| air | latest | Hot reload in development | `go install github.com/air-verse/air@latest` |
| API client | — | Testing endpoints | curl, [Bruno](https://www.usebruno.com/), or Postman |
| DB viewer | — | Inspecting data (optional) | TablePlus or DBeaver |

Optional, for observability: Elasticsearch, Kibana, and Filebeat run as a Docker Compose profile and need about 1.5–2 GB of extra RAM.

Check your setup:

```bash
go version          # go1.22 or newer
docker compose version
goose --version
sqlc version
```

---

## Getting started

```bash
git clone https://github.com/irfanandriansyah1997/bang-backend.git
cd bang-backend
```

### 0. Install git hooks (once per clone)

```bash
brew install lefthook
make hooks
```

After this, every `git commit` runs checks first (like husky in Node.js) and is **blocked if any fail**:

| Hook | Runs when | Checks |
|---|---|---|
| `pre-commit` | staged `.go` files | `gofmt`, `go vet`, `go build`, `go test` |
| `commit-msg` | every commit | message is `<type>(<feature>): <message>`, feature required, ≤ 72 chars |

Run the same code checks by hand any time with `make check`. Hooks are configured in [`lefthook.yml`](lefthook.yml).

### 1. Configure environment *(planned, #2)*

```bash
cp .env.example .env
```

| Variable | Example | Purpose |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `DATABASE_URL` | `postgres://app:app@localhost:5432/fieldservice?sslmode=disable` | Postgres connection |
| `JWT_SECRET` | any long random string | Signs access/refresh tokens |
| `ELASTIC_VERSION` | e.g. `9.x.x` | Only for the observability profile |

Never commit `.env`. It is already in `.gitignore`.

### 2. Start PostgreSQL *(planned, #2)*

```bash
docker compose up -d
```

### 3. Migrate and seed *(planned, #5, #7)*

```bash
make migrate-up
make seed           # 1 admin, 3 technicians, 5 customers, ~10 work orders
```

### 4. Run the API *(planned, #2, #4)*

```bash
make dev            # hot reload with air
# or
make run            # go run ./cmd/api
```

Check that it's up:

```bash
curl http://localhost:8080/health
# {"status":"ok","db":"ok"}
```

From the **Android emulator**, the Mac's localhost is `http://10.0.2.2:8080`.
From the **web app** (Vite on `http://localhost:5173`), CORS is allowed for that origin.

---

## Common commands

| Command | What it does | Status |
|---|---|---|
| `make hooks` | Install git hooks (lefthook) | ✅ available |
| `make check` | gofmt + vet + build + test (same as pre-commit) | ✅ available |
| `make ai-link` | Symlink AI rules/skills from `.ai/` into `.claude/` and `.cursor/` | ✅ available |
| `make run` | Run the API once | planned (#2) |
| `make dev` | Run with hot reload (air) | planned (#2) |
| `make migrate-up` / `make migrate-down` | Apply / roll back migrations (goose) | planned (#5) |
| `make sqlc` | Regenerate Go code from `db/queries/*.sql` | planned (#5) |
| `make seed` | Load sample data | planned (#7) |
| `make test` | `go test ./...` | planned (#2) |

Reset the database completely (deletes all data):

```bash
docker compose down -v
```

### Observability (optional) *(planned)*

```bash
docker compose --profile observability up -d     # + Elasticsearch, Kibana, Filebeat
docker compose --profile observability stop      # stop to free RAM
```

Kibana runs at `http://localhost:5601`. Search logs by `request_id`, which every response returns in the `X-Request-ID` header.

---

## Deployment

**Current target: local only.** The MVP runs on a developer machine with Docker. The Android app is tested on an emulator. No hosted environment exists yet.

When a hosted environment is added, the plan is:

1. **Build** a static binary:
   ```bash
   CGO_ENABLED=0 go build -o bin/api ./cmd/api
   ```
   or package it as a Docker image (multi-stage build → small `distroless`/`alpine` image).
2. **Provision** PostgreSQL 16 (managed service recommended).
3. **Configure** environment variables (`PORT`, `DATABASE_URL`, `JWT_SECRET`) through the platform's secret store, never in the image.
4. **Migrate** before starting the new version: `goose -dir db/migrations postgres "$DATABASE_URL" up`.
5. **Store uploads** (work order photos, ticket attachments) on a persistent volume or object storage. The local version writes to disk.
6. **Ship logs**: the API logs JSON to stdout, so any log collector (Filebeat, the platform's built-in logging) can ingest it.
7. **Health check**: point the platform's health check at `GET /health`.

Production checklist before going live: HTTPS only, a strong `JWT_SECRET`, CORS restricted to the real web origin, rate limiting on public endpoints, Elasticsearch security enabled (it is disabled in the local compose file).

---

## Project structure

```
cmd/api/                 entry point: wires dependencies, starts the server
internal/
  domain/                entities + business rules (no external imports)
  usecase/               one flow per user action + the interfaces it needs
  delivery/http/         router, handlers, middleware, DTOs
  repository/postgres/   database access (pgx + sqlc)
  infrastructure/        config, logger, database connection
db/migrations/           goose migrations
db/queries/              SQL for sqlc
api/                     openapi.yaml (contract for web + Android)
docs/                    architecture, specs, plans
.ai/                     AI rules and skills (single source for Claude Code + Cursor)
```

Every folder has its own `README.md` explaining what belongs there. Start with [`docs/architecture.md`](docs/architecture.md).

---

## Tech stack

- **Language:** Go (`net/http` routing, `log/slog`)
- **Database:** PostgreSQL 16, [pgx](https://github.com/jackc/pgx), [sqlc](https://sqlc.dev), [goose](https://github.com/pressly/goose)
- **Auth:** JWT (access + refresh tokens), bcrypt
- **Validation:** go-playground/validator
- **API contract:** OpenAPI (`api/openapi.yaml`)
- **Observability:** structured JSON logs → Filebeat → Elasticsearch → Kibana
- **Dev:** air (hot reload), Docker Compose via OrbStack

---

## Contributing

- Tasks are GitHub issues, grouped by phase labels (`phase-0-foundation` … `phase-4-call-center`) on the [Project board](https://github.com/users/irfanandriansyah1997/projects/1).
- Follow the layer rules in [`docs/architecture.md`](docs/architecture.md) and the folder READMEs.
- Commit messages use `<type>(<feature>): <message>`, based on [Conventional Commits](https://www.conventionalcommits.org/) with the feature **required**:
  - `feat(docs): adding readme.md with setup steps`
  - `fix(work-order): reject complete without photos`
  - Types: `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `style`, `build`, `ci`, `chore`, `revert`
  - Feature: lowercase kebab-case name of the feature or area (`docs`, `auth`, `customer`, `work-order`, ...)
  - Issue footers: `Refs: #12` (partial), `Closes: #12` (done)
- Run `go vet ./...` and `go test ./...` before opening a PR.
- AI agents (Claude Code, Cursor) read [`AGENTS.md`](AGENTS.md). Rules and skills live in [`.ai/`](.ai/README.md).
