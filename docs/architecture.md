# Architecture — Benah Backend

The Benah backend (Go + PostgreSQL) follows **clean architecture**: business rules sit in the center and don't depend on HTTP, the database, or any framework. Everything technical is an *adapter* around that center.

Each folder also has its own `README.md` with detailed do's and don'ts. This page is the overview.

Spec: [`superpowers/specs/rencana-belajar-fullstack-fde.md`](superpowers/specs/rencana-belajar-fullstack-fde.md)

---

## 1. Folder structure

```
bang-backend/
├── cmd/api/                       entry point: start the app, wire dependencies
├── internal/                      private app code (Go compiler blocks outside imports)
│   ├── domain/                    entities + business rules (no external imports)
│   ├── usecase/                   one flow per user action + the interfaces it needs
│   ├── delivery/http/             router, handlers, middleware, DTOs, error mapping
│   ├── repository/postgres/       implements usecase interfaces with pgx + sqlc
│   └── infrastructure/
│       ├── config/                env → typed Config
│       ├── logger/                slog JSON logger
│       └── database/              pgx connection pool
├── db/
│   ├── migrations/                goose schema migrations
│   └── queries/                   SQL for sqlc
├── api/                           openapi.yaml (contract for web + Android)
└── docs/                          this file, specs, plans
```

---

## 2. The layers

```
            ┌──────────────────────────────────────────────┐
            │ cmd/api  (wires everything at startup)       │
            └──────────────────────────────────────────────┘

  HTTP ──▶ delivery/http ──▶ usecase ──▶ domain
                               │  ▲
                     interface │  │ implements
                               ▼  │
                       repository/postgres ──▶ PostgreSQL

  infrastructure (config, logger, database) = tools used by the outer layers
```

| Layer | Role | React analogy |
|---|---|---|
| `cmd/api` | Power button + wiring (composition root) | `main.tsx` + `<Providers>` |
| `domain` | Business rules, pure and testable | Pure util functions (`canTransition(from, to)`) |
| `usecase` | Order of steps for one user action | Custom hook that receives `api` from outside |
| `delivery/http` | Translate HTTP ↔ usecase | Page/component layer |
| `repository/postgres` | Load/save data | The `fetch`/API client implementation |
| `infrastructure` | Technical tools | `import.meta.env`, logger setup |

### The dependency rule

**Imports point inward only:** `delivery → usecase → domain`.
`repository` implements interfaces that `usecase` defines. `domain` imports nothing from the project.

| This package… | may import | must NOT import |
|---|---|---|
| `domain` | standard library only | everything below |
| `usecase` | `domain` | `delivery`, `repository`, `pgx`, `net/http` |
| `delivery/http` | `usecase`, `domain` (for errors) | `repository`, `pgx` |
| `repository/postgres` | `domain`, `usecase` (interfaces) | `delivery` |
| `infrastructure/*` | external libraries | `delivery`, `usecase` |
| `cmd/api` | **everything** | — |

Why: business rules can be tested without a database or server, and changing a technical detail (HTTP → gRPC, Postgres → something else) touches only one folder.

---

## 3. Do's and don'ts per layer

### `cmd/api`
- ✅ Load config, create logger, open DB, build repo → usecase → handler, start server, graceful shutdown
- ❌ SQL, business rules, handlers, a giant `main.go` (aim for < 100 lines)

### `domain`
- ✅ Entity structs, status enums, rule methods (`wo.Start(now)`), business errors (`ErrInvalidTransition`), accept `now` as a parameter
- ❌ Importing `net/http` / `pgx` / sqlc / JWT / `slog`, `json:`/`db:` tags, DB or file access, HTTP status codes

### `usecase`
- ✅ Define the interfaces it needs (`WorkOrderRepository`, `FileStorage`, `TxManager`, `Clock`), call domain rules, business authorization, transactions via interface, business event logs, `ctx` as first param
- ❌ `*http.Request`, headers/JSON, `pgx`/sqlc types, SQL, HTTP status codes, mutable globals

### `delivery/http`
- ✅ Routes, parse params/JSON/multipart, **format** validation, DTOs with `json:` tags, middleware (request_id, logging, auth, role, CORS, recover), domain error → HTTP mapping
- ❌ Business rules, calling repositories, opening transactions, returning domain structs as JSON, long handlers (decode → call usecase → encode)

### `repository/postgres`
- ✅ Implement usecase interfaces with sqlc/pgx, row ↔ domain mapping, `WithTx`, DB error → domain error, filtering and pagination queries
- ❌ Business rules, returning sqlc types, importing `delivery`, HTTP knowledge, chaining repositories for logic

### `infrastructure`
- ✅ `config.Load()`, `logger.New()`, `database.Connect()`; fail fast at startup
- ❌ Business logic, `os.Getenv` outside `config/`, logging secrets, global singletons

### `db/migrations` and `db/queries`
- ✅ One change per migration with `Up` and `Down`; named sqlc queries; run `make sqlc` after changes
- ❌ **Editing an already-run migration**, editing generated Go code, business rules in SQL

---

## 4. Example: one request through every layer

`POST /api/v1/me/work-orders/{id}/complete` (technician completes a task)

1. **delivery/http**: middleware sets `request_id`, verifies the JWT, checks role `technician`. The handler parses the multipart form (`note`, `customer_approved`, `photos[]`) and checks format (1–5 JPG/PNG files).
2. **usecase** `CompleteWorkOrder.Execute`: loads the WO through `WorkOrderRepository`, checks that it belongs to this technician, saves the photos through `FileStorage`.
3. **domain** `wo.Complete(photos, approved, now)`: status must be `in_progress`, at least one photo, `customer_approved` must be true. Otherwise it returns a domain error.
4. **usecase**: inside a transaction, saves the WO and a `work_order_events` row, then logs `work_order.status_changed`.
5. **repository/postgres**: runs the sqlc queries and maps rows back to domain structs.
6. **delivery/http**: maps the result to JSON (`duration`, `completed_at`, `next_task`), or maps errors (`ErrInvalidTransition` → 409, `ErrPhotoRequired` → 422).

---

## 5. Where does this code go?

1. **If HTTP were replaced by a CLI, would this code still be needed?** No → `delivery`.
2. **If Postgres were replaced by MySQL, would this code change?** Yes → `repository`.
3. **Is this a rule the client/business asked for?** Yes → `domain` (the rule) or `usecase` (the order of steps).

---

## 6. Adding a new feature (checklist)

1. Migration in `db/migrations/` (if the schema changes)
2. Queries in `db/queries/` → `make sqlc`
3. Entity + rules in `domain/` (+ unit tests)
4. Usecase + interfaces in `usecase/` (+ unit tests with fakes)
5. Repository implementation in `repository/postgres/`
6. Handler + DTOs + route in `delivery/http/`
7. Wire it in `cmd/api/main.go`
8. Update `api/openapi.yaml`

---

## 7. Cross-cutting conventions

- **Errors**: `{ "error": { "code": "...", "message": "..." } }`
- **Pagination**: cursor-based, `?limit=15&cursor=<opaque>` → `{ data, next_cursor, prev_cursor }`; cursor = base64 of `(created_at, id)`
- **Time**: ISO 8601, UTC
- **Logging**: `slog` JSON with `request_id` on every line; never log passwords, tokens, or full personal data
- **Context**: every function that does I/O takes `ctx context.Context` as its first parameter
