.PHONY: check hooks ai-link

# Run the same checks as the pre-commit hook, on demand.
check:
	@files=$$(gofmt -l .); if [ -n "$$files" ]; then echo "✗ Not formatted. Run: gofmt -w ."; echo "$$files"; exit 1; fi
	go vet ./...
	go build -o /dev/null ./...
	go test ./...

# Install git hooks from lefthook.yml (once per clone). Requires: brew install lefthook
hooks:
	@command -v lefthook >/dev/null || { echo "lefthook not found. Install it with: brew install lefthook"; exit 1; }
	lefthook install

# Link the single source of AI rules/skills in .ai/ into .claude/ and .cursor/.
# Run after adding, renaming, or deleting a file in .ai/rules or .ai/skills.
ai-link:
	@mkdir -p .claude/rules .claude/skills .cursor/rules .cursor/skills
	@# remove dangling symlinks (sources that were deleted or renamed)
	@find .claude/rules .claude/skills .cursor/rules .cursor/skills -maxdepth 1 -type l ! -exec test -e {} \; -print -delete
	@for f in .ai/rules/*.md; do \
		n=$$(basename $$f .md); \
		ln -sfn ../../.ai/rules/$$n.md .claude/rules/$$n.md; \
		ln -sfn ../../.ai/rules/$$n.md .cursor/rules/$$n.mdc; \
		echo "rule  $$n"; \
	done
	@for d in .ai/skills/*/; do \
		n=$$(basename $$d); \
		ln -sfn ../../.ai/skills/$$n .claude/skills/$$n; \
		ln -sfn ../../.ai/skills/$$n .cursor/skills/$$n; \
		echo "skill $$n"; \
	done
