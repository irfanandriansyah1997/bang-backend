---
name: atomic-conventional-commits
description: Split working-tree changes into small, atomic git commits with Conventional Commits messages. Use whenever the user asks to commit, save work, stage files, prep a branch for review, or split one mixed diff into multiple focused commits — including terse prompts like "commit this", "commit these changes", "split into commits", "ship it", or "wrap up". Also use when the agent notices the working tree mixes unrelated concerns before a commit.
---

# Atomic Conventional Commits

Turn the current working tree into a sequence of **atomic** commits, each with a **Conventional Commits** message. Atomic means one commit captures one logical change that could be reverted on its own without breaking the rest of the tree.

Portable core: the workflow below assumes standard git and Conventional Commits. Conventions specific to this repository are collected in [Repo conventions (bang-backend)](#repo-conventions-bang-backend) near the end. If the repo has its own `AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md`, or `.cursor/rules/*`, read those first and let repo rules override this skill.

## When to use

Trigger this skill when the user says any of:

- "commit this / these changes"
- "split this into commits" / "break this up"
- "prep this branch for PR / review"
- "save my work" / "wrap up" / "ship it"
- The user staged nothing but has a mixed diff
- The user staged part of a change and wants the rest committed sensibly

Also trigger proactively (without waiting to be asked) when the agent is about to commit and the working tree mixes unrelated concerns (e.g. feature + refactor + formatting + dep bump in one diff).

## When NOT to use

- User explicitly asks for one big commit — just run `git commit` directly.
- Amending a commit already pushed to a shared branch — propose a new commit instead (rewriting shared history is destructive).
- Reordering / rewording existing commits — that's `git rebase`, not this skill.
- Merge or rebase in progress — resolve that first; this skill only handles working-tree changes.

## Core principles

1. **Atomic** — one logical change per commit. Test: *reverting this commit alone should leave the tree in a working state.*
2. **Prefer more commits over fewer** — smaller slices are easier to review and revert.
3. **Conventional Commits** — every subject follows `type(scope): summary`.
4. **Never `git add -A` blindly** — stage only what belongs in the next commit; use hunk-level staging for mixed files.
5. **Interactive by default** — propose the plan, show each commit's files + message, and wait for approval before running `git commit`.
6. **Safety first** — scan for secrets, warn on protected branches, refuse destructive rewrites unless the user explicitly opts in.

## Workflow

Copy this checklist and track progress. Do not skip steps.

```
Commit plan:
- [ ] 1. Safety checks (branch + merge state + secrets)
- [ ] 2. Read working-tree state
- [ ] 3. Group changes into atomic slices
- [ ] 4. Present the plan and get approval
- [ ] 5. For each slice: stage → review staged diff → commit
- [ ] 6. Verify tree is clean and summarize
```

### 1. Safety checks

- `git status --short` and `git rev-parse --abbrev-ref HEAD`.
- If HEAD is `main` / `master` / `develop` / `release/*`, **stop and confirm** with the user before proceeding.
- If a merge, rebase, cherry-pick, or bisect is in progress (`.git/MERGE_HEAD`, `.git/REBASE_HEAD`, etc.), stop — this skill does not resolve those.
- Scan the diff for obvious secrets (API keys, tokens, `.env` files, private keys). If any appear, stop and ask the user before staging.

### 2. Read working-tree state

Run these read-only commands and keep the output in mind:

- `git status`
- `git diff --stat`
- `git diff` (unstaged)
- `git diff --cached` (staged, if anything is already staged)
- `git log -n 10 --oneline` — sample recent commits to match the repo's style (scope names, tense, subject length).

### 3. Group changes into atomic slices

Split the diff into a sequence of atomic groups. See [`references/grouping-guide.md`](references/grouping-guide.md) for the full decision tree.

Quick heuristics (priority order):

1. **Feature / concern** — different features, bug fixes, or user-visible behaviors always split.
2. **Scope** — different modules / packages / features usually split, unless one is a trivial hop for the other.
3. **Type** — `feat`, `fix`, `refactor`, `test`, `docs`, `style`, `chore` should not be mixed.
4. **Nature** — behavior change vs. pure formatting vs. dependency bump always split.

Order slices so that intermediate commits still compile and pass tests when possible. Typical safe order: `chore/build` → `refactor` → `feat` / `fix` → `test` → `docs`.

If a single file contains hunks that belong to different slices, plan to use **hunk-level staging** (see [Hunk-level staging without a TTY](#hunk-level-staging-without-a-tty)). Do not commit unrelated hunks together just because they share a file.

### 4. Present the plan and get approval

Before staging anything, show the user a numbered plan, for example:

```
Proposed commits (4):
  1. refactor(auth): extract TokenValidator from AuthService
     files: src/auth/AuthService.kt, src/auth/TokenValidator.kt
  2. feat(auth): support refresh-token rotation
     files: src/auth/AuthService.kt (partial), src/auth/RefreshPolicy.kt
  3. test(auth): cover refresh-token rotation edge cases
     files: src/auth/AuthServiceTest.kt
  4. docs(auth): document refresh-token flow
     files: docs/auth/refresh-tokens.md
```

Wait for explicit approval or edits. If the user changes anything (splits further, merges two, renames a scope), redo the plan and re-confirm.

### 5. For each slice: stage → review → commit

For every slice, in order:

1. **Stage** only that slice.
   - Whole files: `git add path/to/file1 path/to/file2`
   - Mixed files: stage only the relevant hunks with a patch (see [Hunk-level staging without a TTY](#hunk-level-staging-without-a-tty)). A human at a terminal can use `git add -p` instead.
   - Do **not** use `git add -A`, `git add .`, or `git commit -a`.
2. **Review** what will actually be committed: `git diff --cached`. Sanity-check for:
   - no secrets or credentials
   - no stray debug logs, TODOs unrelated to the slice, or commented-out code
   - no unrelated formatting churn from an editor auto-formatter
3. **Write the message** following [`references/commit-message-template.md`](references/commit-message-template.md).
4. **Commit** using a HEREDOC for multi-line messages (keeps formatting exact):

   ```bash
   git commit -m "$(cat <<'EOF'
   type(scope): imperative summary in <= 72 chars

   Optional body explaining what changed and why (not how). Wrap at ~72 chars.
   Reference issues / tickets if the repo convention expects it.
   EOF
   )"
   ```

5. **Do not** add trailers like `Co-authored-by:` or `Generated-by:` unless the user explicitly asks. Do not use `--no-verify` unless the user asks.

If a commit fails (pre-commit hook, lint, tests), fix forward with a **new commit** — do not `--amend` a commit that has already been pushed, and do not amend a commit you did not just create in this session.

### 6. Verify and summarize

- `git status` — expect a clean tree (or only the intentionally leftover changes).
- `git log --oneline -n <N>` — show the new commits.
- Report back to the user: a `slice → SHA → one-line summary` table, plus anything that was intentionally left uncommitted.

## Conventional Commits format

Subject line: `type(scope): imperative summary`

> In this repo the format is `<type>(<feature>): <message>` with the feature **required**. See [Repo conventions](#repo-conventions-bang-backend).

- **Imperative mood** — "add", not "added" or "adds".
- **No trailing period.**
- Aim for ≤ 50 chars, hard cap at 72.
- Scope is optional but recommended when the repo uses scopes; match the scope style of recent commits (`git log --oneline -n 20`).
- Body (optional) is separated by one blank line, wraps at ~72 chars, and answers *what* and *why*, not *how*.
- Footer (optional) carries `BREAKING CHANGE: ...` or issue refs (`Refs: #123`, `Closes: #123`) if the repo uses them.

Full template with worked examples: [`references/commit-message-template.md`](references/commit-message-template.md).

### Commit types

| Type       | Use for                                                             |
|------------|---------------------------------------------------------------------|
| `feat`     | New user-facing capability                                          |
| `fix`      | Bug fix affecting user-facing behavior                              |
| `refactor` | Code change with no behavior change                                 |
| `perf`     | Performance improvement with no behavior change                     |
| `test`     | Adding or updating tests only                                       |
| `docs`     | Documentation only (README, in-repo docs, code comments)            |
| `style`    | Formatting, whitespace, semicolons — no logic change                |
| `build`    | Build system, packaging, dependency version bumps                   |
| `ci`       | CI configuration and pipeline scripts                               |
| `chore`    | Housekeeping that doesn't fit above (tooling, ignores, etc.)        |
| `revert`   | Reverts a previous commit; body cites the reverted SHA              |

Some repos also use `security`. If recent history uses it, follow suit.

### Scope

- Derive from the primary module / package / feature touched (`auth`, `api`, `ui`, `payments`, ...).
- Match the casing and separators used in recent commits (e.g. all lowercase, or `kebab-case`).
- Drop the scope entirely if the change legitimately spans 3+ unrelated modules — usually a signal to split further.

## Hunk-level staging without a TTY

AI agents run commands without an interactive terminal, so `git add -p`, `git add -i`, and `git rebase -i` hang or fail. Stage hunks non-interactively instead:

```bash
git diff -U0 -- path/to/file > /tmp/slice.patch   # 1. dump the file's diff
# 2. edit /tmp/slice.patch: delete the hunks (@@ ... blocks) that belong to other slices
git apply --cached --unidiff-zero /tmp/slice.patch  # 3. stage only what is left
git diff --cached -- path/to/file                   # 4. confirm
```

For a **new** (untracked) file that must be split, run `git add -N path/to/file` first so it shows up in `git diff`.
To unstage a mistake: `git restore --staged path/to/file`, then retry.

## Repo conventions (bang-backend)

Specific to this repository (Go clean architecture; see `docs/architecture.md`). These win over the generic guidance above.

### Commit format (required)

```
<type>(<feature>): <message>
```

Examples: `feat(docs): adding readme.md with setup steps`, `fix(work-order): reject complete without photos`, `test(health): cover db down case`.

- **type**: required. One of `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `style`, `build`, `ci`, `chore`, `revert`.
- **feature**: **required**. Lowercase kebab-case name of the *feature or area* the change belongs to, not the code layer. A commit without `(feature)` is rejected by the `commit-msg` hook.
- **message**: required. A short description of the change, ≤ 72 chars for the whole subject line. Both "add …" and "adding …" are accepted; keep one style within a branch.
- **`!`** after `(feature)` marks a breaking change: `feat(auth)!: drop v1 login`.

The `commit-msg` hook (`scripts/hooks/commit-msg.sh`) enforces this; generic guidance above about optional scopes or imperative mood does not apply here.

### Feature names

Use the feature or area from the spec / GitHub issue. Reuse an existing name before inventing a new one (`git log --oneline` shows what's been used).

| Feature | Covers |
|---|---|
| `docs` | README files, `docs/**` |
| `ai` | `.ai/**`, `.claude/**`, `.cursor/**`, `AGENTS.md`, `CLAUDE.md` |
| `hooks` | `lefthook.yml`, `scripts/hooks/**` |
| `setup` | `go.mod`, `Makefile`, `docker-compose.yml`, `.gitignore`, project tooling |
| `infra` | config, logger, database connection |
| `health` | `GET /health` |
| `http` | shared HTTP base: router, middleware, error format, pagination helpers |
| `auth` | login, refresh, `/me`, JWT, role middleware |
| `technician` | technician CRUD, deactivate + reassign |
| `customer` | customer CRUD |
| `work-order` | work order CRUD, status state machine, assign/cancel/start/complete |
| `dashboard` | dashboard summary |
| `support-ticket` | Call Center tickets and attachments |

The `type` still says *what kind* of change it is, so a docs-only change to a feature is `docs(work-order): ...`, and a test-only change is `test(customer): ...`.

### Slicing a feature across layers

Keep the **same feature name** for every commit of a feature and split one commit per layer, in dependency order, so each commit still compiles:

1. `feat(customer): add customers table migration`
2. `feat(customer): add customer queries` (commit generated sqlc code together with the `.sql` that produced it)
3. `feat(customer): add Customer entity and validation` (domain + its tests)
4. `feat(customer): add customer usecases` (usecase + ports + its tests)
5. `feat(customer): add postgres customer repository`
6. `feat(customer): expose /customers endpoints` (handler, DTOs, route, wiring in `main.go`)
7. `docs(customer): document customer endpoints in openapi`

Keep tests in the same commit as the code they test (Go tests live next to the code).

### Issue references

Tasks live in GitHub issues on `irfanandriansyah1997/bang-backend`. Put the reference in the footer:

- `Refs: #12` for partial progress on an issue
- `Closes: #12` in the commit that finishes the issue's acceptance criteria

### Pre-commit checks

When `.go` files are staged, run `make check` (gofmt, vet, build, test) before committing, and stop on failure. Skip this when only docs or config change.

The lefthook `pre-commit` and `commit-msg` hooks enforce the same checks and the Conventional Commits format. If a hook blocks a commit, fix the cause and commit again. Never pass `--no-verify`.

### Branch

`main` is the default branch. A new repository with no commits yet may get its first commit on `main` once the user confirms. After that, prefer a feature branch per issue (`feat/12-customer-crud`).

## Guardrails

- **Never** stage or commit files matching `.env*`, `*.pem`, `*.key`, `id_rsa*`, `credentials*`, or files the user hasn't explicitly shown you, without asking.
- **Never** run `git push`, `git push --force`, `git reset --hard`, `git rebase -i`, or `git commit --amend` on your own — only when the user explicitly asks.
- **Never** change `git config`.
- Respect any repo-specific rules in `AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md`, or `.cursor/rules/*` (branch naming, required trailers, ticket prefixes, disallowed types). Repo rules win over this skill.

## Anti-patterns

- One "misc changes" commit that touches unrelated areas.
- `git add -A` followed by `git commit -m "wip"` or `git commit -m "update"`.
- Mixing formatting churn with logic changes in the same commit.
- Commit messages that describe the code (`change the loop`) instead of the intent (`avoid O(n²) scan on cache lookup`).
- Amending an already-pushed commit to "fix" a review comment on a shared branch.

## Deliverable

At the end of the run, report:

- The final list of commits: `#`, SHA, subject, one-line rationale.
- Any files intentionally left uncommitted and why.
- Any deferred concerns worth a follow-up commit or PR.

## References

- [`references/grouping-guide.md`](references/grouping-guide.md) — decision tree for slicing a mixed diff into atomic groups.
- [`references/commit-message-template.md`](references/commit-message-template.md) — subject/body/footer template plus type + scope guidance.
- [`references/examples.md`](references/examples.md) — worked before/after examples of splitting a messy diff.
