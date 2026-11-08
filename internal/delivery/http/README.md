# `internal/delivery/http/` — the receptionist

The entry point from outside. Its job is to **translate**: HTTP → usecase input, and usecase output → HTTP response. It makes **no business decisions**.

> React analogy: your page/component layer. It takes input, calls logic, shows the result, but the logic itself lives elsewhere.

## ✅ Do

- Router and endpoint registration (`/api/v1/...`)
- Handlers: read path params, query strings, JSON bodies, multipart forms
- Validate **format** (required fields, valid email, file ≤ 5 MB, JPG/PNG only)
- Request/response DTO structs with `json:"..."` tags
- Middleware: `X-Request-ID`, request logging, JWT auth, role check, CORS, panic recovery
- Map domain errors to HTTP: `ErrInvalidTransition` → 409, `ErrNotFound` → 404, validation → 422
- Write errors in the standard shape: `{ "error": { "code": "...", "message": "..." } }`

## ❌ Don't

- Business rules (`if wo.Status == "assigned" ...`)
- Query the database or call repositories directly
- Validate **business** rules (whether a status can change belongs in domain)
- Return domain structs directly as JSON (always map to a response DTO)
- Open DB transactions
- Write long handlers. A handler should be: **decode → call usecase → encode**

## Example

```go
// internal/delivery/http/work_order_handler.go
func (h *WorkOrderHandler) Complete(w http.ResponseWriter, r *http.Request) {
    in, err := parseCompleteForm(r) // read multipart: note, customer_approved, photos[]
    if err != nil {
        writeError(w, err)
        return
    }
    out, err := h.complete.Execute(r.Context(), in)
    if err != nil {
        writeError(w, err) // e.g. ErrInvalidTransition → 409
        return
    }
    writeJSON(w, http.StatusOK, toCompleteResponse(out))
}
```

## Suggested files (later)

```
router.go            all routes in one place
middleware/          request_id, logging, auth, role, cors, recover
errors.go            domain error → HTTP status + error code
response.go          writeJSON / writeError helpers
pagination.go        cursor encode/decode
*_handler.go         one file per resource (work_order, technician, customer, ...)
```
