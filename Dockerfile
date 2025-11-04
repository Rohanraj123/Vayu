# =======================================================
# Build stage - Compile the Vayu binary
# =======================================================
FROM golang:1.24-alpine AS builder

# Enable modules and set working dir
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary statically (no CGO)
RUN CGO_ENABLED=0 go build -o vayu ./cmd/vayu

# ============================================
# Runtime stage — minimal final image
# ============================================
FROM alpine:latest

# Create a non-root user for security
RUN adduser -D -g '' vayuuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/vayu /app/vayu

# Use non-root user
USER vayuuser

# Expose port
EXPOSE 8080

# Default command (users can override the config path)
# The image expects a config file to be mounted at /app/config.yaml
CMD ["./vayu", "./config.yaml"]
