include .env
export

fmt:
	go fmt ./...

build:
	@mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

docker-up:
	docker-compose -f ./infra/docker-compose.yml up -d --build

docker-down:
	docker-compose -f ./infra/docker-compose.yml down

docker-clean:
	docker-compose -f ./infra/docker-compose.yml down -v

docker-restart:
	$(MAKE) docker-clean
	$(MAKE) docker-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	$(MAKE) migrate-up

sqlc:
	sqlc generate

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path db/migrations -database "$(DB_URL)" down

