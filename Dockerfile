
# # Build the application from source
# FROM golang:1.23.7-alpine3.21 AS build-stage

# RUN apk add --no-cache \
#     ca-certificates \
#     curl \
#     gcc \
#     musl-dev \
#     && curl -sSL https://www.amazontrust.com/repository/AmazonRootCA1.pem -o /usr/local/share/ca-certificates/AmazonRootCA1.crt \
#     && update-ca-certificates



# WORKDIR /app


# # pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod tidy

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/server/main.go


# FROM alpine:3.21
# # Copy CA certificates from build stage
# COPY --from=builder /etc/ssl/certs /etc/ssl/certs
# COPY --from=builder /api /api



# # Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# # Deploy the application binary into a lean image
# FROM scratch AS build-release-stage
# WORKDIR /

# COPY --from=build-stage /api /api

# EXPOSE 8080
# EXPOSE 6379

# ENTRYPOINT ["/api"]

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

EXPOSE 8080
EXPOSE 6379

ENTRYPOINT ["/api"]
