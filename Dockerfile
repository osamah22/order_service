# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Enable cgo
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Install build dependencies for cgo
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Install swag for swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs
RUN swag init -g ./cmd/api/main.go -o ./docs

# Build the binary
RUN go build -o order-service ./cmd/api

# Stage 2: Minimal runtime
FROM alpine:latest

# Install sqlite3 library for runtime
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy binary and docs
COPY --from=builder /app/order-service .
COPY --from=builder /app/docs ./docs

EXPOSE 7070

CMD ["./order-service"]
