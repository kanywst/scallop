FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o scallop ./cmd/scallop

# Create a minimal image
FROM alpine:latest

# Install Docker client
RUN apk add --no-cache docker-cli

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/scallop /app/scallop

# Set the entrypoint
ENTRYPOINT ["/app/scallop"]

# Default command
CMD ["--help"]
