.PHONY: build

build:
	go build -o scanner ./cmd/scanner/main.go

.PHONY: start

start:
	go run ./cmd/scanner/main.go

.PHONY: migrate_up

migrate_up:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose down

.PHONY: test

test:
	go test -v ./...

.PHONY: mock

mock:
	cd ./internal/service/; go generate;
	cd ./internal/store/; go generate; 