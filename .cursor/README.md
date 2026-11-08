# `.cursor/` — Cursor configuration

**Shared rules and skills are not edited here.** Their source is [`../.ai/`](../.ai/README.md). The entries below are symlinks created by `make ai-link`.

```
AGENTS.md                  (repo root) shared always-on rules; Cursor reads it automatically
.cursor/
├── rules/
│   ├── go-style.mdc       → ../../.ai/rules/go-style.md
│   └── database.mdc       → ../../.ai/rules/database.md
└── skills/
    ├── add-feature                    → ../../.ai/skills/add-feature
    └── atomic-conventional-commits    → ../../.ai/skills/atomic-conventional-commits
```

Rule symlinks are named `.mdc` because Cursor ignores plain `.md` files in `rules/`. The target is the same `.md` file Claude uses.

## Cursor-only things

Put a real (non-symlink) `.mdc` file or skill folder directly here when it should apply to Cursor only. `make ai-link` never touches it.

## Notes

- Cursor reads `description`, `globs`, and `alwaysApply`, and ignores Claude's `paths`.
- If a symlinked rule doesn't show up in Cursor's context, check Cursor Settings → Rules. If Cursor doesn't follow symlinks, switch `make ai-link` to copy files instead.
