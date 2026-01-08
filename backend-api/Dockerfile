# Stage 1: Build
FROM golang:1.25.5-alpine3.23 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose port 3000
EXPOSE 3000

# Run the application
CMD ["./main"]
