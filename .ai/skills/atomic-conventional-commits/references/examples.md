# Worked examples

Each example shows a mixed working tree and the atomic commit plan the skill should produce.

---

## Example 1 — Feature + refactor + tests + docs

### Before

`git status`:

```
 M src/auth/AuthService.kt
 M src/auth/AuthServiceTest.kt
?? src/auth/TokenValidator.kt
?? src/auth/RefreshPolicy.kt
 M docs/auth/refresh-tokens.md
```

Inside `AuthService.kt` the diff has two logical groups:

- Extract `TokenValidator` (moves parsing / expiry checks out).
- Add `refresh-token rotation` (new behavior, uses `RefreshPolicy`).

### Plan (4 commits)

```
1. refactor(auth): extract TokenValidator from AuthService
   files: src/auth/AuthService.kt (partial), src/auth/TokenValidator.kt

2. feat(auth): support refresh-token rotation
   files: src/auth/AuthService.kt (partial), src/auth/RefreshPolicy.kt

3. test(auth): cover refresh-token rotation edge cases
   files: src/auth/AuthServiceTest.kt

4. docs(auth): document refresh-token flow
   files: docs/auth/refresh-tokens.md
```

Staging uses hunk-level staging on `src/auth/AuthService.kt` for commits 1 and 2 (a patch + `git apply --cached` for agents, `git add -p` for humans), so different hunks in the same file go to different commits.

Order matters: the refactor lands first so the feature commit only shows the new-behavior hunks.

---

## Example 2 — Bug fix + unrelated formatting

### Before

`git status`:

```
 M src/payments/ReceiptBuilder.kt
 M src/payments/ReceiptRenderer.kt
```

`ReceiptBuilder.kt` fixes a null-safety bug. `ReceiptRenderer.kt` is only reformatted by an editor auto-formatter.

### Plan (2 commits)

```
1. fix(payments): guard against null merchant id in receipt builder
   files: src/payments/ReceiptBuilder.kt

2. style(payments): reformat ReceiptRenderer
   files: src/payments/ReceiptRenderer.kt
```

Do **not** ship the formatting churn in the same commit as the fix — it would make reverting the fix noisy.

---

## Example 3 — Dependency bump that unblocks a feature

### Before

`git status`:

```
 M build.gradle.kts
 M gradle/libs.versions.toml
 M src/ui/HomeScreen.kt
?? src/ui/NewChartCard.kt
```

The Gradle changes bump a charting library to a version that ships `NewChartCard`. `HomeScreen.kt` uses the new card.

### Plan (2 commits)

```
1. build(deps): bump charts library to 2.4.0
   files: build.gradle.kts, gradle/libs.versions.toml

2. feat(home): show revenue chart on home screen
   files: src/ui/HomeScreen.kt, src/ui/NewChartCard.kt
```

The bump lands first so the tree still builds after commit 1. The feature commit is a clean, focused change.

---

## Example 4 — When to keep things together

### Before

`git status`:

```
 M src/config/FeatureFlags.kt
 M src/checkout/CheckoutController.kt
```

Adding a new checkout flow gated by a feature flag. The controller change is unusable without the flag entry, and vice versa.

### Plan (1 commit)

```
1. feat(checkout): add express-checkout flow behind feature flag
   files: src/config/FeatureFlags.kt, src/checkout/CheckoutController.kt
```

Splitting would leave one intermediate commit that either can't compile or ships dead code. Keep it as one commit and mention the flag in the body.

---

## Example 5 — Revert

### Before

The most recent commit `1a2b3c4d5e6f` introduced a color change that fails contrast checks.

### Plan (1 commit)

```
1. revert(ui): revert "feat(ui): highlight active tab with accent color"

   This reverts commit 1a2b3c4d5e6f. The accent color fails WCAG AA
   contrast on the light theme; will retry with a corrected palette.
```

Prefer a new revert commit over `git reset` on a branch that others may already have pulled.

---

## Example 6 — Go feature across clean-architecture layers (this repo)

### Before

`git status`:

```
?? db/migrations/00002_create_customers.sql
?? db/queries/customers.sql
?? internal/repository/postgres/sqlc/customers.sql.go
?? internal/domain/customer.go
?? internal/domain/customer_test.go
?? internal/usecase/customer.go
?? internal/usecase/customer_test.go
?? internal/repository/postgres/customer_repo.go
?? internal/delivery/http/customer_handler.go
 M internal/delivery/http/router.go
 M cmd/api/main.go
 M api/openapi.yaml
```

The work implements GitHub issue #12 (customer CRUD).

### Plan (7 commits)

```
1. feat(customer): add customers table migration
   files: db/migrations/00002_create_customers.sql

2. feat(customer): add customer queries
   files: db/queries/customers.sql, internal/repository/postgres/sqlc/customers.sql.go

3. feat(customer): add Customer entity and validation
   files: internal/domain/customer.go, internal/domain/customer_test.go

4. feat(customer): add create/list/update usecases
   files: internal/usecase/customer.go, internal/usecase/customer_test.go

5. feat(customer): add postgres customer repository
   files: internal/repository/postgres/customer_repo.go

6. feat(customer): expose /customers endpoints
   files: internal/delivery/http/customer_handler.go, internal/delivery/http/router.go, cmd/api/main.go
   footer: Refs: #12

7. docs(customer): document customer endpoints in openapi
   files: api/openapi.yaml
   footer: Closes: #12
```

Every commit uses the same feature name, `customer`, so `git log --oneline | grep "(customer)"` shows the whole feature. Each commit compiles because the layers land inward-out (DB → domain → usecase → adapters). Generated sqlc code is committed together with the query that produced it. Wiring in `main.go` rides along with the handler, because the route is dead code without it.

---

## Anti-example — "wip: various fixes"

`git status` shows unrelated edits across 8 files. The user prompts "commit this". The wrong response is:

```bash
git add -A
git commit -m "wip: various fixes"
```

The right response is to run this skill: read the diff, propose 3–5 atomic commits, and get approval before staging anything.
