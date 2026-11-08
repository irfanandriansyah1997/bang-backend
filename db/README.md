# `db/` — database schema and SQL

| Folder | Tool | Purpose |
|---|---|---|
| [`migrations/`](migrations/README.md) | goose | Versioned schema changes (create/alter tables) |
| [`queries/`](queries/README.md) | sqlc | Hand-written SQL that sqlc turns into type-safe Go code |

These are plain `.sql` files, not Go code, so they live outside `internal/`.
