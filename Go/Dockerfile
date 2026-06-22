# --- Stage 1: Build the Go binary ---
FROM golang:1.26-alpine AS builder

# Install git if required for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the binary statically (ensures it runs seamlessly on alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# --- Stage 2: Final lightweight runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy static assets (HTML, JS, CSS) needed by Gin
COPY --from=builder /app/static ./static

# Expose the application port
EXPOSE 3000

# Run the binary
CMD ["./main"]
