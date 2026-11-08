<!--
PR title format: <type>(<feature>): <message>
e.g. feat(customer): add customer CRUD endpoints
Keep every section short and in bullet points. Write "None" for sections that don't apply.
-->

## Summary
<!-- First bullet = kind of change (feat / fix / refactor / test / docs / build / ci / chore) + one-line why.
     Then one bullet per notable change. -->
- **Type:** feat: <why this PR exists, in one line>
- <change 1>
- <change 2>

## Related issue
<!-- "Closes #N" closes the issue on merge; "Refs #N" only links it. -->
Closes #

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
- None

## How to test
<!-- Steps a reviewer can run locally. Include curl examples for endpoints. -->
```bash
make check
```

## Observability
<!-- New logs / business events (e.g. work_order.status_changed). -->
- [ ] Logs include `request_id`; no passwords, tokens, or full personal data

## Risk & rollback
<!-- What could break? How do we undo it (revert commit, down migration)? -->
- Risk:
- Rollback:

## Checklist
- [ ] `make check` passes locally (gofmt, vet, build, test)
- [ ] Tests added/updated for new logic (domain + usecase at minimum)
- [ ] Follows `docs/architecture.md`: imports point inward, no business rules in handlers/SQL, no sqlc types outside `repository/postgres`
- [ ] Errors use `{ "error": { "code", "message" } }`; list endpoints use cursor pagination
- [ ] Commits follow `<type>(<feature>): <message>`
- [ ] No secrets, `.env`, or credentials committed
- [ ] Docs / READMEs updated if behavior or setup changed
