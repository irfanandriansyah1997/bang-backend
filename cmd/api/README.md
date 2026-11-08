# `cmd/api/` — the power button + the wiring

This folder holds `main.go`, the **entry point** of the API server. It runs when you do `go run ./cmd/api`.

Its only job is to **start the app and connect all the pieces** (this is called the *composition root*). It is the only place that is allowed to know about every package.

> React analogy: `main.tsx` + your `<Providers>` tree. It assembles everything but contains no logic.

Go convention: every executable program lives in `cmd/<program-name>/`. If we later need another program (e.g. a seeder), we add `cmd/seed/`.

## ✅ Do

- Load config, create the logger, open the DB connection
- Build repositories → inject them into usecases → inject usecases into handlers
- Register the router and start the HTTP server
- Handle graceful shutdown (SIGINT / SIGTERM)
- Exit with an error (`os.Exit(1)`) if startup fails

## ❌ Don't

- Write SQL queries
- Write business rules (`if status == ...`)
- Write HTTP handlers here
- Let `main.go` grow huge (aim for < 100 lines; move setup into helper functions)

## Example

```go
// cmd/api/main.go (sketch)
cfg := config.Load()
db := database.Connect(cfg.DatabaseURL)

woRepo := postgres.NewWorkOrderRepo(db)        // repository
woUC   := usecase.NewCompleteWorkOrder(woRepo) // usecase
h      := http.NewWorkOrderHandler(woUC)       // delivery

server.Run(cfg.Port, h)
```

## Import rule

`cmd/api` may import **everything**. That is exactly its job.
