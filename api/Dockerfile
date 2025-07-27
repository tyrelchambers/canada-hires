# Use an official Golang runtime as a parent image
FROM golang:1.24.1-alpine AS builder
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./

# Use Alpine for runtime with non-root user
FROM alpine:3.19 AS prod

# Install CA certificates, Node.js, npm, Chromium and create non-root user
RUN apk --no-cache add \
    ca-certificates \
    nodejs \
    npm \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ttf-freefont \
    && addgroup -g 1001 -S appgroup \
    && adduser -S appuser -u 1001 -G appgroup

# Set Puppeteer to use the installed Chromium
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Set working directory and install puppeteer
WORKDIR /app
RUN npm init -y && npm install puppeteer

# Copy the built binary and directories
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/migrations /app/migrations
COPY --from=builder --chown=appuser:appgroup /app/scraper /app/scraper

# Install scraper dependencies
WORKDIR /app/scraper
RUN npm install

# Return to main app directory
WORKDIR /app

# Change ownership of npm files to appuser
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the port your application is listening on
EXPOSE 8000

# Set the command to run your application
ENTRYPOINT ["/app/main"]