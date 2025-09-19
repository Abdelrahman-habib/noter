# Build stage
FROM golang:1.25.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata make bash curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum go.tool.mod ./

# Download dependencies
RUN go mod download

# Download tool dependencies (for development tools)
RUN go mod download -modfile=go.tool.mod

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Build the application
RUN go build -o bin/noter ./cmd/web

# Final stage
FROM golang:1.25.1-alpine

# Install runtime dependencies
RUN apk --no-cache add ca-certificates bash curl

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/noter ./
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy necessary files
COPY --from=builder /app/db ./db
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/tls ./tls
COPY --from=builder /app/scripts/entrypoint.sh ./scripts/

# Make entrypoint executable
RUN chmod +x scripts/entrypoint.sh

# Create a non-root user
RUN addgroup -g 1001 -S noter && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G noter -g noter noter

# Change ownership
RUN chown -R noter:noter /app

# Switch to non-root user
USER noter

# Expose port
EXPOSE 4444

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f https://localhost:4444/ || exit 1

# Default command
CMD ["./noter"]