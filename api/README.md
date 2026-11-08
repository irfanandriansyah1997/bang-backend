# `api/` — the API contract

Will hold `openapi.yaml`: the **contract** between this backend and its clients (web admin in React, Android app in Kotlin).

Written **before** the web/Android work starts, so all three sides agree on endpoints, request/response shapes, and error codes.

## ✅ Do

- Describe every `/api/v1` endpoint, including auth (bearer JWT) and which role can call it
- Define shared schemas once: the error shape `{ "error": { "code", "message" } }` and cursor pagination `{ data, next_cursor, prev_cursor }`
- Update the contract **in the same PR** that changes an endpoint
- Lint it (e.g. `redocly lint api/openapi.yaml`)

## ❌ Don't

- Let the contract drift from the real code
- Put implementation details here (SQL, internal struct names)
