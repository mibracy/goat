# Build Stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o goat .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/goat .

# Copy the index.html file
COPY index.html .

# Expose the port the application listens on
EXPOSE 8420

# Command to run the application
CMD ["./goat"]
