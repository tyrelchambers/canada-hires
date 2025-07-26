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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./

# Use Alpine for runtime with non-root user
FROM alpine:3.19 AS prod

# Install CA certificates, Node.js, npm and create non-root user
RUN apk --no-cache add ca-certificates nodejs npm && \
    addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

RUN apt get update \
    && apt get install -y wget gnupg \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt get update \
    && apt get install -y google-chrome-stable fonts-ipafont-gothic fonts-wqy-zenhei fonts-thai-tlwg fonts-kacst fonts-freefont-ttf libxss1 \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

RUN npm init -y &&  \
    npm i puppeteer \
    # Add user so we don't need --no-sandbox.
    # same layer as npm install to keep re-chowned files from using up several hundred MBs more space
    && groupadd -r appgroup && useradd -r -g appuser -G audio,video appuser \
    && mkdir -p /home/appuser/Downloads \
    && chown -R appuser:appuser /home/appuser \
    && chown -R appuser:appuser /node_modules \
    && chown -R appuser:appuser /package.json \
    && chown -R appuser:appuser /package-lock.json

# Copy the built binary and directories
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/migrations /app/migrations
COPY --from=builder --chown=appuser:appgroup /app/scraper /app/scraper

# Switch to non-root user
USER appuser

# Expose the port your application is listening on
EXPOSE 8000

# Set the command to run your application
ENTRYPOINT ["/app/main"]
