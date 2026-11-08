# `internal/delivery/` — ways the outside world talks to the app

Each sub-folder is one "door" into the app:

- `http/` — the REST API (`/api/v1`), used by the web admin and the Android app

If we ever add another door (gRPC, CLI, message queue consumer), it becomes a new sibling folder like `delivery/grpc/`. **Usecases and domain don't change** when a new door is added. That is the whole point of separating them.

See [`http/README.md`](http/README.md) for the rules.
