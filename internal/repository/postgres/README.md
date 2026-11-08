# `internal/repository/postgres/` — the warehouse clerk

Implements the interfaces from `usecase/` using PostgreSQL. **Only this folder knows about SQL, pgx, and sqlc.**

sqlc-generated code will live in `postgres/sqlc/` (generated from `db/queries/*.sql`; never edit it by hand).

## ✅ Do

- Implement `FindByID`, `Save`, `List`, ... using sqlc / pgx
- Convert sqlc rows ↔ domain structs (`toDomain(row)`, `fromDomain(wo)`)
- Implement transactions (`WithTx(ctx, fn)`) for the `TxManager` interface
- Translate DB errors into domain errors (`pgx.ErrNoRows` → `domain.ErrNotFound`, unique violation → `domain.ErrEmailTaken`)
- Implement filtering and cursor pagination queries

## ❌ Don't

- Business rules (`if status == ... { return error }`). Store what you're given
- Return sqlc types to the usecase; always return domain structs
- Import `delivery/http`
- Log requests or know anything about HTTP
- Chain several repositories together to implement logic (that's the usecase's job)

## Example

```go
// internal/repository/postgres/work_order_repo.go
type WorkOrderRepo struct{ q *sqlc.Queries }

func (r *WorkOrderRepo) FindByID(ctx context.Context, id int64) (*domain.WorkOrder, error) {
    row, err := r.q.GetWorkOrder(ctx, id)
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    return toDomainWorkOrder(row), nil
}
```

## Testing

Integration tests against a real Postgres (the one from `docker-compose.yml`), not mocks.
