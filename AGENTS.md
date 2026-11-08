# AGENTS.md — shared instructions for AI coding agents

This file is the **single source of truth** for project rules that every AI agent should follow.
Cursor reads `AGENTS.md` automatically; Claude Code loads it via `@AGENTS.md` in `CLAUDE.md`.

## Project

Benah backend: Go + PostgreSQL REST API (`/api/v1`) for a field service / work order app
(web admin in React, Android technician app in Kotlin).

- Spec: `docs/superpowers/specs/rencana-belajar-fullstack-fde.md`
- Architecture: `docs/architecture.md` (read before adding code)
- Tasks: GitHub Project https://github.com/users/irfanandriansyah1997/projects/1

## The developer

New to Go, ~7 years of frontend (React/TypeScript). When explaining code, compare Go concepts with React concepts and keep explanations short.

## Architecture rules (must follow)

- Clean architecture. Imports point inward only: `delivery → usecase → domain`.
- `domain` imports only the standard library: no `net/http`, `pgx`, sqlc, JWT, `slog`, and no `json:` / `db:` tags.
- `usecase` defines the interfaces it needs; `repository/postgres` implements them.
- Business rules live in `domain` (rules) and `usecase` (order of steps), never in handlers or SQL.
- Handlers are thin: decode → call usecase → encode. Map domain errors to HTTP status in one place.
- Never return sqlc types outside `repository/postgres`.
- `cmd/api/main.go` is the only place that wires dependencies.
- Each folder's `README.md` has its own do's and don'ts. Follow them.

## Conventions

- Errors: `{ "error": { "code": "...", "message": "..." } }`
- Pagination: cursor-based, `?limit=15&cursor=<opaque>` → `{ data, next_cursor, prev_cursor }`
- Time: ISO 8601, UTC. Pass `now` into domain methods instead of calling `time.Now()` there.
- Logging: `log/slog` JSON with `request_id`; never log passwords, tokens, or full personal data.
- Every I/O function takes `ctx context.Context` as its first parameter.

## Workflow

- Never edit a migration that has already been run; add a new one.
- After changing `db/queries/*.sql`, run `make sqlc`; never edit generated code.
- Update `api/openapi.yaml` in the same change as an endpoint.
- Run `go test ./...` and `go vet ./...` before saying a task is done.
- AI rules/skills: edit only in `.ai/`, then run `make ai-link`. Never edit the symlinks in `.claude/` or `.cursor/`.
