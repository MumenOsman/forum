# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Install gcc and musl-dev to compile CGO for SQLite
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o web ./cmd/web

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the generated binary from the builder stage
COPY --from=builder /app/web .

# Copy the static and html directories needed for the UI
COPY --from=builder /app/ui ./ui

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./web"]
