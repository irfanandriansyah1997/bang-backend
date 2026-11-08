# `internal/infrastructure/` — the toolbox

Shared **technical** tools used by the other layers. No business logic lives here.

| Folder | Purpose |
|---|---|
| [`config/`](config/README.md) | Read environment variables into a typed `Config` struct |
| [`logger/`](logger/README.md) | Set up `slog` to write JSON logs (for Kibana) |
| [`database/`](database/README.md) | Open the Postgres connection pool |

Later this may also hold `storage/` (save photos/attachments to disk) and `auth/` (sign/verify JWT).

## ✅ Do

- Wrap external libraries behind small, simple constructors (`config.Load()`, `logger.New(cfg)`, `database.Connect(ctx, url)`)
- Fail fast at startup when something required is missing

## ❌ Don't

- Business logic of any kind
- Import `usecase` or `delivery`
- Use global singletons that other packages reach into. Create them in `cmd/api` and pass them down
