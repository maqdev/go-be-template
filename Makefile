TEST_OPTS=--race
LINT_OPTS=

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: build
codegen:
	@go generate ./...

.PHONY: build
build:
	@./build.sh github.com/maqdev/go-be-template/config

.PHONY: test
test:
	@go test $(TEST_OPTS) ./...

.PHONY: lint
lint:
	@golangci-lint run -v --timeout 30m --exclude-use-default $(LINT_OPTS)

.PHONY: _enable_lint_fix
_enable_lint_fix:
	@$(eval LINT_OPTS=--fix)

.PHONY: lint-n-fix
lint-n-fix: _enable_lint_fix lint

.PHONY: init
init:
	@docker-compose up -d

.PHONY: reinit
reinit:
	@docker-compose up -d --force-recreate
