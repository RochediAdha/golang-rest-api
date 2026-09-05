.PHONY: run test tidy seed

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

seed:
	go run ./cmd/seed $(name)
