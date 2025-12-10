# Stage 1 — Build the Go binary
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Enable static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

RUN go mod download

# Copy the entire project
COPY . .

# Build the Go API
RUN go build -ldflags="-s -w" -o server ./cmd/api

# Stage 2 — Final ultra-slim image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder stage
COPY --from=builder /app/server /server

# Expose your Go API port
EXPOSE 3000

# Run the server
CMD ["/server"]
