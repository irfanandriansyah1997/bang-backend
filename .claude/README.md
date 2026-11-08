# `.claude/` — Claude Code configuration

**Shared rules and skills are not edited here.** Their source is [`../.ai/`](../.ai/README.md). The entries below are symlinks created by `make ai-link`.

```
CLAUDE.md                  (repo root) always loaded; imports @AGENTS.md
.claude/
├── rules/
│   ├── go-style.md        → ../../.ai/rules/go-style.md
│   └── database.md        → ../../.ai/rules/database.md
└── skills/
    ├── add-feature                    → ../../.ai/skills/add-feature
    └── atomic-conventional-commits    → ../../.ai/skills/atomic-conventional-commits
```

## Claude-only things

Put a real (non-symlink) file or folder directly here when it should apply to Claude only. `make ai-link` never touches it.

## Notes

- Claude reads only the `paths` field from a rule's frontmatter.
- Every `.md` in `rules/` is loaded as a rule, so don't put a README in there.
- Run `/memory` in Claude Code to see which rules are loaded.
