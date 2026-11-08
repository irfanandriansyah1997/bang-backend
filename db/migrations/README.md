# `db/migrations/` — schema history

Numbered SQL files run in order by **goose** to create or change tables. Think of it as version control for the database schema.

## ✅ Do

- One change per file, named with goose's numbering: `00001_create_users.sql`, `00002_create_customers.sql`
- Write both directions in each file:

  ```sql
  -- +goose Up
  CREATE TABLE users (...);

  -- +goose Down
  DROP TABLE users;
  ```

- Make sure `Down` really undoes `Up` (test with up → down → up)
- Add indexes and constraints (unique, foreign keys, checks) in migrations

## ❌ Don't

- **Edit a migration that has already been run.** Create a new migration for any change
- Put seed/test data here (seeding is a separate step)
- Write Go code here

## Commands (via Makefile, later)

```
make migrate-up
make migrate-down
```
