# Start with the Go base image
FROM golang:1.18-alpine as builder

WORKDIR /app

# Cache dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o url-shortener .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/url-shortener .

# Command to run the application
CMD ["./url-shortener"]
