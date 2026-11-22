# ADR 0001: Log the raw error text of 5xx responses

- **Status:** proposed
- **Date:** 2026-11-22
- **Source:** PR #39, review thread `PRRT_kwDOVnGdm866iR4k` on `internal/delivery/http/errors.go:27`
- **Area:** delivery/http
- **Applies to:** Whether `Responder.Error` writes the full `err.Error()` text into the log when a request fails with a 5xx status.
- **Reason source:** reviewer

## Context

When the server fails with a 500, the HTTP error mapper writes the full error text into the log, so the operator can see the real cause. The automated review flagged this as a minor risk: an error message can carry data that should not be in logs, such as a value from the request, and the project rule says never to log passwords, tokens, or personal data.

The developer chose to keep the raw text on purpose. The protection is meant to live earlier in the stack: repositories and usecases are where errors are created and wrapped, and their messages are checked in the lower environments before anything reaches production.

## Decision

`Responder.Error` may log `err.Error()` for a 5xx response. Error messages created in `repository` and `usecase` must never contain passwords, tokens, or personal data, and the reviewer must check this when they add or change those error messages.

Reason given by the reviewer (`irfanandriansyah1997`):

> Thanks for the review. We are keeping `err.Error()` in the 5xx log on purpose, and the protection is placed earlier in the stack:
>
> - **Repository and usecase** are where errors are created and wrapped. There, secrets and personal data must never be put into error messages. This is checked during sanity testing on the lower environments before anything reaches production.
> - **HTTP layer** (`Responder.Error`) logs the error as it is. It does not sanitize again, so the log shows the real cause without extra cleanup at the end.
>
> Reason: the lower environments are where this is verified, so the raw text there is safe to inspect, and the HTTP layer should not second-guess what the domain and repository already produced.
>
> Note for reviewers: the test `TestResponder_Error_HidesDetailsOf5xx` uses a fake database error on purpose. It checks that the text is not sent to the client and that the 5xx is logged once with `error_code`. The fake text does appear in the log line, which is expected for this design.

## Consequences

- Operators see the real cause of a 500 in the logs, so incidents are faster to debug.
- Risk remains: if a repository or usecase puts sensitive data into an error message, that data reaches the logs. The protection depends on those layers, not on this one.
- The reviewer must still check that new error messages in `repository` and `usecase` do not contain secrets or personal data, and that the logging in `delivery/http` does not add request bodies.

## Scope and limits

- This decision does **not** cover blockers or majors. Those must still be fixed.
- Only minor findings can be accepted here. Nits are never reported, so they need no decision.
