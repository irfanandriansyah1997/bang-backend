# `db/queries/` — SQL for sqlc

Hand-written, named SQL queries. **sqlc** reads these files plus the migrations and generates type-safe Go functions into `internal/repository/postgres/sqlc/`.

> React analogy: like `openapi-typescript` generating types from a contract. You write the source (SQL), the tool writes the typed code.

## ✅ Do

- One file per table/resource: `work_orders.sql`, `users.sql`, `customers.sql`
- Name each query with an annotation:

  ```sql
  -- name: GetWorkOrder :one
  SELECT * FROM work_orders WHERE id = $1;

  -- name: ListWorkOrders :many
  SELECT * FROM work_orders
  WHERE (created_at, id) < ($1, $2)
  ORDER BY created_at DESC, id DESC
  LIMIT $3;
  ```

- Use `:one`, `:many`, `:exec`, `:execrows` correctly
- Run `make sqlc` after every change and commit the generated code

## ❌ Don't

- Write Go code here
- Edit the generated Go files by hand (they are overwritten on the next `make sqlc`)
- Put business rules in SQL (e.g. deciding valid status transitions). Keep that in `domain`
