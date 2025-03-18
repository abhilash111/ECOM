build:
	@go build -o bin/ecom cmd/main.go

.PHONY: build test run migration migrate-up migrate-down

test:
	@go test -v ./...

run:
	@go run cmd/main.go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down