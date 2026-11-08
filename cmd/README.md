# `cmd/` — executable programs

Go convention: each sub-folder here is **one program** with its own `main.go`, built with `go build ./cmd/<name>`.

| Folder | Program |
|---|---|
| [`api/`](api/README.md) | The HTTP API server |

Add a new folder only for a new program (e.g. `cmd/seed/` for a data seeder). Keep each `main.go` thin: it wires things together and calls code from `internal/`.
