# Build the application from source
FROM golang:1.23.7-alpine3.21 AS build-stage

RUN apk add --no-cache \
    ca-certificates \
    curl \
    gcc \
    musl-dev \
    && curl -sSL https://www.amazontrust.com/repository/AmazonRootCA1.pem -o /usr/local/share/ca-certificates/AmazonRootCA1.crt \
    && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/server/main.go


# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...


# Final image for running the app
FROM alpine:3.21 AS build-release-stage

# Install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /

# Copy the binary and certificates
COPY --from=build-stage /api /api
COPY --from=build-stage /etc/ssl/certs /etc/ssl/certs
# Copy the dev.env file from /infra directory
COPY --from=build-stage /app/infra/dev.env /infra/dev.env


EXPOSE 8080
EXPOSE 6379

ENTRYPOINT ["/api"]
