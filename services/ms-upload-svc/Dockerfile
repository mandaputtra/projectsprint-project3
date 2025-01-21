# Use the official Golang image for ARM64
FROM --platform=linux/arm64 golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary for ARM64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ms-upload-svc .

# Use a minimal Alpine image for the final stage
FROM alpine:latest AS final

# Copy the binary from the builder stage
COPY --from=builder /app/ms-upload-svc /ms-upload-svc

# Set the entry point (optional)
ENTRYPOINT ["/ms-upload-svc"]