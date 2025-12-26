# Stage 1: Builder
FROM golang:1.25.4-alpine AS builder

# Install git. 
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# -ldflags="-s -w" reduces binary size by stripping debug information
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/api/main.go

# Stage 2: Runner
FROM alpine:latest

# Install necessary runtime packages
# - ca-certificates: Required to make HTTPS requests
# - tzdata: Required for Timezone handling
# - bash: For the entrypoint script
RUN apk add --no-cache ca-certificates tzdata bash

WORKDIR /app

# Create directories for persistent data (SQLite/Media)
RUN mkdir -p /app/media /app/database

# Copy the Pre-built binary from the builder stage
COPY --from=builder /app/main .

# Copy the entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Expose the port (matches the one in .env.docker)
EXPOSE 5005

# Define volumes to tell Docker these are mount points
VOLUME ["/app/media", "/app/database"]

# Set the entrypoint
ENTRYPOINT ["./entrypoint.sh"]