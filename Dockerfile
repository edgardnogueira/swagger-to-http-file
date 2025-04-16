FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:3.19

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/swagger-to-http-file /usr/local/bin/

# Create a non-root user and switch to it
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set entrypoint
ENTRYPOINT ["swagger-to-http-file"]

# Default command
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="Swagger to HTTP File Converter"
LABEL org.opencontainers.image.description="A CLI tool that converts Swagger/OpenAPI JSON documents into .http files for easy API testing"
LABEL org.opencontainers.image.source="https://github.com/edgardnogueira/swagger-to-http-file"
LABEL org.opencontainers.image.licenses="MIT"
