# Go review checklist

Things to look for in the diff. Strict lint (`.golangci.yml`) already catches many mechanical problems; this list focuses on what a reviewer still has to judge. Each item: what to look for → why it matters → how to fix, explained for a Go beginner with a React/TS background.

## Errors

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| `if err != nil { log...; return nil }`: error logged then swallowed | The caller thinks it succeeded; the bug surfaces far away | Return the error (wrapped); log once at the edge (HTTP error mapper) | major |
| `return err` from another package with no context | "connection refused" with no clue which step failed | `fmt.Errorf("ping database: %w", err)` | major |
| `fmt.Errorf("...: %v", err)` | `%v` loses the original error, so `errors.Is` stops working | Use `%w` | major |
| `err.Error() == "..."` or `strings.Contains(err.Error(), ...)` | Breaks when the message changes | `errors.Is(err, domain.ErrX)` / `errors.As` | major |
| `panic(...)` in non-`main` code | Crashes the whole server for one bad request | Return an error | blocker |
| Domain error turned into a generic error before reaching HTTP | The mapper can't return 404/409; everything becomes 500 | Wrap with `%w` so `errors.Is` still matches | major |

## Context

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| I/O function without `ctx context.Context` as the first parameter | Can't be cancelled or time out; request IDs don't reach logs | Add `ctx` first and pass it down | major |
| `context.Background()` / `context.TODO()` inside handlers or usecases | Drops the request's deadline and `x_request_id` | Use the `ctx` you were given | major |
| `context.WithTimeout` / `WithCancel` without `defer cancel()` | Leaks a timer and goroutine | `defer cancel()` right after creating it | major |
| Logging with a context-free call where a `ctx` exists | `x_request_id` / SDC missing from that line | `log.Info(ctx, ...)`: our logger always takes `ctx` | minor |

## Resources

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| `rows` / `resp.Body` / files / tx not closed on every path | Connection pool runs dry; server hangs | `defer rows.Close()` right after the error check | major |
| `rows.Err()` not checked after the loop | A mid-stream DB error looks like "no more rows" | `if err := rows.Err(); err != nil { ... }` | major |
| Transaction with an early `return` that skips `Rollback` | Locks stay held | `defer tx.Rollback(ctx)` (no-op after commit), or use `TxManager.WithTx` | major |

## Data and security

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| SQL built with `+` or `fmt.Sprintf` from user input | SQL injection | Parameters (`$1`) via sqlc / pgx | blocker |
| Password, token, JWT secret, full phone/address in logs, `Payload`, `SDC`, or error messages | Leaks personal data and credentials into Kibana | Log IDs only; redact | blocker |
| Hard-coded secret or real credential in code or `.env.example` | Public repo, so it's leaked | Env var + placeholder in `.env.example` | blocker |
| Comparing passwords or tokens with `==` | Timing attacks | `bcrypt.CompareHashAndPassword` / `subtle.ConstantTimeCompare` | major |
| Uploaded file type trusted from the filename | Anything can be uploaded | Check size and `http.DetectContentType` | major |
| SDC value built from a number without converting | Won't compile (SDC is string-only), or a mapping conflict if someone works around it | `strconv.Itoa(id)` | minor |

## Concurrency

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| Goroutine started without a way to stop it (`ctx`, `WaitGroup`) | Leaks on shutdown | Tie it to `ctx` and wait for it | major |
| Map or slice written from several goroutines | Data race, random crashes (`go test -race` may catch it) | `sync.Mutex` or channels | blocker |
| Package-level variables that change at runtime | Shared hidden state; tests interfere with each other | Pass dependencies through constructors | major |

## Go idioms

| Look for | Why it matters | Fix | Severity |
|---|---|---|---|
| Interface defined next to its implementation (in `repository/`) | Couples usecase to the database package | Define it in `usecase`, where it's used | major |
| Big interfaces (5+ methods) | Hard to fake in tests | Split into small ones | minor |
| Returning an interface from a constructor | Hides the concrete type for no reason | Return the concrete struct (`*Repo`) | minor |
| Pointer to a pointer, or pointer to an interface | Almost always a mistake | Use the value or the interface directly | major |
| `time.Now()` inside domain | Tests depend on the wall clock | Take `now time.Time` as a parameter | major |
| Exported identifiers without a doc comment | `revive` will flag it; readers don't know the contract | One-line comment starting with the name | minor |
