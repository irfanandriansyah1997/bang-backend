---
description: Rules for the domain and usecase layers (business logic)
globs: internal/domain/**,internal/usecase/**
alwaysApply: false
paths:
  - "internal/domain/**"
  - "internal/usecase/**"
---

# Domain and usecase rules

## Domain (`internal/domain`)

- Imports the **standard library only** (`errors`, `time`, `strings`, ...). No `net/http`, `pgx`, sqlc, JWT, `slog`, or other app packages.
- No struct tags (`json:"..."`, `db:"..."`); those belong to the outer layers.
- Business rules are methods on the entity: `wo.Start(now)`, `wo.Complete(photos, approved, now)`. They return a domain error instead of `bool`.
- Business errors are package-level values: `var ErrInvalidTransition = errors.New("invalid status transition")`. Callers check them with `errors.Is`.
- **Never call `time.Now()` here.** Take `now time.Time` as a parameter so tests are predictable.
- No I/O of any kind: no DB, files, network, or logging.

## Usecase (`internal/usecase`)

- One usecase per user action (`CompleteWorkOrder`, `DeactivateTechnician`), with its own `Input` / `Output` structs.
- **Defines the interfaces (ports) it needs**, small (1–3 methods): `WorkOrderRepository`, `FileStorage`, `TxManager`, `Clock`. `repository/postgres` implements them.
- `ctx context.Context` is the first parameter of every method; pass it on to every port call.
- Orders the steps and checks business authorization ("does this work order belong to this technician?"); the rule itself stays in domain.
- Multi-step writes that must succeed or fail together go through `TxManager.WithTx(ctx, fn)`.
- Wrap errors from ports with context: `fmt.Errorf("find work order %d: %w", id, err)`. Return domain errors unchanged so the HTTP mapper can recognize them.
- Logs **business events** only (`work_order.status_changed` with `logger.SDC{"work_order_id": ..., "from": ..., "to": ...}`), not request logs.
- Never imports `net/http`, `delivery`, `repository`, `pgx`, or sqlc types; never writes SQL; never chooses HTTP status codes.
- No package-level mutable state; dependencies come in through the constructor (`NewCompleteWorkOrder(repo, storage, clock, log)`).
