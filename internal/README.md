# `internal/` — the app's private room

`internal` is a **special folder name in Go**. Packages inside it can only be imported by code in this repository. If another project tries to import them, the **Go compiler refuses**. It is enforced by the toolchain, not just a convention.

- **Put here:** all application code (domain, usecase, delivery, repository, infrastructure).
- **Not here:** code meant to be reused by other projects (that usually goes in `pkg/`; we don't need it now).

## The layers

```
HTTP request
   ↓
delivery/http        → receive request, parse JSON/form, call usecase, write response
   ↓
usecase              → orchestrate the flow: load WO, check owner, check photos, save, log event
   ↓
domain               → business rules: "can complete only if in_progress + has photos"
   ↑
repository/postgres  → load/save data (called by usecase through an interface)
infrastructure       → tools: config, logger, DB connection
```

Example flow, "technician completes a task" (`POST /me/work-orders/{id}/complete`):
handler parses the multipart form → `CompleteWorkOrder` usecase loads the WO via the repository → domain method `wo.Complete(...)` checks the rule → usecase saves via the repository → handler returns JSON.

## Import cheat sheet

Dependencies point **inward only**: `delivery → usecase → domain`.

| This package… | may import | must NOT import |
|---|---|---|
| `domain` | standard library only | everything below |
| `usecase` | `domain` | `delivery`, `repository`, `pgx`, `net/http` |
| `delivery/http` | `usecase`, `domain` (for errors) | `repository`, `pgx` |
| `repository/postgres` | `domain`, `usecase` (interfaces) | `delivery` |
| `infrastructure/*` | external libraries | `delivery`, `usecase` |
| `cmd/api` | **everything** | — (wiring is its job) |

## Not sure where code belongs? Ask 3 questions

1. **If HTTP were replaced by a CLI, would this code still be needed?** No → `delivery`.
2. **If Postgres were replaced by MySQL, would this code change?** Yes → `repository`.
3. **Is this a rule the client/business asked for?** Yes → `domain` (the rule) or `usecase` (the order of steps).
