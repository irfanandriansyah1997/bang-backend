---
name: add-feature
description: Add a new backend feature or endpoint end-to-end following this repo's clean architecture (migration, sqlc query, domain, usecase, repository, HTTP handler, wiring, OpenAPI). Use when the user asks to implement a GitHub issue, a new endpoint, or a new business rule.
---

# Add a feature

Follow `docs/architecture.md` and the folder READMEs. Work through the steps in order and skip the ones the feature doesn't need.

1. **Read the issue / spec**: the acceptance criteria in the GitHub issue and the matching section of `docs/superpowers/specs/rencana-belajar-fullstack-fde.md`.
2. **Migration** in `db/migrations/` (only if the schema changes).
3. **Queries** in `db/queries/`, then run `make sqlc`.
4. **Domain** in `internal/domain/`: entity, rule methods, business errors, plus table-driven unit tests.
5. **Usecase** in `internal/usecase/`: define the interfaces it needs, the `Input`/`Output` structs, the flow, plus unit tests with fake implementations.
6. **Repository** in `internal/repository/postgres/`: implement the interfaces, map rows ↔ domain, map DB errors.
7. **HTTP** in `internal/delivery/http/`: DTOs, a thin handler, the route, and error mapping.
8. **Wire** it in `cmd/api/main.go`.
9. **Contract**: update `api/openapi.yaml`.
10. **Verify**: run `go vet ./...` and `go test ./...`, and hit the endpoint with curl.

## When done

Explain what was added in each layer in 1–2 lines per layer, comparing it to a React concept where it helps. The developer is new to Go.
