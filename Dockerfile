# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with version info
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-s -w -X main.version=docker' \
    -o gowsay .

# Final stage
FROM scratch

# Copy CA certificates for HTTPS (if needed in future)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /build/gowsay /gowsay

# Expose default port
EXPOSE 9000

# Set default command (server mode)
ENTRYPOINT ["/gowsay"]
CMD ["serve"]
