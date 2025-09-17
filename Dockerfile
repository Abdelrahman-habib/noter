# Build stage
FROM golang:1.25.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/snippetbox ./cmd/web

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/snippetbox .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy necessary files
COPY --from=builder /app/db ./db
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/tls ./tls

# Create a non-root user
RUN addgroup -g 1001 -S snippetbox && \
    adduser -S -D -H -u 1001 -h /root -s /sbin/nologin -G snippetbox -g snippetbox snippetbox

# Change ownership
RUN chown -R snippetbox:snippetbox /root

# Switch to non-root user
USER snippetbox

# Expose port
EXPOSE 4444

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider https://localhost:4444/ || exit 1

# Default command
CMD ["./snippetbox"]