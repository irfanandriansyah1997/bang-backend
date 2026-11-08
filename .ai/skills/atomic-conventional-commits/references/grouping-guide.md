# Grouping guide

How to slice a mixed working tree into atomic commits.

## The atomic test

A commit is atomic when **reverting only that commit leaves the tree in a working state**. If reverting one slice would break another slice, they are entangled — either merge them into one commit or reorder / re-split until each stands alone.

Intermediate commits that don't compile or don't pass tests are a smell. Order slices so the tree stays green after each one whenever possible.

## Decision tree

Ask these in order. As soon as one answers "yes, they differ", split.

1. **Different concerns / features?**
   - Two user-visible behaviors, two bug fixes, or a feature plus a mechanical refactor → split.
2. **Different scopes / modules?**
   - Changes in `auth/` and `payments/` with no functional coupling → split.
   - Exception: a one-line call-site change forced by the primary slice can ride along.
3. **Different Conventional Commit types?**
   - `feat` + `fix`, `refactor` + `feat`, `test` + `feat` (usually) → split.
   - Tests written *for* the feature are often okay in the same commit if the repo prefers it; check recent history.
4. **Different natures of change?**
   - Behavior change vs. pure formatting → always split.
   - Behavior change vs. dependency bump → always split.
   - Refactor vs. rename → usually split; a rename that enables the refactor can ride along.
5. **Would reverting one break the other?**
   - Yes → they are entangled; either keep together or reorder so the dependency is committed first.

## Priority order for slicing

When a diff has several axes of variation, split by the highest applicable level:

1. **Feature / concern** (highest priority)
2. **Scope** (module / package)
3. **Type** (`feat`, `fix`, `refactor`, ...)
4. **Nature** (behavior vs. formatting vs. deps)

Within the same feature and scope, prefer to keep implementation + its immediate tests together **only if** the repo convention allows it; otherwise split `feat` from `test`.

## Recommended ordering of slices

When you have to pick an order, this ordering usually keeps the tree green after every commit:

1. `chore` / `build` — tooling, ignores, dependency bumps required by later slices
2. `refactor` — preparatory, non-behavior-changing moves
3. `feat` or `fix` — the actual change
4. `test` — tests for the change (if split out)
5. `docs` — documentation
6. `style` — pure formatting, if any

Reorder if a later slice logically depends on an earlier one.

## Hunk-level staging

If a single file mixes hunks from different slices:

- **AI agent (no TTY):** dump the diff with `git diff -U0 -- <file> > /tmp/slice.patch`, delete the hunks that belong to other slices, then `git apply --cached --unidiff-zero /tmp/slice.patch`. For untracked files, `git add -N <file>` first.
- **Human at a terminal:** `git add -p <file>` and accept only the hunks for the current slice; press `s` to split a hunk further, `e` to edit it manually.
- To undo a staging mistake: `git restore --staged <file>` and retry.
- Verify with `git diff --cached` before committing.

Never commit unrelated hunks together just because they share a file.

## Escape hatches

- If a slice is genuinely inseparable (e.g. a behavior change that only makes sense with a matching config update), keep it as one commit and note *why* in the body.
- If the diff is so tangled that slicing takes longer than the change itself, propose a single commit but flag the entanglement in the message body and suggest a follow-up refactor.
- If the user pushes back on a split, defer to the user — this skill proposes, the user decides.

## Quick smell tests

Split further if any of these apply to a proposed commit:

- Subject needs "and" or a comma to describe it.
- Message body has more than one "This commit also ..." sentence.
- Changed files list touches 3+ unrelated modules.
- Diff mixes renames, formatting, and logic changes.
- You catch yourself wanting to write "misc", "cleanup", "various", or "wip".
