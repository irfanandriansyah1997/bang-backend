# `internal/domain/` — the business brain (kept sterile)

Entities (structs) and business rules that are true **no matter where the data comes from**: HTTP, CLI, or a test.

> React analogy: a pure TypeScript function like `canTransition(from, to)` that you can test without rendering any component.

## ✅ Do

- Entity structs: `WorkOrder`, `User`, `Customer`, `SupportTicket`
- Status / enum types: `StatusOpen`, `StatusAssigned`, `StatusInProgress`, ...
- Methods that enforce rules: `wo.Assign(techID)`, `wo.Start(now)`, `wo.Complete(photos, approved, now)`
- Business errors: `ErrInvalidTransition`, `ErrPhotoRequired`, `ErrNotFound`, `ErrForbidden`
- Import only light standard library packages (`time`, `errors`, `strings`)
- Accept `now time.Time` as a parameter instead of calling `time.Now()`, so tests are predictable

## ❌ Don't

- Import `net/http`, `pgx`, sqlc, JWT, or `slog`
- Add `json:"..."` or `db:"..."` struct tags (that's the outer layers' concern)
- Touch the database or the file system
- Know about HTTP status codes (409, 422). Return a domain error; `delivery` maps it

## Example

```go
// internal/domain/work_order.go
var ErrInvalidTransition = errors.New("invalid status transition")

func (w *WorkOrder) Start(now time.Time) error {
    if w.Status != StatusAssigned {
        return ErrInvalidTransition
    }
    w.Status = StatusInProgress
    w.StartedAt = &now
    return nil
}
```

## Testing

Pure unit tests, no DB, no server. Table-driven tests are ideal for the status state machine.
