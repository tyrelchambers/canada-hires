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

# Download Playwright browsers in a separate stage
FROM golang:1.24.1-bullseye AS playwright-installer
RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libwayland-client0 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxrandr2 \
    libxss1 \
    libxtst6 \
    wget \
    xvfb

WORKDIR /playwright-setup
# Install the exact version from go.mod
COPY go.mod go.sum ./
RUN go mod download
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /playwright-setup/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER}
# Final runtime stage
FROM ubuntu:22.04 AS prod

# Install only runtime dependencies for Playwright
RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libwayland-client0 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxrandr2 \
    libxss1 \
    libxtst6 \
    xvfb \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -m appuser

# Set working directory
WORKDIR /app

# Copy the built binary and directories
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/migrations /app/migrations

# Copy Playwright browsers and drivers from installer stage
COPY --from=playwright-installer /go/bin/playwright /app

RUN apt-get update && apt-get install -y ca-certificates \
    # Install dependencies and all browsers (or specify one)
    && ./playwright install --with-deps \
    && rm -rf /var/lib/apt/lists/*


# Switch to non-root user
USER appuser

# Expose the port your application is listening on
EXPOSE 8000

# Set the command to run your application
ENTRYPOINT ["/app/main"]
