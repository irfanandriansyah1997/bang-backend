---
description: Migration, sqlc query, and repository rules
globs: db/**/*.sql,internal/repository/**
alwaysApply: false
paths:
  - "db/**/*.sql"
  - "internal/repository/**"
---

# Database rules

- Migrations (goose): one change per file, always with `-- +goose Up` and `-- +goose Down`.
- Never edit a migration that has already been run; create a new one.
- Queries (sqlc): name every query (`-- name: GetWorkOrder :one`), one file per table.
- After changing queries, run `make sqlc`. Never edit the generated Go code.
- Cursor pagination: `WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC, id DESC LIMIT n+1`.
- No business rules in SQL; keep them in `internal/domain`.
- Repositories map DB errors to domain errors (`pgx.ErrNoRows` → `domain.ErrNotFound`).
