# `internal/infrastructure/logger/` — structured logging

Sets up Go's built-in `log/slog` to write **JSON** logs to stdout and `logs/app.log`. Filebeat ships that file to Elasticsearch, and you search it in Kibana.

## ✅ Do

- Create the logger with `slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), ...)`
- Use the standard fields from the spec: `request_id`, `client`, `method`, `route`, `status`, `duration_ms`, `user_id`, `role`, `error_code`, `event`
- Provide helpers to attach/read the logger or `request_id` from `context.Context`

## ❌ Don't

- **Log passwords, tokens, or full personal data** (phone numbers, addresses)
- Use `fmt.Println` / `log.Println` for app logs. Always use `slog`
- Log the same error at every layer. Log it once, where it's handled (usually the HTTP error handler)
