---
name: check-pr
description: Decide whether a pull request's code is ready to merge. Lint and unit tests act as a gate, then a review covers architecture, common Go mistakes, primary test cases (tested vs missing), and repo conventions, ending in a verdict of Ready, Needs changes, or Blocked, with every finding explained for a Go beginner. Use when the user asks "cek PR ini", "review PR #N", "is this PR ready to merge", "check-pr", or when the @claude-code-review workflow runs it with --ci.
# Cheapest Claude model; enough for this workflow. Applies for the rest of the turn.
model: haiku
---

# Check a PR

Answer one question with evidence: **is this PR's code good enough to merge?** (Whether it finishes its GitHub issue is a different question; that's `review-task`.)

The developer is **new to Go** (strong React/TypeScript). Every finding must say **what** is wrong, **why** it matters, and **how** to fix it, with a short code example. Be strict on correctness and kind in tone.

## Inputs and modes

```
/check-pr <PR number or URL>          local mode
/check-pr <PR number> --ci            CI mode (called by .github/workflows/claude-review.yml)
```

With no PR given in local mode, use the open PR for the current branch (`gh pr view --json number`).

| | Local mode | CI mode (`--ci`) |
|---|---|---|
| Gate (lint + test) | You run it, in a temp worktree | Already done: the workflow only starts you when Lint and Test are green on the PR head |
| Code | Temp worktree of the PR head | The runner's checkout (already the PR head) |
| Output | Report in chat, in the user's language | **Post to the PR automatically**, in English |
| Posting | Only if the user asks | Always: one inline comment per finding (with code suggestions) + one summary comment |

## Workflow

```
- [ ] 1. Collect the PR (diff, files, commits, body, linked issue)
- [ ] 2. Gate: lint + unit test (local mode only)
- [ ] 3. Review: a. architecture  b. Go mistakes  c. primary test cases  d. repo conventions
- [ ] 4. Grade findings and pick the verdict
- [ ] 5. Report (local: chat; CI: summary comment + inline comments)
```

### 1. Collect the PR

```bash
R=irfanandriansyah1997/bang-backend
gh pr view <PR> -R $R --json number,title,body,headRefName,headRefOid,baseRefName,files,commits,closingIssuesReferences
gh pr diff <PR> -R $R
```

Focus on **changed lines**, plus the code they call when needed to judge them. Don't review untouched files.

### 2. Gate (local mode only)

Use a temp worktree so the user's checkout is never touched:

```bash
WT=<scratchpad>/check-pr-<PR>
git fetch -q origin pull/<PR>/head:check/pr-<PR>
git worktree add -q "$WT" check/pr-<PR>
(cd "$WT" && make lint && make test)
```

If either fails, **stop here**. Report each failure in plain words: the error, what it means, and the fix, with a code example. Then clean up (`git worktree remove --force "$WT"; git branch -D check/pr-<PR>`). Don't review code that doesn't pass lint and tests.

If lint and test both pass, keep the worktree for step 3 and remove it at the end.

### 3. Review

Read the rules first: `docs/architecture.md`, `AGENTS.md`, and the `.ai/rules/*.md` files whose `paths` match the changed files.

**a. Architecture.** Run these on the changed packages; any hit is a finding:

```bash
# domain imports only the standard library
grep -rnE '"(net/http|github.com/jackc/pgx|log/slog)|bang-backend/internal/(usecase|delivery|repository|infrastructure)' internal/domain
# usecase doesn't import delivery, repository, pgx, net/http
grep -rnE '"net/http"|pgx|bang-backend/internal/(delivery|repository)' internal/usecase
# sqlc types stay inside repository/postgres
grep -rn 'postgres/sqlc' --include='*.go' . | grep -v '^./internal/repository/postgres/'
# os.Getenv only in config
grep -rn 'os.Getenv' --include='*.go' . | grep -v 'internal/infrastructure/config/'
# struct tags in domain
grep -rnE '`(json|db):' internal/domain
```

Also check by reading: thin handlers, business rules only in domain/usecase, dependencies wired only in `cmd/api/main.go`.

**b. Go mistakes.** Walk the diff with [`references/go-review-checklist.md`](references/go-review-checklist.md). Lint already caught the mechanical ones; look for what lint can't see (wrong logic, missing `ctx` pass-through, error swallowed into a log, secrets in `Payload`, races).

**c. Primary test cases.** For each new or changed function with logic, list its cases with [`references/test-cases-checklist.md`](references/test-cases-checklist.md) and mark each **tested ✅** or **missing ❌** by reading the `_test.go` files. Every ❌ gets: why it matters, and an example test row the developer can paste.

**d. Repo conventions.**
- Every commit subject passes `sh scripts/hooks/commit-msg.sh` (write the subject to a temp file and run it)
- PR body follows `.github/pull_request_template.md` (no leftover `<...>` placeholders; sections that don't apply say "None")
- Endpoint changes come with an `api/openapi.yaml` change
- New migrations have both `-- +goose Up` and `-- +goose Down`, and no already-merged migration was edited
- `.env.example` is updated for new env vars; no secrets or `.env` committed

### 4. Grade findings and pick the verdict

**Report only what can cause an error, a regression, or a breaking change** (for API clients, the database, env/config, or public function signatures), plus security and data problems. Too many comments bury the ones that matter, so style, naming, docs, and "nice to have" tests are **not reported at all**.

| Severity | Use for | Reported? |
|---|---|---|
| **blocker** | Architecture rule broken, secret or personal data in logs or code, SQL injection, error ignored on an important path, a breaking change nobody planned for, new logic with no tests at all | ✅ inline comment |
| **major** | A likely bug or runtime error, a primary test case missing, `ctx` not passed on, resource not closed, error not wrapped where the caller can't tell what failed, OpenAPI not updated, migration without `Down`, breaking change to an API, column, or env variable | ✅ inline comment |
| **minor** | A real but small risk of an error or regression that is quick to fix (e.g. a missing test that could hide a bug, a default that is unsafe outside local) | ✅ inline comment |
| **nit** | Style, naming, doc comments, long functions, formatting, optional tests, wording in the PR body | ❌ never reported, not even in the summary |

If you are not sure that something can really cause an error or a regression, it is a nit: leave it out.

**Verdict:** any blocker → **Blocked** · else any major → **Needs changes** · else **Ready** (minor findings are still reported). The workflow turns the verdict into the PR status: **Ready → success, Needs changes or Blocked → failure**.

Don't invent findings to fill the report. A clean PR is **Ready** with "no findings". When unsure whether something is a real problem, say so, and mark it minor at most.

### 5. Report

Findings go **on the lines they are about**, one comment per finding. The developer is new to Go and reads each problem next to its code. One giant comment loses that context. Templates: [`references/report-template.md`](references/report-template.md).

**Two kinds of comments:**

| Comment | Where | Contains |
|---|---|---|
| **Inline comment**, one per finding | On the exact line (or line range) in the PR diff | Severity, what, why it matters, and a **code suggestion** |
| **Summary comment**, one per PR | The PR conversation | Verdict, counts, test-cases table, conventions, and findings that can't go on a diff line |

**Inline comment rules:**
- **Every reported finding (blocker, major, minor) gets its own inline comment**, as long as its line is part of the PR diff. Don't bundle several findings into one comment. If the same problem appears on several lines, comment on each line. From the 2nd occurrence on, the comment can be short and point to the first.
- **Add a `suggestion` block whenever the fix is a direct change to the commented lines**, so the developer can click *Commit suggestion* on GitHub:

  ````markdown
  ```suggestion
  	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
  ```
  ````

  - The block **replaces exactly the commented line(s)**. For a multi-line fix, comment on the whole range (start line → end line) and put the full replacement in the block.
  - Keep the original indentation (Go uses tabs) and write code that compiles and passes `make lint`.
  - When the fix touches other lines, adds a new file, or is a test to add, use a normal ` ```go ` block instead and say where it goes.
- Only comment on lines that are in the diff (added, or context lines in a changed hunk). Findings about unchanged code, missing files, or the PR as a whole go in the summary.
- Order of posting: blocker → major → minor. Up to **10 inline comments**; if there are more minors, put the rest in the summary as a list with `file:line` and no extra comment.

**Summary comment rules:** short, **does not repeat** what the inline comments already say, and links to them by listing `file:line` + title. Its first two lines are machine-readable markers:

```
<!-- check-pr -->
<!-- check-pr-verdict: <ready|needs-changes|blocked>; commit: <full 40-character head sha> -->
```

The workflow reads the second line to set the **Claude Review** status (ready → success, otherwise failure) and checks that `commit` is the PR head it started with. Write the verdict exactly as above, with the full SHA from `headRefOid`.

**Local mode:**
- Show the report in chat in the user's language. Use the same structure, with each finding under its `file:line`, plus the suggestion.
- Offer to post it to the PR, as inline comments + summary, in English.
- Remove the temp worktree and branch.

**CI mode:**
1. Post the **inline comments** with the `mcp__github_inline_comment__create_inline_comment` tool, one call per finding.
2. Post or update the **summary comment**. If a comment with the marker already exists, update it instead of adding another:

   ```bash
   ID=$(gh api repos/$R/issues/<PR>/comments --paginate \
     --jq '[.[] | select(.body | startswith("<!-- check-pr -->"))][0].id // empty')
   if [ -n "$ID" ]; then
     gh api -X PATCH repos/$R/issues/comments/$ID -F body=@<report.md>
   else
     gh pr comment <PR> -R $R --body-file <report.md>
   fi
   ```

3. Write everything in **English**, with the beginner-friendly explanations.

## Never

- Approve, request changes as a formal review, merge, push, or edit files in the PR. This skill only comments.
- Re-run CI or re-trigger workflows.
- Paste secrets, tokens, or `.env` content into a comment, even when reporting that one leaked: point to the file and line instead.
- Review the whole repository: only what the PR changes.
