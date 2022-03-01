.PHONY: build

build:
	go build -o scanner ./cmd/scanner/main.go

.PHONY: start

start:
	go run ./cmd/scanner/main.go