# Architecture Decision Records (ADR)

An ADR records one decision that a reviewer accepted, so the same finding is not raised again in every review.

## Folder rules

- Files here are named `NNNN-short-title.md`. Start from [`0000-template.md`](0000-template.md).
- Only `.md` files live here. The review checks this folder, and a PR that adds anything else here is reviewed as code.

## Status

| Status | Meaning | Used by the review? |
|---|---|---|
| `proposed` | Written from a review, waiting for a person to merge it | No |
| `accepted` | Merged to `main`, and the team agreed to it | Yes, when it is on `main` |
| `superseded` | Replaced by a newer ADR (link to it) | No |

The review reads only `accepted` ADRs from `main`. An ADR added in the same PR is treated as a proposal for that PR, so a PR cannot approve its own code.

## Who writes what

- An AI may only create a new file with status `proposed`. It never changes `accepted` ADRs.
- A person decides whether an ADR is accepted, by merging its PR.

## How an ADR is added from a review

1. A person replies to a **minor** review comment, with `/skip-suggestion` as the last line of the reply. The text above that line is the reason.
2. The next `@claude-code-review` run (`check-pr --ci`) checks the ADRs on the PR branch. If none covers the point yet, it writes a new `proposed` ADR here. The reason comes from the reply, or the AI writes one when the reply has none (marked as such).
3. The workflow step **Commit ADR** in `.github/workflows/claude-review.yml` (a script, not the AI) checks that the run added only new `docs/adr/NNNN-*.md` files with `Status: proposed` and changed nothing else. It then commits them as `github-actions[bot]` and pushes to the PR branch with the `ADR_PUSH_TOKEN` secret.
4. The push changes only `docs/adr/*.md`, so **Claude Review outdated** carries the verdict of the reviewed commit over to the new commit. No new review is needed.

Blockers and majors never get an ADR: they must be fixed.

`ADR_PUSH_TOKEN` is a fine-grained personal access token with **Contents: Read and write** on this repository only. A push with the default `GITHUB_TOKEN` would not start Lint, Test, or the outdated workflow, so the PR could not pass its required checks. Without the secret, the step leaves a comment and commits nothing.
