---
description: Go style conventions for this repo
globs: cmd/**/*.go,internal/**/*.go
alwaysApply: false
paths:
  - "cmd/**/*.go"
  - "internal/**/*.go"
---

# Go style

- Format with `gofmt`; keep `go vet ./...` clean.
- Return errors, don't panic. Wrap with context: `fmt.Errorf("find work order %d: %w", id, err)`.
- Compare errors with `errors.Is` / `errors.As`, never by string.
- Small interfaces (1–3 methods), defined by the package that uses them.
- Constructors are named `NewXxx` and return concrete types.
- Accept `ctx context.Context` as the first parameter of any function doing I/O.
- No package-level mutable state; pass dependencies in through constructors.
- Tests live next to the code (`xxx_test.go`); prefer table-driven tests.
- Package names are short, lowercase, singular (`workorder`, not `workOrders`).
