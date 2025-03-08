# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mailbox-api ./main.go

# Final stage
FROM alpine:3.16

WORKDIR /app

# Install required packages
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Copy binary from builder stage
COPY --from=builder /app/mailbox-api .

# Copy migrations and data
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/data ./data

# Copy and set permissions for entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Create directory for logs
RUN mkdir -p /var/log/mailbox-api

# Expose port
EXPOSE 8080

# Use entrypoint script
ENTRYPOINT ["docker-entrypoint.sh"]

# Command to run
CMD ["./mailbox-api"]