include ./configs/.config.env

.PHONY: build

build:
	go build -o scanner ./cmd/scanner/main.go

.PHONY: run

run:
	go run ./cmd/scanner/main.go

.PHONY: migrate_up

migrate_up:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose down

.PHONY: test

test:
	go test -v ./...

.PHONY: mock

mock:
	cd ./internal/service/; go generate;
	cd ./internal/store/; go generate; 

.PHONY: docker

docker:
	docker build -t tg_scanner .
