# `.ai/` — single source of AI rules and skills

**Edit rules and skills here only.** Claude Code and Cursor read them through symlinks:

```
.ai/rules/go-style.md          ◀── .claude/rules/go-style.md    (symlink)
                               ◀── .cursor/rules/go-style.mdc   (symlink)
.ai/skills/add-feature/        ◀── .claude/skills/add-feature   (symlink)
                               ◀── .cursor/skills/add-feature   (symlink)
```

Always-on, project-wide rules live in [`../AGENTS.md`](../AGENTS.md). Cursor reads it natively, and Claude imports it from `CLAUDE.md`.

## Add or change a rule

1. Create `.ai/rules/<topic>.md` with **both** frontmatters (each tool ignores the other's fields):

   ```markdown
   ---
   description: HTTP handler rules          # Cursor
   globs: internal/delivery/http/**         # Cursor (comma-separated)
   alwaysApply: false                       # Cursor
   paths:                                   # Claude (YAML list, quoted)
     - "internal/delivery/http/**"
   ---

   # HTTP rules
   - ...
   ```

2. Run `make ai-link`.

Gotchas:
- **Don't start `globs` with `*`** (e.g. `**/*.go`). In YAML a leading `*` means "alias", which breaks the frontmatter for Claude. Use `cmd/**/*.go,internal/**/*.go`, or quote it if Cursor accepts that.
- Keep `globs` and `paths` matching the same files.
- A rule with `alwaysApply: true` and no `paths` is always on in both tools. Usually that belongs in `AGENTS.md` instead.

## Add or change a skill

1. Create `.ai/skills/<skill-name>/SKILL.md` (folder name = `name`):

   ```markdown
   ---
   name: skill-name
   description: What it does and WHEN to use it.
   ---

   # Steps
   1. ...
   ```

2. Run `make ai-link`.

Extra files (references, templates, scripts) go in the same folder.

## `make ai-link`

- Creates or refreshes all symlinks in `.claude/` and `.cursor/`
- Removes symlinks whose source was deleted or renamed
- Never touches real (non-symlink) files, so tool-specific extras can still live directly in `.claude/` or `.cursor/`

Editing an existing file in `.ai/` needs **no** command; both tools see the change immediately. Run `make ai-link` only when you add, rename, or delete a file.
