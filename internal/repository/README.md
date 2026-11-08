# `internal/repository/` — data storage implementations

Each sub-folder implements the interfaces defined in `usecase/` for one storage technology:

- `postgres/` — PostgreSQL via pgx + sqlc

If storage ever changes (another DB, an in-memory version for tests), it's a new sibling folder. **Usecases don't change**, because they only know the interface.

See [`postgres/README.md`](postgres/README.md) for the rules.
