---
name: create-pr
description: Open a pull request (or rewrite an existing PR's description) using this repo's PR template. Short bullet-list summary with the change type on top, filled sections, no "Type of change" or commit list. Links the PR to its GitHub issue (task) if there is one, assigns it, copies the labels, and moves the task to In review. Use when the user asks to "open a PR", "bikin PR", "create pull request", "update PR description", "rapikan deskripsi PR", or after commits are ready on a feature branch.
---

# Create a PR

Turn the commits on the current branch into a pull request whose description follows [`.github/pull_request_template.md`](../../../.github/pull_request_template.md). The reader is a busy reviewer (or the developer, who is new to Go), so **keep it short and in bullet points**.

Reference example: PR #19 (`feat(setup): init go module, entry point, hooks, and CI`), with the two differences below.

## Format rules

- **No "Type of change" section.** The kind of change goes in the **first bullet of Summary**: `- **Type:** feat + ci: <one-line why>`.
- **Summary is a bullet list**, one bullet per notable change. No paragraphs.
- **No commit list** at the bottom. Reviewers see commits in the Commits tab.
- Every template section stays, in template order. If one doesn't apply, write `- None` (or `None` for tables) instead of deleting it or leaving placeholders.
- Remove the HTML hint comments (`<!-- ... -->`) from the final body.
- Prefer bullets and small tables over prose in every section.
- **Title** = `<type>(<feature>): <message>`, ≤ 72 chars, using the same format as commits (see the `atomic-conventional-commits` skill). Pick the type/feature of the main change; mention secondary ones in the Summary Type bullet.
- End the body with the attribution line if the session requires one.

## Workflow

```
- [ ] 1. Safety checks
- [ ] 2. Collect what changed + find the related task
- [ ] 3. Verify locally
- [ ] 4. Draft title + body
- [ ] 5. Show the draft and get approval
- [ ] 6. Push + create (or edit) the PR, assign, label, move the task to In review
- [ ] 7. Check CI
```

### 1. Safety checks

```bash
git status --short                       # must be clean; if not, commit first (/atomic-conventional-commits)
git rev-parse --abbrev-ref HEAD          # must NOT be main
git fetch -q && git log --oneline origin/main..HEAD   # must have commits
git rev-parse --abbrev-ref @{u} 2>/dev/null   # upstream; if it is origin/main, push with: git push -u origin <branch>
gh pr list --head "$(git branch --show-current)" --state open   # existing PR? then edit it instead of creating
```

Stop and tell the user if the tree is dirty, the branch is `main`, or there are no commits.

### 2. Collect what changed

```bash
git log --format='%h %s%n%b' origin/main..HEAD   # subjects, bodies, Closes/Refs footers
git diff --stat origin/main...HEAD
git diff origin/main...HEAD                      # read it; don't summarize from file names only
```

- **Related task**: see [Linking the PR to its task](#linking-the-pr-to-its-task).
- **Layers**: map changed paths to the "Changes by layer" table (`internal/domain`, `internal/usecase`, `internal/delivery/http`, `internal/repository/postgres`, `internal/infrastructure`, `cmd/api`). For changes outside these layers (tooling, CI, docs, AI config), add a **Tooling** table right below, as in PR #19.
- **API / DB / Config**: look for changes in `api/openapi.yaml`, `internal/delivery/http`, `db/migrations`, `db/queries`, `.env.example`, `internal/infrastructure/config`.

### Linking the PR to its task

Every PR should say which GitHub issue (task) it delivers. Find it in this order:

1. `Closes: #N` / `Refs: #N` footers in the branch's commits
2. The branch name: `feat/12-customer-crud` → #12
3. An issue whose acceptance criteria match the diff: search open issues, using the same phase label when you can guess it

   ```bash
   gh issue list -R irfanandriansyah1997/bang-backend --state open --json number,title,labels \
     --jq '.[] | "#\(.number) [\(.labels | map(.name) | join(","))] \(.title)"'
   ```

4. Otherwise **ask the user**: "Is this PR for a task? Give me the issue link or number, or say none."

Confirm the match with the user when it came from step 3 (a guess), or when the diff only covers part of the issue.

**`Closes` or `Refs`?** Use `Closes #N` only if the PR meets every acceptance criterion and DoD item of the issue (read it with `gh issue view N`). Otherwise use `Refs #N` and say what's left in the Summary.

**When a task is found**, include these in the approval step (5) and apply them in step 6:

| What | How | Why |
|---|---|---|
| Link in the body | `Closes #N` / `Refs #N` under **Related issue** | Shows in the issue's *Development* panel; `Closes` closes it on merge |
| Assignee | `--assignee @me` | The PR has an owner |
| Labels | Copy the issue's labels (e.g. `phase-2-auth-master-data`) with `--label` | Filtering PRs by phase works like issues |
| Board status | Move the **issue** to `In review` on Project 1 | The board shows the task is waiting for review |

Moving the issue on the board:

```bash
ITEM=$(gh project item-list 1 --owner irfanandriansyah1997 --limit 200 --format json \
  --jq '.items[] | select(.content.number==<N>) | .id')
gh project field-list 1 --owner irfanandriansyah1997 --format json \
  --jq '.fields[] | select(.name=="Status") | {id, options}'      # Status field id + "In review" option id
gh project item-edit --id "$ITEM" --project-id "$(gh project view 1 --owner irfanandriansyah1997 --format json --jq .id)" \
  --field-id <status-field-id> --single-select-option-id <in-review-option-id>
```

**When there is no task** (the user says none): write `- None` under **Related issue**, still assign `@me`, and skip labels and the board.

### 3. Verify locally

```bash
make check            # gofmt, vet, build, test; stop and report if it fails
```

Run anything else the "How to test" section will claim (e.g. `go run ./cmd/api`, a `curl`). Only tick checklist items you actually verified, and leave the rest unticked with a short note.

### 4. Draft title and body

Write the body to a file in the scratchpad (not the repo), using this shape:

````markdown
## Summary
- **Type:** <type[s]>: <why this PR exists, one line>
- <notable change 1>
- <notable change 2>

## Related issue
Closes #<N>

## Changes by layer
| Layer | Change |
|---|---|
| `<layer>` | <one line> |

<!-- only if needed -->
### Tooling
| Area | Change |
|---|---|
| <area> | <one line> |

## API changes
None   <!-- or the Method / Path / Role / Change table + OpenAPI checkbox -->

## Database changes
- None

## Config / environment
- None

## How to test
```bash
make check
<commands with expected output in comments>
```

## Observability
- <new logs/events, or "Not applicable yet">

## Risk & rollback
- Risk: <low/medium/high + why>
- Rollback: <revert the merge commit / down migration>

## Checklist
- [x] `make check` passes locally (gofmt, vet, build, test)
- [ ] ...   <!-- tick only what you verified -->
````

### 5. Show the draft and get approval

Show the user:

- the title and the full body
- the **related task** (`Closes #N` / `Refs #N` / none) and how it was found
- the assignee, the labels, and the board move (`#N → In review`) that will be applied

Wait for an explicit OK or edits. Pushing, opening a PR, and moving the board are visible to others, so **never do step 6 without approval**.

### 6. Push and create or edit the PR

```bash
git push -u origin "$(git branch --show-current)"
gh pr create --base main --head "$(git branch --show-current)" --title "<title>" --body-file <draft.md> \
  --assignee @me --label <issue-label>          # one --label per issue label; omit if no task
# existing PR instead:
gh pr edit <PR> --title "<title>" --body-file <draft.md> --add-assignee @me --add-label <issue-label>
```

Then, if there is a task, move the issue to **In review** (see [Linking the PR to its task](#linking-the-pr-to-its-task)) and confirm the link landed:

```bash
gh pr view <PR> --json closingIssuesReferences,assignees,labels   # Closes #N shows up in closingIssuesReferences
```

`closingIssuesReferences` only lists `Closes`/`Fixes`/`Resolves` links, so it stays empty for `Refs #N`. That's expected.

### 7. Check CI

Read the PR's checks once (`gh pr checks <N>`, or the host's PR status tool if available) and report them: the `main` ruleset requires **Lint** and **Test** to pass before merge. Don't poll in a loop. If a check fails, show the failing step and suggest the fix.

Finish by giving the user the PR link and suggesting `/review-task <PR> <issue>` before merging.
