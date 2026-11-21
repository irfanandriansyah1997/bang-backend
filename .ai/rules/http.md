---
description: Rules for HTTP handlers, middleware, DTOs, and error mapping
globs: internal/delivery/http/**
alwaysApply: false
paths:
  - "internal/delivery/http/**"
---

# HTTP rules

- **Thin handlers:** decode → call one usecase → encode. No business rules, no SQL, no repository calls.
- **DTOs are separate from domain types.** Request/response structs with `json:"..."` tags live here; map to and from domain in small functions (`toCustomerResponse`). Never encode a domain struct directly.
- **Validate format here, rules in domain.** Required fields, email format, file size/type, `limit` range → 422 from the handler. "Can this status change?" belongs to domain.
- **One error mapper.** Domain errors become HTTP status + error code in a single function (`writeError`). Handlers never pick status codes for domain errors themselves.

  | Domain error | Status | `error.code` |
  |---|---|---|
  | `ErrNotFound` | 404 | `not_found` |
  | `ErrForbidden` | 403 | `forbidden` |
  | `ErrInvalidTransition` | 409 | `INVALID_STATUS_TRANSITION` |
  | validation | 422 | `validation_failed` |
  | anything else | 500 | `internal_error` (log it, don't leak details) |

- **Error shape:** always `{ "error": { "code": "...", "message": "..." } }`.
- **Lists use cursor pagination:** `?limit=15&cursor=<opaque>` (default 15, max 20) → `{ "data": [...], "next_cursor": "...", "prev_cursor": "..." }`.
- **Context:** pass `r.Context()` to the usecase. Middleware stores `x_request_id` / SDC with `logger.WithRequestID` / `logger.WithSDC`; handlers don't build their own.
- **Logging:** log a 5xx once, in the error mapper, with `logger.SDC{"error_code": code}`. Don't log request bodies with passwords, tokens, or personal data.
- **Routes:** all under `/api/v1`, registered in one router file. Role checks are middleware, not `if` statements inside handlers.
- **Every endpoint change updates `api/openapi.yaml` in the same PR.**
- **Tests:** `httptest` per handler covering the success status and every error path's status + `error.code`.
