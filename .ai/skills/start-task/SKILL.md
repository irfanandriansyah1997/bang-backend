---
name: start-task
description: Brief the developer on a GitHub issue before they start it. Writes a plan file in docs/superpowers/plans/ with what to build, step by step with full source code, and the acceptance criteria and Definition of Done, then gives a short summary in chat. Keeps that file updated when requirements change. Use when the user shares an issue link or number and wants to understand or start the task themselves, e.g. "start issue #2", "jelasin task ini", "apa yang harus dikerjain di <issue link>", "mau kerjain issue berikutnya". Do NOT write the code into the repo; use add-feature for that.
# Planning and writing code needs a stronger model; Sonnet is half the price of Opus.
model: sonnet
---

# Start a task

The developer is **new to Go** (strong React/TypeScript background) and wants to **type the code themselves**. They asked to be "spoon-fed": give them the full source code for every file, explain it simply, and tell them exactly where it goes and how to check it. Explain in the language the user writes in (usually Indonesian); keep code, commands, and file names as-is.

**Do not create or edit code files in the repo.** The brief, with all source code, goes into **one plan file** (see [Where the brief goes](#where-the-brief-goes)), which the user reads and types from. Editing the issue or the board is fine **after asking** (step 5).

### Where the brief goes

- Write the full brief to `docs/superpowers/plans/YYYY-MM-DD-issue-<N>-<slug>.md` (today's date, e.g. `2026-11-08-issue-2-infrastructure.md`). It's a doc, not code, so creating it is allowed.
- **In chat, only reply with a short summary**: the file link, the goal, the list of step titles, any warnings (blocked dependency, a tool that isn't running), and what to do first. Don't paste the full brief into chat.
- **Keep the file current.** When the user asks for changes later (a new requirement, a different API, the next batch of steps), update the plan file in place: replace the affected step with the new full version, add a line to the **Changelog** section at the bottom, and tell the user in chat which steps changed. Don't make them hunt through old messages.
- The plan file has a **Progress** checklist (one box per step) near the top so the user can tick steps as they go.
- Code in the plan file must be the exact code you compiled and tested in the scratch module. Copy it from the scratch files rather than retyping it.
- If a plan file for this issue already exists, update it instead of creating a second one.
- Commit the plan file together with the user's work on the feature branch (`docs(<feature>): add plan for #N`) only when the user commits, and include it in the commit plan that `/atomic-conventional-commits` proposes.

### Teaching rules (always)

- **Full source code, not fragments.** For every file, give the complete content (package line, imports, all functions) so it compiles as-is. If a later step changes an existing file, show the **whole new version** of that file and say what changed.
- **Say where it goes.** Before each code block: the exact path (`internal/domain/work_order.go`), whether it's a **new** file or an **edit**, and which layer it belongs to.
- **Explain simply.** After each code block, a short "what this does" in plain words, then a table or bullets explaining each new Go concept **the first time it appears**, with a React/TS comparison.
- **Small steps.** One file or one concept per step. Every step must compile and has a check.
- **Show expected output.** Each check lists the exact command and what the terminal should print, including "no output = success" where that applies.
- **Warn about traps.** 1–2 common beginner mistakes per step with the error message they'll see and how to fix it.
- **Tests too.** Give the full `_test.go` file for new logic, using table-driven tests, and explain how to run one test (`go test -run Name -v ./path`).
- **No jargon without a definition.** If you use a term (port, adapter, DTO, pointer, receiver, context), define it in one line.

## Input

A GitHub issue URL (`https://github.com/irfanandriansyah1997/bang-backend/issues/N`) or a number (`#N`). If none is given, list open issues with status `Ready` on the project board and ask which one.

## Workflow

```
- [ ] 1. Read the issue, its labels, and its board status
- [ ] 2. Read the context (spec, architecture, folder READMEs, current code)
- [ ] 3. Check dependencies and what's already done
- [ ] 4. Write the brief to the plan file (goal → steps with full code → DoD), then a short chat summary
- [ ] 5. Offer to update the issue / board
```

### 1. Read the issue

```bash
gh issue view <N> -R irfanandriansyah1997/bang-backend --json number,title,state,labels,body,url
gh project item-list 1 --owner irfanandriansyah1997 --limit 200 --format json \
  --jq '.items[] | select(.content.number==<N>) | {status, priority, size}'
```

If the issue is closed, say so and stop unless the user still wants a brief.

### 2. Read the context

- `docs/architecture.md`: layer rules and the "add a feature" checklist
- The `README.md` of every folder the task touches (do's and don'ts)
- The matching part of `docs/superpowers/specs/rencana-belajar-fullstack-fde.md` (feature table, data model, endpoints, business rules)
- Current code in the touched folders (`ls`, read files) so the steps build on what exists

### 3. Check dependencies and current state

- Issues this one depends on (e.g. #4 `/health` needs #2 infrastructure). Look at earlier phase labels and the issue body. If a dependency is still open, warn the user and suggest doing it first.
- For each acceptance criterion, check whether the repo **already** satisfies it (file exists, command passes). Mark those as already done in the brief so the user doesn't redo them.

### 4. Write the brief

Write it to the plan file using this structure:

````markdown
# Plan: Issue #N: <title>
> Issue: <issue URL> · Branch: `feat/<N>-<short-slug>` · Last updated: <YYYY-MM-DD>

**Goal:** <one or two sentences, plain language>
**Depends on:** #X (done ✅ / not done ⚠️)
**Already done in the repo:** <criteria that already pass, or "nothing yet">
**Before you start:** <tools to install/start, e.g. Docker; `git switch -c feat/<N>-<slug>` from an up-to-date main>

## Progress
- [ ] Step 1: <title>
- [ ] Step 2: <title>
- [ ] ...

## What you'll build
<a small tree of the files this task creates/edits, marked NEW / EDIT>
```
internal/domain/
  work_order.go        NEW
  work_order_test.go   NEW
cmd/api/main.go        EDIT
```

## Steps
### Step 1: <what, in plain words>
**File:** `path/to/file.go` (NEW, layer: `domain`)
**Why:** <1–2 lines, with a React/TS comparison when it helps>

```go
// FULL file content: package, imports, every function.
// Inline comments on the lines that introduce a new concept.
```

**What this does:** <2–4 plain sentences>

| Go concept | Meaning | In React/TS |
|---|---|---|
| `func (w *WorkOrder) Start()` | method with a pointer receiver | a class method that changes `this` |

**Check:**
```bash
go build ./...      # (no output = success)
```
**Common mistakes:**
- <mistake> → you'll see `<error message>` → fix: <how>

### Step 2: ...
(each step that edits an existing file shows the WHOLE updated file and says "changed: …")

## Definition of Done
- [ ] <each acceptance criterion, made checkable>
- [ ] `make check` passes (gofmt, vet, build, test)
- [ ] Tests for new logic (domain + usecase at minimum)
- [ ] Follows docs/architecture.md (imports inward, thin handlers, no sqlc types outside repository)
- [ ] `api/openapi.yaml` updated (if endpoints changed)
- [ ] Commits use `<type>(<feature>): <message>`; last commit or PR body has `Closes #N`
- [ ] PR opened with the template, CI Lint + Test green

## When you're stuck
<2–3 example prompts the user can send to get hints rather than full answers>

## Changelog
- <YYYY-MM-DD>: plan created
````

Then reply in chat with the short summary described in [Where the brief goes](#where-the-brief-goes).

Guidelines for steps and source code:

- Order steps so each one can be run and checked before the next (inner layers first: domain → usecase → repository → http → wiring).
- Follow the [teaching rules](#teaching-rules-always): full files, path + NEW/EDIT, plain explanation, concept table, check with expected output, common mistakes.
- Make sure the code you show is correct. When it's non-trivial, compile and test it in a scratch module outside the repo before showing it, and never write it into the repo.
- Code must follow `docs/architecture.md` and the folder READMEs, because the user will copy it as the pattern for later tasks.
- Explain new Go concepts (interfaces, pointers, `context`, error returns, `defer`, struct tags) the first time they appear in this brief.
- Only include DoD items that apply (no OpenAPI item for a tooling task).
- Long tasks: write the first 3–5 steps in full, list the rest as titles in Progress and as a short "Next batch" section, and offer to continue. When the user says continue, fill those steps into the **same** plan file.

### 5. Offer to update the issue and board

Ask before each of these; never do them silently:

- If the issue has no **Definition of Done** section, offer to add the DoD from the brief to the issue body.
- If the board status is `Ready` or `Backlog`, offer to move it to `In progress`:

  ```bash
  ITEM=$(gh project item-list 1 --owner irfanandriansyah1997 --limit 200 --format json \
    --jq '.items[] | select(.content.number==<N>) | .id')
  gh project field-list 1 --owner irfanandriansyah1997 --format json   # find Status field + "In progress" option ids
  gh project item-edit --id "$ITEM" --project-id <project-id> --field-id <status-field-id> \
    --single-select-option-id <in-progress-option-id>
  ```

When the user finishes, point them to `/atomic-conventional-commits` to commit and `/review-task` to cross-check the PR.
