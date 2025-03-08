.PHONY: test test-unit test-integration test-all build clean

test-unit:
	go test ./cmd/... ./internal/...

test-integration:
	go test ./tests/integration/... -timeout 30s

test-all: test-unit test-integration

test:
	go test -short ./...

build:
	go build -o jokecli

clean:
	rm -f jokecli
	rm -f coverage.out 