build:
	@go build -o bin/ecom cmd/server/main.go

test:
	@go test -v ./...
run:
	@go run cmd/server/main.go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down