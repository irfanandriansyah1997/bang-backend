---
name: adr-from-review
description: Turn the minor findings a reviewer skipped with /skip-suggestion on a PR into a proposed ADR file in docs/adr/, so future reviews do not raise them again. Use when the user asks to "buat ADR dari review", "adr from review", or after a check-pr run where a minor was skipped. check-pr --ci follows its CI mode (write the file only; the workflow commits it). Writes only docs/adr/NNNN-*.md with status proposed, never code, and never blockers or majors.
model: sonnet
---

# Write an ADR from a review

Turn a decision a reviewer already made (a minor finding skipped with `/skip-suggestion`) into a **proposed** ADR. Read [`docs/adr/README.md`](../../../docs/adr/README.md) and [`docs/adr/0000-template.md`](../../../docs/adr/0000-template.md) first, and follow them exactly.

## Hard limits

- Write **only** `docs/adr/NNNN-short-title.md`. Nothing else: no code, no workflow, no config, no skill.
- Status is always `proposed`. Never `accepted`. Never change an existing ADR.
- Only **minor** findings. Never a blocker or a major: those are not decisions a reviewer can skip.
- Never put a secret, token, password, or personal data into the ADR.
- Never push, merge, or open a PR without the user's explicit OK.

## Modes

| | Local mode (`/adr-from-review <PR>`) | CI mode (from `check-pr --ci`) |
|---|---|---|
| Who runs it | The user, on their machine | The `check-pr` review in `.github/workflows/claude-review.yml` |
| Branch check (step 0) | Yes | Skip: the runner's checkout is already the PR head |
| Find skipped minors (step 1) | Yes | Skip: `check-pr` step 1b already found them |
| Existing ADR check (step 1b) | Yes | **Yes**, with `Glob docs/adr/*.md` + `Read` (the checkout is the PR head) |
| Path rule (step 2) | Yes | Skip: the workflow's "Commit ADR" step checks it |
| Number, write (steps 3–4) | Yes | **Yes** |
| Commit (step 5) | You commit, after the user's OK | **Never.** Only write the file. The workflow validates it, commits it as `github-actions[bot]`, and pushes it |

## 0. Check the branch (local mode)

The ADR is committed to the PR's own branch, so the PR head must be checked out here:

```bash
gh pr view <PR> -R irfanandriansyah1997/bang-backend --json headRefName --jq .headRefName
git branch --show-current
```

If the two names differ, stop and tell the user which branch to switch to.

## 1. Find the skipped minors (local mode)

```bash
R=irfanandriansyah1997/bang-backend
gh api repos/$R/issues/<PR>/comments --paginate --jq '.[] | select(.body|startswith("<!-- check-pr -->")) | .body' | tail -1
```

The summary's section **Skipped by the reviewer (`/skip-suggestion`)** lists the candidates. Confirm each one in its thread:

```bash
gh api graphql -f query='query($o:String!,$n:String!,$p:Int!){repository(owner:$o,name:$n){pullRequest(number:$p){reviewThreads(first:100){nodes{isResolved path line comments(first:50){nodes{author{login} body}}}}}}}' \
  -F o=irfanandriansyah1997 -F n=bang-backend -F p=<PR>
```

A finding counts as skipped only when the **last non-empty line of a reply by a person** is exactly `/skip-suggestion` (the same rule as `check-pr`). Lines starting with `>` do not count, and replies by `claude[bot]` do not count. If the summary says a major was skipped, stop and tell the user: majors are never skipped.

## 1b. Check for an existing ADR first

Before writing anything, see whether a decision on the same point already exists. Read the ADRs on the current branch (the PR head). The engineer rebases from the latest `main` before the PR is reviewed, so the branch already has every ADR that is on `main`:

```bash
git ls-tree --name-only HEAD docs/adr/ | grep -E '/[0-9]{4}-.*\.md$' | while read f; do echo "== $f"; git show "HEAD:$f"; done
```

For each ADR, compare its **Area**, **Applies to**, and **Decision** with the skipped finding (same file or layer, same kind of problem):

| What you find | What to do |
|---|---|
| An `accepted` ADR (it is on `main`, so the rebased branch has it) that covers the same point | **Do not write a new ADR.** Tell the user the existing ADR number and file, and stop. The finding is already decided. |
| A `proposed` ADR in this PR that covers the same point | **Do not write a second one.** Tell the user which file already proposes it. |
| A `superseded` ADR, or one about a different point | Not a match. Continue. |
| Nothing similar | Continue to step 2. |

When two ADRs are close but not the same, write the new one and link to the old one in its Context, so the reader sees both.

## 2. Check the path rule before writing (local mode)

Run from the PR branch, and stop if anything fails:

```bash
git fetch -q origin main
git diff --name-only origin/main...HEAD
```

Every changed file must match `docs/adr/[^/]+\.md`. If the PR also changes code, that is fine for the review, but this skill must not commit that code. It only writes its own ADR file.

## 3. Pick the number

The number is the next free four-digit number. Look at both `main` and the PR branch, so two proposals do not take the same number:

```bash
git ls-tree --name-only origin/main docs/adr/ | sed 's|docs/adr/||' | grep -E '^[0-9]{4}-'
git ls-tree --name-only HEAD docs/adr/ | sed 's|docs/adr/||' | grep -E '^[0-9]{4}-'
```

Take the highest number in both lists and add 1. Never reuse a number. When one run writes several ADRs, give each the next number (0004, 0005, ...).

## 4. Write the ADR

Copy `0000-template.md` and fill every section:

- **Status:** `proposed`
- **Date:** today (`YYYY-MM-DD`)
- **Source:** `PR #<number>, review thread <link to the thread>`
- **Area:** the folder or layer, for example `delivery/http`
- **Applies to:** one sentence: what the decision covers.
- **Reason source:** `reviewer` or `AI` (see [The reason](#the-reason)).
- **Context:** describe the skipped finding in plain words, so a developer new to Go can follow it. Do not paste the code.
- **Decision:** a rule a reviewer can check, in one or two sentences, followed by the reason.
- **Consequences:** what gets easier, what risk remains, and what the reviewer must still check.
- **Scope and limits:** keep the template's text. It says blockers and majors are not covered.

Name the file `NNNN-short-title.md` (lowercase letters, digits, and dashes only).

### The reason

The reason comes from the reply in which the person wrote `/skip-suggestion`:

1. Take the lines **above** the `/skip-suggestion` line.
2. Drop empty lines and lines that start with `>` (quoted text).
3. What is left is the reviewer's reason.

| What is left | Reason source | What to write under **Decision** |
|---|---|---|
| Some text | `reviewer` | `Reason given by the reviewer (@<login>):` and the text, as a quote. Fix only obvious typos; do not change its meaning. |
| Nothing (the reply was only the token) | `AI` | `Reason written by the AI (the reviewer gave none):` and one or two sentences on why the finding is acceptable for now: the risk is small, and what would make it worth fixing later. Do not invent facts about the team's plans. |

If the reason contains a secret, token, or personal data, leave that part out and write `(removed: sensitive data)`.

## 5. Check and commit (local mode)

```bash
git add docs/adr/NNNN-short-title.md
git diff --cached --name-only | grep -vE '^docs/adr/[^/]+\.md$' && echo "STOP: only docs/adr/*.md may be committed" && exit 1
git commit -m "docs(adr): propose NNNN <short title>" -m "Refs: #<issue>"
```

Show the user the file and the commit message. Push only after the user says yes.

**CI mode:** do not run `git add` or `git commit` (the tools are blocked anyway). Write the file and stop. The workflow step *Commit ADR* checks that only new `docs/adr/NNNN-*.md` files with `Status: proposed` were added, then commits and pushes them.

## 6. Hand off (local mode)

Tell the user:
- the file path and the decision, in one line
- that the status is `proposed` until the PR is merged to `main`
- that the review reads it only after it is `accepted` on `main`
- that the push is ADR-only, so the workflow carries the previous verdict over and **no new review is needed**. The PR status stays the same as before the ADR was added.
