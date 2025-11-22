include .env
export

build:
	@mkdir -p bin
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

docker-up:
	docker-compose -f ./infra/docker-compose.yml up -d

docker-down:
	docker-compose -f ./infra/docker-compose.yml down

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
