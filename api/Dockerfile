# Use an official Golang runtime as a parent image
FROM golang:1.24.1-alpine AS builder
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the containe
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./apps/server/.

# Use Alpine for runtime with non-root user
FROM alpine:3.19 AS prod

# Install CA certificates and create non-root user
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Copy the built binary and directories
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/migrations /app/migrations

# Switch to non-root user
USER appuser

# Expose the port your application is listening on
EXPOSE 8000

# Set the command to run your application
ENTRYPOINT ["/app/main"]
