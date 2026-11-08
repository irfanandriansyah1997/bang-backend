# `internal/infrastructure/config/` — app settings

Reads environment variables (`PORT`, `DATABASE_URL`, `JWT_SECRET`, ...) into one typed `Config` struct.

> React analogy: `import.meta.env` in Vite, but read once and validated.

## ✅ Do

- One `Config` struct with every setting the app needs
- Sensible defaults for local dev (e.g. `PORT=8080`)
- Validate at startup and return a clear error if a required value is missing
- Keep `.env.example` at the repo root in sync with the fields here

## ❌ Don't

- Call `os.Getenv` anywhere else in the codebase. Read it here once, then pass values down
- Commit real secrets (`.env` is git-ignored; `.env.example` holds placeholders)
- Log the full config (it contains secrets)

## Example

```go
type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}

func Load() (Config, error) { ... }
```
