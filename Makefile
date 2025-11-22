run:
	go run ./cmd/api

docker-up:
	docker-compose -f ./infra/docker-compose.yml up -d

docker-down:
	docker-compose -f ./infra/docker-compose.yml down