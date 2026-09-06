.PHONY: run test tidy seed docker-up docker-down docker-seed

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

seed:
	go run ./cmd/seed $(name)

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-seed:
	docker compose run --rm --entrypoint /app/seed api $(name)
