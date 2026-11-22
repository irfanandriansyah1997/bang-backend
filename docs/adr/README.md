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
