# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set environment for SQLite compilation
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with SQLite build tags for musl
RUN go build -tags "sqlite_omit_load_extension" -o notifypipe ./cmd/notifypipe

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

# Copy binary from builder
COPY --from=builder /build/notifypipe .

# Copy static files
COPY --from=builder /build/web ./web

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./notifypipe"]
