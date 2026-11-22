# ADR NNNN: <short title in plain words>

- **Status:** proposed
- **Date:** YYYY-MM-DD
- **Source:** PR #<number>, review thread <link or "summary comment">
- **Area:** <for example: delivery/http, repository, usecase, CI>
- **Applies to:** <what this decision covers, in one sentence>
- **Reason source:** <reviewer | AI (the reviewer gave no reason)>

## Context

What happened, and why a decision is needed. Write it so a developer who is new to Go can follow it.
Describe the finding that was skipped, not the code that was reviewed.

## Decision

What the team accepts or does, in one or two sentences. Be specific: a rule that a reviewer can check.

Then the reason, in one of these two forms:

> Reason given by the reviewer (@login): <the reviewer's own words from the `/skip-suggestion` reply>

> Reason written by the AI (the reviewer gave none): <why the finding is acceptable for now, and what would make it worth fixing later>

## Consequences

- What becomes easier.
- What becomes harder, or what risk remains.
- What the reviewer must still check, even with this decision.

## Scope and limits

- This decision does **not** cover blockers or majors. Those must still be fixed.
- Only minor findings can be accepted here. Nits are never reported, so they need no decision.

---

## Rules for writing this file

- File name: `docs/adr/NNNN-short-title.md`, with the next free four-digit number. Do not reuse a number.
- New ADRs start as `proposed`. Only a person changes the status to `accepted` (after the PR is merged).
- Copy this template. Do not change its headings.
- One decision per file.
- Never put a secret, token, password, or personal data in an ADR.
