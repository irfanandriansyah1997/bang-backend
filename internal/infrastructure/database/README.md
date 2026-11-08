# `internal/infrastructure/database/` — the DB connection

Opens and configures the PostgreSQL connection pool (`pgxpool`) and checks it with a ping at startup.

Local connection string (from `docker-compose.yml`):
`postgres://app:app@localhost:5432/fieldservice?sslmode=disable`

## ✅ Do

- `Connect(ctx, url) (*pgxpool.Pool, error)` that pings before returning
- Pool settings (max connections, timeouts) in one place
- Close the pool on graceful shutdown

## ❌ Don't

- Write business queries here (that's `repository/postgres/`)
- Open a new connection per request. Reuse the one pool created in `cmd/api`
