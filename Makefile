.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build binary
	go build -o ghashdeep main.go

.PHONY: test-integration
test-integration: ## run integration tests
	go test -v -tags integration ./...

.PHONY: update-deps
update-deps: ## update dependency libraries
	go get -u \
		github.com/cespare/xxhash \
		github.com/lmittmann/tint \
		github.com/oklog/run \
		github.com/spf13/cobra \
		github.com/stretchr/testify \
		github.com/zeebo/blake3
	go mod tidy
