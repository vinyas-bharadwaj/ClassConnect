# Stage 1 — Build the Go binary
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates openssl

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

# Generate SSL certificate and key
RUN openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost" && \
    chmod 644 server.key server.crt

# Build the Go API
RUN go build -ldflags="-s -w" -o server ./cmd/api

# Stage 2 — Final ultra-slim image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder stage
COPY --from=builder /app/server /app/server

# Copy SSL certificate and key
COPY --from=builder /app/server.crt /app/server.crt
COPY --from=builder /app/server.key /app/server.key

# Expose your Go API port
EXPOSE 3000

# Run the server
CMD ["/app/server"]
