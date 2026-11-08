---
name: review-task
description: Cross-check a pull request against its GitHub issue. Verify every acceptance criterion and Definition of Done item with evidence (files, commands, CI), report what's done and what's missing, then, with the user's OK, comment on the PR or tick the issue checklist. Use when the user gives a PR link and an issue link, or asks "cek PR ini sudah sesuai task?", "review task #N", "apakah issue sudah selesai di PR ini", "cross check PR".
---

# Review a task (PR ↔ issue)

Decide, with evidence, whether a PR fully delivers its issue. You are a strict but friendly reviewer for a developer who is new to Go. Report in the language the user writes in (usually Indonesian). Write PR comments in English.

## Inputs

1. **PR**: URL or number (`https://github.com/irfanandriansyah1997/bang-backend/pull/N`)
2. **Issue**: URL or number

If only the PR is given, use its linked issues (`closingIssuesReferences`). If neither links an issue, ask the user for one.

## Workflow

```
- [ ] 1. Collect the requirements (issue)
- [ ] 2. Collect the evidence (PR: diff, commits, CI, body)
- [ ] 3. Verify locally on the PR head
- [ ] 4. Judge every item → ✅ / ❌ / ⚠️
- [ ] 5. Report to the user
- [ ] 6. With the user's OK: comment on the PR, or tick the issue checklist
```

### 1. Collect the requirements

```bash
gh issue view <ISSUE> -R irfanandriansyah1997/bang-backend --json number,title,state,body,labels
```

Extract every checkbox under **Acceptance criteria** and **Definition of Done**. If the issue has no DoD, use the standard one:

- `make check` passes (gofmt, vet, build, test)
- Tests for new logic (domain + usecase at minimum)
- Follows `docs/architecture.md`
- `api/openapi.yaml` updated when endpoints change
- Commits follow `<type>(<feature>): <message>`; `Closes #N` in a commit or the PR body
- CI **Lint** and **Test** green

### 2. Collect the evidence

```bash
R=irfanandriansyah1997/bang-backend
gh pr view <PR> -R $R --json number,title,state,body,headRefName,baseRefName,mergeStateStatus,closingIssuesReferences,files,commits,statusCheckRollup
gh pr diff <PR> -R $R
gh pr checks <PR> -R $R
```

### 3. Verify locally on the PR head

Use a temporary worktree so the user's checkout is never touched:

```bash
WT=<scratchpad>/pr-<PR>
git fetch origin pull/<PR>/head:review/pr-<PR>
git worktree add "$WT" review/pr-<PR>
(cd "$WT" && make check && go run ./cmd/api)   # plus any command an item needs
git worktree remove --force "$WT"
git branch -D review/pr-<PR>
```

Run the specific commands each criterion implies (e.g. `curl` an endpoint, `goose up/down`, `go test -run`). If something needs services that aren't running (Postgres), say which items you could not verify and why. Never guess a result.

Architecture checks on the changed Go files:

```bash
# domain must not import external or app packages
grep -rnE '"(net/http|github.com/jackc/pgx|log/slog)|bang-backend/internal/(usecase|delivery|repository|infrastructure)' internal/domain
# usecase must not import delivery, repository, pgx, net/http
grep -rnE '"net/http"|pgx|bang-backend/internal/(delivery|repository)' internal/usecase
# no sqlc types outside repository/postgres
grep -rn 'postgres/sqlc' --include='*.go' . | grep -v '^./internal/repository/postgres/'
# os.Getenv only in config
grep -rn 'os.Getenv' --include='*.go' . | grep -v 'internal/infrastructure/config/'
```

Commit format: run every commit subject in the PR through `scripts/hooks/commit-msg.sh`.

### 4. Judge every item

| Mark | Meaning |
|---|---|
| ✅ | Verified with evidence (file + line, command output, CI result) |
| ❌ | Missing or wrong: say exactly what is missing and where it should go |
| ⚠️ | Could not verify (needs a service, ambiguous criterion) or partly done |

Also flag, even if not in the checklist:

- A change that is **out of scope** for the issue (belongs to another issue)
- Missing tests for new domain/usecase logic
- Secrets, debug prints, commented-out code, TODOs without an issue
- PR template sections left as unfilled placeholders

### 5. Report to the user

```markdown
## Review: PR #<PR> ↔ Issue #<ISSUE>
**Verdict:** ✅ Ready to merge / ❌ Not done yet (<k> items missing) / ⚠️ Needs your check

| # | Item | Evidence | Status |
|---|---|---|---|
| 1 | <criterion> | <file:line / command → output / CI> | ✅ |
| 2 | <criterion> | <what is missing> | ❌ |

### What to do next
1. <concrete fix, with file and a short snippet if the user needs it>
2. ...

### Also noticed
- <out-of-scope / quality notes, if any>
```

If anything is ❌, tell the user clearly that the task is **not done yet** and what to do, step by step, so they can fix it themselves.

### 6. Follow-up actions (ask first, each one)

Never post, edit, approve, or merge without the user's explicit OK in chat.

- **Not done**: offer to post the review on the PR:

  ```bash
  gh pr comment <PR> -R irfanandriansyah1997/bang-backend --body-file <report.md>
  ```

  Comment shape: a verdict line, then the table limited to ❌ / ⚠️ items, then the "What to do next" list. Use English and keep it short.

- **All done**: offer to tick the verified boxes on the issue (edit the body; only tick ✅ items) and add a short validation table, as was done for #1. Leave items that depend on the merge (e.g. "merged to main", "moved to Done") unticked with a note.

- Do **not** use `gh pr review --approve` or merge the PR. That stays with the user.

After a fix is pushed, the user can run this skill again on the same PR.
