
# Build the application from source
FROM golang:1.23.7-alpine3.21 AS build-stage
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...


# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
