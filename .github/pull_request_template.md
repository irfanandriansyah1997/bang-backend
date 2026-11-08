<!--
PR title format: <type>(<feature>): <message>
e.g. feat(customer): add customer CRUD endpoints
Delete any section that doesn't apply (write "None" rather than leaving it blank).
-->

## Summary
<!-- What does this PR do and why? 1–3 sentences. -->

## Related issue
<!-- "Closes #N" closes the issue on merge; "Refs #N" only links it. -->
Closes #

## Type of change
- [ ] `feat`: new feature / endpoint
- [ ] `fix`: bug fix
- [ ] `refactor`: no behavior change
- [ ] `test`: tests only
- [ ] `docs`: documentation only
- [ ] `build` / `ci` / `chore`: tooling, CI, dependencies

## Changes by layer
<!-- One line per touched layer. Remove the rows you didn't touch. -->
| Layer | Change |
|---|---|
| `domain` | |
| `usecase` | |
| `delivery/http` | |
| `repository/postgres` | |
| `infrastructure` | |
| `cmd/api` (wiring) | |

## API changes
<!-- New/changed endpoints, request/response shape, status codes, error codes. -->
| Method | Path | Role | Change |
|---|---|---|---|
| | | | |

- [ ] `api/openapi.yaml` updated
- [ ] Breaking change for web / Android clients (describe below)

## Database changes
<!-- Migrations, new tables/columns/indexes, seed data. -->
- Migrations:
- [ ] `Down` migration tested (up → down → up)
- [ ] Seed data updated if needed
- [ ] Safe on existing data (no destructive change without a plan)

## Config / environment
<!-- New env vars, changed defaults, new dependencies. Update .env.example. -->
None

## How to test
<!-- Steps a reviewer can run locally. Include curl examples for endpoints. -->
```bash
make check
```

## Observability
<!-- New logs / business events (e.g. work_order.status_changed). Confirm no secrets or personal data are logged. -->
- [ ] Logs include `request_id`; no passwords, tokens, or full personal data

## Risk & rollback
<!-- What could break? How do we undo it (revert commit, down migration)? -->

## Checklist
- [ ] `make check` passes locally (gofmt, vet, build, test)
- [ ] Tests added/updated for new logic (domain + usecase at minimum)
- [ ] Follows `docs/architecture.md`: imports point inward, no business rules in handlers/SQL, no sqlc types outside `repository/postgres`
- [ ] Errors use `{ "error": { "code", "message" } }`; list endpoints use cursor pagination
- [ ] Commits follow `<type>(<feature>): <message>`
- [ ] No secrets, `.env`, or credentials committed
- [ ] Docs / READMEs updated if behavior or setup changed
