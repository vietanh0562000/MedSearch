# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o medsearch main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/medsearch .

# Copy any other needed files (e.g., config, migrations, etc.)
# COPY --from=builder /app/app/config ./app/config

# Expose the port your API runs on (default 8080)
EXPOSE 8080

# Default command (can be overridden)
CMD ["./medsearch", "api"]