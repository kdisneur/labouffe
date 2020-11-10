GIT_BIN := git
GIT_SHA := $(shell $(GIT_BIN) rev-parse HEAD)

GO_BIN := go
GO_FMT_BIN := gofmt
GO_LINT_BIN := $(GO_BIN) run ./vendor/golang.org/x/lint/golint
GO_REFLEX_BIN := $(GO_BIN) run ./vendor/github.com/cespare/reflex
GO_STATICCHECK_BIN := $(GO_BIN) run ./vendor/honnef.co/go/tools/cmd/staticcheck

DEPLOY_FOLDER := $(shell mktemp -d /tmp/livesite.XXXX)

$(DEPLOY_FOLDER)/.git:
	$(GIT_BIN) worktree add --force $(DEPLOY_FOLDER) refs/heads/live-site

live-reload:
	@echo "+ $@"
	$(GO_BIN) run ./vendor/github.com/cosmtrek/air

generate-site: $(DEPLOY_FOLDER)/.git
	@echo "+ $@"
	$(GO_BIN) run . -output $(DEPLOY_FOLDER); \
	if [ -z "$$($(GIT_BIN) -C $(DEPLOY_FOLDER) status -s)" ]; then \
		echo "no changes to commit"; \
	else \
		$(GIT_BIN) -C $(DEPLOY_FOLDER) add -A .; \
		$(GIT_BIN) -C $(DEPLOY_FOLDER) commit --message="New version from SHA $(GIT_SHA)"; \
		$(GIT_BIN) -C $(DEPLOY_FOLDER) push; \
		echo "site generated: $(DEPLOY_FOLDER)"; \
	fi;

test: test-style test-unit

test-style: test-fmt test-lint test-vet test-staticcheck

test-unit:
	@echo "+ $@"
	@$(GO_BIN) test -race ./...

test-fmt:
	@echo "+ $@"
	@test -z "$$($(GO_FMT_BIN) -l -e -s main.go ./internal | tee /dev/stderr)" || \
	  ( >&2 echo "=> please format Go code with '$(GO_FMT_BIN) -s -w .'" && false)

test-lint:
	@echo "+ $@"
	@test -z "$$($(GO_LINT_BIN) internal . | tee /dev/stderr )"

test-tidy:
	@echo "+ $@"
	@$(GO_BIN) mod tidy
	@test -z "$$($(GIT_BIN) status --short go.mod go.sum | tee /dev/stderr)" || \
	  ( >&2 echo "=> please tidy the Go modules with '$(GO_BIN) mod tidy'" && false)

test-vet:
	@echo "+ $@"
	@$(GO_BIN) vet ./...

test-staticcheck:
	@echo "+ $@"
	@$(GO_STATICCHECK_BIN) ./...
