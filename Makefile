TEST_OPTS=--race

.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: lint-n-fix
lint-n-fix: $(eval LINT_OPTS=--fix) lint

.PHONY: lint
lint:
	golangci-lint run -v --timeout 30m --exclude-use-default $(LINT_OPTS)

.PHONY: test
test:
	go test $(TEST_OPTS) ./...

.PHONY: build
build:
	./build.sh github.com/maqdev/go-be-template/config