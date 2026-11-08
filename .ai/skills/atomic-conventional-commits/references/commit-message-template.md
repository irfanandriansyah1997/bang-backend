# Commit message template

The Conventional Commits shape used by this skill.

> **bang-backend:** the subject must be `<type>(<feature>): <message>`. The feature is required, in lowercase kebab-case, and names the feature/area (e.g. `feat(docs): adding readme.md with setup steps`). "add" and "adding" are both accepted. The `commit-msg` hook rejects anything else. Where this file says scope is optional, that does not apply to this repo.

## Shape

```
<type>(<scope>): <imperative summary>
                                       <-- blank line
<body: what changed and why, wrap ~72 chars>
                                       <-- blank line
<footer: BREAKING CHANGE / Refs / Closes>
```

- **Subject** (first line): required, imperative mood, no trailing period, ≤ 50 chars ideal, hard cap 72.
- **Blank line** between subject and body, and between body and footer.
- **Body**: optional but recommended when the change is non-obvious. Explain *what* and *why*, not *how* — the diff already shows how.
- **Footer**: optional. Use `BREAKING CHANGE: <description>` for breaking changes. Use `Refs: <id>` or `Closes: <id>` if the repo uses issue trackers.

## Subject rules

- Imperative mood: "add", "fix", "remove" — not "added", "adds", "fixed".
- No trailing period.
- Present tense.
- Lowercase after the colon, unless the first word is a proper noun or acronym.
- No emoji unless the repo convention already uses them.

## Type

Pick exactly one:

| Type       | Use for                                                             |
|------------|---------------------------------------------------------------------|
| `feat`     | New user-facing capability                                          |
| `fix`      | Bug fix affecting user-facing behavior                              |
| `refactor` | Code change with no behavior change                                 |
| `perf`     | Performance improvement with no behavior change                     |
| `test`     | Adding or updating tests only                                       |
| `docs`     | Documentation only                                                  |
| `style`    | Formatting, whitespace — no logic change                            |
| `build`    | Build system, packaging, dependency bumps                           |
| `ci`       | CI configuration and pipeline scripts                               |
| `chore`    | Housekeeping that doesn't fit above                                 |
| `revert`   | Reverts a previous commit; body cites the reverted SHA              |

If recent history uses extra types (e.g. `security`), follow suit.

## Scope

- Derive from the primary module / package / feature touched.
- Match the casing / separator convention of recent commits — check `git log --oneline -n 20`.
- Prefer short, stable scopes (`auth`, `ui`, `api`, `db`) over paths (`src/services/auth`).
- Drop the scope if the change spans 3+ unrelated areas — usually a signal to split further.

## Body

Write the body when any of these are true:

- Change touches 2+ files
- Diff is longer than ~30 lines
- Intent is not obvious from the diff
- The change was driven by a bug report, incident, or ticket that adds context

Answer in this order:

1. **What** changed (one or two lines)
2. **Why** it changed (motivation, constraint, tradeoff)
3. **Notable non-obvious choices** (only if any)

Do not narrate the diff step by step. Do not paste in stack traces or logs unless they're the reason for the change.

## Footer

- `BREAKING CHANGE: <description>` — required when the change breaks an API or contract. Also acceptable as `!` after the type: `feat(api)!: drop v1 endpoints`.
- `Refs: <id>` — soft reference to an issue.
- `Closes: <id>` — closes the issue on merge, if the repo uses that convention.
- Do **not** add `Co-authored-by:` or `Generated-by:` trailers unless the user explicitly asks.

## Multi-line commits with HEREDOC

Use a HEREDOC to keep formatting exact and avoid shell quoting issues:

```bash
git commit -m "$(cat <<'EOF'
feat(auth): support refresh-token rotation

Rotate the refresh token on every use so a leaked token cannot be
replayed indefinitely. The old token is invalidated server-side once
the new one is issued.

Refs: AUTH-482
EOF
)"
```

## Worked examples

### Feature (single file, obvious intent)

```
feat(ui): add empty-state illustration to inbox
```

### Feature (needs a body)

```
feat(auth): support refresh-token rotation

Rotate the refresh token on every use so a leaked token cannot be
replayed indefinitely. The old token is invalidated server-side once
the new one is issued.

Refs: AUTH-482
```

### Bug fix

```
fix(payments): guard against null merchant id in receipt builder

`ReceiptBuilder` assumed `merchantId` was always present, which is not
true for guest checkout. Fall back to the guest tag instead of
throwing.

Closes: PAY-311
```

### Refactor (no behavior change)

```
refactor(auth): extract TokenValidator from AuthService

Move token parsing and expiry checks into a dedicated collaborator so
`AuthService` only orchestrates. No behavior change; existing tests
continue to pass unchanged.
```

### Breaking change

```
feat(api)!: drop v1 payments endpoints

BREAKING CHANGE: `/v1/payments/*` endpoints have been removed. All
callers must migrate to `/v2/payments/*`, which uses the new envelope
schema documented in docs/api/payments-v2.md.
```

### Revert

```
revert(ui): revert "feat(ui): highlight active tab with accent color"

This reverts commit 1a2b3c4d5e6f. The accent color fails WCAG AA
contrast on the light theme; will retry with a corrected palette.
```
