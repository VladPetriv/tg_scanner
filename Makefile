.PHONY: build

build:
	go build -o scanner ./cmd/scanner/main.go

.PHONY: run

run:
	go run ./cmd/scanner/main.go

.PHONY: lint
lint:
	golangci-lint run ./...
