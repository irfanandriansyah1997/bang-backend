# Primary test cases checklist

For each new or changed function with logic, list the cases below that apply, then mark each one **tested ✅** or **missing ❌** by reading its `_test.go`. The goal is to show the developer which cases they tend to overlook, so explain each ❌.

Rules for writing tests: `.ai/rules/testing.md`.

## Every function

| Case | Ask |
|---|---|
| Happy path | Is the normal, valid input tested with an exact expected result? |
| Invalid input | Empty string, zero, negative, too long, wrong format, unknown enum value? |
| Edge / boundary | First and last valid values, one past each (e.g. 1 and 5 photos OK, 0 and 6 rejected); empty and nil collections |

## By layer

### `internal/domain` (rule methods)
| Case | Example |
|---|---|
| Every allowed transition | `open → assigned`, `assigned → in_progress`, `in_progress → completed` |
| **Every forbidden transition** | `open → completed`, `completed → anything`, `cancelled → anything` (beginners usually test only one) |
| Each rule that can say "no" | no photos → `ErrPhotoRequired`; `customer_approved=false` → error |
| State after success | `StartedAt` set to the given `now`; status changed |
| State after failure | Nothing changed when the method returned an error |

### `internal/usecase`
| Case | Example |
|---|---|
| Happy path with fake ports | Saves once, returns the right output |
| Not found | Fake repo returns `ErrNotFound` → usecase returns it (wrapped, still `errors.Is`-able) |
| Each dependency failing | Repo save fails, storage upload fails, tx fails → error returned, nothing half-written |
| Authorization | Another technician's work order → `ErrForbidden` |
| Domain rule rejection passes through | Domain returns `ErrInvalidTransition` → usecase returns it unchanged |
| Side effects | Event recorded / business log written exactly once on success, never on failure |

### `internal/delivery/http`
| Case | Example |
|---|---|
| Success status and body shape | 200/201 with the documented JSON |
| Malformed request | Bad JSON, missing required field → 422 `validation_failed` |
| Each domain error mapped | `ErrNotFound` → 404, `ErrForbidden` → 403, `ErrInvalidTransition` → 409, unknown → 500 |
| Auth | No token → 401; wrong role → 403 |
| Pagination (list endpoints) | Default limit, max limit, `next_cursor` present or absent |

### `internal/repository/postgres` (behind `TEST_DATABASE_URL`)
| Case | Example |
|---|---|
| Round trip | Save then find returns the same domain value |
| Not found | `pgx.ErrNoRows` → `domain.ErrNotFound` |
| Constraint violation | Duplicate email → `domain.ErrEmailTaken` |
| Pagination order and cursor | Second page starts after the cursor with no duplicates |

### `internal/infrastructure`
| Case | Example |
|---|---|
| Defaults | Missing optional env → default value |
| Required value missing | Missing `JWT_SECRET` → clear error |
| Output format | Logger writes valid JSON with the required fields |
| Failure | Unreachable DB → error, not panic |

## Reporting a missing case

For each ❌, give:
1. **Why it matters**, in one sentence: what bug would slip through.
2. **An example test row** to add to the existing table:

```go
{name: "from completed", status: StatusCompleted, wantErr: ErrInvalidTransition},
```

Missing happy-path or rule-violation cases are **major**. A function with logic and **no test at all** is a **blocker**. Other missing cases are **major** if they guard a business rule or money/data integrity, otherwise **minor**.
