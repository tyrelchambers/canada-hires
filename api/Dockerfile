# Use an official Golang runtime as a parent image
FROM golang:1.24.1-alpine AS builder
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code (includes integrated scraper)
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the main Go application (now includes integrated scraper)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./

# Build the CLI commands
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o job_scrape_cli ./cmd/job_scrape
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o reddit_post_cli ./cmd/reddit_post

# Final runtime stage using Alpine for minimal footprint
FROM golang:1.24-alpine AS prod

# Install runtime dependencies for chromedp and make
RUN apk update && apk add --no-cache \
    ca-certificates \
    chromium \
    make \
    && rm -rf /var/cache/apk/*

# Create non-root user (Alpine Linux compatible)
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -h /home/appuser -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy the built binary and directories
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/job_scrape_cli /app/job_scrape_cli
COPY --from=builder --chown=appuser:appgroup /app/reddit_post_cli /app/reddit_post_cli
COPY --from=builder --chown=appuser:appgroup /app/migrations /app/migrations
COPY --from=builder --chown=appuser:appgroup /app/makefile /app/makefile


# Set environment variable for chromium path
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Switch to non-root user for runtime
USER appuser

# Expose the port your application is listening on
EXPOSE 8000

# Set the command to run your application
ENTRYPOINT ["/app/main"]
