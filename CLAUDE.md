# CLAUDE.md

Shared project rules live in AGENTS.md (also used by Cursor):

@AGENTS.md

## Claude Code specifics

- Scoped rules and skills: single source in `.ai/`, symlinked into `.claude/rules/` and `.claude/skills/` by `make ai-link` (see `.ai/README.md`).
- Claude-only skills/rules can live directly in `.claude/` as real files.
