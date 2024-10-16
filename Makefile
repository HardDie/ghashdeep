.PHONY: build
build:
	go build -o ghashdeep main.go

.PHONY: test-integration
test-integration:
	go test -v -tags integration ./...
