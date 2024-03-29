TEST_OPTS=--race
LINT_OPTS=

LINTER_VERSION=v1.57.2
SQLC_VERSION=v1.26.0
PROTOBUF_DOCKER=jaegertracing/protobuf:v0.5.0

USER_ID = $(shell id -u)
GROUP_ID = $(shell id -g)

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINTER_VERSION)
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: build
codegen:
	@go generate ./...

.PHONY: proto
proto:
	@mkdir -p ./gen/proto
	@docker run --rm -v "${PWD}":"/data/" -w "/data/" --user "$(USER_ID):$(GROUP_ID)" "${PROTOBUF_DOCKER}" --go_out="./gen/proto" --proto_path "./" "./proto/**/*.proto"

.PHONY: build
build:
	@./build.sh github.com/maqdev/go-be-template/config

.PHONY: build-all
build-all: codegen sqlc build

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
	@docker-compose run migrations-testdata # this will create db and run migrations

.PHONY: deinit
deinit:
	@docker-compose down

.PHONY: reinit
reinit: deinit init

.PHONY: clean
clean:
	rm -rf ./gen/*

#.PHONY: mocks
#mocks:
	#mockery --dir=abc --name=ABC --output ./gen/mocks/abc --recursive --with-expecter
