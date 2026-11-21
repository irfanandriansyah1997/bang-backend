.PHONY: check lint fmt test tools run dev migrate-up migrate-down sqlc seed hooks ai-link

# One golangci-lint version for everyone: make, the git hook, VS Code, and CI.
# To upgrade, change it here AND in .github/workflows/ci.yml (golangci-lint-action `version:`).
GOLANGCI_LINT_VERSION := v2.15.0
BIN := $(CURDIR)/bin
# The real binary has the version in its name, so a version bump installs the new one.
GOLANGCI_LINT_BIN := $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
# Stable path for VS Code (.vscode/settings.json) and lefthook; a symlink to the binary above.
GOLANGCI_LINT := $(BIN)/golangci-lint

# Pinned version of air (hot reload), installed into ./bin by `make dev`.
AIR_VERSION := v1.67.4
AIR_BIN := $(BIN)/air-$(AIR_VERSION)
AIR := $(BIN)/air

# Load .env (if it exists) into the environment of the next command.
# Values with spaces or & must be quoted in .env.
LOAD_ENV := set -a; [ ! -f .env ] || . ./.env; set +a;

# Run everything CI runs: lint, build, test.
check: lint test

# Strict lint (config: .golangci.yml). Includes gofmt/goimports and go vet.
lint: tools
	$(GOLANGCI_LINT) run ./...

# Fix formatting and import order in place.
fmt: tools
	$(GOLANGCI_LINT) fmt ./...

# Install the pinned golangci-lint into ./bin (once, or after a version bump).
tools: $(GOLANGCI_LINT_BIN)
	@ln -sfn $(notdir $(GOLANGCI_LINT_BIN)) $(GOLANGCI_LINT)

$(GOLANGCI_LINT_BIN):
	@echo "installing golangci-lint $(GOLANGCI_LINT_VERSION) into ./bin"
	@mkdir -p $(BIN)
	GOBIN=$(BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@mv $(BIN)/golangci-lint $@

# Build everything, then run tests with the race detector and coverage.
test:
	go build -o /dev/null ./...
	go test -race -count=1 -cover ./...

# Run the API with the settings from .env. Stop it with Ctrl+C.
run:
	@$(LOAD_ENV) go run ./cmd/api

# Run the API and restart it on every .go change (config: .air.toml).
dev: $(AIR_BIN)
	@ln -sfn $(notdir $(AIR_BIN)) $(AIR)
	@$(LOAD_ENV) $(AIR)

$(AIR_BIN):
	@echo "installing air $(AIR_VERSION) into ./bin"
	@mkdir -p $(BIN)
	@rm -f $(AIR)
	GOBIN=$(BIN) go install github.com/air-verse/air@$(AIR_VERSION)
	@mv $(BIN)/air $@

# Placeholders: they fail on purpose, so nobody thinks they ran.
migrate-up migrate-down:
	@echo "make $@ is not ready yet: migrations come with issue #5"; exit 1

sqlc:
	@echo "make sqlc is not ready yet: it comes with issue #5"; exit 1

seed:
	@echo "make seed is not ready yet: it comes with issue #7"; exit 1

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
