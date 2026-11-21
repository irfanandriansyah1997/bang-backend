.PHONY: check lint fmt test hooks ai-link

GOLANGCI_LINT ?= golangci-lint

# Run everything CI runs: lint, build, test.
check: lint test

# Strict lint (config: .golangci.yml). Includes gofmt/goimports and go vet.
lint:
	@command -v $(GOLANGCI_LINT) >/dev/null || { echo "golangci-lint not found. Install it with: brew install golangci-lint"; exit 1; }
	$(GOLANGCI_LINT) run ./...

# Fix formatting and import order in place.
fmt:
	@command -v $(GOLANGCI_LINT) >/dev/null || { echo "golangci-lint not found. Install it with: brew install golangci-lint"; exit 1; }
	$(GOLANGCI_LINT) fmt ./...

# Build everything, then run tests with the race detector and coverage.
test:
	go build -o /dev/null ./...
	go test -race -count=1 -cover ./...

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
