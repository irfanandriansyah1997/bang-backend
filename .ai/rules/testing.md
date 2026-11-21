---
description: How to write Go tests in this repo and which cases every test must cover
globs: cmd/**/*_test.go,internal/**/*_test.go
alwaysApply: false
paths:
  - "cmd/**/*_test.go"
  - "internal/**/*_test.go"
---

# Testing rules

## Primary cases (must be covered for every new function with logic)

| Case | Question to ask | Example |
|---|---|---|
| Happy path | Does the normal input give the expected result? | `Start()` from `assigned` → `in_progress` |
| Invalid input | What happens with empty, zero, wrong format, or out-of-range values? | empty `JWT_SECRET`, `limit=0`, unknown status |
| Business rule violations | Is every rule that can say "no" tested? | `Complete()` without photos → `ErrPhotoRequired` |
| Dependency errors | What if the repository / storage / clock returns an error? | fake repo returns `ErrNotFound` → usecase returns it wrapped |
| Edge cases | Boundaries, nil/empty collections, duplicates, the last item | 1 photo and 5 photos OK, 0 and 6 rejected |
| Authorization (usecase) | Can the wrong user do it? | technician completes someone else's work order → `ErrForbidden` |

Not every row applies to every function; skip a row only when it truly can't happen.

## How to write them

- **Table-driven tests** with `t.Run(tt.name, ...)`; one row per case, name rows after the case (`"missing secret"`, not `"case 3"`).
- Test names: `TestThing_Condition` (`TestLoadFrom`, `TestWorkOrder_Start_FromOpen`).
- Assert with `got` / `want`; on failure print both: `t.Errorf("Port = %q, want %q", got, want)`.
- Compare errors with `errors.Is(err, domain.ErrX)`, never by message text.
- Use `t.Fatal` / `t.Fatalf` when the rest of the test can't continue, `t.Errorf` otherwise.
- Helpers call `t.Helper()` first.
- Use `t.TempDir()` for files and `t.Setenv()` for env vars; never touch real paths or the real environment.
- **No `time.Sleep`.** Pass `now` or a fake clock instead.
- **Fakes over mocks:** a small struct in the test file that implements the usecase's interface.
- Tests that need Postgres run only when `TEST_DATABASE_URL` is set; otherwise `t.Skip(...)`. CI must pass without a database.
- Keep tests next to the code (`x.go` → `x_test.go`, same package).

## Where tests are required

- `internal/domain`: every rule method. No DB, no HTTP.
- `internal/usecase`: every usecase, using fake ports.
- `internal/delivery/http`: status code + error code per error path (`httptest`).
- `internal/repository/postgres`: integration tests behind `TEST_DATABASE_URL`.
- `internal/infrastructure`: config parsing, logger output, connection errors.
