.PHONY: build

build:
	go build -o scanner ./cmd/scanner/main.go

.PHONY: run

run:
	go run ./cmd/scanner/main.go

.PHONY: docker

docker:
	docker build -t tg_scanner .
