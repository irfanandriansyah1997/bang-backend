# `internal/usecase/` — the director of the flow

One usecase = **one action a user performs**: `CompleteWorkOrder`, `Login`, `DeactivateTechnician`, ...
It decides the **order of steps**, but not the technical details of how data is stored or how requests arrive.

> React analogy: a custom hook that receives its `api` through props/context. In tests you pass a fake one.

## ✅ Do

- **Define the interfaces (ports) it needs**: `WorkOrderRepository`, `FileStorage`, `TxManager`, `Clock`
  - In Go, interfaces belong to the *consumer*, so they live here, not in `repository/`
- Call domain methods for business rules
- Check business authorization ("does this WO belong to the logged-in technician?")
- Open transactions through an interface (`uc.tx.WithTx(ctx, func(ctx) error { ... })`)
- Log business events (`work_order.status_changed`)
- Take `context.Context` as the first parameter of every method
- Define its own `Input` / `Output` structs

## ❌ Don't

- Import `net/http`, use `*http.Request`, read headers/cookies/JSON
- Import `pgx` or sqlc-generated types
- Write SQL
- Decide HTTP status codes
- Keep mutable global state (package-level variables that change)
- Return sqlc structs to the handler

## Example

```go
// internal/usecase/complete_work_order.go
type WorkOrderRepository interface { // interface defined HERE, implemented in repository/postgres
    FindByID(ctx context.Context, id int64) (*domain.WorkOrder, error)
    Save(ctx context.Context, wo *domain.WorkOrder) error
}

type CompleteWorkOrder struct {
    repo  WorkOrderRepository
    clock Clock
}

func (uc *CompleteWorkOrder) Execute(ctx context.Context, in Input) (*Output, error) {
    wo, err := uc.repo.FindByID(ctx, in.WorkOrderID)
    if err != nil {
        return nil, err
    }
    if wo.TechnicianID != in.UserID {
        return nil, domain.ErrForbidden
    }
    if err := wo.Complete(in.Photos, in.CustomerApproved, uc.clock.Now()); err != nil {
        return nil, err
    }
    if err := uc.repo.Save(ctx, wo); err != nil {
        return nil, err
    }
    return toOutput(wo), nil
}
```

## Testing

Unit test with **fake** implementations of the interfaces (a simple struct in the test file). No DB needed.
