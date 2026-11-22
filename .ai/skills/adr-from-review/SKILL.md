---
name: adr-from-review
description: Turn the minor findings a reviewer skipped with /skip-suggestion on a PR into a proposed ADR file in docs/adr/, so future reviews do not raise them again. Use when the user asks to "buat ADR dari review", "adr from review", or after a check-pr run where a minor was skipped. Writes only docs/adr/NNNN-*.md with status proposed, never code, and never blockers or majors.
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

## 0. Check the branch

The ADR is committed to the PR's own branch, so the PR head must be checked out here:

```bash
gh pr view <PR> -R irfanandriansyah1997/bang-backend --json headRefName --jq .headRefName
git branch --show-current
```

If the two names differ, stop and tell the user which branch to switch to.

## 1. Find the skipped minors

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

## 2. Check the path rule before writing

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

Take the highest number in both lists and add 1. Never reuse a number.

## 4. Write the ADR

Copy `0000-template.md` and fill every section:

- **Status:** `proposed`
- **Date:** today (`YYYY-MM-DD`)
- **Source:** `PR #<number>, review thread <link to the thread>`
- **Area:** the folder or layer, for example `delivery/http`
- **Applies to:** one sentence: what the decision covers.
- **Context:** describe the skipped finding in plain words, so a developer new to Go can follow it. Do not paste the code.
- **Decision:** a rule a reviewer can check, in one or two sentences.
- **Consequences:** what gets easier, what risk remains, and what the reviewer must still check.
- **Scope and limits:** keep the template's text. It says blockers and majors are not covered.

Name the file `NNNN-short-title.md` (lowercase, dashes).

## 5. Check and commit

```bash
git add docs/adr/NNNN-short-title.md
git diff --cached --name-only | grep -vE '^docs/adr/[^/]+\.md$' && echo "STOP: only docs/adr/*.md may be committed" && exit 1
git commit -m "docs(adr): propose NNNN <short title>" -m "Refs: #<issue>"
```

Show the user the file and the commit message. Push only after the user says yes.

## 6. Hand off

Tell the user:
- the file path and the decision, in one line
- that the status is `proposed` until the PR is merged to `main`
- that the review reads it only after it is `accepted` on `main`
- that the push is ADR-only, so the workflow carries the previous verdict over and **no new review is needed**. The PR status stays the same as before the ADR was added.
