# Dockerfile

# Use the official Go image as the base image
FROM golang:1.24-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Install dockerize (a wait-for-it utility)
# This command directly installs dockerize into /usr/local/bin/
RUN apk add --no-cache curl \
    && wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-alpine-linux-amd64-v0.6.1.tar.gz \
    && tar -xzf dockerize-alpine-linux-amd64-v0.6.1.tar.gz -C /usr/local/bin \
    && rm dockerize-alpine-linux-amd64-v0.6.1.tar.gz


# --- REST OF YOUR DOCKERFILE REMAINS THE SAME ---

# Copy go.mod and go.sum to the workspace and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o my-go-crud-app ./main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy dockerize from the builder stage
COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/dockerize

# Copy the compiled binary from the builder stage
COPY --from=builder /app/my-go-crud-app .

# Expose the port your application listens on
EXPOSE 8080

# Command to run the executable, using dockerize to wait for postgres
CMD ["dockerize", "-wait", "tcp://postgresDB:5432", "-timeout", "30s", "./my-go-crud-app"]